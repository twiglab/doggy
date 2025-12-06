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

func New(client mqtt.Client) *MQTTAction {
	return &MQTTAction{client: client}
}

func (c *MQTTAction) HandleData(ctx context.Context, data human.DataMix) error {
	var bb bytes.Buffer
	bb.Grow(1024)

	err := json.MarshalWrite(&bb, &data)
	if err != nil {
		return err
	}

	token := c.client.Publish(pushTopic(data.Head.UUID, data.Type), 0, false, bb)

	return token.Error()
}

func pushTopic(uuid, typ string) string {
	return "dcp/" + uuid + "/" + typ
}
