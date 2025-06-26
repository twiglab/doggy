package out

import (
	"context"

	"github.com/twiglab/doggy/pkg/oc"
)

type Accumulator interface {
	Sum(context.Context, *oc.SumArg, *oc.SumReply) error
	MutiSum(context.Context, *oc.MutiSumArg, *oc.MutiSumReply) error
}

type Outer interface {
	Accumulator
}

type UnimplOut struct {
}

func (UnimplOut) Sum(_ context.Context, _ *oc.SumArg, _ *oc.SumReply) error {
	panic("not implement")
}

func (UnimplOut) MutiSum(_ context.Context, _ *oc.MutiSumArg, _ *oc.MutiSumReply) error {
	panic("not implement")
}
