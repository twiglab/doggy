package out

import (
	"context"
	"net/http"

	"github.com/twiglab/doggy/pkg/oc"
)

type Accumulator interface {
	SumOf(context.Context, *oc.SumArgs, *oc.SumReply) error
}

type OutServ struct {
	Accumulator Accumulator
}

func NewOutServ(acc Accumulator) *OutServ {
	return &OutServ{
		Accumulator: acc,
	}
}

func (h *OutServ) Sum(r *http.Request, args *oc.SumArgs, reply *oc.SumReply) error {
	return h.Accumulator.SumOf(r.Context(), args, reply)
}
