package branch

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"

	"got/internal/got/filesystem"
)

var Cmd = &cobra.Command{
	Use: `branch <newbranch>
   branch --list
   branch -d <branchname>`,
	Short:                 "Create and list branches",
	DisableFlagsInUseLine: true,
}

func init() {
	list := Cmd.Flags().Bool("list", false, "list all branches")
	dlt := Cmd.Flags().BoolP("delete", "d", false, "delete branch")
	Cmd.Args = func(cmd *cobra.Command, args []string) error {
		action, err := actionFromFlags(*list, *dlt)
		if err != nil {
			return err
		}
		switch action {
		case actionCreate:
			if len(args) != 1 {
				return errors.New("wrong number of arguments")
			}
		case actionList:
			if len(args) != 0 {
				return errors.New("wrong number of arguments")
			}
		case actionDelete:
			if len(args) != 1 {
				return errors.New("wrong number of arguments")
			}
		}
		return nil
	}
	Cmd.Run = func(cmd *cobra.Command, args []string) {
		runBranch(cmd, args, *list, *dlt)
	}
}

func runBranch(cmd *cobra.Command, args []string, list bool, delete bool) {
	action, err := actionFromFlags(list, delete)
	if err != nil {
		fmt.Println(err)
		return
	}
	g, err := filesystem.NewGot()
	if err != nil {
		fmt.Println(err)
		return
	}
	switch action {
	case actionCreate:
		newBranch := args[0]
		err = g.CreateBranch(newBranch)
		if err != nil {
			fmt.Println(err)
		}
	case actionList:
		branches, err := g.ListBranches()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print(branches)
	case actionDelete:
		branchName := args[0]
		err := g.DeleteBranch(branchName)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func actionFromFlags(list, delete bool) (action, error) {
	if list && !delete {
		return actionList, nil
	}
	if delete && !list {
		return actionDelete, nil
	}
	if !list && !delete {
		return actionCreate, nil
	}
	return -1, errors.New("flags are not compatible")
}

type action int

const (
	actionList action = iota
	actionDelete
	actionCreate
)
