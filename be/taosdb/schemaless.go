package taosdb

import (
	"context"
	"time"

	"github.com/influxdata/line-protocol/v2/lineprotocol"

	"github.com/taosdata/driver-go/v3/ws/schemaless"
	"github.com/twiglab/doggy/be"
	"github.com/twiglab/doggy/pkg/human"
)

type Schemaless struct {
	schemaless *schemaless.Schemaless
}

func NewSchLe(s *schemaless.Schemaless) *Schemaless {
	return &Schemaless{
		schemaless: s,
	}
}

func (c *Schemaless) Name() string {
	return be.TAOS
}

func (s *Schemaless) HandleData(ctx context.Context, data human.DataMix) error {
	switch data.Type {
	case human.COUNT:
		return s.handleCount(ctx, data)
	case human.DENSITY:
		return s.handleDensity(ctx, data)
	case human.QUEUE:
		return s.handleQueue(ctx, data)
	}
	return be.ErrUnimplType
}

func (s *Schemaless) handleCount(_ context.Context, data human.DataMix) error {
	if !be.HasHuman(data.HumanCountIn, data.HumanCountOut) {
		return nil
	}

	var enc lineprotocol.Encoder

	enc.SetPrecision(lineprotocol.Millisecond)
	enc.StartLine(MA_COUNTY)

	enc.AddTag(TAG_DIVICE_ID, data.Head.Code)
	enc.AddTag(TAG_PROJECT, data.Head.Project)
	enc.AddTag(TAG_TYPE, data.Type)
	enc.AddTag(TAG_UUID, data.Head.ID)

	enc.AddField(FIELD_COUNT_IN, lineprotocol.MustNewValue(int64(data.HumanCountIn)))
	enc.AddField(FIELD_COUNT_OUT, lineprotocol.MustNewValue(int64(data.HumanCountOut)))
	enc.EndLine(data.EndTime)

	if err := enc.Err(); err != nil {
		return err
	}

	bs := enc.Bytes()
	line := bytesToStr(bs)

	return s.schemaless.Insert(line, schemaless.InfluxDBLineProtocol, TSDB_SML_TIMESTAMP_MILLI_SECONDS, 0, 0)
}

func (s *Schemaless) handleDensity(_ context.Context, data human.DataMix) error {
	if !be.HasCount(data.HumanCount) {
		return nil
	}

	var enc lineprotocol.Encoder

	enc.SetPrecision(lineprotocol.Millisecond)
	enc.StartLine(MA_DENSITY)

	enc.AddTag(TAG_DIVICE_ID, data.Head.Code)
	enc.AddTag(TAG_PROJECT, data.Head.Project)
	enc.AddTag(TAG_TYPE, data.Type)
	enc.AddTag(TAG_UUID, data.Head.ID)

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

func (s *Schemaless) handleQueue(_ context.Context, data human.DataMix) error {
	if !be.HasCount(data.HumanCount) {
		return nil
	}

	var enc lineprotocol.Encoder

	enc.SetPrecision(lineprotocol.Millisecond)
	enc.StartLine(MA_QUEUE)

	enc.AddTag(TAG_DIVICE_ID, data.Head.Code)
	enc.AddTag(TAG_PROJECT, data.Head.Project)
	enc.AddTag(TAG_TYPE, data.Type)
	enc.AddTag(TAG_UUID, data.Head.ID)

	enc.AddField(FIELD_QUEUE_COUNT, lineprotocol.MustNewValue(int64(data.HumanCount)))
	enc.AddField(FIELD_QUEUE_TIME, lineprotocol.MustNewValue(int64(data.QueueTime)))
	enc.EndLine(time.Now())

	if err := enc.Err(); err != nil {
		return err
	}

	bs := enc.Bytes()
	line := bytesToStr(bs)

	return s.schemaless.Insert(line, schemaless.InfluxDBLineProtocol, TSDB_SML_TIMESTAMP_MILLI_SECONDS, 0, 0)
}
