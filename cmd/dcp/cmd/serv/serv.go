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
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigType("yaml")
		viper.SetConfigName("dcp")
	}

	viper.AutomaticEnv()

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
	fmt.Println("backend:", conf.BackendConf.Use)
	fmt.Println("muti-sub :", conf.SubsConf.Muti)
	fmt.Println("--------------------")
}

func fullMux(conf AppConf) http.Handler {
	ctx := buildAll(context.Background(), conf)
	mux := FullHandler(ctx, conf)
	return mux
}

func servCmd() {
	var conf AppConf

	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	printConf(conf)

	/*
		mux := fullMux(conf)

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
	*/

	RunWithConf(conf)
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

func RunWithConf(conf AppConf) {
	mux := fullMux(conf)

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
