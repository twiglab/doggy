package meta

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/twiglab/doggy/cmd/dcp/utils"
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
	Example: "dcp meta post --addr 1.2.3.4 --url https://5.6.7.8:8080/meta/upload",
}

var metadataURL string

func init() {
	MetaCmd.AddCommand(MetaPostCmd)
	MetaPostCmd.Flags().StringVar(&metadataURL, "url", "", "平台URL")
}

func metaPost(cmd *cobra.Command, args []string) {
	var (
		err  error
		port int = 80
		url  *url.URL
	)

	if _, _, err = utils.VerifyAddr(addr); err != nil {
		log.Fatal("no ip")
	}

	dev, _ := holo.OpenDevice(addr, user, pwd)
	defer dev.Close()

	if url, err = url.Parse(metadataURL); err != nil {
		log.Fatal(err)
	}

	if url.Port() != "" {
		if port, err = strconv.Atoi(url.Port()); err != nil {
			log.Fatal(err)
		}
	}

	var resp *holo.CommonResponseID
	resp, err = dev.PostMetadataSubscription(context.Background(), holo.SubscriptionReq{
		Address:     url.Hostname(),
		Port:        port,
		TimeOut:     0,
		HttpsEnable: 1,
		MetadataURL: metadataURL,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %d, code: %d, str: %s\n", resp.ID, resp.CommonResponse.StatusCode, resp.CommonResponse.StatusString)
}
