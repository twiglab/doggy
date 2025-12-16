package pf

import "time"

//go:generate go tool msgp
type Channel struct {
	SN      string    `msg:"s"`
	IpAddr  string    `msg:"p"`
	RegTime time.Time `msg:"r"`

	UUID string `msg:"i"`
	Code string `msg:"c"`

	X string `msg:"x"`
	Y string `msg:"y"`
	Z string `msg:"z"`
}
