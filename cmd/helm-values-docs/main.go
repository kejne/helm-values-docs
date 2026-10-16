package main

import (
	"os"

	"github.com/kejne/helm-values-docs/internal/adapters/in/cli"
	"github.com/kejne/helm-values-docs/internal/adapters/out/filesystem"
	"github.com/kejne/helm-values-docs/internal/adapters/out/jsonschema"
	"github.com/kejne/helm-values-docs/internal/adapters/out/markdown"
	"github.com/kejne/helm-values-docs/internal/adapters/out/yamlvalues"
	"github.com/kejne/helm-values-docs/internal/app"
)

func main() {
	service := app.New(jsonschema.Loader{}, yamlvalues.Adapter{}, markdown.Adapter{}, filesystem.TemplateAdapter{}, filesystem.Store{})
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr, service))
}
