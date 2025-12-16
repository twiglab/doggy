package pf

import (
	"context"
)

type Camera interface {
	Setup(ctx context.Context) error
	Close() error

	SerialNumber() string
	IpAddr() string
	/*
		ChannelUseData(channelID string) (ChannelUserData, error)
	*/
}

type DeviceResolver[C Camera, R any] interface {
	Resolve(ctx context.Context, data R) (C, error)
}
