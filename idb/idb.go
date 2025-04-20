package idb

import "github.com/twiglab/doggy/holo"

const (
	TYPE_12 = "12"
	TYPE_15 = "15"

	TAG_UUID      = "uuid"
	TAG_TYPE      = "type"
	TAG_DIVICE_ID = "device_id"

	FIELD_DENSITY_COUNT = "count"
	FIELD_DENSITY_RATIO = "ratio"

	FIELD_COUNT_IN  = "in"
	FIELD_COUNT_OUT = "out"

	MA_DENSITY = "density"
	MA_COUNTY  = "count"
)

type IdbConf struct {
	URL    string
	Token  string
	Org    string
	Bucket string
}

func hasHumanCount(data holo.HumanMix) bool {
	if data.HumanCountIn == 0 && data.HumanCountOut == 0 {
		return false
	}
	return true
}
