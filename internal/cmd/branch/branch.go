package branch

import (
	"fmt"

	"github.com/spf13/cobra"

	"got/internal/got/filesystem"
)

var Cmd = &cobra.Command{
	Use: `branch <newbranch>
   branch --list`,
	Short:                 "Create and list branches",
	DisableFlagsInUseLine: true,
}

func init() {
	list := Cmd.Flags().Bool("list", false, "list all branches")
	Cmd.Run = func(cmd *cobra.Command, args []string) {
		runBranch(cmd, args, *list)
	}
}

func runBranch(cmd *cobra.Command, args []string, list bool) {
	g, err := filesystem.NewGot()
	if err != nil {
		fmt.Println(err)
		return
	}
	if list {
		branches, err := g.ListBranches()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print(branches)
		return
	}
	newBranch := args[0]
	err = g.CreateBranch(newBranch)
	if err != nil {
		fmt.Println(err)
		return
	}
}
