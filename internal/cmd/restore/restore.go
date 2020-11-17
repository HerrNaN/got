package restore

import (
	"fmt"

	"github.com/spf13/cobra"

	"got/internal/got/filesystem"
)

var Cmd = &cobra.Command{
	Use:   "restore [--staged] <file>",
	Short: "Restore specified file in the working tree from HEAD or Index",
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

	filename := args[0]

	if staged {
		err := g.Unstage(filename)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		//g.Discard(filename)
	}
}
