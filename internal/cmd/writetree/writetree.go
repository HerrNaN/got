package writetree

import (
	"fmt"

	"got/internal/got/filesystem"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "write-tree",
	Short: "Create tree object from the current index",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		g, err := filesystem.NewGot()
		if err != nil {
			fmt.Println(err)
			return
		}
		sum, err := g.WriteTree()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(sum)
	},
}
