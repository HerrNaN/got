package add

import (
	"fmt"

	"github.com/spf13/cobra"

	"got/internal/got/filesystem"
)

var Cmd = &cobra.Command{
	Use:   "add <file>",
	Short: "Add file into the index",
	Args:  cobra.ExactArgs(1),
	Run:   runAdd,
}

func runAdd(cmd *cobra.Command, args []string) {
	g, err := filesystem.NewGot()
	if err != nil {
		fmt.Println(err)
		return
	}
	f := args[0]
	g.AddToIndex(f)
}
