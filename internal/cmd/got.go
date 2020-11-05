package cmd

import (
	"fmt"

	"got/internal/cmd/readtree"

	"got/internal/cmd/writetree"

	"got/internal/cmd/updateindex"

	"got/internal/cmd/catfile"

	"github.com/spf13/cobra"

	"got/internal/cmd/hashobject"
	gotInit "got/internal/cmd/init"
)

var GotCmd = &cobra.Command{
	Short: "The basics of git implemented in go",
}

func init() {
	GotCmd.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println("Nothing")
	}
	GotCmd.AddCommand(gotInit.Cmd)
	GotCmd.AddCommand(hashobject.Cmd)
	GotCmd.AddCommand(catfile.Cmd)
	GotCmd.AddCommand(updateindex.Cmd)
	GotCmd.AddCommand(writetree.Cmd)
	GotCmd.AddCommand(readtree.Cmd)
}
