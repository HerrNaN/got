package readtree

import (
	"fmt"

	"got/internal/objects"

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
		id, err := objects.IdFromString(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		err = g.ReadTree(id)
		if err != nil {
			fmt.Println(err)
		}
	},
}
