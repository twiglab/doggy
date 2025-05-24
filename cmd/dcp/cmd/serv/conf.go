package serv

import (
	"database/sql"
	"log"
	"os"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/spf13/cobra"
	"github.com/twiglab/doggy/orm"
	"github.com/twiglab/doggy/orm/ent"
	"github.com/twiglab/doggy/taosdb"
	"gopkg.in/yaml.v3"
)

type KeepliveJobConf struct {
	Crontab string `yaml:"crontab" mapstructure:"crontab"`

	MetadataURL string `yaml:"metadata-url" mapstructure:"metadata-url"`
	Addr        string `yaml:"addr" mapstructure:"addr"`
	Port        int    `yaml:"port" mapstructure:"port"`
}

type LoggerConf struct {
	Level   string `yaml:"level" mapstructure:"level"`
	LogFile string `yaml:"log-file" mapstructure:"log-file"`
}

type TaosDBConf struct {
	Addr     string `yaml:"addr" mapstructure:"addr"`
	Port     int    `yaml:"port" mapstructure:"port"`
	Protocol string `yaml:"protocol" mapstructure:"protocol"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
	DBName   string `yaml:"dbname" mapstructure:"dbname"`
}

type InfluxDBConf struct {
	URL    string `yaml:"url" mapstructure:"url"`
	Token  string `yaml:"token" mapstructure:"token"`
	Org    string `yaml:"org" mapstructure:"org"`
	Bucket string `yaml:"bucket" mapstructure:"bucket"`
}

type BackendConf struct {
	Use          string       `yaml:"use" mapstructure:"use"`
	InfluxDBConf InfluxDBConf `yaml:"influx-db" mapstructure:"influx-db"`
	TaosDBConf   TaosDBConf   `yaml:"taos" mapstructure:"taos"`
}

type JobConf struct {
	Disable  bool            `yaml:"disable" mapstructure:"disable"`
	Keeplive KeepliveJobConf `yaml:"keeplive" mapstructure:"keeplive"`
}

type AutoRegConf struct {
	MetadataURL string `yaml:"metadata-url" mapstructure:"metadata-url"`
	Addr        string `yaml:"addr" mapstructure:"addr"`
	Port        int    `yaml:"port" mapstructure:"port"`
}

type CameraDBConf struct {
	CsvCameraDB CsvCameraDB `yaml:"csvdb" mapstructure:"csvdb"`
}

type CsvCameraDB struct {
	CameraUser string `yaml:"camera-user" mapstructure:"camera-user"`
	CameraPwd  string `yaml:"camera-pwd" mapstructure:"camera-pwd"`
	CsvFile    string `yaml:"csv-file" mapstructure:"csv-file"`
}

type ServerConf struct {
	Addr       string `yaml:"addr" mapstructure:"addr"`
	CertFile   string `yaml:"cert-file" mapstructure:"cert-file"`
	KeyFile    string `yaml:"key-file" mapstructure:"key-file"`
	ForceHttps int    `yaml:"force-https" mapstructure:"force-https"`
}

type DBConf struct {
	Name string `yaml:"name" mapstructure:"name"`
	DSN  string `yaml:"dsn" mapstructure:"dsn"`
}

type AppConf struct {
	ID           string       `yaml:"id" mapstructure:"id"`
	LoggerConf   LoggerConf   `yaml:"log" mapstructure:"log"`
	ServerConf   ServerConf   `yaml:"server" mapstructure:"server"`
	BackendConf  BackendConf  `yaml:"backend" mapstructure:"backend"`
	CameraDBConf CameraDBConf `yaml:"camera-db" mapstructure:"camera-db"`
	AutoRegConf  AutoRegConf  `yaml:"auto-reg" mapstructure:"auto-reg"`
	DBConf       DBConf       `yaml:"db" mapstructure:"db"`
	JobConf      JobConf      `yaml:"job" mapstructure:"job"`
}

func MustEntClient(dbconf DBConf) *ent.Client {
	c, err := orm.OpenClient(dbconf.Name, dbconf.DSN)
	if err != nil {
		log.Fatal("ent client create error: ", err)
	}
	return c
}

func MustOpenTaosDB(conf AppConf) *sql.DB {
	db, err := taosdb.OpenDB(taosdb.Config{
		Addr:     conf.BackendConf.TaosDBConf.Addr,
		Port:     conf.BackendConf.TaosDBConf.Port,
		Protocol: conf.BackendConf.TaosDBConf.Protocol,
		Username: conf.BackendConf.TaosDBConf.Username,
		Password: conf.BackendConf.TaosDBConf.Password,
		DBName:   conf.BackendConf.TaosDBConf.DBName,
	})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func MustIdb(conf InfluxDBConf) *influxdb3.Client {
	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:         conf.URL,
		Token:        conf.Token,
		Organization: conf.Org,
		Database:     conf.Bucket,
	})
	if err != nil {
		log.Fatal("idb create error: ", err)
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
	conf := AppConf{
		ID: "dcp",

		LoggerConf: LoggerConf{
			Level:   "INFO",
			LogFile: "dcp.log",
		},

		ServerConf: ServerConf{
			Addr:       "0.0.0.0:10005",
			CertFile:   "repo/server.crt",
			KeyFile:    "repo/server.key",
			ForceHttps: 1,
		},

		CameraDBConf: CameraDBConf{
			CsvCameraDB: CsvCameraDB{
				CameraUser: "ApiAdmin",
				CameraPwd:  "AAaa1234%%",
				CsvFile:    "repo/cameradb.csv",
			},
		},

		AutoRegConf: AutoRegConf{
			Addr:        "127.0.0.1",
			Port:        10005,
			MetadataURL: "https://127.0.0.1:10005/pf/upload",
		},

		BackendConf: BackendConf{
			Use: "taos",
			InfluxDBConf: InfluxDBConf{
				URL:    "url",
				Token:  "token",
				Org:    "org",
				Bucket: "bucket",
			},
			TaosDBConf: TaosDBConf{
				Addr:     "127.0.0.1",
				Port:     6041,
				Protocol: "ws",
				Username: "root",
				Password: "taosdata",
				DBName:   "taosdb",
			},
		},

		DBConf: DBConf{
			Name: "sqlite3",
			DSN:  "repo/dcp.db",
		},

		JobConf: JobConf{
			Keeplive: KeepliveJobConf{
				Crontab:     "*/10 * * * *",
				Addr:        "127.0.0.1",
				Port:        10005,
				MetadataURL: "https://127.0.0.1:10005/pf/upload",
			},
		},
	}

	enc := yaml.NewEncoder(os.Stdout)
	defer enc.Close()
	enc.SetIndent(2)
	enc.Encode(conf)
}
