package meta

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/twiglab/doggy/holo"
)

// MetaCmd represents the serv command
var MetaGetCmd = &cobra.Command{
	Use:   "get",
	Short: "查询设备元数据订阅信息",
	Long:  `查询设备元数据订阅信息`,
	Run: func(cmd *cobra.Command, args []string) {
		metaGet()
	},
	Example: "dcp meta get --addr 1.2.3.4",
}

func init() {
	MetaCmd.AddCommand(MetaGetCmd)
}

func metaGet() {
	dev, _ := holo.ConnectDevice(addr, user, pwd)
	defer dev.Close()

	subs, err := dev.GetMetadataSubscription(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(subs)
}
