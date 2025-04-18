package meta

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/twiglab/doggy/holo"
)

// MetaCmd represents the serv command
var MetaPostCmd = &cobra.Command{
	Use:   "post",
	Short: "设置设备元数据订阅信息",
	Long:  `使用配置文件启动dcp服务`,
	Run: func(cmd *cobra.Command, args []string) {
		metaPost(cmd, args)
	},
	Example: "dcp meta post 1.2.3.4",
}

var metadataURL string

func init() {
	MetaCmd.AddCommand(MetaPostCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// servCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// servCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// MetaCmd.PersistentFlags().StringVar(&username, "username", "ApiAdmin", "摄像头认证用户名")
	// MetaCmd.PersistentFlags().StringVar(&password, "password", "Aaa1234%%", "摄像头认证用户密码")
	MetaPostCmd.Flags().StringVar(&metadataURL, "url", "", "平台URL")

}

func metaPost(cmd *cobra.Command, args []string) {
	if net.ParseIP(ip) == nil {
		log.Fatal("no ip")
	}

	dev, _ := holo.OpenDevice(ip, username, password)
	defer dev.Close()

	u, err := url.Parse(metadataURL)
	if err != nil {
		log.Fatal(err)
	}

	port := 80

	if u.Port() != "" {
		var e error
		if port, e = strconv.Atoi(u.Port()); e != nil {
			log.Fatal(e)
		}
	}

	resp, err := dev.PostMetadataSubscription(context.Background(), holo.SubscriptionReq{
		Address:     u.Hostname(),
		Port:        port,
		TimeOut:     0,
		HttpsEnable: 1,
		MetadataURL: metadataURL,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}
