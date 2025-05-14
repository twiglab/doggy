package fake

import (
	"github.com/spf13/cobra"
	"github.com/twiglab/doggy/cmd/dcp/cmd/fake/d3252"
)

// DbCmd represents the serv command
var FakeCmd = &cobra.Command{
	Use:   "fake",
	Short: "模拟器",
	Long:  `模拟器`,
}

func init() {
	FakeCmd.AddCommand(d3252.D3252Cmd)
}
