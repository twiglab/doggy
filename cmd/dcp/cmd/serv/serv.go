package serv

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var backendOnly bool

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
	ServCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file")
	ServCmd.Flags().BoolVarP(&backendOnly, "backend-only", "b", false, "backend only")
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
	fmt.Println("backend-only:", backendOnly)
	fmt.Println("backend:", conf.BackendConf.Use)
	fmt.Println("sub :", conf.SubsConf.Muti)
	fmt.Println("--------------------")
}

func backendMux(conf AppConf) (context.Context, http.Handler) {
	_, ctx := buildRootlogger(context.Background(), conf)
	_, ctx = buildBackend(ctx, conf)
	mux := BackendHandler(ctx, conf)
	return ctx, mux
}

func fullMux(conf AppConf) (context.Context, http.Handler) {
	ctx := buildAll(context.Background(), conf)
	mux := FullHandler(ctx, conf)
	return ctx, mux
}

func servCmd() {
	conf := AppConf{}

	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	printConf(conf)

	var (
		mux http.Handler

		// ctx context.Context
	)
	if backendOnly {
		_, mux = backendMux(conf)
	} else {
		_, mux = fullMux(conf)
		/*
			if conf.JobConf.Enable != 0 {
				s := buildAllJob(ctx, conf)
				s.Start()
			}
		*/
	}

	svr := &http.Server{
		Addr:         conf.ServerConf.Addr,
		Handler:      mux,
		IdleTimeout:  90 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := runSvr(svr, conf.ServerConf); err != nil {
		log.Fatal(err)
	}
}

func runSvr(s *http.Server, sc ServerConf) error {
	if sc.ForceHttps == 0 {
		return s.ListenAndServe()
	}
	if sc.CertFile == "" || sc.KeyFile == "" {
		return errors.New("no cert and key file")
	}
	return s.ListenAndServeTLS(sc.CertFile, sc.KeyFile)
}
