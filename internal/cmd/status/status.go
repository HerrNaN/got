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
	headType, err := g.HeadType()
	if err != nil {
		fmt.Println(err)
		return
	}
	s, err := g.Status()
	if err != nil {
		fmt.Println(err)
		return
	}

	switch headType {
	case filesystem.HeadTypeEmpty:
		fmt.Println("No commits yet")
	case filesystem.HeadTypeRef:
		ref, err := g.HeadAsRef()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("On branch %s\n", ref.Name())
	}
	fmt.Println(s)
}
