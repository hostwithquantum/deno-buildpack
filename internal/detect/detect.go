package detect

import (
	"github.com/hostwithquantum/deno-buildpack/internal/meta"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func Detect(logs scribe.Emitter) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		plan := packit.DetectResult{
			Plan: packit.BuildPlan{
				Requires: []packit.BuildPlanRequirement{},
			},
		}

		logs.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		logs.Process("Running detection...")
		logs.Process("Checking working directory: %s", context.WorkingDir)

		finder := meta.Factory()
		if err := finder.Find(context.WorkingDir); err != nil {
			return packit.DetectResult{}, packit.Fail.WithMessage("%s", err.Error())
		}

		if !finder.HasMatch() {
			logs.Subprocess("Not a deno app")
			return packit.DetectResult{}, packit.Fail.WithMessage("no deno configuration files found")
		}

		logs.Detail("Detected deno")
		logs.Detail("Found matches: %#v", finder.GetMatches())

		v := meta.VersionFactory(logs)
		var requirements = []packit.BuildPlanRequirement{}

		for _, path := range finder.GetMatches() {
			denoVersion, err := v.GetVersionByFile(path)
			if err != nil {
				continue
			}

			requirements = append(requirements, packit.BuildPlanRequirement{
				Name: "deno",
				Metadata: map[string]any{
					"version":        denoVersion,
					"version_source": path,
				},
			})

			break
		}

		plan.Plan.Requires = requirements
		return plan, nil
	}
}
