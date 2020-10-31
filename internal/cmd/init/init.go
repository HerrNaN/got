package init

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"got/internal/got"
)

var Cmd = &cobra.Command{
	Use:   "init",
	Short: "Initialized a got repository",
	Run: func(cmd *cobra.Command, args []string) {
		wd, _ := os.Getwd()
		_, err := os.Stat(got.GotRootDir)
		if os.IsNotExist(err) {
			os.Mkdir(got.GotRootDir, os.ModePerm)
			fmt.Printf("Repository initialized in %s/%s\n", wd, got.GotRootDir)
			return
		}
		fmt.Printf("Repository already exists for %s\n", wd)
	},
}
