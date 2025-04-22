package holo

import "fmt"

type DeviceAutoRegisterData struct {
	DeviceName    string            `json:"DeviceName"`
	Manufacturer  string            `json:"Manufacturer"`
	DeviceType    string            `json:"DeviceType"`
	SerialNumber  string            `json:"SerialNumber"`
	DeviceVersion DeviceVersionData `json:"DeviceVersion"`
	IpAddr        string            `json:"IpAddr"` // for SEC 9.0.0 +
}

type DeviceVersionData struct {
	Software string `json:"Software"`
	Uboot    string `json:"Uboot"`
	Kernel   string `json:"Kernel"`
	Hardware string `json:"Hardware"`
}

type RebootResp struct {
	Code int    `json:"HSErrorCode"`
	Msg  string `json:"HSErrorMsg"`
}

func (r RebootResp) Error() string {
	return fmt.Sprintf("code = %d, msg = %s", r.Code, r.Msg)
}

func (r RebootResp) IsErr() bool {
	return r.Code != 0
}

type DeviceID struct {
	UUID     string `json:"UUID"`
	DeviceID string `json:"deviceID"`
}
type DeviceIDList struct {
	IDs []DeviceID `json:"IDs"`
}
