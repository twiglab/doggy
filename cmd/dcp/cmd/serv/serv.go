package serv

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/twiglab/doggy/hx"
)

var cfgFile string

var ServCmd = &cobra.Command{
	Use:   "serv",
	Short: "启动dcp服务",
	Long:  `使用配置文件启动dcp服务`,
	Run: func(cmd *cobra.Command, args []string) {
		servCmd()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	ServCmd.Flags().StringVarP(&cfgFile, "config", "c", "dcp.yaml", "config file")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {

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
	enc := yaml.NewEncoder(os.Stdout)
	defer enc.Close()
	enc.SetIndent(2)
	fmt.Println("--------------------")
	enc.Encode(conf)
	fmt.Println("--------------------")
}

func servCmd() {
	conf := AppConf{}

	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	printConf(conf)

	ctx := buildCtx(conf)

	mux := chi.NewMux()
	mux.Use(middleware.Recoverer)
	mux.Mount("/", pfTestHandle())
	mux.Mount("/pf", pfHandle(ctx, conf))
	mux.Mount("/admin", pageHandle(ctx, conf))

	mux.Mount("/debug", middleware.Profiler())

	mux.Mount("/jsonrpc", outHandle(ctx, conf))

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
