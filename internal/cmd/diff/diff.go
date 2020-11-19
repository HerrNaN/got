package diff

import (
	"fmt"

	"github.com/spf13/cobra"

	"got/internal/got/filesystem"
)

var Cmd = &cobra.Command{
	Use:   "diff <path>...",
	Short: "Show the diff between your working directory and the the index",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runDiff(cmd, args)
	},
}

func runDiff(cmd *cobra.Command, args []string) {
	g, err := filesystem.NewGot()
	if err != nil {
		fmt.Println(err)
		return
	}
	diffs, err := g.DiffPath(args...)
	if err != nil {
		fmt.Println(err)
		return
	}
	if diffs != "" {
		fmt.Println(diffs)
	}
}
