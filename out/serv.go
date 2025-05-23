package out

import (
	"context"
	"net/http"

	"github.com/twiglab/doggy/pkg/oc"
)

type Accumulator interface {
	SumOf(context.Context, *oc.AreaArg, *oc.Reply) error
}

type Densimeter interface {
	CollectOf(context.Context, *oc.AreaArg, *oc.Reply) error
}

type Outer interface {
	Accumulator
	Densimeter
}

type UnimplOut struct {
}

func (*UnimplOut) SumOf(_ context.Context, _ *oc.AreaArg, _ *oc.Reply) error {
	panic("not implement")
}

func (*UnimplOut) CollectOf(_ context.Context, _ *oc.AreaArg, _ *oc.Reply) error {
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

func (h *OutServ) Sum(r *http.Request, args *oc.AreaArg, reply *oc.Reply) error {
	return h.Outer.SumOf(r.Context(), args, reply)
}

func (h *OutServ) Collect(r *http.Request, args *oc.AreaArg, reply *oc.Reply) error {
	return h.Outer.CollectOf(r.Context(), args, reply)
}
