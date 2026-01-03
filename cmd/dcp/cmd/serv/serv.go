package serv

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var vp *viper.Viper = viper.New()

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
		vp.SetConfigFile(cfgFile)
	} else {
		vp.AddConfigPath(".")
		vp.SetConfigName("dcp")
		vp.SetConfigType("toml")
	}

	vp.AutomaticEnv()

	if err := vp.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", vp.ConfigFileUsed())
	}
}

func fullMux(conf AppConf) http.Handler {
	ctx := buildAll(context.Background(), conf)
	mux := FullHandler(ctx, conf)
	return mux
}

func servCmd() {

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
