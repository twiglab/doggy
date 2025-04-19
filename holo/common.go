package holo

import "fmt"

// 1.2.5
type CommonResponse struct {
	RequestUrl   string `json:"RequestURL"`
	StatusCode   int    `json:"StatusCode"`
	StatusString string `json:"StatusString"`
}

func (r CommonResponse) Error() string {
	return fmt.Sprintf("url = %s, code = %d, msg = %s", r.RequestUrl, r.StatusCode, r.StatusString)
}

func (r CommonResponse) String() string {
	return fmt.Sprintf("url = %s, code = %d, msg = %s", r.RequestUrl, r.StatusCode, r.StatusString)
}

func (r CommonResponse) IsErr() bool {
	// 5.2 响应码
	return r.StatusCode != 0
}

type CommonResponseID struct {
	CommonResponse
	ID int `json:"ID,omitempty"`
}
