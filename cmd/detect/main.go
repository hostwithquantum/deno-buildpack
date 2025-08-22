package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/env"
	"github.com/hostwithquantum/deno-buildpack/internal/detect"
	"github.com/hostwithquantum/deno-buildpack/internal/meta"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func main() {
	logEmitter := scribe.NewEmitter(os.Stdout).WithLevel(os.Getenv("BP_LOG_LEVEL"))

	var appEnv meta.AppEnv
	if err := env.Parse(&appEnv); err != nil {
		fmt.Fprintln(os.Stdout, fmt.Errorf("failed getting environment: %s", err))
		os.Exit(1)
	}

	packit.Detect(detect.Detect(logEmitter, appEnv))
}
