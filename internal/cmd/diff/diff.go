package diff

import (
	"fmt"

	"github.com/spf13/cobra"

	"got/internal/got/filesystem"
)

var Cmd = &cobra.Command{
	Use:   "diff [--cached] <path>",
	Short: "Show the diff between your working directory and the the index",
	Args:  cobra.ExactArgs(1),
}

func init() {
	cached := Cmd.Flags().Bool("cached", false, "Show diff between your index and HEAD")
	Cmd.Run = func(cmd *cobra.Command, args []string) {
		runDiff(cmd, args, *cached)
	}
}

func runDiff(cmd *cobra.Command, args []string, cached bool) {
	g, err := filesystem.NewGot()
	if err != nil {
		fmt.Println(err)
		return
	}

	var diffs string
	if cached {
		diffs, err = g.DiffIndexPath(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		diffs, err = g.DiffPath(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if diffs != "" {
		fmt.Print(diffs)
	}
}
