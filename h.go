package doggy

import (
	"context"
	"fmt"
)

type M map[string]any

type HoleHandl struct {
}

func (h *HoleHandl) HandleRegister(ctx context.Context, data DeviceRegisterData) error {
	fmt.Println(data)
	return nil
}

func (h *HoleHandl) HandleMeta(ctx context.Context, data M) error {
	fmt.Println(data)
	return nil
}
