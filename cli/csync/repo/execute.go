package repo

import (
	"flag"
	"fmt"
	"os"

	"github.com/DirtyHairy/csync/cli"
)

func Execute(name string, arguments []string) error {
	flags := flag.NewFlagSet(name, flag.ExitOnError)

	commands := []cli.Command{
		{
			Name:        "list",
			Description: "list repositories",
		},
		{
			Name:        "edit",
			Description: "edit repositories",
		},
		{
			Name:        "add",
			Description: "add an existing repository",
		},
		{
			Name:        "create",
			Description: "create a new repository",
		},
		{
			Name:        "remove",
			Description: "remove a repository",
		},
	}

	usage := cli.Usage{
		Usage:       fmt.Sprintf("usage: %s [options] command [command options]", name),
		Description: "Manage csync repositories.",
		Flags:       flags,
		Commands:    commands,
	}

	flags.Usage = usage.Print

	flags.Parse(arguments)

	err := cli.DispatchCommand(usage.Commands, flags.Args(), name)
	if _, ok := err.(cli.CommandDispatchError); ok && err != nil {
		fmt.Printf("ERROR: %v\n", err)
		usage.Print()
		os.Exit(1)
	}

	return err
}
