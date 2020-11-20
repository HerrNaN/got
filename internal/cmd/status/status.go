package status

import (
	"fmt"

	"github.com/spf13/cobra"

	"got/internal/got/filesystem"
)

var Cmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of index and working tree",
	Args:  cobra.NoArgs,
	Run:   runStatus,
}

func runStatus(cmd *cobra.Command, args []string) {
	g, err := filesystem.NewGot()
	if err != nil {
		fmt.Println(err)
		return
	}
	head, err := g.Head()
	if err != nil {
		fmt.Println(err)
		return
	}
	s, err := g.Status()
	if err != nil {
		fmt.Println(err)
		return
	}

	if head == nil {
		fmt.Println("No commits yet")
	} else {
		fmt.Printf("HEAD at %s\n", *head)
	}

	fmt.Println(s)
}
