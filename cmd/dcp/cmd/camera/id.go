package camera

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

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
	dev, _ := holo.OpenDevice("", addr, user, pwd)
	defer dev.Close()

	resp, err := dev.GetDeviceID(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("---------------\n")

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(resp)
}
