package holo

import (
	"fmt"
	"net/url"
	"strconv"
)

const (
	HUMMAN_DENSITY = 12 // 密度 2.6.7
	HUMMAN_QUEUE   = 13 // 排队长度 2.6.8
	HUMMAN_COUNT   = 15 // 人数 2.6.9
)

func SubReq(s string) (SubscriptionReq, error) {
	var err error
	var u *url.URL

	if u, err = url.ParseRequestURI(s); err != nil {
		return SubscriptionReq{}, err
	}

	port := 80
	if p := u.Port(); p != "" {
		if port, err = strconv.Atoi(p); err != nil {
			return SubscriptionReq{}, err
		}
	}

	https := 0
	if u.Scheme == "https" {
		https = 1
	}

	return SubscriptionReq{
		Address:     u.Hostname(),
		TimeOut:     0,
		Port:        port,
		HttpsEnable: https,
		MetadataURL: s,
	}, nil
}

func CameraURL(addr, path string, ssl bool) string {
	if ssl {
		return "https://" + addr + path
	}
	return "http://" + addr + path
}

type ApiError struct {
	Code int
	Text string
}

func NewApiError(code int, text string) *ApiError {
	return &ApiError{
		Code: code,
		Text: text,
	}
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("code = %d, text = %s", e.Code, e.Text)
}

func CheckErr(comm *CommonResponse, err error) error {
	if err != nil {
		return err
	}
	return comm.Err()
}
