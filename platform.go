package doggy

import (
	"context"
	"fmt"
	"time"

	"github.com/twiglab/doggy/holo"
)

func subUrl(addr string, port int, path string) string {
	return fmt.Sprintf("https://%s:%d%s", addr, port, path)
}

type PlatformConfig struct {
	BaseURL string
	Address string
	Port    int
}

type DeviceResolve struct {
	Username string
	Password string
}

func (d *DeviceResolve) Resolve(data holo.DeviceAutoRegisterData) (*holo.Device, error) {
	return holo.OpenDevice(data.IpAddr, d.Username, d.Password)
}

type M map[string]any

type HoloHandle struct {
	Conf    PlatformConfig
	Resolve *DeviceResolve
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

	device, err := h.Resolve.Resolve(data)
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

	_, err = device.MetadataSubscription(ctx, holo.MetadataSubscriptionReq{
		Address:     h.Conf.Address,
		Port:        h.Conf.Port,
		TimeOut:     0,
		HttpsEnable: 1,
		MetadataURL: "/MetadataEntry",
	})
	return err
}

func (h *HoloHandle) HandleMeta(ctx context.Context, data M) error {
	fmt.Println(data)
	return nil
}
