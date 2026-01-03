package serv

import (
	"context"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/taosdata/driver-go/v3/ws/schemaless"
	"github.com/twiglab/doggy/be"
	"github.com/twiglab/doggy/be/mqttc"
	"github.com/twiglab/doggy/be/taosdb"
)

func backendName2(s string) string {

	if key, _, ok := strings.Cut(s, "_"); ok {
		switch strings.ToLower(key) {
		case be.TAOS:
			return be.TAOS
		case be.MQTT:
			return be.MQTT
		case be.MQTT5:
			return be.MQTT5
		case be.LOG:
			return be.LOG
		}
	}

	return be.NOOP

}
func buildTaos2(ctx context.Context, v *viper.Viper) (*taosdb.Schemaless, error) {
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

func buildMQTT2(ctx context.Context, v *viper.Viper) (*mqttc.MQTTAction, error) {
	c := &mqttc.MQTTAction{}
	return c, nil
}

func buildBackend2(ctx context.Context, v *viper.Viper) (be.MutiAction, context.Context) {
	var act be.MutiAction

	blist := v.GetStringSlice("backend.list")

	for _, bk := range blist {
		switch backendName2(bk) {
		case be.TAOS:
			if h, err := buildTaos2(ctx, subviper(bk, v)); err == nil {
				act.Add(h)
			}
		case be.MQTT:
			if h, err := buildMQTT2(ctx, subviper(bk, v)); err == nil {
				act.Add(h)
			}
		}
	}
	return act, context.WithValue(ctx, keyBackend, act)
}
