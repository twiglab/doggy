package holo

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

// 2.6.4

type SubscripionItemData struct {
	ID int `json:"id"`
}

type SubscripionsData struct {
	Subscripions []SubscripionItemData `json:"subscriptions"`
}

func (s SubscripionsData) IsEmpty() bool {
	return len(s.Subscripions) == 0
}

func (s SubscripionsData) Size() int {
	return len(s.Subscripions)
}
