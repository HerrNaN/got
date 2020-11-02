package writetree

import (
	"fmt"

	"github.com/spf13/cobra"

	"got/internal/got"
	"got/internal/index/file"
	"got/internal/objects/disk"
)

var Cmd = &cobra.Command{
	Use:   "write-tree",
	Short: "Create tree object from the current index",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if !got.IsInitialized() {
			fmt.Println("repository is not initialized")
			return
		}
		i, err := file.ReadFromFile()
		if err != nil {
			i = file.NewIndex()
		}
		g := got.NewGot(disk.NewObjects(), i)
		sum := g.WriteTree()
		fmt.Println(sum)
	},
}
