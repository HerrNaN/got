package cmd

import (
	"fmt"

	"got/internal/cmd/branch"

	"got/internal/cmd/log"

	"github.com/spf13/cobra"

	"got/internal/cmd/add"
	"got/internal/cmd/catfile"
	"got/internal/cmd/commit"
	"got/internal/cmd/diff"
	"got/internal/cmd/hashobject"
	gotInit "got/internal/cmd/init"
	"got/internal/cmd/readtree"
	"got/internal/cmd/restore"
	"got/internal/cmd/status"
	"got/internal/cmd/updateindex"
	"got/internal/cmd/writetree"
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
	GotCmd.AddCommand(add.Cmd)
	GotCmd.AddCommand(status.Cmd)
	GotCmd.AddCommand(commit.Cmd)
	GotCmd.AddCommand(restore.Cmd)
	GotCmd.AddCommand(diff.Cmd)
	GotCmd.AddCommand(log.Cmd)
	GotCmd.AddCommand(branch.Cmd)
}
