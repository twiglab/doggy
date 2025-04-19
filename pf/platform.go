package pf

import (
	"context"
	"log"
	"time"

	"github.com/twiglab/doggy/holo"
)

const (
	HUMMAN_DENSITY = 12
	HUMMAN_COUNT   = 15
)

type Config struct {
	MetadataURL string
	Address     string
	Port        int

	NotMetaAutoSub bool
}

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

type DeviceResolver interface {
	Resolve(ctx context.Context, data holo.DeviceAutoRegisterData) (*holo.Device, error)
}

type SimpleProcess struct {
	Username string
	Password string
}

func NewSimpleProcess(user, pwd string) *SimpleProcess {
	return &SimpleProcess{Username: user, Password: pwd}
}

func (d *SimpleProcess) AutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	now := time.Now().Format(time.RFC3339Nano)
	log.Printf("auto reg sn = %s, ip = %s, time = %s\n", data.SerialNumber, data.IpAddr, now)
	return nil
}

func (d *SimpleProcess) Resolve(ctx context.Context, data holo.DeviceAutoRegisterData) (*holo.Device, error) {
	return holo.OpenDevice(data.IpAddr, d.Username, d.Password)
}

func (d *SimpleProcess) HandleCount(ctx context.Context, common holo.Common, target holo.HumanMix) error {
	start := holo.MilliToTime(target.StartTime, target.TimeZone).Format(time.RFC3339Nano)
	end := holo.MilliToTime(target.EndTime, target.TimeZone).Format(time.RFC3339Nano)
	log.Printf("count in = %d, out = %d, start = %s, end = %s, type = %d\n", target.HumanCountIn, target.HumanCountOut, start, end, target.TargetType)
	return nil
}

func (d *SimpleProcess) HandleDensity(ctx context.Context, common holo.Common, target holo.HumanMix) error {
	now := time.Now().Format(time.RFC3339Nano)
	log.Printf("density count = %d, ration = %d, type = %d, time = %s\n", target.HumanCount, target.AreaRatio, target.TargetType, now)
	return nil
}

type Handle struct {
	Conf Config

	Resolver       DeviceResolver
	Register       DeviceRegister
	CountHandler   CountHandler
	DensityHandler DensityHandler
}

func (h *Handle) metaAutoSub(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	device, err := h.Resolver.Resolve(ctx, data)
	if err != nil {
		return err
	}
	defer device.Close()

	subscriptions, err := device.GetMetadataSubscription(ctx)
	if err != nil {
		return err
	}

	if !subscriptions.IsEmpty() {
		return nil
	}

	resp, err := device.PostMetadataSubscription(ctx, holo.SubscriptionReq{
		Address:     h.Conf.Address,
		Port:        h.Conf.Port,
		TimeOut:     0,
		HttpsEnable: 1,
		MetadataURL: h.Conf.MetadataURL,
	})

	if err != nil {
		return err
	}

	if resp.IsErr() {
		return resp
	}

	return nil
}

func (h *Handle) HandleAutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	if err := h.Register.AutoRegister(ctx, data); err != nil {
		return err
	}

	if !h.Conf.NotMetaAutoSub {
		return h.metaAutoSub(ctx, data)
	}
	return nil
}

func (h *Handle) HandleMetadata(ctx context.Context, data holo.MetadataObjectUpload) error {
	for _, target := range data.MetadataObject.TargetList {
		if target.TargetType == HUMMAN_DENSITY {
			return h.DensityHandler.HandleDensity(ctx, data.MetadataObject.Common, target)
		}
		if target.TargetType == HUMMAN_COUNT {
			return h.CountHandler.HandleCount(ctx, data.MetadataObject.Common, target)
		}
	}
	return nil
}
