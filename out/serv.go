package out

import "net/http"

type SumArgs struct {
	Table     string
	StartTime int64
	EndTime   int64
}

type SumReply struct {
	Total int64 `json:"total"`
}

type OutServ struct{}

func NewOutServ() *OutServ {
	return &OutServ{}
}

func (h *OutServ) Sum(r *http.Request, args *SumArgs, reply *SumReply) error {
	reply.Total = 8888
	return nil
}
