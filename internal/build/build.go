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

		finder := meta.Factory()
		finder.Find(context.WorkingDir)

		if !finder.HasMatch() {
			logger.Process("not a deno app")
			return packit.BuildResult{}, nil
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

		logger.Process("getting deno version")

		v := version.VersionFactory(appEnv, logger)
		denoVersion, err := v.Find(context.WorkingDir)
		if err != nil {
			return packit.BuildResult{}, err
		}

		if denoVersion != "" {
			logger.Detail("discovered %s", denoVersion)
		}

		layerBinPath := filepath.Join(layer.Path, "bin")
		err = os.MkdirAll(layerBinPath, os.ModePerm)
		if err != nil {
			return packit.BuildResult{}, err
		}

		var downloadUrl string
		downloadDest := filepath.Join(layer.Path, "deno.zip")

		logger.Process("download")
		if denoVersion != "latest" {
			logger.Detail("building download for: %s", denoVersion)
			downloadUrl = fmt.Sprintf(
				"https://github.com/denoland/deno/releases/download/%s/deno-x86_64-unknown-linux-gnu.zip",
				denoVersion)
		} else {
			logger.Detail("building download for latest")
			downloadUrl = "https://github.com/denoland/deno/releases/latest/download/deno-x86_64-unknown-linux-gnu.zip"
		}

		logger.Process("installing deno %s", denoVersion)

		logger.Subprocess("downloading deno")
		logger.Subdetail("url: %s", downloadUrl)
		logger.Subdetail("to: %s", downloadDest)

		resp, err := http.Get(downloadUrl)
		if err != nil {
			logger.Detail("download failed")
			return packit.BuildResult{}, err
		}

		defer resp.Body.Close()

		logger.Subprocess("extracting download")
		logger.Subdetail("to: %s", layerBinPath)

		zip := vacation.NewZipArchive(resp.Body).StripComponents(0)
		err = zip.Decompress(layerBinPath)
		if err != nil {
			logger.Detail("failed")
			return packit.BuildResult{}, err
		}

		layer.BuildEnv.Append("PATH", layerBinPath, ":")
		layer.LaunchEnv.Append("PATH", layerBinPath, ":")

		var launchMetadata packit.LaunchMetadata

		denoRunArgs := []string{"run"}

		logger.Process("determine permissions for deno process")

		if appEnv.AllowAll {
			logger.Detail("granting all permissions — this is not very secure")
			denoRunArgs = append(denoRunArgs, "--allow-all")
		} else {
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
		if _, err = os.Stat(filepath.Join(context.WorkingDir, "tsconfig.js")); err == nil {
			fmt.Println("we could run bundle here")
		}

		logger.EnvironmentVariables(layer)

		denoRunArgs = append(denoRunArgs, appEnv.DenoMain)

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
