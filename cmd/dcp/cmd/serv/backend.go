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
	case "http", "HTTP":
		return be.HTTP
	case "log", "LOG":
		return be.LOG
	}
	return be.NOOP
}

func buildTaos2(_ context.Context, v *viper.Viper) (*taosdb.Schemaless, error) {
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

func buildMQTT2(_ context.Context, v *viper.Viper) (*mqttc.MQTTAction, error) {
	c := &mqttc.MQTTAction{}
	return c, nil
}

func buildBackend2(ctx context.Context, v *viper.Viper) (be.MutiAction, context.Context) {
	var act be.MutiAction

	blist := v.GetStringSlice("backend.list")

	for _, bk := range blist {
		switch backendName2(bk) {
		case be.TAOS:
			h, err := buildTaos2(ctx, v)
			if err != nil {
				log.Println(err)
			}
			act.Add(h)
		case be.MQTT:
			h, err := buildMQTT2(ctx, v)
			if err != nil {
				log.Println(err)
			}
			act.Add(h)
		case be.HTTP:
		case be.LOG:
		}
	}
	return act, context.WithValue(ctx, keyBackend, act)
}
