package oc

import "time"

const CALL_SUM = "out.Sum"

func ToMilliTimestamp(t time.Time) int64 {
	return t.UnixMilli()
}

type SumArg struct {
	Start int64    `json:"starte"`
	End   int64    `json:"end"`
	IDs   []string `json:"ids"`
}

type MutiSumArg struct {
	Args []SumArg
}

type SumReply struct {
	InTotal  int64
	OutTotal int64
}

type MutiSumReply struct {
	Replies []SumReply
}
