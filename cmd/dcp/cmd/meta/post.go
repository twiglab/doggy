package meta

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/twiglab/doggy/holo"
)

// MetaCmd represents the serv command
var MetaPostCmd = &cobra.Command{
	Use:   "post",
	Short: "设置设备元数据订阅信息",
	Long:  `设置设备元数据订阅信息`,
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

func metaPost(_ *cobra.Command, _ []string) {
	var (
		err  error
		port = 80
		url  *url.URL
	)

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

	var resp *holo.CommonResponse
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

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(resp)
}
