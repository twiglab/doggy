package meta

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/spf13/cobra"
	"github.com/twiglab/doggy/holo"
)

// MetaCmd represents the serv command
var MetaGetCmd = &cobra.Command{
	Use:   "get",
	Short: "查询设备元数据订阅信息",
	Long:  `使用配置文件启动dcp服务`,
	Run: func(cmd *cobra.Command, args []string) {
		metaGet(cmd, args)
	},
	Example: "dcp meta get --addr 1.2.3.4",
}

func init() {
	MetaCmd.AddCommand(MetaGetCmd)
}

func metaGet(cmd *cobra.Command, args []string) {
	if _, err := net.ResolveTCPAddr("", addr); err != nil {
		log.Fatal("bad addr")
	}

	dev, _ := holo.OpenDevice(addr, user, pwd)
	defer dev.Close()

	subs, err := dev.GetMetadataSubscription(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, sub := range subs.Subscripions {
		fmt.Printf("ID: %d, url: %s", sub.ID, sub.MetadataURL)
	}
}
