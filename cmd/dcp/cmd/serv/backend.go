package serv

import (
	"context"
	"log"
	"time"

	"github.com/spf13/viper"
	"github.com/taosdata/driver-go/v3/ws/schemaless"
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
	url := taosdb.SchemalessURL(
		v.GetString("backend.taos.addr"),
		v.GetInt("backend.taos.port"),
	)
	sc := schemaless.NewConfig(url, 1,
		schemaless.SetDb(v.GetString("backend.taos.dbname")),
		schemaless.SetAutoReconnect(true),
		schemaless.SetUser(v.GetString("backend.taos.username")),
		schemaless.SetPassword(v.GetString("backend.taos.password")),
		schemaless.SetReadTimeout(5*time.Second),
		schemaless.SetWriteTimeout(5*time.Second),
	)
	s, err := schemaless.NewSchemaless(sc)
	if err != nil {
		return err
	}

	a.Add(taosdb.NewSchLe(s))
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
		case be.TAOS:
			if err := buildTaosAction(&acts, v); err != nil {
				log.Println(err)
			}
		case be.MQTT:
			if err := buildMQTTAction(&acts, v); err != nil {
				log.Println(err)
			}
		case be.LOG:
			_ = buildLogAction(&acts, v)
		}
	}
	return acts, context.WithValue(ctx, keyBackend, acts)
}
