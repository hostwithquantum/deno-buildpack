package meta

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

// this needs to be updated when we increase the deno version
const (
	versionDeno    = "v2.1.5"
	versionDefault = "__default__"
)

type Version struct {
	logger scribe.Emitter
}

func VersionFactory(logger scribe.Emitter) *Version {
	return &Version{
		logger: logger,
	}
}

func (v *Version) GetVersionByFile(path string) (string, error) {
	return v.extractVersion(path)
}

// Find
// This determines the deno version to download. This does not validate the version
// in terms of available release.
func (v *Version) Find(ctx packit.BuildContext) (string, error) {
	bp, err := os.Open(filepath.Join(ctx.CNBPath, "buildpack.toml"))
	if err != nil {
		return "", packit.Fail.WithMessage("failed to find buildpack.toml: %s", err)
	}

	var configuration BuildpackConfig
	if err := decode(bp, &configuration); err != nil {
		return "", err
	}

	for _, config := range configuration.Metadata.Configurations {
		switch config.Name {
		case "BP_RUNWAY_DENO_VERSION":
			denoVersion := getEnvWithDefault(config.Name, config.Default)
			if denoVersion != "" && denoVersion != versionDefault && denoVersion != "latest" {
				return fixVersionString(denoVersion), nil
			}
		case "BP_RUNWAY_DENO_FILE_VERSION":
			versionFile := getEnvWithDefault(config.Name, config.Default)
			runtimeFile := filepath.Join(ctx.WorkingDir, versionFile)
			v.logger.Detail("trying %s", runtimeFile)
			if fileExists(runtimeFile) {
				return v.extractVersion(runtimeFile)
			}
		}
	}

	dvmrcFile := filepath.Join(ctx.WorkingDir, DENO_BP_DVMRC_FILE)
	v.logger.Detail("trying %s", dvmrcFile)
	if fileExists(dvmrcFile) {
		return v.extractVersion(dvmrcFile)
	}

	// use our default version
	v.logger.Detail("using default: %s", versionDeno)
	return versionDeno, nil
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func (v *Version) extractVersion(file string) (string, error) {
	cnt, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	if len(cnt) > 0 {
		denoVersion := string(cnt)
		v.logger.Subdetail("discovered %s", denoVersion)
		return fixVersionString(denoVersion), nil
	}

	return versionDeno, nil
}

func fixVersionString(ver string) string {
	if string(ver[0]) != "v" {
		ver = fmt.Sprintf("v%s", ver)
	}
	return ver
}
