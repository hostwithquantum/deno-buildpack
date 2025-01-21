package version_test

import (
	"os"
	"testing"

	"github.com/caarlos0/env"
	"github.com/hostwithquantum/deno-buildpack/internal/meta"
	"github.com/hostwithquantum/deno-buildpack/internal/version"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/stretchr/testify/assert"
)

func TestFindEnvDefault(t *testing.T) {
	var allTheVars meta.AppEnv
	err := env.Parse(&allTheVars)
	assert.NoError(t, err)

	logEmitter := scribe.NewEmitter(os.Stdout).WithLevel("DEBUG")

	v := version.VersionFactory(allTheVars, logEmitter)
	assert.NotNil(t, v)

	detectedVersion, err := v.Find(".")
	assert.NoError(t, err)

	assert.Equal(t, "v2.1.5", detectedVersion)
}
