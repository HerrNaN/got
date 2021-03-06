package add

import (
	"fmt"

	"github.com/spf13/cobra"

	"got/internal/got/filesystem"
)

var Cmd = &cobra.Command{
	Use:   "add <file>...",
	Short: "Add file(s) into the index",
	Args:  cobra.MinimumNArgs(1),
	Run:   runAdd,
}

func runAdd(cmd *cobra.Command, args []string) {
	g, err := filesystem.NewGot()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = g.AddPath(args...)
	if err != nil {
		fmt.Println(err)
	}
}
