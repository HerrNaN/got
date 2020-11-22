package checkout

import (
	"fmt"

	"github.com/spf13/cobra"

	"got/internal/got/filesystem"
)

var Cmd = &cobra.Command{
	Use:   "checkout <branchname>",
	Short: "Switch branches",
	Args:  cobra.ExactArgs(1),
}

func init() {
	Cmd.Run = func(cmd *cobra.Command, args []string) {
		runCheckout(cmd, args)
	}
}

func runCheckout(cmd *cobra.Command, args []string) {
	g, err := filesystem.NewGot()
	if err != nil {
		fmt.Println(err)
		return
	}
	branchName := args[0]
	err = g.Checkout(branchName)
	if err != nil {
		fmt.Println(err)
		return
	}
}
