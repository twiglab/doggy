package idb

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/twiglab/doggy/holo"
)

type IdbPoint struct {
	client influxdb2.Client
	w      api.WriteAPI
}

func NewPointHandle(conf IdbConf) *IdbPoint {
	client := influxdb2.NewClientWithOptions(
		conf.URL,
		conf.Token,
		influxdb2.
			DefaultOptions().
			SetBatchSize(20),
	)
	writeAPI := client.WriteAPI(conf.Org, conf.Bucket)

	return &IdbPoint{
		client: client,
		w:      writeAPI,
	}
}

func (h *IdbPoint) HandleDensity(ctx context.Context, common holo.Common, data holo.HumanMix) error {
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

func (h *IdbPoint) HandleCounty(ctx context.Context, common holo.Common, data holo.HumanMix) error {
	if !hasHumanCount(data) {
		return nil
	}

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
