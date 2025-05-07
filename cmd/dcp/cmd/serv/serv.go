package serv

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/twiglab/doggy/hx"
	"github.com/twiglab/doggy/idb"
	"github.com/twiglab/doggy/job"
	"github.com/twiglab/doggy/orm"

	"github.com/twiglab/doggy/pf"
)

var cfgFile string

// ServCmd represents the serv command
var ServCmd = &cobra.Command{
	Use:   "serv",
	Short: "启动dcp服务",
	Long:  `使用配置文件启动dcp服务`,
	Run: func(cmd *cobra.Command, args []string) {
		servCmd()
	},
}

var plan int = -1

func init() {
	cobra.OnInitialize(initConfig)
	ServCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (default is dcp.yaml)")
	ServCmd.Flags().IntVar(&plan, "plan", -1, "指定方案")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Search config in home directory with name "dcp" (without extension).
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
	fmt.Println(plan)
	fmt.Println("conf.Plan ", conf.Plan)
	fmt.Println("--------------------")
}

func servCmd() {
	conf := AppConf{}

	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	printConf(conf)

	var h *pf.Handle
	crontab := job.NewCron()

	if plan >= 0 {
		conf.Plan = plan
	}

	switch conf.Plan {
	case 1:
		idb3 := idb.NewIdbPoint(MustIdb(conf.InfluxDBConf))
		eh := orm.NewEntHandle(MustEntClient(conf.DBConf))
		fixUser := &pf.FixUserDeviceResolve{User: conf.FixUserConf.CameraUser, Pwd: conf.FixUserConf.CameraPwd}

		keeplive := &pf.KeepLiveJob{
			DeviceLoader:   eh,
			DeviceResolver: fixUser,

			MetadataURL: conf.AutoRegConf.MetadataURL,
			Addr:        conf.AutoRegConf.Addr,
			Port:        conf.AutoRegConf.Port,
		}

		crontab.AddJob(conf.JobConf.Keeplive, keeplive)
		autoSub := &pf.AutoSub{
			DeviceResolver: fixUser,
			UploadHandler:  eh,

			MetadataURL: conf.AutoRegConf.MetadataURL,
			Addr:        conf.AutoRegConf.Addr,
			Port:        conf.AutoRegConf.Port,
		}
		h = pf.NewHandle(
			pf.WithCountHandler(idb3),
			pf.WithDensityHandler(idb3),
			pf.WithDeviceRegister(autoSub),
		)
	case 2:
		idb3 := idb.NewIdbPoint(MustIdb(conf.InfluxDBConf))
		h = pf.NewHandle(
			pf.WithCountHandler(idb3),
			pf.WithDensityHandler(idb3),
		)
	case 5:
		eh := orm.NewEntHandle(MustEntClient(conf.DBConf))
		fixUser := &pf.FixUserDeviceResolve{User: conf.FixUserConf.CameraUser, Pwd: conf.FixUserConf.CameraPwd}

		keeplive := &pf.KeepLiveJob{
			DeviceLoader:   eh,
			DeviceResolver: fixUser,

			MetadataURL: conf.AutoRegConf.MetadataURL,
			Addr:        conf.AutoRegConf.Addr,
			Port:        conf.AutoRegConf.Port,
		}

		crontab.AddJob(conf.JobConf.Keeplive, keeplive)
		autoSub := &pf.AutoSub{
			DeviceResolver: fixUser,
			UploadHandler:  eh,

			MetadataURL: conf.AutoRegConf.MetadataURL,
			Addr:        conf.AutoRegConf.Addr,
			Port:        conf.AutoRegConf.Port,
		}
		h = pf.NewHandle(
			pf.WithDeviceRegister(autoSub),
		)
	default:
		h = pf.NewHandle()
	}

	pfHandle := pf.PlatformHandle(h)

	mux := chi.NewMux()
	mux.Use(middleware.Recoverer)
	mux.Mount("/", pfHandle)
	mux.Mount("/pf", pfHandle)
	mux.Mount("/debug", middleware.Profiler())

	crontab.Start()

	svr := hx.NewServ().SetAddr(conf.ServerConf.Addr).SetHandler(mux)
	if err := runSvr(svr, conf.ServerConf); err != nil {
		log.Fatal(err)
	}
}

func runSvr(s *hx.Svr, sc ServerConf) error {
	if sc.ForceHttps == 0 {
		return s.Run()
	}
	if sc.CertFile == "" || sc.KeyFile == "" {
		return errors.New("no cert and key file")
	}
	return s.RunTLS(sc.CertFile, sc.KeyFile)
}
