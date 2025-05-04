package camera

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/twiglab/doggy/holo"
)

var SetIDCmd = &cobra.Command{
	Use:   "setid",
	Short: "设置机ID",
	Long:  `设置相机ID`,
	Run: func(cmd *cobra.Command, args []string) {
		setid(args)
	},

	Example: "dcp camera setid <uuid1=deviceID1> [uuid2=deviceID2 ...] --addr 1.2.3.4",
}

func init() {
	CameraCmd.AddCommand(SetIDCmd)
}

func fromString(s string) holo.DeviceID {
	uuid, deviceID, ok := strings.Cut(s, ",")
	if !ok {
		log.Fatal("param error")
	}

	return holo.DeviceID{
		UUID:     strings.TrimSpace(uuid),
		DeviceID: strings.TrimSpace(deviceID),
	}
}

func setid(args []string) {
	dev, _ := holo.OpenDevice(addr, user, pwd)
	defer dev.Close()

	var ids []holo.DeviceID
	for _, s := range args {
		ids = append(ids, fromString(s))
	}

	idList := holo.DeviceIDList{IDs: ids}

	resp, err := dev.PutDeviceID(context.Background(), idList)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusString)

}
