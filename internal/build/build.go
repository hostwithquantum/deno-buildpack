package build

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hostwithquantum/deno-buildpack/internal/meta"
	"github.com/hostwithquantum/deno-buildpack/internal/version"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/paketo-buildpacks/packit/v2/vacation"
)

type PackageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

func Build(logger scribe.Emitter, appEnv meta.AppEnv) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {

		logger.Title("%s %s", context.BuildpackInfo.ID, context.BuildpackInfo.Version)

		var metadata map[string]any
		for _, p := range context.Plan.Entries {
			if p.Name != "deno" {
				continue
			}
			metadata = p.Metadata
			break
		}

		if metadata == nil {
			return packit.BuildResult{}, packit.Fail.WithMessage("malformed plan")
		}

		logger.Process("Getting version source")
		path, ok := metadata["version_source"].(string)
		if !ok {
			return packit.BuildResult{}, packit.Fail.WithMessage(
				"broken detection process, expected `version_source`",
			)
		}
		logger.Subprocess("Found: %s", path)

		logger.Process("Fetching deno version")
		v := version.VersionFactory(appEnv, logger)
		denoVersion, err := v.GetVersionByFile(path)
		if err != nil {
			return packit.BuildResult{}, packit.Fail.WithMessage(
				"failed to get get deno version: %s", err,
			)
		}

		if denoVersion != "" {
			logger.Subprocess("Found %q in %s", denoVersion, path)
		}

		layer, err := context.Layers.Get(meta.BPLayerName)
		if err != nil {
			logger.Process("failed to fetch layer: %s", err)
			return packit.BuildResult{}, err
		}
		layer, err = layer.Reset()
		if err != nil {
			logger.Process("failed to reset layer: %s", err)
			return packit.BuildResult{}, err
		}

		layer.Build = false
		layer.Launch = true

		layerBinPath := filepath.Join(layer.Path, "bin")
		if err := os.MkdirAll(layerBinPath, os.ModePerm); err != nil {
			return packit.BuildResult{}, err
		}

		var downloadUrl string
		if denoVersion != "latest" {
			downloadUrl = fmt.Sprintf(
				"https://github.com/denoland/deno/releases/download/%s/deno-x86_64-unknown-linux-gnu.zip",
				denoVersion)
		} else {
			downloadUrl = "https://github.com/denoland/deno/releases/latest/download/deno-x86_64-unknown-linux-gnu.zip"
		}

		logger.Subprocess("Downloading deno %q from Github", denoVersion)

		// FIXME(till): we will eventually need token support here to avoid rate-limiting
		resp, err := http.Get(downloadUrl)
		if err != nil {
			logger.Detail("Download failed")
			return packit.BuildResult{}, err
		}

		defer resp.Body.Close()

		logger.Subprocess("Extracting download")
		logger.Detail("Destination: %q", layerBinPath)

		zip := vacation.NewZipArchive(resp.Body).StripComponents(0)
		if err := zip.Decompress(layerBinPath); err != nil {
			logger.Detail("failed")
			return packit.BuildResult{}, err
		}

		layer.BuildEnv.Append("PATH", layerBinPath, ":")
		layer.LaunchEnv.Append("PATH", layerBinPath, ":")

		var launchMetadata packit.LaunchMetadata

		denoRunArgs := []string{"run"}

		logger.Process("Permission setup for deno app")

		if appEnv.AllowAll {
			logger.Subprocess("Granting all permissions — this is not very secure")
			denoRunArgs = append(denoRunArgs, "--allow-all")
		} else {
			logger.Subprocess("Setting granular permissions")
			if appEnv.AllowEnv != "false" {
				assembleArgs(&denoRunArgs, "--allow-env", appEnv.AllowEnv)
				logger.Detail("Set --allow-env")
			}

			if appEnv.AllowHRTime {
				denoRunArgs = append(denoRunArgs, "--allow-hrtime")
				logger.Detail("Set --allow-hrtime")
			}

			if appEnv.AllowNet != "false" {
				assembleArgs(&denoRunArgs, "--allow-net", appEnv.AllowNet)
				logger.Detail("Set --allow-net")
			}

			if appEnv.AllowFFI {
				denoRunArgs = append(denoRunArgs, "--allow-ffi")
				logger.Detail("Set --allow-ffi")
			}

			if appEnv.AllowRead != "false" {
				assembleArgs(&denoRunArgs, "--allow-read", appEnv.AllowRead)
				logger.Detail("Set --allow-read")
			}

			if appEnv.AllowRun != "false" {
				assembleArgs(&denoRunArgs, "--allow-run", appEnv.AllowRun)
				logger.Detail("Set --allow-run")
			}

			if appEnv.AllowWrite != "false" {
				assembleArgs(&denoRunArgs, "--allow-write", appEnv.AllowWrite)
				logger.Detail("Set --allow-write")
			}
		}

		// run bundle?
		logger.Process("Bundler")
		if _, err = os.Stat(filepath.Join(context.WorkingDir, "tsconfig.js")); err == nil {
			logger.Subprocess("Detected tsconfig.js")
			logger.Detail("Send feedback if we should run the bundler here: support@runway.horse")
		}

		logger.EnvironmentVariables(layer)

		// fall back to main.ts according to `deno init`
		logger.Process("Finding entrypoint")
		if len(appEnv.DenoMain) == 0 {
			logger.Subprocess("Using default (main.ts)")
			denoRunArgs = append(denoRunArgs, "main.ts")
		} else {
			var foundMain bool = false
			for _, f := range appEnv.DenoMain {
				if _, err := os.Stat(filepath.Join(context.WorkingDir, f)); err == nil {
					foundMain = true
					denoRunArgs = append(denoRunArgs, f)
					logger.Subprocess("Using %q", f)
					break
				}
			}

			if !foundMain {
				logger.Subprocess("Unable to determine main file/entrypoint for app, please see BP_RUNWAY_DENO_MAIN")
				return packit.BuildResult{}, fmt.Errorf("unable to find entrypoint")
			}
		}

		launchMetadata.Processes = []packit.Process{
			{
				Type:    "web",
				Command: "deno",
				Args:    denoRunArgs,
				Default: true,
				Direct:  false,
			},
		}

		logger.LaunchProcesses(launchMetadata.Processes)

		return packit.BuildResult{
			Layers: []packit.Layer{layer},
			Launch: launchMetadata,
		}, nil
	}
}

func assembleArgs(args *[]string, opt string, optValue string) {
	if len(optValue) == 0 || optValue == "true" {
		*args = append(*args, opt)
	} else {
		*args = append(*args, fmt.Sprintf("%s=%s", opt, optValue))
	}
}
