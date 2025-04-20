package idb

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/twiglab/doggy/holo"
)

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

type IdbHandle struct {
	client influxdb2.Client
	w      api.WriteAPI
}

func NewIdbHandle(conf IdbConf) *IdbHandle {
	client := influxdb2.NewClientWithOptions(
		conf.URL,
		conf.Token,
		influxdb2.
			DefaultOptions().
			SetBatchSize(20),
	)
	writeAPI := client.WriteAPI(conf.Org, conf.Bucket)

	return &IdbHandle{
		client: client,
		w:      writeAPI,
	}
}

func (h *IdbHandle) HandleDensity(ctx context.Context, common holo.Common, data holo.HumanMix) error {
	p := influxdb2.NewPointWithMeasurement(MA_DENSITY).
		AddField(FIELD_DENSITY_COUNT, data.HumanCount).
		AddField(FIELD_DENSITY_RATIO, data.AreaRatio).
		AddTag(TAG_UUID, common.UUID).
		AddTag(TAG_DIVICE_ID, common.DeviceID).
		AddTag(TAG_TYPE, TYPE_12).
		SetTime(time.Now())

	h.w.WritePoint(p)
	return nil
}

func (h *IdbHandle) HandleCounty(ctx context.Context, common holo.Common, data holo.HumanMix) error {
	time := holo.MilliToTime(data.EndTime, data.TimeZone)

	p := influxdb2.NewPointWithMeasurement(MA_COUNTY).
		AddField(FIELD_COUNT_IN, data.HumanCountIn).
		AddField(FIELD_COUNT_OUT, data.HumanCountOut).
		AddTag(TAG_UUID, common.UUID).
		AddTag(TAG_DIVICE_ID, common.DeviceID).
		AddTag(TAG_TYPE, TYPE_15).
		SetTime(time)

	h.w.WritePoint(p)
	return nil
}
