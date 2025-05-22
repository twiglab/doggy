package idb

import (
	"context"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/twiglab/doggy/pkg/oc"
)

type SumTotal struct {
	InTotal int64
}

type SumParam struct {
	Start int64
	End   int64
	IDs   []string
}

func (h *IdbPoint) SumOfPoints(ctx context.Context, param SumParam) (total SumTotal, err error) {
	var (
		iter *influxdb3.PointValueIterator
		pv   *influxdb3.PointValues
	)

	sql := sumOfPointsSQL(param.IDs)

	iter, err = h.client.QueryPointValueWithParameters(ctx, sql,
		influxdb3.QueryParameters{
			TIME_START: param.Start,
			TIME_END:   param.End,
		},
	)

	if err != nil {
		return
	}

	for {
		pv, err = iter.Next()
		if err == influxdb3.Done {
			break
		}

		if err != nil {
			return
		}

		if pi := pv.GetIntegerField(IN_TOTAL); pi != nil {
			total.InTotal = *pi
		}
	}
	return
}

type IdbOut struct {
	Point *IdbPoint
}

func NewIdbOut(p *IdbPoint) *IdbOut {
	return &IdbOut{Point: p}
}

func (p *IdbOut) SumOf(ctx context.Context, in *oc.SumArgs, out *oc.SumReply) error {
	result, err := p.Point.SumOfPoints(ctx, SumParam{
		Start: in.Start,
		End:   in.End,
		IDs:   in.IDs,
	})
	if err != nil {
		return err
	}

	out.Total = result.InTotal
	return nil
}
