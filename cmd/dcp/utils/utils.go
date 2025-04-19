package utils

import (
	"fmt"
	"net"
	"strconv"
)

func VerifyAddr(addr string) (host string, port int, err error) {
	var portS string

	if host, portS, err = net.SplitHostPort(addr); err != nil {
		return
	}

	if port, err = strconv.Atoi(portS); err != nil {
		return
	}

	if net.ParseIP(host) == nil {
		err = fmt.Errorf("bad host: %s", host)
	}

	return
}
