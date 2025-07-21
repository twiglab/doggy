package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/twiglab/doggy/cmd/dcp/cmd/camera"
	"github.com/twiglab/doggy/cmd/dcp/cmd/fake"
	"github.com/twiglab/doggy/cmd/dcp/cmd/meta"
	"github.com/twiglab/doggy/cmd/dcp/cmd/serv"
	"github.com/twiglab/doggy/cmd/dcp/cmd/version"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dcp",
	Short: "(๑•̀ㅂ•́)و✧ 客流平台",
	Long:  `(๑•̀ㅂ•́)و✧`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(serv.ServCmd)
	rootCmd.AddCommand(meta.MetaCmd)
	rootCmd.AddCommand(camera.CameraCmd)
	rootCmd.AddCommand(version.VersionCmd)
	rootCmd.AddCommand(fake.FakeCmd)
}
