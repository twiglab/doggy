package out

import (
	"context"

	"net/http"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
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

const SERV_OUT = "out"

func OutHandle(outer Outer) http.Handler {
	s := rpc.NewServer()
	s.RegisterCodec(json2.NewCodec(), "application/json")
	s.RegisterService(outer, SERV_OUT)
	return s
}
