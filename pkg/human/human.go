package human

import "time"

const (
	DENSITY = "density"
	QUEUE   = "queue"
	COUNT   = "count"
	UNKNOWN = "unknown"
)

type Head struct {
	SN       string `json:"sn"`
	IpAddr   string `json:"ipAddr"`
	UUID     string `json:"uuid"`
	DeviceID string `json:"deviceID"`
	Project  string `json:"project"`
}

type DataMix struct {
	Head Head `json:"head"`

	Type string `json:"type"`

	HumanCountIn  int       `json:"humanCountIn,omitempty"`
	HumanCountOut int       `json:"humanCountOut,omitempty"`
	BeginTime     time.Time `json:"beginTime,omitzero"`
	EndTime       time.Time `json:"endTime,omitzero"`

	HumanCount int `json:"humanCount,omitempty"`
	AreaRatio  int `json:"areaRatio,omitempty"`

	QueueTime int `json:"queueTime,omitempty"` // 2.6.8

}
