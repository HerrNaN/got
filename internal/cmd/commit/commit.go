package commit

import (
	"fmt"

	"github.com/spf13/cobra"

	"got/internal/got/filesystem"
)

var Cmd = &cobra.Command{
	Use:   "commit -m message",
	Short: "Commit changes in the index",
	Args:  cobra.NoArgs,
}

func init() {
	message := Cmd.Flags().StringP("message", "m", "", "The message to describe the commit")
	Cmd.Run = func(cmd *cobra.Command, args []string) {
		run(cmd, args, *message)
	}
}

func run(cmd *cobra.Command, args []string, message string) {
	g, err := filesystem.NewGot()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = g.Commit(message)
	if err != nil {
		fmt.Println(err)
		return
	}
}
