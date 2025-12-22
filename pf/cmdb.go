package pf

import (
	"context"
)

type ChannelUserData struct {
	UUID string
	Code string

	X string
	Y string
	Z string
}

type Camera interface {
	Setup(ctx context.Context) error
	Close() error

	SerialNumber() string
	IpAddr() string
	ChannelData(channelID string) (ChannelUserData, error)
}

type DeviceResolver[R any] interface {
	Resolve(ctx context.Context, data R) (Camera, error)
}

