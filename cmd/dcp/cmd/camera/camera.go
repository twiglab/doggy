package camera

import (
	"github.com/spf13/cobra"
)

// CameraCmd represents the serv command
var CameraCmd = &cobra.Command{
	Use:   "camera",
	Short: "相机操作",
	Long:  `相机操作`,
	/*
		Run: func(cmd *cobra.Command, args []string) {
			metaGet(cmd, args)
		},
		Example: "dcp meta --addr 1.2.3.4",
	*/
}

var (
	user string
	pwd  string
	addr string
)

func init() {
	CameraCmd.PersistentFlags().StringVarP(&user, "user", "u", "ApiAdmin", "用户名")
	CameraCmd.PersistentFlags().StringVarP(&pwd, "pwd", "p", "Aaa1234%%", "密码")
	CameraCmd.PersistentFlags().StringVar(&addr, "addr", "127.0.0.1:10007", "相机地址(ip:port)")
}
