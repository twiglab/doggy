package pf

import (
	"context"
	"log/slog"
	"time"

	"github.com/twiglab/doggy/holo"
)

type DeviceRegister interface {
	// 对应 2.1.4 自动注册
	AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error
}

type CountHandler interface {
	// 对应 2.6.9 人数上报(过线检测，type = 15)
	HandleCount(ctx context.Context, common holo.Common, data holo.HumanMix) error
}

type DensityHandler interface {
	// 对应2.6.7 密度上报(密度检查，type = 12)
	HandleDensity(ctx context.Context, common holo.Common, data holo.HumanMix) error
}

type Option func(*Handle)

func WithCountHandler(h CountHandler) Option {
	return func(c *Handle) {
		if h != nil {
			c.countHandler = h
		}
	}
}

func WithDensityHandler(h DensityHandler) Option {
	return func(c *Handle) {
		if h != nil {
			c.densityHandler = h
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
	countHandler   CountHandler
	densityHandler DensityHandler
	deviceRegister DeviceRegister
	toucher        Cache[string, time.Time]
	cache          Cache[string, Channel]
}

func NewHandle(opts ...Option) *Handle {
	action := &cameraAction{}
	h := &Handle{
		countHandler:   action,
		densityHandler: action,
		deviceRegister: action,
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

	/*
		if err != nil {
			return err
		}
			if ok && time.Since(last) < 90*time.Second {
				slog.Debug("ignore muti reg",
					slog.String("sn", data.SerialNumber),
					slog.String("ipAddr", data.IpAddr),
				)
				return nil
			}
	*/

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
		switch target.TargetType {
		case holo.HUMMAN_DENSITY:
			if err := h.densityHandler.HandleDensity(ctx, common, target); err != nil {
				slog.ErrorContext(ctx, "HandleMetadata", slog.Int("targetType", target.TargetType), slog.String("errText", err.Error()))
			}
		case holo.HUMMAN_COUNT:
			if err := h.countHandler.HandleCount(ctx, common, target); err != nil {
				slog.ErrorContext(ctx, "HandleMetadata", slog.Int("targetType", target.TargetType), slog.String("errText", err.Error()))
			}
		default:
			slog.ErrorContext(ctx, "HandleMetadata", slog.Int("targetType", target.TargetType), slog.String("errText", "unsupport type"))
		}
	}
	return nil
}

type cameraAction struct{}

func (d *cameraAction) AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	slog.DebugContext(ctx, "receive reg data",
		slog.String("module", "cameraAction"),
		slog.String("sn", data.SerialNumber),
		slog.String("addr", data.IpAddr))
	return nil
}

func (d *cameraAction) HandleCount(ctx context.Context, common holo.Common, target holo.HumanMix) error {
	start := holo.MilliToTime(target.StartTime, target.TimeZone)
	end := holo.MilliToTime(target.EndTime, target.TimeZone)

	c := slog.Group("common", slog.String("uuid", common.UUID), slog.String("deviceID", common.DeviceID))
	da := slog.Group("data", slog.Int("in", target.HumanCountIn), slog.Int("out", target.HumanCountOut))
	t := slog.Group("time", slog.Time("start", start), slog.Time("end", end))

	slog.DebugContext(ctx, "HandleCount", slog.Int("targetType", target.TargetType), c, da, t)

	return nil
}

func (d *cameraAction) HandleDensity(ctx context.Context, common holo.Common, target holo.HumanMix) error {
	c := slog.Group("common", slog.String("uuid", common.UUID), slog.String("deviceID", common.DeviceID))
	da := slog.Group("data", slog.Int("count", target.HumanCount), slog.Int("areaRatio", target.AreaRatio))

	slog.DebugContext(ctx, "HandleDensity", slog.Int("targetType", target.TargetType), c, da)

	return nil
}
