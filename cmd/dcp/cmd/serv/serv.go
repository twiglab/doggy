package serv

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/twiglab/doggy/hx"
)

var cfgFile string
var vp *viper.Viper = viper.New()

var ServCmd = &cobra.Command{
	Use:   "serv",
	Short: "启动dcp服务",
	Long:  `使用配置文件启动dcp服务`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return servCmd()
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

func servCmd() error {

	ctx := buildAll(context.Background(), vp)
	mux := MainHandle(ctx, pid(vp))
	addr := vp.GetString("server.addr")

	s := hx.NewServer(ctx, addr, mux)

	key := vp.GetString("server.key")
	cert := vp.GetString("server.cert")
	https := vp.GetInt("server.https")

	return runSvr(s, https, cert, key)
}

func runSvr(s *http.Server, forceHttps int, cert, key string) error {
	if forceHttps == 0 {
		return s.ListenAndServe()
	}
	if cert == "" || key == "" {
		return errors.New("no cert and key file")
	}
	return s.ListenAndServeTLS(cert, key)
}
