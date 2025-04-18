package serv

import (
	"context"
	"log"
	"net/url"
	"strconv"

	"github.com/twiglab/doggy/orm"
	"github.com/twiglab/doggy/orm/ent"
	"github.com/twiglab/doggy/pf"
)

type PlatformConfig struct {
	MetadataURL string `yaml:"metadata-url" mapstructure:"metadata-url"`
	CameraUser  string `yaml:"camera-user" mapstructure:"camera-user"`
	CameraPwd   string `yaml:"camera-pwd" mapstructure:"camera-pwd"`
}

type ServerConfig struct {
	Addr     string `yaml:"addr" mapstructure:"addr"`
	CertFile string `yaml:"cert-file" mapstructure:"cert-file"`
	KeyFile  string `yaml:"key-file" mapstructure:"key-file"`
}

type DB struct {
	DSN string `yaml:"dsn" mapstructure:"dsn"`
}

type App struct {
	ID             string         `yaml:"id" mapstructure:"id"`
	PlatformConfig PlatformConfig `yaml:"platform" mapstructure:"platform"`
	ServerConf     ServerConfig   `yaml:"server" mapstructure:"server"`
	DB             DB             `yaml:"db" mapstructure:"db"`
}

func MustPlatformConfig(s string) pf.PlatformConfig {
	pc, err := NewPlatformConfig(s)
	if err != nil {
		log.Fatal(err)
	}
	return pc
}

func NewPlatformConfig(s string) (pc pf.PlatformConfig, err error) {
	var u *url.URL

	if u, err = url.Parse(s); err != nil {
		return
	}

	pc.Address = u.Hostname()
	pc.Port = 80
	pc.MetadataURL = s

	if p := u.Port(); p != "" {
		pc.Port, err = strconv.Atoi(p)
	}

	return
}

func MustClient(s string) *ent.Client {
	db, err := orm.FromURL(context.Background(), s)
	if err != nil {
		log.Fatal(err)
	}

	return orm.OpenClient(db)
}
