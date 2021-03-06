package log

import (
	"fmt"
	"strings"

	"got/internal/pkg/terminal"

	"github.com/spf13/cobra"

	"got/internal/got/filesystem"
)

var Cmd = &cobra.Command{
	Use:   "log",
	Short: "List commits that are reachable by following the 'parent' links from HEAD",
	Args:  cobra.NoArgs,
}

func init() {
	n := Cmd.Flags().IntP("number", "n", 0, "show a maximum of n entries")
	Cmd.Run = func(cmd *cobra.Command, args []string) {
		runLog(cmd, args, *n)
	}
}

func runLog(cmd *cobra.Command, args []string, n int) {
	g, err := filesystem.NewGot()
	if err != nil {
		fmt.Println(err)
		return
	}
	log, err := g.Log(n)
	if err != nil {
		fmt.Println(err)
		return
	}
	height, err := terminal.Height()
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(strings.Split(log.String(), "\n")) >= height {
		err := terminal.RunLess(log.String())
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Print(log)
	}

}
