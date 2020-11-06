package status

import (
	"fmt"
	"os"

	"github.com/gookit/color"

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
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("couldn't get working directory")
		return
	}
	staged, unstaged, err := g.Status(wd)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("On branch not implemented")
	fmt.Println("Changes to be committed:")
	for _, s := range staged {
		color.Green.Printf("        modified:   %s\n", s)
	}
	fmt.Println()
	fmt.Println("Changes not staged for commit:")
	for _, u := range unstaged {
		color.Red.Printf("        modified:   %s\n", u)
	}
	fmt.Println()

}
