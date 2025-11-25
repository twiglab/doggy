package pf

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/twiglab/doggy/holo"
)

var ErrUnimplType = errors.New("unsupport type")

type DeviceRegister interface {
	// 对应 2.1.4 自动注册
	AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error
}

type Option func(*Handle)

func WithDataHandler(h DataHandler) Option {
	return func(c *Handle) {
		if h != nil {
			c.dataHandler = h
		}
	}
}

func WithDeviceRegister(h DeviceRegister) Option {
	return func(c *Handle) {
		if h != nil {
			c.deviceRegister = h
		}
	}
}

func WithToucher(t Cache[string, time.Time]) Option {
	return func(c *Handle) {
		if t != nil {
			c.toucher = t
		}
	}
}

func WithCache(cache Cache[string, Channel]) Option {
	return func(c *Handle) {
		if cache != nil {
			c.cache = cache
		}
	}
}

type Handle struct {
	deviceRegister DeviceRegister
	dataHandler    DataHandler
	toucher        Cache[string, time.Time]
	cache          Cache[string, Channel]
}

func NewHandle(opts ...Option) *Handle {
	action := &action{}
	h := &Handle{
		deviceRegister: action,
		dataHandler:    action,
		toucher:        emptyCache[string, time.Time]{},
		cache:          emptyCache[string, Channel]{},
	}

	for _, o := range opts {
		o(h)
	}
	return h
}

func (h *Handle) HandleAutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	ch := data.FirstChannel()
	if _, ok, _ := h.toucher.Get(ctx, ch.UUID); ok {
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

func (h *Handle) HandleMetadata(ctx context.Context, data holo.MetadataObjectUpload) error {
	if err := h.toucher.Set(ctx, data.MetadataObject.Common.UUID, time.Now()); err != nil {
		return err
	}
	var common = data.MetadataObject.Common
	if item, ok, _ := h.cache.Get(ctx, data.MetadataObject.Common.UUID); ok {
		common.DeviceID = item.Code
	}
	for _, target := range data.MetadataObject.TargetList {
		if err := h.dataHandler.HandleData(ctx, UploadeData{Common: common, Target: target}); err != nil {
			slog.ErrorContext(ctx, "HandleData",
				slog.Any("taget", target),
				slog.Any("common", common),
				slog.Any("err", err))
		}
	}
	return nil
}

type action struct{}

func (d *action) AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	slog.DebugContext(ctx, "receive reg data",
		slog.String("module", "cameraAction"),
		slog.String("sn", data.SerialNumber),
		slog.String("addr", data.IpAddr))
	return nil
}

func (d *action) HandleData(ctx context.Context, data UploadeData) error {
	slog.DebugContext(ctx, "handleData", slog.Any("data", data))
	return nil
}
