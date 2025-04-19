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

	NoMetaAutoSub int `yaml:"no-meta-auto-sub" mapstructure:"on-meta-auto-sub"`
}

type ServerConfig struct {
	Addr     string `yaml:"addr" mapstructure:"addr"`
	CertFile string `yaml:"cert-file" mapstructure:"cert-file"`
	KeyFile  string `yaml:"key-file" mapstructure:"key-file"`
}

type DB struct {
	DSN string `yaml:"dsn" mapstructure:"dsn"`
}

type AppConf struct {
	ID             string         `yaml:"id" mapstructure:"id"`
	PlatformConfig PlatformConfig `yaml:"platform" mapstructure:"platform"`
	ServerConf     ServerConfig   `yaml:"server" mapstructure:"server"`
	DB             DB             `yaml:"db" mapstructure:"db"`
}

func MustPfConf(c PlatformConfig) pf.Config {
	pc, err := NewPfConf(c)
	if err != nil {
		log.Fatal(err)
	}
	return pc
}

func NewPfConf(c PlatformConfig) (pc pf.Config, err error) {
	var u *url.URL

	if u, err = url.Parse(c.MetadataURL); err != nil {
		return
	}

	pc.Address = u.Hostname()
	pc.Port = 80
	pc.MetadataURL = c.MetadataURL

	if p := u.Port(); p != "" {
		pc.Port, err = strconv.Atoi(p)
	}

	if c.NoMetaAutoSub != 0 {
		pc.NotMetaAutoSub = true
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
