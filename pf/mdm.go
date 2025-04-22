package pf

import "time"

type CameraUplaod struct {
	SN     string
	IpAddr string
	Last   time.Time
	UUID   string
}

type CameraSetup struct {
	SN       string
	Pos      string // 图纸上标注的名称
	Floor    string // 楼层
	Building string // 建筑物
	Area     string // 区域
	NATAddr  string // 按照配置时设置的NAT地址
	User     string
	Pwd      string
}

type CameraUsing struct {
	SN    string
	UUID  string
	AlgID string // 算法id， 15 过线，12 密度
	Name  string // 用于显示的名称
	Memo  string // 业务备注
	BK    string // 业务Key(uuid+algid), 业务code
}
