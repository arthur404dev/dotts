package main

import (
	"os"

	"github.com/arthur404dev/dotts/cmd/dotts/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
