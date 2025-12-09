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

func (c *MQTTAction) HandleData(ctx context.Context, data human.DataMix) error {
	var bb bytes.Buffer
	bb.Grow(1024)

	err := json.MarshalWrite(&bb, &data)
	if err != nil {
		return err
	}

	pubToken := c.client.Publish("", 0x00, false, bb)
	pubToken.Wait()

	return pubToken.Error()
}

/*
func Main() {
	// 定义 MQTT Broker 地址和端口
	broker := "mqtt://127.0.0.1:1883"
	clientID := "go_mqtt_client"
	// 配置客户端选项
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetKeepAlive(60 * time.Second) // 设置 KeepAlive 时间
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		fmt.Printf("连接丢失: %v\n", err)
	}
	// 创建客户端
	client := mqtt.NewClient(opts)
	// 连接到 Broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("MQTT 连接成功")
	// 调用 Ping 检查函数
	checkMQTTPing(client)
	client.Disconnect(250)
}

func checkMQTTPing(client mqtt.Client) {
	topic := "ping/test"
	message := "ping"
	// 订阅主题
	token := client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("收到消息: %s\n", msg.Payload())
	})
	token.Wait()
	// 发布消息
	pubToken := client.Publish(topic, 1, false, message)
	pubToken.Wait()
	fmt.Println("Ping 消息已发送，等待响应...")
	time.Sleep(2 * time.Second) // 等待响应
	if pubToken.Error() != nil {
		fmt.Printf("Ping 检查失败: %v\n", pubToken.Error())
	} else {
		fmt.Println("Ping 检查成功")
	}
}
*/
