package serv

import (
	"context"
	"log"

	"github.com/spf13/viper"
	"github.com/twiglab/doggy/be"
	"github.com/twiglab/doggy/be/mqttc"
	"github.com/twiglab/doggy/be/taosdb"
)

type adder interface {
	Add(h be.DataHandler)
}

func backendName(key string) string {
	switch key {
	case "taos", "TAOS":
		return be.TAOS
	case "mqtt", "MQTT":
		return be.MQTT
	case "log", "LOG":
		return be.LOG
	}
	return be.NOOP
}

func buildLogAction(a adder, v *viper.Viper) error {
	logDir := v.GetString("backend.log.logdir")
	a.Add(be.NewLogAction(logDir))
	return nil
}

func buildTaosAction(a adder, v *viper.Viper) error {
	dsn := v.GetString("backend.taos.dsn")

	s, err := taosdb.NewSchLe(dsn)
	if err != nil {
		return err
	}

	a.Add(s)
	return nil
}

func buildMQTTAction(a adder, v *viper.Viper) error {
	brokers := v.GetStringSlice("backend.mqtt.brokers")
	clientID := v.GetString("backend.mqtt.cliend-id")
	cli, err := mqttc.BuildMQTTCLient(mqttc.MQTTConf{
		ClientID: clientID,
		Borkers:  brokers,
	})
	if err != nil {
		return err
	}
	a.Add(mqttc.NewMQTTAction(cli))
	return nil
}

func buildBackend(ctx context.Context, v *viper.Viper) (be.MutiAction, context.Context) {
	var acts be.MutiAction

	blist := v.GetStringSlice("backend.use")

	for _, bk := range blist {
		switch backendName(bk) {
		case be.LOG:
			_ = buildLogAction(&acts, v)
		case be.TAOS:
			if err := buildTaosAction(&acts, v); err != nil {
				log.Println(err)
			}
		case be.MQTT:
			if err := buildMQTTAction(&acts, v); err != nil {
				log.Println(err)
			}
		case be.HTTP:
		}
	}
	return acts, context.WithValue(ctx, keyBackend, acts)
}
