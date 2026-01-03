package pf

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/hx"
	"github.com/twiglab/doggy/pkg/human"
)

type DataHandler interface {
	HandleData(ctx context.Context, data human.DataMix) error
}

type DeviceRegister interface {
	// 对应 2.1.4 自动注册
	AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error
}

type Option func(*MainHandle)

func WithDataHandler(h DataHandler) Option {
	return func(c *MainHandle) {
		if h != nil {
			c.dataHandler = h
		}
	}
}

func WithDeviceRegister(h DeviceRegister) Option {
	return func(c *MainHandle) {
		if h != nil {
			c.deviceRegister = h
		}
	}
}

func WithToucher(t TouchCache) Option {
	return func(c *MainHandle) {
		if t != nil {
			c.toucher = t
		}
	}
}

func WithCache(cache ChannelCache) Option {
	return func(c *MainHandle) {
		if cache != nil {
			c.cache = cache
		}
	}
}

type TouchCache = Cache[string, time.Time]
type ChannelCache = Cache[string, ChannelExtra]

type MainHandle struct {
	deviceRegister DeviceRegister
	dataHandler    DataHandler
	toucher        TouchCache
	cache          ChannelCache

	project string
}

func NewMainHandle(project string, opts ...Option) *MainHandle {
	var logA = LogAction{}
	h := &MainHandle{
		deviceRegister: logA,
		dataHandler:    logA,
		toucher:        emptyCache[string, time.Time]{}, // 用于防止摄像头短时间内重复注册
		cache:          emptyCache[string, ChannelExtra]{},

		project: project,
	}

	for _, o := range opts {
		o(h)
	}
	return h
}

func (h *MainHandle) HandleAutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	ch := data.FirstChannel()
	if _, ok, _ := h.toucher.Get(ctx, ch.UUID); ok { // find it, do not register
		return nil
	}

	if err := h.deviceRegister.AutoRegister(ctx, data); err != nil {
		return err
	}
	return h.toucher.Set(ctx, ch.UUID, time.Now())
}

func (h *MainHandle) HandleMetadata(ctx context.Context, data holo.MetadataObjectUpload) error {
	if err := h.toucher.Set(ctx, data.MetadataObject.Common.UUID, time.Now()); err != nil {
		return err
	}

	head := human.Head{
		Project: h.project,
		ID:      data.MetadataObject.Common.UUID,
		Code:    data.MetadataObject.Common.DeviceID,
	}
	if item, ok, _ := h.cache.Get(ctx, data.MetadataObject.Common.UUID); ok {
		head.SN = item.SN
		head.IpAddr = item.IpAddr
		if item.Code != "" {
			head.Code = item.Code
		}
	}

	for _, target := range data.MetadataObject.TargetList {
		data := human.DataMix{
			Head: head,

			Type: dataType(target.TargetType),

			HumanCountIn:  target.HumanCountIn,
			HumanCountOut: target.HumanCountOut,
			BeginTime:     holo.MilliToTime(target.StartTime, target.TimeZone),
			EndTime:       holo.MilliToTime(target.EndTime, target.TimeZone),

			HumanCount: target.HumanCount,
			AreaRatio:  target.AreaRatio,

			QueueTime: target.QueueTime,
		}

		if err := h.dataHandler.HandleData(ctx, data); err != nil {
			slog.ErrorContext(ctx, "HandleData",
				slog.Any("data", data),
				slog.Any("err", err))
		}
	}
	return nil
}

func dataType(t int) string {
	switch t {
	case holo.HUMMAN_COUNT:
		return human.COUNT
	case holo.HUMMAN_DENSITY:
		return human.DENSITY
	case holo.HUMMAN_QUEUE:
		return human.QUEUE
	default:
		return human.UNKNOWN
	}
}

func DeviceAutoRegisterUpload(h *MainHandle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var data holo.DeviceAutoRegisterData
		if err := hx.BindAndClose(r, &data); err != nil {
			_ = hx.JsonTo(http.StatusInternalServerError,
				holo.CommonResponseFailedError(r.URL.Path, err), w)
			return
		}

		ctx := r.Context()

		if err := h.HandleAutoRegister(ctx, data); err != nil {
			slog.ErrorContext(ctx, "handleAutoReg error",
				slog.String("module", "hapi"),
				slog.Any("error", err))

			_ = hx.JsonTo(http.StatusInternalServerError,
				holo.CommonResponseFailedError(r.URL.Path, err), w)
			return
		}

		_ = hx.JsonTo(http.StatusOK, holo.CommonResponseOK(r.URL.Path), w)
	}
}

func MetadataEntryUpload(h *MainHandle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data holo.MetadataObjectUpload
		if err := hx.BindAndClose(r, &data); err != nil {
			slog.ErrorContext(r.Context(), "error-01",
				slog.String("method", "MetadataEntryUpload"),
				slog.Any("error", err))

			http.Error(w, "error-01", http.StatusInternalServerError)
			return
		}

		if err := h.HandleMetadata(r.Context(), data); err != nil {
			slog.ErrorContext(r.Context(), "error-02",
				slog.String("method", "MetadataEntryUpload"),
				slog.Any("error", err))

			http.Error(w, "error-02", http.StatusInternalServerError)
			return
		}

		hx.NoContent(w)
	}
}

func PlatformHandle(h *MainHandle) http.Handler {
	r := chi.NewRouter()
	r.Put("/nat", DeviceAutoRegisterUpload(h))
	r.Put("/1", DeviceAutoRegisterUpload(h))

	r.Post("/upload", MetadataEntryUpload(h))
	r.Post("/2", MetadataEntryUpload(h))
	return r
}
