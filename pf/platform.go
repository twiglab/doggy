package pf

import (
	"context"
	"fmt"
	"time"

	"github.com/twiglab/doggy/holo"
)

func platformURL(base string, path string) string {
	return base + path
}

type PlatformConfig struct {
	BaseURL string
	Address string
	Port    int
}

type DeviceResolver interface {
	Resolve(ctx context.Context, data holo.DeviceAutoRegisterData) (*holo.Device, error)
}

type DeviceResolve struct {
	Username string
	Password string
}

func (d *DeviceResolve) Resolve(ctx context.Context, data holo.DeviceAutoRegisterData) (*holo.Device, error) {
	return holo.OpenDevice(data.IpAddr, d.Username, d.Password)
}

type M map[string]any

type HoloHandle struct {
	Conf     PlatformConfig
	Resolver DeviceResolver
	// log  *slog.Logger
}

func (h *HoloHandle) HandleAutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
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

	_, err = device.MetadataSubscription(ctx, holo.SubscriptionReq{
		Address:     h.Conf.Address,
		Port:        h.Conf.Port,
		TimeOut:     0,
		HttpsEnable: 1,
		//MetadataURL: "/MetadataEntry",
		MetadataURL: platformURL(h.Conf.BaseURL, "/MetadataEntry"),
	})
	return err
}

func (h *HoloHandle) HandleMetadata(ctx context.Context, data holo.MetadataObjectUpload) error {
	for _, target := range data.MetadataObject.TargetList {
		if target.TargetType == 12 {
			return nil
		}
		if target.TargetType == 15 {
			if target.HumanCountIn == 0 && target.HumanCountOut == 0 {
				return nil
			}
			return h.handleHuman15(ctx, data.MetadataObject.Common, target)
		}
	}
	return nil
}

func (h *HoloHandle) handleHuman15(ctx context.Context, common holo.Common, target holo.HumanMix) error {
	fmt.Println(target.HumanCountIn, target.HumanCountOut)
	fmt.Println(time.UnixMilli(target.StartTime).Format(time.RFC3339), time.UnixMilli(target.EndTime).Format(time.RFC3339))
	return nil
}
