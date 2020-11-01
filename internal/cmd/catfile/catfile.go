package catfile

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"got/internal/objects/disk"
)

var Cmd = &cobra.Command{
	Use:                   "cat-file { -t | -s | -p } object",
	DisableFlagsInUseLine: true,
	Short:                 "Provide content or type and size information for repository objects",
	Args:                  cobra.ExactArgs(1),
}

type flag string

const (
	flagType        = "type"
	flagSize        = "size"
	flagPrettyPrint = "pretty-print"
)

func init() {
	showType := Cmd.Flags().BoolP(flagType, "t", false, "Show type of object")
	showSize := Cmd.Flags().BoolP(flagSize, "s", false, "Show size of object")
	prettyPrintContent := Cmd.Flags().BoolP(flagPrettyPrint, "p", false, "Pretty-print content of object")
	Cmd.Run = func(cmd *cobra.Command, args []string) {
		run(cmd, args, *showType, *showSize, *prettyPrintContent)
	}
}

func run(cmd *cobra.Command, args []string, showType bool, showSize bool, prettyPrint bool) {
	flagUsed, err := flagsAreCompatible(showType, showSize, prettyPrint)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	objs := disk.NewObjects()
	sum := args[0]
	o, err := objs.Get(sum)
	if err != nil {
		fmt.Printf("object %s not found", sum)
		return
	}
	switch flagUsed {
	case flagType:
		fmt.Println(o.Type)
	case flagSize:
		fmt.Println(o.Size)
	case flagPrettyPrint:
		fmt.Println(o.Bs)
	}
}

func flagsAreCompatible(showType, showSize, prettyPrint bool) (flag, error) {
	if !showSize && !showType && !prettyPrint {
		return "", errors.New("one of -s, -t or -p flag must be used")
	}
	if showType && showSize {
		return "", errors.New("-s and -t are not compatible")
	}
	if showType && prettyPrint {
		return "", errors.New("-p and -t are not compatible")
	}
	if showSize && prettyPrint {
		return "", errors.New("-s and -p are not compatible")
	}
	if showSize {
		return flagSize, nil
	}
	if showType {
		return flagType, nil
	}
	return flagPrettyPrint, nil
}
