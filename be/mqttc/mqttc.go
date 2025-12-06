package mqttc

import (
	"bytes"
	"context"
	"encoding/json/v2"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/twiglab/doggy/pkg/human"
)

type MQTTAction struct {
	client mqtt.Client
}

func New() *MQTTAction {
	return &MQTTAction{}
}

func (c *MQTTAction) HandleData(ctx context.Context, data human.DataMix) error {
	var bb bytes.Buffer
	bb.Grow(1024)

	err := json.MarshalWrite(&bb, &data)
	if err != nil {
		return err
	}

	token := c.client.Publish("test/topic", 0, false, bb.Bytes())
	return token.Error()
}
