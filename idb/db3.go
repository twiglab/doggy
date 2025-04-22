package idb

import (
	"context"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/twiglab/doggy/holo"
)

type IdbPoint3 struct {
	client *influxdb3.Client
}

func NewIdb3(client *influxdb3.Client) *IdbPoint3 {
	return &IdbPoint3{client: client}
}

func (h *IdbPoint3) HandleDensity(ctx context.Context, common holo.Common, data holo.HumanMix) error {
	p := influxdb3.NewPointWithMeasurement(MA_DENSITY).
		SetField(FIELD_DENSITY_COUNT, data.HumanCount).
		SetField(FIELD_DENSITY_RATIO, data.AreaRatio).
		SetTag(TAG_UUID, common.UUID).
		SetTag(TAG_DIVICE_ID, common.DeviceID).
		SetTag(TAG_TYPE, TYPE_12).
		SetTimestamp(time.Now())

	return h.client.WritePoints(ctx, []*influxdb3.Point{p})
}

func (h *IdbPoint3) HandleCount(ctx context.Context, common holo.Common, data holo.HumanMix) error {
	if !hasHumanCount(data) {
		return nil
	}

	time := holo.MilliToTime(data.EndTime, data.TimeZone)

	p := influxdb3.NewPointWithMeasurement(MA_COUNTY).
		SetField(FIELD_COUNT_IN, data.HumanCountIn).
		SetField(FIELD_COUNT_OUT, data.HumanCountOut).
		SetTag(TAG_UUID, common.UUID).
		SetTag(TAG_DIVICE_ID, common.DeviceID).
		SetTag(TAG_TYPE, TYPE_15).
		SetTimestamp(time)

	return h.client.WritePoints(ctx, []*influxdb3.Point{p})
}
