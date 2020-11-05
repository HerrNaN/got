package hashobject

import (
	"fmt"

	"got/internal/got/filesystem"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "hash-object [-w] <file>...",
	Short: "Compute object ID and optionally creates a blob from a file",
	Args:  cobra.ExactArgs(1),
}

func init() {
	write := Cmd.Flags().BoolP("write", "w", false, "Actually write the object into the database")
	Cmd.Run = func(cmd *cobra.Command, args []string) {
		run(cmd, args, *write)
	}
}

func run(cmd *cobra.Command, args []string, write bool) {
	g, err := filesystem.NewGot()
	if err != nil {
		fmt.Println(err)
		return
	}
	f := args[0]
	sum := g.HashFile(f, write)
	fmt.Println(sum)
}
