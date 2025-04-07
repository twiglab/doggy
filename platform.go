package doggy

import (
	"context"
	"fmt"
	"time"

	"github.com/twiglab/doggy/holo"
)


type M map[string]any

type HoloHandle struct {
	Conf holo.DeviceCommonConfig
	// log  *slog.Logger
}

func (h *HoloHandle) HandleAutoRegister(ctx context.Context, data holo.DeviceAutoRegisterData) error {
	fmt.Println(time.Now())
	fmt.Println(data.DeviceName)
	fmt.Println(data.Manufacturer)
	fmt.Println(data.DeviceType)
	fmt.Println(data.SerialNumber)
	fmt.Println(data.DeviceVersion)
	fmt.Println(data.IpAddr)
	fmt.Println("------------------------")

	device := holo.NewDevice(data.IpAddr, h.Conf)
	defer device.Close()

	subscriptions, err := device.GetMetadataSubscription(ctx)
	if err != nil {
		return err
	}

	if subscriptions.Size() != 0 {
		return nil
	}

	_, err = device.MetadataSubscription(ctx, holo.MetadataSubscriptionReq{
		TimeOut:     0,
		HttpsEnable: 1,
	})
	return err
}

func (h *HoloHandle) HandleMeta(ctx context.Context, data M) error {
	fmt.Println(data)
	return nil
}
