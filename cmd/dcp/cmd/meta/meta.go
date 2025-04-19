package meta

import (
	"fmt"
	"net"
	"strconv"

	"github.com/spf13/cobra"
)

// MetaCmd represents the serv command
var MetaCmd = &cobra.Command{
	Use:   "meta",
	Short: "设备元数据订阅操作",
	Long:  `设备元数据订阅操作`,
	Run: func(cmd *cobra.Command, args []string) {
		metaGet(cmd, args)
	},
	Example: "dcp meta --addr 1.2.3.4",
}

var (
	user string
	pwd  string
	addr string
)

func init() {
	MetaCmd.PersistentFlags().StringVar(&user, "user", "ApiAdmin", "摄像头认证用户名")
	MetaCmd.PersistentFlags().StringVar(&pwd, "pwd", "Aaa1234%%", "摄像头认证用户密码")
	MetaCmd.PersistentFlags().StringVar(&addr, "addr", "", "摄像头地址含端口(如：1.2.3.4:80)")
}

func verifyAddr(addr string) (host string, port int, err error) {
	var portS string

	if host, portS, err = net.SplitHostPort(addr); err != nil {
		return
	}

	if port, err = strconv.Atoi(portS); err != nil {
		return
	}

	if net.ParseIP(host) == nil {
		err = fmt.Errorf("bad host: %s", host)
	}

	return
}
