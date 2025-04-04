package doggy

type MetaCommonData struct {
	UUID     string
	DeviceID string
}

type HumanCountData struct {
	TargetType           int // 元数据类型，2.6.9取值为15
	HumanCountIn         int
	HumanCountOut        int
	StartTime            int
	EndTime              int
	TimeZone             int
	DayLightSavingOffset int //夏令时偏移（秒）
}

type HumanCountUploadData struct {
	Common     MetaCommonData
	TargetList []HumanCountData
}

type DeviceRegisterData struct {
	DeviceName    string
	Manufacturer  string
	DeviceType    string
	SerialNumver  string
	IpAddr        string
	DeviceVersion DeviceVersionData
}

type DeviceVersionData struct {
	Software string
	Uboot    string
	Kernel   string
	Hardware string
}
