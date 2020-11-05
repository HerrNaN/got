package init

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"got/internal/got"
	"got/internal/got/filesystem"
)

var Cmd = &cobra.Command{
	Use:   "init",
	Short: "Initialized a got repository",
	Run: func(cmd *cobra.Command, args []string) {
		wd, _ := os.Getwd()
		err := filesystem.Initialize(wd)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Repository initialized in %s/%s", wd, got.RootDir)
	},
}
