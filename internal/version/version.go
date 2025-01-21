package version

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hostwithquantum/deno-buildpack/internal/meta"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

// this needs to be updated when we increase the deno version
const (
	versionDeno    = "v2.1.5"
	versionDefault = "__default__"
)

type Version struct {
	env    meta.AppEnv
	logger scribe.Emitter
}

func VersionFactory(env meta.AppEnv, logger scribe.Emitter) *Version {
	return &Version{
		env:    env,
		logger: logger,
	}
}

// Find
// This determines the deno version to download. This does not validate the version
// in terms of available release.
func (v *Version) Find(workDir string) (string, error) {
	denoVersion := v.env.DenoVersion
	if denoVersion != "" && denoVersion != versionDefault {
		return fixVersionString(denoVersion), nil
	}

	runtimeFile := filepath.Join(workDir, v.env.DenoFileVersion)
	v.logger.Detail("trying %s", runtimeFile)

	if fileExists(runtimeFile) {
		return v.extractVersion(runtimeFile)
	}

	dvmrcFile := filepath.Join(workDir, meta.DENO_BP_DVMRC_FILE)
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
