package serv

import (
	"log"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/twiglab/doggy/orm"
	"github.com/twiglab/doggy/orm/ent"
)

type InfluxDBConf struct {
	URL    string
	Token  string
	Org    string
	Bucket string
}

type AutoRegConf struct {
	MetadataURL string `yaml:"metadata-url" mapstructure:"metadata-url"`
	Addr        string `yaml:"addr" mapstructure:"addr"`
	Port        int    `yaml:"port" mapstructure:"port"`
}

type FixUserConf struct {
	CameraUser string `yaml:"camera-user" mapstructure:"camera-user"`
	CameraPwd  string `yaml:"camera-pwd" mapstructure:"camera-pwd"`
}

type ServerConfig struct {
	Addr     string `yaml:"addr" mapstructure:"addr"`
	CertFile string `yaml:"cert-file" mapstructure:"cert-file"`
	KeyFile  string `yaml:"key-file" mapstructure:"key-file"`
}

type DB struct {
	Name string `yaml:"name" mapstructure:"name"`
	DSN  string `yaml:"dsn" mapstructure:"dsn"`
}

type AppConf struct {
	ID           string       `yaml:"id" mapstructure:"id"`
	Plan         int          `yaml:"plan" mapstructure:"plan"` // 启动方案 0 debug
	ServerConf   ServerConfig `yaml:"server" mapstructure:"server"`
	InfluxDBConf InfluxDBConf `yaml:"influx-db" mapstructure:"influx-db"`
	FixUserConf  FixUserConf  `yaml:"fix-user" mapstructure:"fix-user"`
	AutoRegConf  AutoRegConf  `yaml:"auto-reg" mapstructure:"auto-reg"`
	DBConf       DB           `yaml:"db" mapstructure:"db"`
}

func MustEntClient(dbconf DB) *ent.Client {
	c, err := orm.OpenClient(dbconf.Name, dbconf.DSN)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func MustIdb(conf InfluxDBConf) *influxdb3.Client {
	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:         conf.URL,
		Token:        conf.Token,
		Organization: conf.Org,
		Database:     conf.Bucket,
	})
	if err != nil {
		log.Fatal(err)
	}
	return client
}
