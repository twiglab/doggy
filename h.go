package doggy

import (
	"context"
	"fmt"
	"time"

	"github.com/twiglab/doggy/holo"
)

type M map[string]any

type HoleHandl struct {
}

func (h *HoleHandl) HandleRegister(ctx context.Context, data holo.DeviceRegisterData) error {
	fmt.Println(time.Now())
	fmt.Println(data.DeviceName)
	fmt.Println(data.Manufacturer)
	fmt.Println(data.DeviceType)
	fmt.Println(data.SerialNumber)
	fmt.Println(data.DeviceVersion)
	fmt.Println(data.IpAddr)
	fmt.Println("------------------------")
	return nil
}

func (h *HoleHandl) HandleMeta(ctx context.Context, data M) error {
	fmt.Println(data)
	return nil
}
