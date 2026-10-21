package main

import (
	"os"

	"github.com/smedijan/helm-values-docs/internal/adapters/in/cli"
	"github.com/smedijan/helm-values-docs/internal/adapters/out/filesystem"
	"github.com/smedijan/helm-values-docs/internal/adapters/out/jsonschema"
	"github.com/smedijan/helm-values-docs/internal/adapters/out/markdown"
	"github.com/smedijan/helm-values-docs/internal/adapters/out/yamlvalues"
	"github.com/smedijan/helm-values-docs/internal/app"
)

func main() {
	service := app.New(jsonschema.Loader{}, yamlvalues.Adapter{}, markdown.Adapter{}, filesystem.TemplateAdapter{}, filesystem.Store{})
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr, service))
}
