package holo

import "fmt"

// 1.2.5
type CommonResponse struct {
	RequestUrl   string `json:"RequestURL"`
	StatusCode   int    `json:"StatusCode"`
	StatusString string `json:"StatusString"`
}

func (r CommonResponse) Error() string {
	return fmt.Sprintf("url = %s, code = %d, str = %s", r.RequestUrl, r.StatusCode, r.StatusString)
}

func (r CommonResponse) String() string {
	return fmt.Sprintf("url = %s, code = %d, str = %s", r.RequestUrl, r.StatusCode, r.StatusString)
}

func (r CommonResponse) IsErr() bool {
	// 5.2 响应码
	return r.StatusCode != 0
}

type CommonResponseID struct {
	CommonResponse
	ID int `json:"ID,omitempty"`
}

// 2.6.9
type MetaCommonData struct {
	UUID     string `json:"UUID"`
	DeviceID string `json:"deviceID"`
}

type MetaHumanCountData struct {
	TargetType           int `json:"targetType"` // 元数据类型，2.6.9取值为15
	HumanCountIn         int `json:"humanCountIn"`
	HumanCountOut        int `json:"humanCountOut"`
	StartTime            int `json:"startTime"`
	EndTime              int `json:"endTime"`
	TimeZone             int `json:"timeZone"`
	DayLightSavingOffset int `json:"dayLightSavingOffset"` //夏令时偏移（秒）
}

type HumanCountUploadData struct {
	Common     MetaCommonData       `json:"common"`
	TargetList []MetaHumanCountData `json:"targetList"`
}

type DeviceRegisterData struct {
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

// 2.6.1
type MetadataSubscriptionReq struct {
	Address     string `json:"address"`
	Port        int    `json:"port"`
	TimeOut     int    `json:"timeOut"`
	HttpsEnable int    `json:"httpsEnable"`
	MetadataURL string `json:"metadataURL"`
	DigUserName string `json:"digUserName,omitempty"`
	DigUserPwd  string `json:"digUserPwd,omitempty"`
}
