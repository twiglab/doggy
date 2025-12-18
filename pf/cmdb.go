package pf

import (
	"context"

	"github.com/twiglab/doggy/cmdb"
)

type Camera interface {
	Setup(ctx context.Context) error
	Close() error

	SerialNumber() string
	IpAddr() string
	ChannelData(channelID string) (cmdb.ChannelUserData, error)
}

type DeviceResolver[C Camera, R any] interface {
	Resolve(ctx context.Context, data R) (C, error)
}
