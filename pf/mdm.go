package pf

import "time"

type CameraSetup struct {
	SN     string `json:"sn"`
	IpAddr string `json:"ip_addr"`
	Last   time.Time
	UUID   string
	User   string
	Pwd    string
}

type CameraPos struct {
	SN       string
	Pos      string // 图纸上标注的名称
	Floor    string // 楼层
	Building string // 建筑物
	Area     string // 区域
}

type CameraUsing struct {
	SN       string
	UUID     string
	DeviceID string
	AlgID    string // 算法id， 15 过线，12 密度
	Name     string // 用于显示的名称
	Memo     string // 业务备注
	BK       string // 业务Key(uuid+algid), 业务code
}
