package camera

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/twiglab/doggy/cmd/dcp/utils"
	"github.com/twiglab/doggy/holo"
)

// MetaCmd represents the serv command
var RebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "重启相机",
	Long:  `重启相机`,
	Run: func(cmd *cobra.Command, args []string) {
		reboot(cmd, args)
	},
	Example: "dcp camera reboot --addr 1.2.3.4",
}

func init() {
	CameraCmd.AddCommand(RebootCmd)
}

func reboot(cmd *cobra.Command, args []string) {

	if _, _, err := utils.VerifyAddr(addr); err != nil {
		log.Fatal("no ip")
	}

	dev, _ := holo.OpenDevice(addr, user, pwd)
	defer dev.Close()

	resp, err := dev.Reboot(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("success, code: %d, msg: %s\n", resp.Code, resp.Msg)
}
