package oc

import "time"

const (
	SumCall     = "out.Sum"
	MuitSumCall = "out.MutiSum"
)

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

func OneSumArg(start, end int64, ids []string) MutiSumArg {
	return MutiSumArg{
		Args: []SumArg{
			{Start: start, End: end, IDs: ids},
		},
	}
}

type SumReply struct {
	InTotal  int64
	OutTotal int64
}

type MutiSumReply struct {
	Replies []SumReply
}
