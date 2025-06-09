package taosdb

import (
	"context"
	"log/slog"
	"time"

	"github.com/influxdata/line-protocol/v2/lineprotocol"

	"github.com/taosdata/driver-go/v3/ws/schemaless"
	"github.com/twiglab/doggy/holo"
)

type Schemaless struct {
	schemaless *schemaless.Schemaless
}

func NewSchLe(s *schemaless.Schemaless) *Schemaless {
	return &Schemaless{
		schemaless: s,
	}
}

func (s *Schemaless) HandleCount(ctx context.Context, common holo.Common, data holo.HumanMix) error {
	if !hasCount(data.HumanCountIn, data.HumanCountOut) {
		return nil
	}

	start := holo.MilliToTime(data.StartTime, data.TimeZone)
	end := holo.MilliToTime(data.EndTime, data.TimeZone)

	slog.InfoContext(ctx, "HandleCount", slog.Int(TAG_TYPE, data.TargetType),
		slog.Group("common", slog.String(TAG_UUID, common.UUID), slog.String(TAG_DIVICE_ID, common.DeviceID)),
		slog.Group("data", slog.Int(FIELD_COUNT_IN, data.HumanCountIn), slog.Int(FIELD_COUNT_OUT, data.HumanCountOut)),
		slog.Group("time", slog.Time("start", start), slog.Time("end", end)),
	)

	var enc lineprotocol.Encoder

	enc.SetPrecision(lineprotocol.Millisecond)
	enc.StartLine(MA_COUNTY)
	enc.AddTag(TAG_DIVICE_ID, common.DeviceID)
	enc.AddTag(TAG_TYPE, TYPE_15)
	enc.AddTag(TAG_UUID, common.UUID)
	enc.AddField(FIELD_COUNT_IN, lineprotocol.MustNewValue(int64(data.HumanCountIn)))
	enc.AddField(FIELD_COUNT_OUT, lineprotocol.MustNewValue(int64(data.HumanCountOut)))
	enc.EndLine(end)

	if err := enc.Err(); err != nil {
		return err
	}

	bs := enc.Bytes()
	line := bytesToStr(bs)

	return s.schemaless.Insert(line, schemaless.InfluxDBLineProtocol, TSDB_SML_TIMESTAMP_MILLI_SECONDS, 0, 0)
}

func (s *Schemaless) HandleDensity(ctx context.Context, common holo.Common, data holo.HumanMix) error {
	slog.InfoContext(ctx, "HandleCount", slog.Int(TAG_TYPE, data.TargetType), slog.Time("now", time.Now()),
		slog.Group("common", slog.String(TAG_UUID, common.UUID), slog.String(TAG_DIVICE_ID, common.DeviceID)),
		slog.Group("data", slog.Int(FIELD_DENSITY_COUNT, data.HumanCount), slog.Int(FIELD_DENSITY_RATIO, data.AreaRatio)),
	)

	var enc lineprotocol.Encoder

	enc.SetPrecision(lineprotocol.Millisecond)
	enc.StartLine(MA_DENSITY)
	enc.AddTag(TAG_DIVICE_ID, common.DeviceID)
	enc.AddTag(TAG_TYPE, TYPE_12)
	enc.AddTag(TAG_UUID, common.UUID)
	enc.AddField(FIELD_DENSITY_COUNT, lineprotocol.MustNewValue(int64(data.HumanCount)))
	enc.AddField(FIELD_DENSITY_RATIO, lineprotocol.MustNewValue(int64(data.AreaRatio)))
	enc.EndLine(time.Now())

	if err := enc.Err(); err != nil {
		return err
	}

	bs := enc.Bytes()
	line := bytesToStr(bs)

	return s.schemaless.Insert(line, schemaless.InfluxDBLineProtocol, TSDB_SML_TIMESTAMP_MILLI_SECONDS, 0, 0)
}
