package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type Usage struct {
	Usage       string
	Description string
	Commands    []Command

	Out   io.Writer
	Flags *flag.FlagSet
}

func (u *Usage) Print() {
	out := u.Out
	if out == nil {
		out = os.Stderr
	}

	u.Flags.SetOutput(out)

	fmt.Fprintf(out, "\n%s\n\n", u.Usage)

	if u.Description != "" {
		fmt.Fprintf(out, "%s\n\n", u.Description)
	}

	if u.Commands != nil && len(u.Commands) > 0 {
		fmt.Fprintf(out, "Available commands:\n\n")
		for _, command := range u.Commands {
			fmt.Fprintf(out, "  %s\n    \t%s\n", command.Name, command.Description)
		}
		fmt.Fprintln(out)
	}

	fmt.Fprintf(out, "Available options:\n\n")
	u.Flags.PrintDefaults()

	fmt.Fprintf(out, "  -help, --help, -h\n    \tshow this help message\n\n")
}
