package serv

import (
	"fmt"
	"log"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/telemetry"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/twiglab/doggy/hx"

	"github.com/twiglab/doggy/pf"
)

var cfgFile string

// ServCmd represents the serv command
var ServCmd = &cobra.Command{
	Use:   "serv",
	Short: "启动dcp服务",
	Long:  `使用配置文件启动dcp服务`,
	Run: func(cmd *cobra.Command, args []string) {
		serv(cmd, args)
	},

	PreRun: func(cmd *cobra.Command, args []string) {
		telemetry.Start(telemetry.Config{
			ReportCrashes: true,
		})
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	ServCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (default is dcp.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name "dcp" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("dcp")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func printConf(conf AppConf) {
	fmt.Println("--------------------")
	fmt.Println("( ͡° ᴥ ͡° ʋ)")
	fmt.Println("--------------------")
}

func serv(cmd *cobra.Command, args []string) {
	conf := AppConf{}

	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	printConf(conf)

	process := pf.NewSimpleProcess(conf.PlatformConfig.CameraUser, conf.PlatformConfig.CameraPwd)
	// eh := orm.NewEntHandle(MustClient(config.DB.DSN))

	h := &pf.Handle{
		Conf: MustPfConf(conf.PlatformConfig),

		Resolver:       process,
		Register:       process,
		CountHandler:   process,
		DensityHandler: process,
	}

	pfHandle := pf.PlatformHandle(h)

	mux := chi.NewMux()
	mux.Use( /* middleware.Logger, */ middleware.Recoverer)
	mux.Mount("/", pfHandle)
	mux.Mount("/pf", pfHandle)

	svr := hx.NewServ().SetAddr(conf.ServerConf.Addr).SetHandler(mux)
	if err := runSvr(svr, conf.ServerConf.CertFile, conf.ServerConf.KeyFile); err != nil {
		log.Fatal(err)
	}
}

func runSvr(s *hx.Svr, cert, key string) error {
	if cert == "" || key == "" {
		return s.Run()
	}
	return s.RunTLS(cert, key)
}
