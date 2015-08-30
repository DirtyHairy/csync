package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/DirtyHairy/csync/cli"
	"github.com/DirtyHairy/csync/cli/csync/push"
	"github.com/DirtyHairy/csync/cli/csync/repo"
	"github.com/DirtyHairy/csync/lib"
)

func Execute(name string, arguments []string) error {
	flags := flag.NewFlagSet(name, flag.ExitOnError)

	commands := []cli.Command{
		{
			Name:        "push",
			Description: "push changes (unidirectional)",
			Dispatcher:  cli.CommandDispatchFunction(push.Execute),
		},
		{
			Name:        "repo",
			Description: "manage repositories",
			Dispatcher:  cli.CommandDispatchFunction(repo.Execute),
		},
	}

	usage := cli.Usage{
		Usage:       fmt.Sprintf("usage: %s [options] command [command options]", name),
		Description: "Sync data between local and remote repositories.",
		Flags:       flags,
		Commands:    commands,
	}

	flags.Usage = usage.Print

	showVersion := false
	flags.BoolVar(&showVersion, "version", false, "show program version")

	flags.Parse(arguments)

	if showVersion {
		fmt.Printf("csync version %s\n", lib.VERSION)
		return nil
	}

	err := cli.DispatchCommand(commands, flags.Args(), name)
	if _, ok := err.(cli.CommandDispatchError); ok && err != nil {
		fmt.Printf("ERROR: %v\n", err)
		usage.Print()
		os.Exit(1)
	}

	return err
}
