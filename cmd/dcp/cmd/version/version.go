package version

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/twiglab/doggy"
)

// DbCmd represents the serv command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "查看版本信息",
	Long:  `查看版本信息`,
	Run: func(cmd *cobra.Command, args []string) {
		version()
	},
	Example: "dcp version",
}

func version() {
	ver := doggy.Version()
	fmt.Println("-------------")
	fmt.Printf("Version: %s\n", ver.Version)
	fmt.Printf("Go Version: %s\n", ver.GoVersion)
	fmt.Printf("OS/Arch: %s\n", ver.OsArch)

	fmt.Printf("Git Commit: %s\n", ver.GitCommit)
	fmt.Printf("Build Time: %s\n", ver.BuildTime)
	fmt.Printf("Backends: %s\n", ver.Backends)
	fmt.Println("-------------")
}
