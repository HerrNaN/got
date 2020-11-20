package catfile

import (
	"errors"
	"fmt"

	"got/internal/got/filesystem"

	"github.com/spf13/cobra"

	"got/internal/objects"
)

var Cmd = &cobra.Command{
	Use:                   "cat-file { -t | -p } object",
	DisableFlagsInUseLine: true,
	Short:                 "Provide content or type and size information for repository objects",
	Args:                  cobra.ExactArgs(1),
}

type flag string

const (
	flagType        = "type"
	flagPrettyPrint = "pretty-print"
)

func init() {
	showType := Cmd.Flags().BoolP(flagType, "t", false, "Show type of object")
	prettyPrintContent := Cmd.Flags().BoolP(flagPrettyPrint, "p", false, "Pretty-print content of object")
	Cmd.Run = func(cmd *cobra.Command, args []string) {
		run(cmd, args, *showType, *prettyPrintContent)
	}
}

func run(cmd *cobra.Command, args []string, showType bool, prettyPrint bool) {
	flagUsed, err := flagsAreCompatible(showType, prettyPrint)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	g, err := filesystem.NewGot()
	if err != nil {
		fmt.Println(err)
		return
	}
	id, err := objects.IdFromString(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	t, err := g.Objects.TypeOf(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	var o objects.Object
	switch t {
	case objects.TypeBlob:
		o, _ = g.Objects.GetBlob(id)
	case objects.TypeTree:
		o, _ = g.Objects.GetTree(id)
	default:
		fmt.Println("no object found")
		return
	}
	switch flagUsed {
	case flagType:
		fmt.Println(o.Type())
	case flagPrettyPrint:
		fmt.Println(o.Content())
	}
}

func flagsAreCompatible(showType, prettyPrint bool) (flag, error) {
	if !showType && !prettyPrint {
		return "", errors.New("one of -t or -p flag must be used")
	}
	if showType && prettyPrint {
		return "", errors.New("-t and -p are not compatible")
	}
	if showType {
		return flagType, nil
	}
	return flagPrettyPrint, nil
}
