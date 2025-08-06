package detect

import (
	"github.com/hostwithquantum/deno-buildpack/internal/meta"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func Detect(logs scribe.Emitter) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		logs.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		logs.Process("Running detection...")
		logs.Process("Checking working directory: %s", context.WorkingDir)

		finder := meta.Factory()
		if err := finder.Find(context.WorkingDir); err != nil {
			return packit.DetectResult{}, err
		}

		if !finder.HasMatch() {
			logs.Subprocess("Not a deno app")
			return packit.DetectResult{}, packit.Fail.WithMessage("no deno configuration files found")
		}

		logs.Detail("Detected deno")
		logs.Detail("Found matches: %#v", finder.GetMatches())

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "deno-buildpack",
						Metadata: map[string]any{
							"launch": true,
							"build":  false,
						},
					},
				},
			},
		}, nil
	}
}
