package meta

import (
	"github.com/spf13/cobra"
)

// MetaCmd represents the serv command
var MetaCmd = &cobra.Command{
	Use:   "meta",
	Short: "元数据操作",
	Long:  `元数据操作`,
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
	MetaCmd.PersistentFlags().StringVarP(&user, "user", "u", "ApiAdmin", "用户名")
	MetaCmd.PersistentFlags().StringVarP(&pwd, "pwd", "p", "Aaa1234%%", "密码")
	MetaCmd.PersistentFlags().StringVar(&addr, "addr", "", "像机地址(ip:port)")
}
