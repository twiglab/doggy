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
	token   string
}

func NewCloudAction(token, baseURL string) *CloudAction {
	clinet := resty.New().
		SetBaseURL(baseURL).
		SetAuthToken(token)

	return &CloudAction{
		client:  clinet,
		baseURL: baseURL,
		token:   token,
	}
}

func (c *CloudAction) HandleData(ctx context.Context, x UploadeData) error {
	data := holo.MetadataObjectUpload{
		MetadataObject: holo.MetadataObject{
			Common:     x.Common,
			TargetList: []holo.HumanMix{x.Target},
		},
	}
	r := c.client.R().
		SetContext(ctx).
		SetBody(data)
	_, err := r.Post("/cloud")
	return err
}
