package camera

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/twiglab/doggy/holo"
)

// MetaCmd represents the serv command
var RebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "重启相机",
	Long:  `重启相机`,
	Run: func(cmd *cobra.Command, args []string) {
		reboot()
	},

	Example: "dcp camera reboot --addr 1.2.3.4",
}

func init() {
	CameraCmd.AddCommand(RebootCmd)
}

func reboot() {
	dev, _ := holo.ConnectDevice(addr, user, pwd)
	defer dev.Close()

	resp, err := dev.Reboot(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(resp)
}
