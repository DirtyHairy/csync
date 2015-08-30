package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if err := Execute(filepath.Base(os.Args[0]), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}
