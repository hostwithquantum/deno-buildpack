package build

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hostwithquantum/deno-buildpack/internal/meta"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/paketo-buildpacks/packit/v2/vacation"
)

type PackageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

func Build(logger scribe.Emitter) packit.BuildFunc {
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
		v := meta.VersionFactory(logger)
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

		logger.Process("Permission setup for deno app")

		runArgs, err := meta.Config(context, logger)
		if err != nil {
			return packit.BuildResult{}, err
		}

		// run bundle?
		logger.Process("Bundler")
		if _, err = os.Stat(filepath.Join(context.WorkingDir, "tsconfig.js")); err == nil {
			logger.Subprocess("Detected tsconfig.js")
			logger.Detail("Send feedback if we should run the bundler here: support@runway.horse")
		}

		logger.EnvironmentVariables(layer)

		launchMetadata.Processes = []packit.Process{
			{
				Type:    "web",
				Command: "deno",
				Args:    runArgs,
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
