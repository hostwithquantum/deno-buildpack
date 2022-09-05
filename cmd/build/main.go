package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/env"
	"github.com/hostwithquantum/deno-buildpack/internal/build"
	"github.com/hostwithquantum/deno-buildpack/internal/meta"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func main() {
	logEmitter := scribe.NewEmitter(os.Stdout).WithLevel(os.Getenv("BP_LOG_LEVEL"))

	var allTheVars meta.AppEnv
	err := env.Parse(&allTheVars)
	if err != nil {
		fmt.Fprintln(os.Stdout, fmt.Errorf("failed getting environment: %s", err))
		os.Exit(1)
	}

	packit.Build(build.Build(logEmitter, allTheVars))
}
