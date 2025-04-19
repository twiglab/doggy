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
	Short: "查询设备元数据订阅信息",
	Long:  `使用配置文件启动dcp服务`,
	Run: func(cmd *cobra.Command, args []string) {
		metaGet(cmd, args)
	},
	Example: "dcp meta 1.2.3.4",
}

var (
	username string
	password string
	addr     string
)

func init() {
	MetaCmd.PersistentFlags().StringVar(&username, "username", "ApiAdmin", "摄像头认证用户名")
	MetaCmd.PersistentFlags().StringVar(&password, "password", "Aaa1234%%", "摄像头认证用户密码")
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
