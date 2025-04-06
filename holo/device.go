package holo

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
