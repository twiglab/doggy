package holo

type DeviceAutoRegisterData struct {
	DeviceName    string            `json:"DeviceName"`
	Manufacturer  string            `json:"Manufacturer"`
	DeviceType    string            `json:"DeviceType"`
	SerialNumber  string            `json:"SerialNumber"`
	DeviceVersion DeviceVersionData `json:"DeviceVersion"`
	IpAddr        string            `json:"IpAddr"`                // for SEC 9.0.0 +
	ChannelInfo   []Channel         `json:"ChannelInfo,omitempty"` // for SEC 11.0.0 +
}

func (d DeviceAutoRegisterData) FirstChannel() Channel {
	if len(d.ChannelInfo) <= 0 {
		return Channel{}
	}
	return d.ChannelInfo[0]
}

type DeviceVersionData struct {
	Software string `json:"Software"`
	Uboot    string `json:"Uboot"`
	Kernel   string `json:"Kernel"`
	Hardware string `json:"Hardware"`
}

type Channel struct {
	ChannelID int    `json:"ChannelId"`
	UUID      string `json:"UUID"`
	DeviceID  string `json:"DeviceId"`
}

type DeviceID struct {
	UUID     string `json:"UUID"`
	DeviceID string `json:"deviceID"`
}

type DeviceIDList struct {
	IDs []DeviceID `json:"IDs"`
}

type SysBaseInfo struct {
	PlatformType string `json:"platformType"` // 设备平台
	BarCode      string `json:"barCode"`      // 设备BarCode
	BomCode      string `json:"bomCode"`      // 设备BomCode
	DrvCode      string `json:"drvCode"`      // 设备款型名

	Manufacturer   string `json:"manufacturer"`   // 厂商名称。SDC 11.1.0版本新增
	SoftVersion    string `json:"softVersion"`    // 软件包版本信息
	KernelVersion  string `json:"kernelVersion"`  // 内核版本，SDC 11.1.0版本新增
	HardVersion    string `json:"hardVersion"`    // 硬件版本，SDC 11.1.0版本新增
	FullDeviceType string `json:"fullDeviceType"` // 完整设备型号
}
