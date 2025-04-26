package serv

import (
	"log"
	"net/url"
	"strconv"

	"github.com/twiglab/doggy/orm"
	"github.com/twiglab/doggy/orm/ent"
)

type InfluxDBConf struct {
	URL    string
	Token  string
	Org    string
	Bucket string
}

type DbUploadConfig struct {
	MetadataURL string `yaml:"metadata-url" mapstructure:"metadata-url"`
	CameraUser  string `yaml:"camera-user" mapstructure:"camera-user"`
	CameraPwd   string `yaml:"camera-pwd" mapstructure:"camera-pwd"`

	IpAddr string
	Port   int

	DB DB `yaml:"db" mapstructure:"db"`
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
	ID             string `yaml:"id" mapstructure:"id"`
	DbUploadConfig DbUploadConfig
	ServerConf     ServerConfig `yaml:"server" mapstructure:"server"`
	IdbConf        InfluxDBConf
}

func NewUploadConfx(c DbUploadConfig) (pc orm.UploadConfig, err error) {
	var u *url.URL

	if u, err = url.Parse(c.MetadataURL); err != nil {
		return
	}

	pc.User = c.CameraUser
	pc.Pwd = c.CameraPwd

	pc.Address = u.Hostname()
	pc.Port = 80
	pc.MetadataURL = c.MetadataURL

	if p := u.Port(); p != "" {
		pc.Port, err = strconv.Atoi(p)
	}

	return
}

func MustOpenClient(dbconf DB) *ent.Client {
	c, err := orm.OpenClient(dbconf.Name, dbconf.DSN)
	if err != nil {
		log.Fatal(err)
	}
	return c
}
