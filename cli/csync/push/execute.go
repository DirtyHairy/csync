package push

import (
	"flag"
	"fmt"
	"os"

	"github.com/DirtyHairy/csync/cli"
	"github.com/DirtyHairy/csync/lib/cmd/push"
)

func Execute(name string, arguments []string) error {
	flags := flag.NewFlagSet(name, flag.ExitOnError)

	usage := cli.Usage{
		Usage:       fmt.Sprintf("usage: %s [options] source_repo target_repo", name),
		Description: "Unidirectional sync between source_repo and target_repo.",
		Flags:       flags,
	}

	flags.Usage = usage.Print

	flags.Parse(arguments)

	if flags.NArg() != 2 {
		usage.Print()
		os.Exit(1)
	}

	return push.Execute(push.Config{
		SourceRepoId: flags.Arg(0),
		TargetRepoId: flags.Arg(1),
	})
}
