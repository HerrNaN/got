package readtree

import (
	"fmt"

	"github.com/spf13/cobra"

	"got/internal/got/filesystem"
)

var Cmd = &cobra.Command{
	Use:   "read-tree <object>",
	Short: "Reads tree information into the index",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		g, err := filesystem.NewGot()
		if err != nil {
			fmt.Println(err)
			return
		}
		sum := args[0]
		err = g.ReadTree(sum)
		if err != nil {
			fmt.Println(err)
		}
	},
}
