package serv

import (
	"database/sql"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/twiglab/doggy/orm"
	"github.com/twiglab/doggy/orm/ent"
	"github.com/twiglab/doggy/taosdb"
	"gopkg.in/yaml.v3"
)

type Sub struct {
	MetadataURL string `yaml:"metadata-url" mapstructure:"metadata-url"`
	Addr        string `yaml:"addr" mapstructure:"addr"`
	Port        int    `yaml:"port" mapstructure:"port"`
	TimeOut     int    `yaml:"time-out" mapstructure:"time-out"`
	HttpsEnable int    `yaml:"https-enable" mapstructure:"https-enable" `
}

type SubsConf struct {
	Main    Sub   `yaml:"main" mapstructure:"main"`
	Backups []Sub `yaml:"backups" mapstructure:"backups"`
}

type KeepliveJobConf struct {
	Crontab string `yaml:"crontab" mapstructure:"crontab"`
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

type BackendConf struct {
	Use        string     `yaml:"use" mapstructure:"use"`
	TaosDBConf TaosDBConf `yaml:"taos" mapstructure:"taos"`
}

type JobConf struct {
	Enable   int             `yaml:"enable" mapstructure:"enable"`
	Keeplive KeepliveJobConf `yaml:"keeplive" mapstructure:"keeplive"`
}

type AutoRegConf struct {
	MutiSub int `yaml:"muti-sub" mapstructure:"muti-sub"`
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
	SubsConf     SubsConf     `yaml:"subs" mapstructure:"subs"`
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
		SubsConf: SubsConf{
			Main: Sub{
				Addr:        "127.0.0.1",
				Port:        10005,
				TimeOut:     0,
				HttpsEnable: 1,
				MetadataURL: "https://127.0.0.1:10005/pf/upload",
			},

			Backups: []Sub{
				{
					Addr:        "127.0.0.1",
					Port:        10005,
					TimeOut:     0,
					HttpsEnable: 1,
					MetadataURL: "https://127.0.0.1:10005/pf/upload",
				},
				{
					Addr:        "127.0.0.1",
					Port:        10005,
					TimeOut:     0,
					HttpsEnable: 1,
					MetadataURL: "https://127.0.0.1:10005/pf/upload",
				},
			},
		},
		CameraDBConf: CameraDBConf{
			CsvCameraDB: CsvCameraDB{
				CameraUser: "ApiAdmin",
				CameraPwd:  "AAaa1234%%",
				CsvFile:    "repo/cameradb.csv",
			},
		},

		AutoRegConf: AutoRegConf{
			MutiSub: 1,
		},

		BackendConf: BackendConf{
			Use: "taos",
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
			Enable: 0,
			Keeplive: KeepliveJobConf{
				Crontab: "*/10 * * * *",
			},
		},
	}

	enc := yaml.NewEncoder(os.Stdout)
	defer enc.Close()
	enc.SetIndent(2)
	enc.Encode(conf)
}
