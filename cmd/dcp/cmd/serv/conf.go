package serv

import (
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type TenantConf struct {
	TenantID string `yaml:"tenant-id" mapstructure:"tenant-id"`
}

type SubsConf struct {
	Muti    int      `yaml:"muti" mapstructure:"muti"`
	Main    string   `yaml:"main" mapstructure:"main"`
	Backups []string `yaml:"backups" mapstructure:"backups"`
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

type CameraDBConf struct {
	CsvCameraDB CsvCameraDB `yaml:"csvdb" mapstructure:"csvdb"`
}

type CsvCameraDB struct {
	CameraUser string `yaml:"camera-user" mapstructure:"camera-user"`
	CameraPwd  string `yaml:"camera-pwd" mapstructure:"camera-pwd"`
}

type ServerConf struct {
	Addr       string `yaml:"addr" mapstructure:"addr"`
	CertFile   string `yaml:"cert-file" mapstructure:"cert-file"`
	KeyFile    string `yaml:"key-file" mapstructure:"key-file"`
	ForceHttps int    `yaml:"force-https" mapstructure:"force-https"`
}

type EtcdConf struct {
	URLs []string `yaml:"urls" mapstructure:"urls"`
}

type AppConf struct {
	ID           string       `yaml:"id" mapstructure:"id"`
	LoggerConf   LoggerConf   `yaml:"log" mapstructure:"log"`
	ServerConf   ServerConf   `yaml:"server" mapstructure:"server"`
	TenantConf   TenantConf   `yaml:"tenant" mapstructure:"tenant"`
	SubsConf     SubsConf     `yaml:"subs" mapstructure:"subs"`
	BackendConf  BackendConf  `yaml:"backend" mapstructure:"backend"`
	CameraDBConf CameraDBConf `yaml:"camera-db" mapstructure:"camera-db"`
	EtcdConf     EtcdConf     `yaml:"etcd" mapstructure:"etcd"`
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
	enc := yaml.NewEncoder(os.Stdout)
	defer enc.Close()
	enc.SetIndent(2)
	enc.Encode(defaultConfig())
}

func defaultConfig() AppConf {
	return AppConf{
		ID: "dcp",

		LoggerConf: LoggerConf{
			Level:   "debug",
			LogFile: "console",
		},

		ServerConf: ServerConf{
			Addr:       "0.0.0.0:10005",
			CertFile:   "repo/server.crt",
			KeyFile:    "repo/server.key",
			ForceHttps: 1,
		},
		TenantConf:TenantConf{
			TenantID:"0000000000",
		},
		SubsConf: SubsConf{
			Muti: 1,
			Main: "https://127.0.0.1:10005/pf/upload",

			Backups: []string{
				"https://127.0.0.1:10005/pf/upload",
			},
		},
		CameraDBConf: CameraDBConf{
			CsvCameraDB: CsvCameraDB{
				CameraUser: "ApiAdmin",
				CameraPwd:  "AAaa1234%%",
			},
		},

		BackendConf: BackendConf{
			Use: "none",
			TaosDBConf: TaosDBConf{
				Addr:     "127.0.0.1",
				Port:     6041,
				Protocol: "ws",
				Username: "root",
				Password: "taosdata",
				DBName:   "taosdb",
			},
		},

		EtcdConf: EtcdConf{
			URLs: []string{"127.0.0.1:2379"},
		},
	}

}
