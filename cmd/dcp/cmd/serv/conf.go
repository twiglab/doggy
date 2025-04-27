package serv

import (
	"log"
	"os"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/spf13/cobra"
	"github.com/twiglab/doggy/orm"
	"github.com/twiglab/doggy/orm/ent"
	"gopkg.in/yaml.v3"
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

var ConfCmd = &cobra.Command{
	Use:   "config",
	Short: "生成配置文件",
	Long:  `生成配置文件`,
	Run: func(cmd *cobra.Command, args []string) {
		confCmd()
	},
	Example: "dcp serv config",
}

func init() {
	ServCmd.AddCommand(ConfCmd)
}

func confCmd() {
	conf := AppConf{ID:"dcp"}

	enc := yaml.NewEncoder(os.Stdout)
	defer enc.Close()
	enc.SetIndent(2)
	enc.Encode(conf)
}
