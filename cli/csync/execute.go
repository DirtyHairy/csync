package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/DirtyHairy/csync/cli"
	"github.com/DirtyHairy/csync/cli/csync/push"
	"github.com/DirtyHairy/csync/cli/csync/repo"
	"github.com/DirtyHairy/csync/lib/environment"
)

func bootstrap() environment.MutableEnvironment {
	var err error

	env := environment.New()
	err = env.Load()

	if err == nil {
		err = env.Save()
	}

	if err != nil {
		fmt.Printf("reading/creating the config failed: %v\n", err)
		os.Exit(1)
	}

	return env
}

func Execute(name string, arguments []string) error {
	var err error

	env := bootstrap()

	defer func() {
		if err := env.Save(); err != nil {
			panic(err)
		}
	}()

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
		fmt.Printf("csync version %s\n", env.Version())
		return nil
	}

	err = cli.DispatchCommand(commands, flags.Args(), name)
	if _, ok := err.(cli.CommandDispatchError); ok && err != nil {
		fmt.Printf("ERROR: %v\n", err)
		usage.Print()
		os.Exit(1)
	}

	return err
}
