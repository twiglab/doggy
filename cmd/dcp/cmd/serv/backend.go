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

func backendName2(key string) string {
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

func buildTaosAction(_ context.Context, v *viper.Viper) (*taosdb.Schemaless, error) {
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
		return nil, err
	}

	return taosdb.NewSchLe(s), nil
}

func buildMQTTAction(_ context.Context, v *viper.Viper) (*mqttc.MQTTAction, error) {
	brokers := v.GetStringSlice("backend.mqtt.brokers")
	clientID := v.GetString("backend.mqtt.cliend-id")
	cli, err := mqttc.BuildMQTTCLient(mqttc.MQTTConf{
		ClientID: clientID,
		Borkers:  brokers,
	})
	if err != nil {
		return nil, err
	}
	return mqttc.NewMQTTAction(cli), nil
}

func buildLogAction(_ context.Context, v *viper.Viper) (be.LogAction, error) {
	logDir := v.GetString("backend.log.logdir")
	return be.NewLogAction(logDir), nil
}

func buildBackend(ctx context.Context, v *viper.Viper) (be.MutiAction, context.Context) {
	var acts be.MutiAction

	blist := v.GetStringSlice("backend.use")

	for _, bk := range blist {
		switch backendName2(bk) {
		case be.TAOS:
			sch, err := buildTaosAction(ctx, v)
			if err != nil {
				log.Println(err)
			} else {
				acts.Add(sch)
			}
		case be.MQTT:
			mqc, err := buildMQTTAction(ctx, v)
			if err != nil {
				log.Println(err)
			} else {
				acts.Add(mqc)
			}
		case be.LOG:
			la, _ := buildLogAction(ctx, v)
			acts.Add(la)
		}
	}
	return acts, context.WithValue(ctx, keyBackend, acts)
}
