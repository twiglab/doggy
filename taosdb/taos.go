package taosdb

import (
	"fmt"
)

const (
	TYPE_12 = "12"
	TYPE_15 = "15"

	TAG_UUID      = "uuid"
	TAG_TYPE      = "type"
	TAG_DIVICE_ID = "device_id"

	FIELD_DENSITY_COUNT = "human_count"
	FIELD_DENSITY_RATIO = "human_ratio"

	FIELD_COUNT_IN  = "human_in"
	FIELD_COUNT_OUT = "human_out"

	MA_DENSITY = "density"
	MA_COUNTY  = "count"
)

const (
	TSDB_SML_TIMESTAMP_NANO_SECONDS = "ns"
)

type Config struct {
	Addr     string
	Port     int
	Protocol string
	Username string
	Password string
	DBName   string
}

func db(conf Config) (string, string) {
	dns := fmt.Sprintf("%s:%s@ws(%s:%d)/%s", conf.Username, conf.Password, conf.Addr, conf.Port, conf.DBName)
	return "taosWS", dns
}

func schemalessURL(conf Config) string {
	return fmt.Sprintf("ws://%s:%d", conf.Addr, conf.Port)
}

func bytesToStr(bs []byte) string {
	//return unsafe.String(&bs[0], len(bs))
	return string(bs)
}

func hasCount(in, out int) bool {
	if in == 0 && out == 0 {
		return false
	}
	return true
}
