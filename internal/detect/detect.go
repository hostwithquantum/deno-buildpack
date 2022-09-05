package detect

import (
	"github.com/hostwithquantum/deno-buildpack/internal/meta"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func Detect(logs scribe.Emitter) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		logs.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		finder := meta.Factory()
		return packit.DetectResult{}, finder.Find(context.WorkingDir)
	}
}
