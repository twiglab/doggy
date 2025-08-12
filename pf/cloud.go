package pf

import (
	"context"

	"github.com/twiglab/doggy/holo"
	"resty.dev/v3"
)

type DataHandler interface {
	HandleData(ctx context.Context, data UploadeData) error
}

type UploadeData struct {
	Common holo.Common
	Target holo.HumanMix
	Tenant Tenant
}

type CloudAction struct {
	client  *resty.Client
	baseURL string
	key     string
}

func NewCloudAction(key, baseURL string) *CloudAction {
	clinet := resty.New().
		SetBaseURL(baseURL).
		SetAuthToken(key)

	return &CloudAction{
		client:  clinet,
		baseURL: baseURL,
		key:     key,
	}
}

func (c *CloudAction) HandleData(ctx context.Context, data UploadeData) error {
	upload := holo.MetadataObjectUpload{
		MetadataObject: holo.MetadataObject{
			Common:     data.Common,
			TargetList: []holo.HumanMix{data.Target},
		},
	}
	r := c.client.R().
		SetContext(ctx).
		SetBody(upload)
	_, err := r.Post("/cloud")
	return err
}
