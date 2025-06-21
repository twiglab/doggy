package out

import (
	"context"
	"net/http"

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

func (*UnimplOut) Sum(_ context.Context, _ *oc.SumArg, _ *oc.SumReply) error {
	panic("not implement")
}

func (*UnimplOut) MutiSum(_ context.Context, _ *oc.MutiSumArg, _ *oc.MutiSumReply) error {
	panic("not implement")
}

type OutServ struct {
	Outer Outer
}

func NewOutServ(out Outer) *OutServ {
	return &OutServ{
		Outer: out,
	}
}

func (h *OutServ) Sum(r *http.Request, args *oc.SumArg, reply *oc.SumReply) error {
	return h.Outer.Sum(r.Context(), args, reply)
}

func (h *OutServ) MutiSum(r *http.Request, in *oc.MutiSumArg, out *oc.MutiSumReply) error {
	return h.Outer.MutiSum(r.Context(), in, out)
}
