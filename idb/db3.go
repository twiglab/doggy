package idb

import (
	"context"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3/batching"
	"github.com/twiglab/doggy/holo"
)

type IdbPoint struct {
	client  *influxdb3.Client
	batcher *fixed
}

func emitCallbackFn(client *influxdb3.Client) func([]*influxdb3.Point) {
	return func(ps []*influxdb3.Point) {
		client.WritePoints(context.Background(), ps)
	}
}

func NewIdbPoint(client *influxdb3.Client) *IdbPoint {
	batch := NewFixed(batching.WithEmitCallback(emitCallbackFn(client)))
	return &IdbPoint{client: client, batcher: batch}
}

func (h *IdbPoint) HandleDensity2(ctx context.Context, common holo.Common, data holo.HumanMix) error {
	p := influxdb3.NewPointWithMeasurement(MA_DENSITY).
		SetField(FIELD_DENSITY_COUNT, data.HumanCount).
		SetField(FIELD_DENSITY_RATIO, data.AreaRatio).
		SetTag(TAG_UUID, common.UUID).
		SetTag(TAG_DIVICE_ID, common.DeviceID).
		SetTag(TAG_TYPE, TYPE_12).
		SetTimestamp(time.Now())

	return h.batcher.AddOne(p)
}

func (h *IdbPoint) HandleDensity(ctx context.Context, common holo.Common, data holo.HumanMix) error {
	p := influxdb3.NewPointWithMeasurement(MA_DENSITY).
		SetField(FIELD_DENSITY_COUNT, data.HumanCount).
		SetField(FIELD_DENSITY_RATIO, data.AreaRatio).
		SetTag(TAG_UUID, common.UUID).
		SetTag(TAG_DIVICE_ID, common.DeviceID).
		SetTag(TAG_TYPE, TYPE_12).
		SetTimestamp(time.Now())

	return h.client.WritePoints(ctx, []*influxdb3.Point{p})
}

func (h *IdbPoint) HandleCount(ctx context.Context, common holo.Common, data holo.HumanMix) error {
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
