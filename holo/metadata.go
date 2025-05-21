package holo

// 2.6.1
type SubscriptionReq struct {
	Address     string `json:"address"`
	Port        int    `json:"port"`
	TimeOut     int    `json:"timeOut"`
	HttpsEnable int    `json:"httpsEnable"`
	MetadataURL string `json:"metadataURL"`
	DigUserName string `json:"digUserName,omitempty"`
	DigUserPwd  string `json:"digUserPwd,omitempty"`
}

// 2.6.9
type Common struct {
	UUID     string `json:"UUID"`
	DeviceID string `json:"deviceID"`
}

type HumanMix struct {
	TargetType           int   `json:"targetType"` // 元数据类型，2.6.9为15, 2.6.7为12
	HumanCountIn         int   `json:"humanCountIn"`
	HumanCountOut        int   `json:"humanCountOut"`
	StartTime            int64 `json:"startTime"`
	EndTime              int64 `json:"endTime"`
	TimeZone             int64 `json:"timeZone"`
	DayLightSavingOffset int64 `json:"dayLightSavingOffset"` //夏令时偏移（秒）

	HumanCount int `json:"humanCount"`
	AreaRatio  int `json:"areaRatio"`
}

type MetadataObject struct {
	Common     Common     `json:"common"`
	TargetList []HumanMix `json:"targetList"`
}

type MetadataObjectUpload struct {
	MetadataObject MetadataObject `json:"metadataObject"`
}

// 2.6.4

type SubscriptionItem struct {
	ID          int    `json:"id"`
	MetadataURL string `json:"metadataURL"`
}

type Subscripions struct {
	Subscriptions []SubscriptionItem `json:"subscriptions,omitempty"`
}

func (s Subscripions) IsEmpty() bool {
	return len(s.Subscriptions) == 0
}

func (s Subscripions) Size() int {
	return len(s.Subscriptions)
}
