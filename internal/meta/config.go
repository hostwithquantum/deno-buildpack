package meta

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

var (
	// this "maps" environment variable to deno CLI option
	envOption = map[string]string{
		"BP_RUNWAY_DENO_PERM_ENV":    "--allow-env",
		"BP_RUNWAY_DENO_PERM_HRTIME": "--allow-hrtime",
		"BP_RUNWAY_DENO_PERM_NET":    "--allow-net",
		"BP_RUNWAY_DENO_PERM_FFI":    "--allow-ffi",
		"BP_RUNWAY_DENO_PERM_READ":   "--allow-read",
		"BP_RUNWAY_DENO_PERM_RUN":    "--allow-run",
		"BP_RUNWAY_DENO_PERM_WRITE":  "--allow-write",
		"BP_RUNWAY_DENO_PERM_ALL":    "--allow-all",
	}
)

func Config(ctx packit.BuildContext, logger scribe.Emitter) ([]string, error) {
	bp, err := os.Open(filepath.Join(ctx.CNBPath, "buildpack.toml"))
	if err != nil {
		return []string{}, packit.Fail.WithMessage("failed to find buildpack.toml: %s", err)
	}

	var configuration BuildpackConfig
	if err := decode(bp, &configuration); err != nil {
		return []string{}, err
	}

	var (
		runArgs  = []string{"run"}
		mainFile string
	)

configServe:
	for _, config := range configuration.Metadata.Configurations {
		if config.Name != "BP_RUNWAY_DENO_SERVE" {
			continue
		}
		v := getEnvWithDefault(config.Name, config.Default)
		if v == "true" {
			runArgs = []string{"serve"}
			break configServe
		}
	}

configMain:
	for _, config := range configuration.Metadata.Configurations {
		switch config.Name {
		case "BP_RUNWAY_DENO_MAIN":
			logger.Process("Finding entrypoint")
			mainFiles := getEnvAsStringSlice(config.Name, config.Default)
			for _, path := range mainFiles {
				path = strings.TrimSpace(path) // trim whitespace from comma-separated values
				if _, err := os.Stat(filepath.Join(ctx.WorkingDir, path)); err == nil {
					mainFile = path
					logger.Subprocess("Using %q", path)
					break configMain
				}
			}

			logger.Subprocess("Unable to determine main file/entrypoint for app, please see BP_RUNWAY_DENO_MAIN")
			return runArgs, fmt.Errorf("unable to find entrypoint")
		}
	}

configPermAll:
	for _, config := range configuration.Metadata.Configurations {
		switch config.Name {
		case "BP_RUNWAY_DENO_PERM_ALL": // this overrides all
			v := getEnvWithDefault(config.Name, config.Default)
			if v == "true" {
				logger.Subprocess("Granting all permissions — this is not very secure")
				assembleArgs(&runArgs, envOption[config.Name], "")
				runArgs = append(runArgs, mainFile)
				return runArgs, nil
			}
			break configPermAll
		}
	}

	for _, config := range configuration.Metadata.Configurations {
		switch config.Name {
		case "BP_RUNWAY_DENO_PERM_HRTIME", "BP_RUNWAY_DENO_PERM_FFI":
			v := getEnvWithDefault(config.Name, config.Default)
			if v == "true" {
				assembleArgs(&runArgs, envOption[config.Name], "")
				logger.Detail("Set %s", envOption[config.Name])
			}
		case "BP_RUNWAY_DENO_PERM_ENV", "BP_RUNWAY_DENO_PERM_NET", "BP_RUNWAY_DENO_PERM_READ", "BP_RUNWAY_DENO_PERM_RUN", "BP_RUNWAY_DENO_PERM_WRITE":
			v := getEnvWithDefault(config.Name, config.Default)
			if v != "false" {
				assembleArgs(&runArgs, envOption[config.Name], v)
				logger.Detail("Set %s=%s", envOption[config.Name], v)
			}
		}
	}

	runArgs = append(runArgs, mainFile)

	return runArgs, nil
}

// getEnvWithDefault returns the environment variable value or the default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsStringSlice returns the environment variable as a comma-separated slice
func getEnvAsStringSlice(key, defaultValue string) []string {
	value := getEnvWithDefault(key, defaultValue)
	if value == "" {
		return []string{}
	}
	return strings.Split(value, ",")
}

func assembleArgs(args *[]string, opt string, optValue string) {
	if len(optValue) == 0 || optValue == "true" {
		*args = append(*args, opt)
	} else {
		*args = append(*args, fmt.Sprintf("%s=%s", opt, optValue))
	}
}
