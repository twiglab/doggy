package pf

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/pkg/human"
)

var ErrUnimplType = errors.New("unsupport type")

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

func WithToucher(t Cache[string, time.Time]) Option {
	return func(c *MainHandle) {
		if t != nil {
			c.toucher = t
		}
	}
}

func WithCache(cache Cache[string, Channel]) Option {
	return func(c *MainHandle) {
		if cache != nil {
			c.cache = cache
		}
	}
}

type MainHandle struct {
	deviceRegister DeviceRegister
	dataHandler    DataHandler
	toucher        Cache[string, time.Time]
	cache          Cache[string, Channel]

	project string
}

func NewMainHandle(project string, opts ...Option) *MainHandle {
	action := &noopAction{}
	h := &MainHandle{
		deviceRegister: action,
		dataHandler:    action,
		toucher:        emptyCache[string, time.Time]{}, // 用于防止摄像头短时间内重复注册
		cache:          emptyCache[string, Channel]{},

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
		slog.Error("AutoReg error",
			slog.Any("data", data),
			slog.Any("error", err))
		return err
	}
	return h.toucher.Set(ctx, ch.UUID, time.Now())
}

func (h *MainHandle) HandleMetadata(ctx context.Context, data holo.MetadataObjectUpload) error {
	if err := h.toucher.Set(ctx, data.MetadataObject.Common.UUID, time.Now()); err != nil {
		return err
	}

	head := human.Head{
		Project:  h.project,
		UUID:     data.MetadataObject.Common.UUID,
		DeviceID: data.MetadataObject.Common.DeviceID,
	}
	if item, ok, _ := h.cache.Get(ctx, data.MetadataObject.Common.UUID); ok {
		head.SN = item.SN
		head.IpAddr = item.IpAddr
		head.DeviceID = item.Code
	}

	for _, target := range data.MetadataObject.TargetList {
		data := human.DataMix{
			Head: head,

			Type: dataType(target.TargetType),

			HumanCountIn:  target.HumanCountIn,
			HumanCountOut: target.HumanCountOut,
			BeginTime:     human.MilliToTime(target.StartTime, target.TimeZone),
			EndTime:       human.MilliToTime(target.EndTime, target.TimeZone),

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
