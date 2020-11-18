package restore

import (
	"fmt"

	"github.com/spf13/cobra"

	"got/internal/got/filesystem"
)

var Cmd = &cobra.Command{
	Use:   "restore [--staged] <file>",
	Short: "Restore specified file in the working tree from HEAD or Index",
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	staged := Cmd.Flags().Bool("staged", false, "unstage changes from index")
	Cmd.Run = func(cmd *cobra.Command, args []string) {
		restoreRun(cmd, args, *staged)
	}
}

func restoreRun(cmd *cobra.Command, args []string, staged bool) {
	g, err := filesystem.NewGot()
	if err != nil {
		fmt.Println(err)
		return
	}

	if staged {
		err := g.UnstagePath(args...)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err := g.DiscardPath(args...)
		if err != nil {
			fmt.Println(err)
		}
	}
}
