package mqttc

import (
	"context"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/twiglab/doggy/be"
	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/pf"
)

type MQTTAction struct {
	client mqtt.Client
}

func New() *MQTTAction {
	return &MQTTAction{}
}

func (c *MQTTAction) HandleData(ctx context.Context, data pf.UploadeData) error {
	switch data.Target.TargetType {
	case holo.HUMMAN_COUNT:
		return c.handleCount(ctx, data.Common, data.Target)
	case holo.HUMMAN_DENSITY:
		return c.handleDensity(ctx, data.Common, data.Target)
	case holo.HUMMAN_QUEUE:
		return c.handleQueue(ctx, data.Common, data.Target)
	}
	return pf.ErrUnimplType
}

func (c *MQTTAction) handleCount(ctx context.Context, common holo.Common, data holo.HumanMix) error {
	if !be.HasHuman(data.HumanCountIn, data.HumanCountOut) {
		return nil
	}

	start := holo.MilliToTime(data.StartTime, data.TimeZone)
	end := holo.MilliToTime(data.EndTime, data.TimeZone)

	slog.DebugContext(ctx, "HandleCount", slog.Any("data", data), slog.Any("common", common),
		slog.Group("time", slog.Time("start", start), slog.Time("end", end)),
	)

	token := c.client.Publish("test/topic", 0, false, "Hello MQTT")
	return token.Error()
}
func (c *MQTTAction) handleDensity(ctx context.Context, common holo.Common, data holo.HumanMix) error {
	if !be.HasCount(data.HumanCount) {
		return nil
	}

	slog.DebugContext(ctx, "HandleQueue", slog.Any("common", common), slog.Any("data", data))

	token := c.client.Publish("test/topic", 0, false, "Hello MQTT")
	return token.Error()
}

func (c *MQTTAction) handleQueue(ctx context.Context, common holo.Common, data holo.HumanMix) error {
	if !be.HasCount(data.HumanCount) {
		return nil
	}

	slog.DebugContext(ctx, "HandleQueue", slog.Any("common", common), slog.Any("data", data))

	token := c.client.Publish("test/topic", 0, false, "Hello MQTT")
	return token.Error()
}
