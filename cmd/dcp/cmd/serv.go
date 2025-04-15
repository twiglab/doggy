/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/twiglab/doggy/cmd/dcp/conf"
	"github.com/twiglab/doggy/hx"

	"github.com/twiglab/doggy/pf"
)

// servCmd represents the serv command
var servCmd = &cobra.Command{
	Use:   "serv",
	Short: "启动dcp服务",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		serv(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(servCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// servCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// servCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func serv(cmd *cobra.Command, args []string) {
	config := conf.App{}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	h := &pf.HoloHandle{
		Conf: pf.PlatformConfig{
			Address: config.PlatformConfig.Address,
			Port:    config.PlatformConfig.Port,
		},

		Resolver: &pf.DeviceResolve{
			Username: config.CommonDeviceConfig.Username,
			Password: config.CommonDeviceConfig.Password,
		},
	}

	pfHandle := pf.PlatformHandle(h)

	mux := chi.NewMux()
	mux.Use(middleware.Logger, middleware.Recoverer)
	mux.Mount("/", pfHandle)
	mux.Mount("/pf", pfHandle)

	svr := hx.NewServ().SetAddr(config.Web.Addr).SetHandler(mux)
	if err := runSvr(svr, config.Web.CertFile, config.Web.KeyFile); err != nil {
		log.Fatal(err)
	}
}

func runSvr(s *hx.Svr, cert, key string) error {
	if cert == "" || key == "" {
		return s.Run()
	}
	return s.RunT(cert, key)
}
