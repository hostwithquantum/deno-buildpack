package meta_test

import (
	"os"
	"testing"

	"github.com/hostwithquantum/deno-buildpack/internal/meta"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/stretchr/testify/assert"
)

func TestFindEnvDefault(t *testing.T) {
	logEmitter := scribe.NewEmitter(os.Stdout).WithLevel("DEBUG")

	v := meta.VersionFactory(logEmitter)
	assert.NotNil(t, v)

	detectedVersion, err := v.Find(packit.DetectContext{
		CNBPath:    "../../",
		WorkingDir: ".",
	})
	assert.NoError(t, err)

	assert.Equal(t, "v2.1.5", detectedVersion)
}
