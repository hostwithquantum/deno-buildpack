package version

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hostwithquantum/deno-buildpack/internal/meta"
	"github.com/paketo-buildpacks/packit/v2/scribe"
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

func (v *Version) Find(workDir string) (string, error) {
	denoVersion := v.env.DenoVersion
	if denoVersion != "" {
		return fixVersionString(denoVersion), nil
	}

	if v.env.DenoFileVersion != "" {
		runtimeFile := filepath.Join(workDir, v.env.DenoFileVersion)
		v.logger.Detail("trying %s", runtimeFile)

		if _, err := os.Stat(runtimeFile); err != nil {
			cnt, err := ioutil.ReadFile(filepath.Join(workDir, v.env.DenoFileVersion))
			if err != nil {
				return "", err
			}
			if len(cnt) > 0 {
				denoVersion = string(cnt)
				v.logger.Subdetail("discovered: %s", denoVersion)

				return fixVersionString(denoVersion), nil
			}
		}
	}

	v.logger.Detail("trying %s", meta.DENO_BP_DVMRC_FILE)
	if _, err := os.Stat(filepath.Join(workDir, meta.DENO_BP_DVMRC_FILE)); err != nil {
		cnt, err := ioutil.ReadFile(filepath.Join(workDir, meta.DENO_BP_DVMRC_FILE))
		if err != nil {
			return "", err
		}

		if len(cnt) > 0 {
			denoVersion = string(cnt)
			v.logger.Subdetail("discovered %s", denoVersion)
			return fixVersionString(denoVersion), nil
		}
	}

	// use latest
	return "latest", nil
}

func fixVersionString(ver string) string {
	if string(ver[0]) != "v" {
		ver = fmt.Sprintf("v%s", ver)
	}
	return ver
}
