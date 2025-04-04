package doggy

import (
	"context"
	"fmt"
	"log/slog"
)

type HoleHandl struct {
}

func (h *HoleHandl) HandleRegister(ctx context.Context, data DeviceRegisterData) error {
	slog.DebugContext(ctx, "register", "data", data)
	fmt.Println(data)
	return nil
}
