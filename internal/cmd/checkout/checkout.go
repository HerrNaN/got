package checkout

import (
	"fmt"

	"github.com/spf13/cobra"

	"got/internal/got/filesystem"
)

var Cmd = &cobra.Command{
	Use:   "checkout {<branchname> | -b <newbranch>}",
	Short: "Switch branches",
	Args:  cobra.ExactArgs(1),
}

func init() {
	create := Cmd.Flags().BoolP("", "b", false, "create a new branch and switch to it")
	Cmd.Run = func(cmd *cobra.Command, args []string) {
		runCheckout(cmd, args, *create)
	}
}

func runCheckout(cmd *cobra.Command, args []string, create bool) {
	g, err := filesystem.NewGot()
	if err != nil {
		fmt.Println(err)
		return
	}
	branchName := args[0]
	err = g.Checkout(branchName, create)
	if err != nil {
		fmt.Println(err)
		return
	}
}
