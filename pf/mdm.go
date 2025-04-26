package pf

import "context"

type CameraSetup struct {
	SN     string
	IpAddr string

	// 文档规定最多6个uuid
	UUID1 string
	UUID2 string

	User string
	Pwd  string
}

type CameraPoint struct {
	SN       string
	Pos      string // 图纸上标注的名称
	Floor    string // 楼层
	Building string // 建筑物
	Area     string // 区域
}

type CameraUsing struct {
	SN    string
	UUID  string
	AlgID string // 算法id， 15 过线，12 密度
	Name  string // 用于显示的名称
	BK    string // 业务Key(如果可能请设置成device_id), 业务code
}

type CameraUsingGetter interface {
	GetCameraUsings(context.Context) map[string]CameraUsing
}

type CameraPosGetter interface {
	GetCameraPoints(context.Context) map[string]CameraPoint
}
