package oc

type SumArgs struct {
	Table string   `json:"table omitempty"`
	Start int64    `json:"starte"`
	End   int64    `json:"end"`
	IDs   []string `json:"ids"`
	UUIDs []string `json:"uuids omitempty"`
}

type SumReply struct {
	Total int64 `json:"total"`
}
