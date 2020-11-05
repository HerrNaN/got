package updateindex

import (
	"fmt"

	"got/internal/got/filesystem"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "update-index [--add] file",
	Short: "Register file contents in the working tree to the index",
	Args:  cobra.RangeArgs(1, 2),
}

func init() {
	add := Cmd.Flags().Bool("add", false, "If a specified file isn’t in the index already then it’s added. Default behaviour is to ignore new files.")
	Cmd.Run = func(cmd *cobra.Command, args []string) {
		run(cmd, args, *add)
	}
}

func run(cmd *cobra.Command, args []string, add bool) {
	g, err := filesystem.NewGot()
	if err != nil {
		fmt.Println(err)
		return
	}
	filename := args[0]
	if !add && !g.Index.HasEntryFor(filename) {
		return
	}
	g.AddToIndex(filename)
}
