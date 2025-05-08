package pf

import (
	"context"
	"log"
	"time"

	"github.com/twiglab/doggy/holo"
)

type DeviceLoader interface {
	All(context.Context) ([]CameraUpload, error)
}

type emptyDeviceLoad struct{}

func (*emptyDeviceLoad) All(ctx context.Context) (cameras []CameraUpload, err error) { return }

var EmptyDeviceLoad = &emptyDeviceLoad{}

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
		c.countHandler = h
	}
}

func WithDensityHandler(h DensityHandler) Option {
	return func(c *Handle) {
		c.densityHandler = h
	}
}

func WithDeviceRegister(h DeviceRegister) Option {
	return func(c *Handle) {
		c.deviceRegister = h
	}
}

type Handle struct {
	countHandler   CountHandler
	densityHandler DensityHandler
	deviceRegister DeviceRegister
}

func NewHandle(opts ...Option) *Handle {
	action := &cameraAction{}
	h := &Handle{
		countHandler:   action,
		densityHandler: action,
		deviceRegister: action,
	}

	for _, o := range opts {
		o(h)
	}
	return h
}

func (h *Handle) HandleAutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	return h.deviceRegister.AutoRegister(ctx, data)
}

func (h *Handle) HandleMetadata(ctx context.Context, data holo.MetadataObjectUpload) error {
	for _, target := range data.MetadataObject.TargetList {
		switch target.TargetType {
		case holo.HUMMAN_DENSITY:
			if err := h.densityHandler.HandleDensity(ctx, data.MetadataObject.Common, target); err != nil {
				log.Println(err)
			}
		case holo.HUMMAN_COUNT:
			if err := h.countHandler.HandleCount(ctx, data.MetadataObject.Common, target); err != nil {
				log.Println(err)
			}
		default:
			log.Println("unsupport type ", target.TargetType)
		}
	}
	return nil
}

type cameraAction struct {
}

func (d *cameraAction) AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	log.Printf("auto reg sn = %s, ip = %s\n", data.SerialNumber, data.IpAddr)
	return nil
}

func (d *cameraAction) HandleCount(ctx context.Context, common holo.Common, target holo.HumanMix) error {
	start := holo.MilliToTime(target.StartTime, target.TimeZone).Format(time.RFC3339Nano)
	end := holo.MilliToTime(target.EndTime, target.TimeZone).Format(time.RFC3339Nano)
	log.Printf("count in = %d, out = %d, start = %s, end = %s, type = %d\n", target.HumanCountIn, target.HumanCountOut, start, end, target.TargetType)
	return nil
}

func (d *cameraAction) HandleDensity(ctx context.Context, common holo.Common, target holo.HumanMix) error {
	now := time.Now().Format(time.RFC3339Nano)
	log.Printf("density count = %d, ration = %d, type = %d, time = %s\n", target.HumanCount, target.AreaRatio, target.TargetType, now)
	return nil
}
