package hashobject

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"

	"got/internal/objects"

	"got/internal/objects/disk"
)

var Cmd = &cobra.Command{
	Use:   "hash-object [-w] <file>...",
	Short: "Compute object ID and optionally creates a blob from a file",
	Args:  cobra.ExactArgs(1),
}

func init() {
	write := Cmd.Flags().BoolP("write", "w", false, "Actually write the object into the database")
	Cmd.Run = func(cmd *cobra.Command, args []string) {
		run(cmd, args, *write)
	}

}

func run(cmd *cobra.Command, args []string, write bool) {
	f := args[0]
	bs, _ := ioutil.ReadFile(f)
	objs := disk.NewObjects()
	sum := objs.HashObject(bs, write, objects.TypeBlob)
	fmt.Println(sum)
}
