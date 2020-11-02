package updateindex

import (
	"fmt"

	"got/internal/got"

	"github.com/spf13/cobra"

	"got/internal/index/file"
	"got/internal/objects/disk"
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
	if !got.IsInitialized() {
		fmt.Println("repository is not initialized")
		return
	}
	i, err := file.ReadFromFile()
	if err != nil {
		i = file.NewIndex()
	}
	g := got.NewGot(disk.NewObjects(), i)
	filename := args[0]
	if !add && !g.Index.HasEntryFor(filename) {
		return
	}
	sum := g.HashFile(filename, true)
	g.AddToIndex(sum, filename)
}
