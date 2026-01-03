package pf

import (
	"context"

	"github.com/twiglab/doggy/holo"
)

type UserDataCache = Cache[string, ChannelUserData]

type HoloCameraSetup struct {
	MainSub holo.SubscriptionReq
	Backups []holo.SubscriptionReq
	Muti    int
}

type HoloCamera struct {
	device  *holo.Device
	setup   HoloCameraSetup
	regData holo.DeviceAutoRegisterData

	cache UserDataCache
}

func (c *HoloCamera) ChannelData(ctx context.Context, ChannelID string) (data ChannelUserData, err error) {
	data, _, err = c.cache.Get(ctx, ChannelID)
	return
}

func (c *HoloCamera) IpAddr() string {
	return c.regData.IpAddr
}

func (c *HoloCamera) SerialNumber() string {
	return c.regData.SerialNumber
}

func (c *HoloCamera) Close() error {
	return c.device.Close()
}

func (c *HoloCamera) Setup(ctx context.Context) error {
	if c.setup.Muti > 0 {
		subs, err := c.device.GetMetadataSubscription(ctx)
		if err != nil {
			return err
		}
		if len(subs.Subscriptions) == 0 {
			res, err := c.device.PostMetadataSubscription(ctx, c.setup.MainSub)
			if err := holo.CheckErr(res, err); err != nil {
				return err
			}
			if c.setup.Muti > 1 {
				for _, sub := range c.setup.Backups {
					res, err := c.device.PostMetadataSubscription(ctx, sub)
					if err := holo.CheckErr(res, err); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

type CameraDB struct {
	User     string
	Pwd      string
	UseHttps bool
}

func NewCamereaDB() *CameraDB {
	return &CameraDB{}
}

func (r *CameraDB) Resolve(ctx context.Context, data holo.DeviceAutoRegisterData) (*HoloCamera, error) {
	c, err := holo.OpenDevice(data.IpAddr, r.User, r.Pwd, r.UseHttps)
	if err != nil {
		return nil, err
	}

	return &HoloCamera{
		device:  c,
		regData: data,
	}, nil
}
