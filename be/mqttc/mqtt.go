package mqttc

import (
	"bytes"
	"context"
	"encoding/json/v2"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/doggy/be"
	"github.com/twiglab/doggy/pkg/human"
)

func pushTopic(uuid, typ string) string {
	return "dcp/" + uuid + "/" + typ
}

type MQTTAction struct {
	client mqtt.Client
}

func NewMQTTAction(client mqtt.Client) *MQTTAction {
	return &MQTTAction{client: client}
}

func (c *MQTTAction) Name() string {
	return be.MQTT
}

func (c *MQTTAction) HandleData(ctx context.Context, data human.DataMix) error {
	var bb bytes.Buffer
	bb.Grow(1024)

	err := json.MarshalWrite(&bb, &data)
	if err != nil {
		return err
	}

	topic := pushTopic(data.Head.ChannelID, data.Type)

	pubToken := c.client.Publish(topic, 0x00, false, bb)
	pubToken.Wait()

	return pubToken.Error()
}

type MQTTConf struct {
	ClientID string
	Borkers  []string
}

func BuildMQTTCLient(conf MQTTConf) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.SetClientID(conf.ClientID)

	for _, b := range conf.Borkers {
		opts.AddBroker(b)
	}

	client := mqtt.NewClient(opts)
	// 连接到 Broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}
