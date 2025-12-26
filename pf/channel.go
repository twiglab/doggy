package pf

import "time"

//go:generate go tool msgp
type ChannelExtra struct {
	SN      string    `msg:"s"`
	IpAddr  string    `msg:"i"`

	RegTime time.Time `msg:"r"`

	UUID string `msg:"u"`
	Code string `msg:"c"`

	X string `msg:"x"`
	Y string `msg:"y"`
	Z string `msg:"z"`
}
