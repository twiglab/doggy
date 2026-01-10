package pf

import (
	"context"
)

type ChannelUserData struct {
	SN   string
	UUID string
	Code string

	X string
	Y string
	Z string
}

type DeviceResolver[R any, C any] interface {
	Resolve(ctx context.Context, data R) (C, error)
}
