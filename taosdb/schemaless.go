package taosdb

import (
	"context"
	"time"

	"github.com/influxdata/line-protocol/v2/lineprotocol"

	comm "github.com/taosdata/driver-go/v3/common"
	"github.com/taosdata/driver-go/v3/ws/schemaless"
	"github.com/twiglab/doggy/holo"
)

type Schemaless struct {
	schemaless *schemaless.Schemaless
}

func NewSchemaless(conf Config) (*Schemaless, error) {
	url := schemalessURL(conf)

	sc := schemaless.NewConfig(url, 1,
		schemaless.SetDb(conf.DBName),
		schemaless.SetAutoReconnect(true),
		schemaless.SetUser(conf.Username),
		schemaless.SetPassword(conf.Password),
		schemaless.SetReadTimeout(5*time.Second),
		schemaless.SetWriteTimeout(5*time.Second),
	)

	s, err := schemaless.NewSchemaless(sc)
	if err != nil {
		return nil, err
	}

	return &Schemaless{schemaless: s}, nil
}

func (s *Schemaless) HandleDensity(ctx context.Context, common holo.Common, data holo.HumanMix) error {
	var enc lineprotocol.Encoder

	enc.SetPrecision(lineprotocol.Microsecond)
	enc.StartLine(MA_DENSITY)
	enc.AddTag(TAG_UUID, common.UUID)
	enc.AddTag(TAG_DIVICE_ID, common.DeviceID)
	enc.AddTag(TAG_TYPE, TYPE_12)
	enc.AddField(FIELD_DENSITY_COUNT, lineprotocol.MustNewValue(data.HumanCount))
	enc.AddField(FIELD_DENSITY_RATIO, lineprotocol.MustNewValue(data.AreaRatio))
	enc.EndLine(time.Now())

	if err := enc.Err(); err != nil {
		return nil
	}

	bs := enc.Bytes()
	line := bytesToStr(bs)

	return s.schemaless.Insert(line, schemaless.InfluxDBLineProtocol, "ms", 0, comm.GetReqID())
}

func (s *Schemaless) HandleCount(ctx context.Context, common holo.Common, data holo.HumanMix) error {
	if !hasCount(data.HumanCountIn, data.HumanCountOut) {
		return nil
	}

	var enc lineprotocol.Encoder

	enc.SetPrecision(lineprotocol.Microsecond)
	enc.StartLine(MA_COUNTY)
	enc.AddTag(TAG_UUID, common.UUID)
	enc.AddTag(TAG_DIVICE_ID, common.DeviceID)
	enc.AddTag(TAG_TYPE, TYPE_15)
	enc.AddField(FIELD_COUNT_IN, lineprotocol.MustNewValue(data.HumanCountIn))
	enc.AddField(FIELD_COUNT_OUT, lineprotocol.MustNewValue(data.HumanCountOut))
	enc.EndLine(holo.MilliToTime(data.EndTime, data.TimeZone))

	if err := enc.Err(); err != nil {
		return nil
	}

	bs := enc.Bytes()
	line := bytesToStr(bs)

	return s.schemaless.Insert(line, schemaless.InfluxDBLineProtocol, "ms", 0, comm.GetReqID())
}
