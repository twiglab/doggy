package oc

type AreaArg struct {
	Start int64    `json:"starte"`
	End   int64    `json:"end"`
	IDs   []string `json:"ids"`
}

type Reply struct {
	ValueA int64 `json:"value_a,omitempty"`
	ValueB int64 `json:"value_b,omitempty"`
	ValueC int64 `json:"value_c,omitempty"`
	ValueD int64 `json:"value_d,omitempty"`
	ValueE int64 `json:"value_e,omitempty"`
}
