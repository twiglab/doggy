/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package serv

import (
	"fmt"
	"log"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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
}

func init() {
	cobra.OnInitialize(initConfig)
	//rootCmd.AddCommand(ServCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// servCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// servCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	ServCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dcp.yaml)")

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

		// Search config in home directory with name ".dcp" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".dcp")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func serv(cmd *cobra.Command, args []string) {
	config := App{}

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

	svr := hx.NewServ().SetAddr(config.ServerConf.Addr).SetHandler(mux)
	if err := runSvr(svr, config.ServerConf.CertFile, config.ServerConf.KeyFile); err != nil {
		log.Fatal(err)
	}
}

func runSvr(s *hx.Svr, cert, key string) error {
	if cert == "" || key == "" {
		return s.Run()
	}
	return s.RunTLS(cert, key)
}
