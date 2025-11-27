package holo

import (
	"fmt"
)

// 1.2.5
type ResponseStatus struct {
	RequestUrl   string `json:"RequestURL"`
	StatusCode   int    `json:"StatusCode"`
	StatusString string `json:"StatusString"`
	ID           int    `json:"ID,omitempty"` // for 2.6.x
}

type CommonResponse struct {
	ResponseStatus ResponseStatus `json:"ResponseStatus"`
}

func (c CommonResponse) Err() error {
	if c.ResponseStatus.StatusCode == 0 {
		return nil
	}

	return &CommonError{
		Code: c.ResponseStatus.StatusCode,
		Text: c.ResponseStatus.StatusString,
	}
}

func NewCommonResponse(url string, code int, str string) *CommonResponse {
	return &CommonResponse{
		ResponseStatus: ResponseStatus{
			RequestUrl:   url,
			StatusCode:   code,
			StatusString: str,
		},
	}
}

func CommonResponseOK(url string) *CommonResponse {
	return NewCommonResponse(url, 0, "OK")
}

func CommonResponseFailed(url string) *CommonResponse {
	return NewCommonResponse(url, -1, "FAILED")
}

func CommonResponseFailedText(url, text string) *CommonResponse {
	return NewCommonResponse(url, -1, text)
}

func CommonResponseFailedError(url string, err error) *CommonResponse {
	return NewCommonResponse(url, -1, err.Error())
}

type CommonError struct {
	Code int
	Text string
}

func (ce *CommonError) Error() string {
	return fmt.Sprintf("code = %d, text = %s", ce.Code, ce.Text)
}
