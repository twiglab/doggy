package holo

import (
	"fmt"
	"time"
)

const (
	HUMMAN_DENSITY = 12
	HUMMAN_COUNT   = 15
)

func cameraURL(addr, path string) string {
	url := "https://" + addr + path
	return url
}

func MilliToTime(milli int64, tz int64) time.Time {
	return time.UnixMilli(milli)
}

type ApiError struct {
	Code int
	Text string
	error
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
