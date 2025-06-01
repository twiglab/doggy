package camera

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/twiglab/doggy/holo"
)

var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "获取相机基本信息",
	Long:  `获取相机基本信息`,
	Run: func(cmd *cobra.Command, args []string) {
		getSysBaseInfo()
	},

	Example: "dcp camera info --addr 1.2.3.4",
}

func init() {
	CameraCmd.AddCommand(InfoCmd)
}

func getSysBaseInfo() {
	dev, _ := holo.OpenDevice("", addr, user, pwd)
	defer dev.Close()

	info, err := dev.GetSysBaseInfo(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(info)
}
