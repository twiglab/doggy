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
var MetaCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "清除设备元数据订阅信息",
	Long:  `清除设备元数据订阅信息`,
	Run: func(cmd *cobra.Command, args []string) {
		metaCLean()
	},
	Example: "dcp meta clean --addr 1.2.3.4",
}

func init() {
	MetaCmd.AddCommand(MetaCleanCmd)
}

func metaCLean() {
	var (
		err error
	)

	dev, _ := holo.ConnectDevice(addr, user, pwd)
	defer dev.Close()

	var resp *holo.CommonResponse
	resp, err = dev.DeleteMetadataSubscription(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	if err := resp.Err(); err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(resp)
}
