package oc

import (
	"context"

	"github.com/imroc/req/v3"
)

type SumArgs struct {
	Table string   `json:"table omitempty"`
	Start int64    `json:"starte"`
	End   int64    `json:"end"`
	IDs   []string `json:"ids"`
	UUIDs []string `json:"uuids omitempty"`
}

type SumReply struct {
	Total int64 `json:"total"`
}

func ocURL(addr string) string {
	url := "https://" + addr + "/out"
	return url
}

type Client struct {
	baseURL string
	client  *req.Client
}

func New(addr string) *Client {
	base := ocURL(addr)

	c := req.C().
		EnableInsecureSkipVerify().
		SetCommonRetryCount(3).
		SetBaseURL(base)

	return &Client{
		client:  c,
		baseURL: base,
	}
}

func (c *Client) Sum(ctx context.Context, args *SumArgs) (reply *SumReply, err error) {
	_, err = c.client.R().
		SetContext(ctx).
		SetBody(args).
		SetSuccessResult(&reply).
		Post("/sum")
	return
}
