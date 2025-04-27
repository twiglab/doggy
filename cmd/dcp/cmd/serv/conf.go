package serv

import (
	"log"

	"github.com/twiglab/doggy/orm"
	"github.com/twiglab/doggy/orm/ent"
)

type InfluxDBConf struct {
	URL    string
	Token  string
	Org    string
	Bucket string
}

type AutoSubConfig struct {
	MetadataURL string `yaml:"metadata-url" mapstructure:"metadata-url"`
	Addr        string
	Port        int
}

type FixUserConfig struct {
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
	ID         string       `yaml:"id" mapstructure:"id"`
	ServerConf ServerConfig `yaml:"server" mapstructure:"server"`
	IdbConf    InfluxDBConf
}

func MustOpenClient(dbconf DB) *ent.Client {
	c, err := orm.OpenClient(dbconf.Name, dbconf.DSN)
	if err != nil {
		log.Fatal(err)
	}
	return c
}
