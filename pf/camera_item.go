package pf

import "time"

//go:generate msgp
type CameraItem struct {
	SN       string    `msg:"s"`
	IpAddr   string    `msg:"p"`
	LastTime time.Time `msg:"t"`

	UUID string `msg:"i"`
	Code string `msg:"c"`

	User string `msg:"u"`
	Pwd  string `msg:"x"`
}
