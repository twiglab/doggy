package taosdb

import (
	"fmt"
	"unsafe"
)

const (
	TAG_CODE    = "code"
	TAG_PROJECT = "project"
	TAG_SN      = "sn"
	// TAG_TYPE    = "type"
	TAG_UUID = "uuid"

	TAG_X = "x"
	TAG_Y = "y"
	TAG_Z = "z"

	FIELD_DENSITY_COUNT = "human_count"
	FIELD_DENSITY_RATIO = "human_ratio"

	FIELD_QUEUE_COUNT = "human_count"
	FIELD_QUEUE_TIME  = "queue_time"

	FIELD_COUNT_IN  = "human_in"
	FIELD_COUNT_OUT = "human_out"

	MA_DENSITY = "stb_density"
	MA_COUNTY  = "stb_count"
	MA_QUEUE   = "stb_queue"
)

const (
	TSDB_SML_TIMESTAMP_SECONDS       = "s"
	TSDB_SML_TIMESTAMP_MILLI_SECONDS = "ms"
	TSDB_SML_TIMESTAMP_MICRO_SECONDS = "us"
	TSDB_SML_TIMESTAMP_NANO_SECONDS  = "ns"
)

func TaosDSN(user, pwd, addr string, port int, dbname string) (string, string) {
	dns := fmt.Sprintf("%s:%s@ws(%s:%d)/%s", user, pwd, addr, port, dbname)
	return "taosWS", dns
}

func SchemalessURL(addr string, port int) string {
	return fmt.Sprintf("ws://%s:%d", addr, port)
}

func bytesToStr(bs []byte) string {
	return unsafe.String(&bs[0], len(bs))
}
