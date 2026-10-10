package main

import (
	"os"

	"github.com/kejne/helm-values-docs/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
