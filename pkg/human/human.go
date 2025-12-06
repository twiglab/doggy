package human

import "time"

const (
	DENSITY = 12 // 密度 2.6.7
	QUEUE   = 13 // 排队长度 2.6.8
	COUNT   = 15 // 人数 2.6.9
)

type Head struct {
	SN       string
	IpAddr   string
	UUID     string `json:"UUID"`
	DeviceID string `json:"deviceID"`
	Project  string `json:"project"`
}

type DataMix struct {
	Head Head `json:"head"`

	TargetType int `json:"targetType"`

	HumanCountIn  int       `json:"humanCountIn,omitempty"`
	HumanCountOut int       `json:"humanCountOut,omitempty"`
	BeginTime     time.Time `json:"beginTime,omitzero"`
	EndTime       time.Time `json:"endTime,omitzero"`

	HumanCount int `json:"humanCount,omitempty"`
	AreaRatio  int `json:"areaRatio,omitempty"`

	QueueTime int `json:"queueTime,omitempty"` // 2.6.8

}

func MilliToTime(milli int64, tz int64) time.Time {
	return time.UnixMilli(milli)
}
