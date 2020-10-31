package cmd

import (
	"fmt"

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
}
