package taosdb

import (
	"fmt"
	"unsafe"
)

const (
	TYPE_12 = "12"
	TYPE_13 = "13"
	TYPE_15 = "15"

	TAG_UUID      = "uuid"
	TAG_TYPE      = "type"
	TAG_DIVICE_ID = "device_id"
	TAG_TENANT_ID = "tenant_id"

	FIELD_DENSITY_COUNT = "human_count"
	FIELD_DENSITY_RATIO = "human_ratio"

	FIELD_QUEUE_COUNT = "human_count"
	FIELD_QUEUE_TIME  = "queue_time"

	FIELD_COUNT_IN  = "human_in"
	FIELD_COUNT_OUT = "human_out"

	MA_DENSITY = "s_density"
	MA_COUNTY  = "s_count"
	MA_QUEUE   = "s_queue"
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
