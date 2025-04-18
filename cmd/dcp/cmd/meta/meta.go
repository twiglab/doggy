package meta

import (
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
	ip string
)

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// servCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// servCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	MetaCmd.PersistentFlags().StringVar(&username, "username", "ApiAdmin", "摄像头认证用户名")
	MetaCmd.PersistentFlags().StringVar(&password, "password", "Aaa1234%%", "摄像头认证用户密码")
	MetaCmd.PersistentFlags().StringVar(&ip, "ip", "", "摄像头ip地址")

}
