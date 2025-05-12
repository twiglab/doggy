package out

import (
	"context"
	"net/http"
)

type Accumulator interface {
	SumOf(context.Context, *SumArgs, *SumReply) error
}

type SumArgs struct {
	Table string   `json:"table"`
	Start int64    `json:"starte"`
	End   int64    `json:"end"`
	IDs   []string `json:"ids"`
}

type SumReply struct {
	Total int64 `json:"total"`
}

type OutServ struct {
	Accumulator Accumulator
}

func NewOutServ() *OutServ {
	return &OutServ{}
}

func (h *OutServ) Sum(r *http.Request, args *SumArgs, reply *SumReply) error {
	return h.Accumulator.SumOf(r.Context(), args, reply)
}
