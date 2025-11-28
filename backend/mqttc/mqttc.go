package mqttc

import (
	"context"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/twiglab/doggy/backend"
	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/pf"
)

type MQTTAction struct {
	client mqtt.Client
}

func (c *MQTTAction) HandleData(ctx context.Context, data pf.UploadeData) error {
	switch data.Target.TargetType {
	case holo.HUMMAN_COUNT:
		return c.handleCount(ctx, data.Common, data.Target)
	case holo.HUMMAN_QUEUE:
	case holo.HUMMAN_DENSITY:
	}
	return pf.ErrUnimplType
}

func (c *MQTTAction) handleCount(ctx context.Context, common holo.Common, data holo.HumanMix) error {
	if !backend.HasHuman(data.HumanCountIn, data.HumanCountOut) {
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

/*
func X() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://broker.emqx.io:1883")
	opts.SetClientID("go_mqtt_client")
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Connected to MQTT Broker")

	// Subscribe to a topic
	client.Subscribe("test/topic", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message: %s\n", msg.Payload())
	})

	// Publish a message
	client.Publish("test/topic", 0, false, "Hello MQTT")

}
*/
