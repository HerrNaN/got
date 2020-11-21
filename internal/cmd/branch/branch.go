package branch

import (
	"fmt"

	"github.com/spf13/cobra"

	"got/internal/got/filesystem"
)

var Cmd = &cobra.Command{
	Use:   "branch <newbranch>",
	Short: "Create branches",
	Args:  cobra.ExactArgs(1),
}

func init() {
	Cmd.Run = func(cmd *cobra.Command, args []string) {
		runBranch(cmd, args)
	}
}

func runBranch(cmd *cobra.Command, args []string) {
	g, err := filesystem.NewGot()
	if err != nil {
		fmt.Println(err)
		return
	}
	newBranch := args[0]
	err = g.CreateBranch(newBranch)
	if err != nil {
		fmt.Println(err)
		return
	}
}
