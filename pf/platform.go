package pf

import (
	"context"
	"fmt"
	"time"

	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/orm"
)

type PlatformConfig struct {
	MetadataURL string
	Address     string
	Port        int
}

type DeviceResolver interface {
	Resolve(ctx context.Context, data holo.DeviceAutoRegisterData) (*holo.Device, error)
}

type DeviceResolve struct {
	Username string
	Password string
}

func NewDeviceResolve(user, pwd string) *DeviceResolve {
	return &DeviceResolve{Username: user, Password: pwd}
}

func (d *DeviceResolve) Resolve(ctx context.Context, data holo.DeviceAutoRegisterData) (*holo.Device, error) {
	return holo.OpenDevice(data.IpAddr, d.Username, d.Password)
}

type M map[string]any

type Handle struct {
	Conf     PlatformConfig
	Resolver DeviceResolver
	Op       *orm.DBOP
}

func (h *Handle) HandleAutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	device, err := h.Resolver.Resolve(ctx, data)
	if err != nil {
		return err
	}

	defer device.Close()

	subscriptions, err := device.GetMetadataSubscription(ctx)
	if err != nil {
		return err
	}

	if subscriptions.Size() != 0 {
		return nil
	}

	_, err = device.PostMetadataSubscription(ctx, holo.SubscriptionReq{
		Address:     h.Conf.Address,
		Port:        h.Conf.Port,
		TimeOut:     0,
		HttpsEnable: 1,
		MetadataURL: h.Conf.MetadataURL,
	})
	return err
}

func (h *Handle) HandleMetadata(ctx context.Context, data holo.MetadataObjectUpload) error {
	for _, target := range data.MetadataObject.TargetList {
		if target.TargetType == 12 {
			return h.handleHuman12(ctx, data.MetadataObject.Common, target)
		}
		if target.TargetType == 15 {
			return h.handleHuman15(ctx, data.MetadataObject.Common, target)
		}
	}
	return nil
}

func (h *Handle) handleHuman15(ctx context.Context, common holo.Common, target holo.HumanMix) error {
	if target.HumanCountIn == 0 && target.HumanCountOut == 0 {
		return nil
	}

	fmt.Println(target.HumanCountIn, target.HumanCountOut)
	fmt.Println(time.UnixMilli(target.StartTime).Format(time.RFC3339), time.UnixMilli(target.EndTime).Format(time.RFC3339))

	return nil
}

func (h *Handle) handleHuman12(ctx context.Context, common holo.Common, target holo.HumanMix) error {
	fmt.Println(target.HumanCount)
	return nil
}
