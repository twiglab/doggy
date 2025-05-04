package camera

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/twiglab/doggy/holo"
)

var IDCmd = &cobra.Command{
	Use:   "id",
	Short: "获取相机ID",
	Long:  `获取相机ID`,
	Run: func(cmd *cobra.Command, args []string) {
		getid()
	},

	Example: "dcp camera id --addr 1.2.3.4",
}

func init() {
	CameraCmd.AddCommand(IDCmd)
}

func getid() {
	dev, _ := holo.OpenDevice(addr, user, pwd)
	defer dev.Close()

	resp, err := dev.GetDeviceID(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	for _, id := range resp.IDs {
		fmt.Printf("uuid = deviceID\n")
		fmt.Printf("---------------\n")
		fmt.Printf("%s = %s\n", id.UUID, id.DeviceID)
	}
}
