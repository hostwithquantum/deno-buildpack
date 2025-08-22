package detect_test

import (
	"os"
	"testing"

	"github.com/hostwithquantum/deno-buildpack/internal/detect"
	"github.com/hostwithquantum/deno-buildpack/internal/meta"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/stretchr/testify/assert"
)

func TestDetect(t *testing.T) {
	logEmitter := scribe.NewEmitter(os.Stdout).WithLevel("debug")

	testCases := []struct {
		IsDeno     bool
		SamplePath string
	}{
		{
			SamplePath: "../../samples/not-deno",
			IsDeno:     false,
		},
		{
			SamplePath: "../../samples/deno",
			IsDeno:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.SamplePath, func(t *testing.T) {
			detectFunct := detect.Detect(logEmitter, meta.AppEnv{})
			res, err := detectFunct(packit.DetectContext{
				WorkingDir: tc.SamplePath,
				CNBPath:    "../../",
				Platform:   packit.Platform{},
			})

			if tc.IsDeno {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			} else {
				assert.Error(t, err)
				assert.Empty(t, res)
			}

			// t.Logf("Error: %#v", err)
			// t.Logf("Res: %#v", res)
		})
	}
}
