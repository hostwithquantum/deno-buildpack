package meta

import (
	"io"

	"github.com/BurntSushi/toml"
	"github.com/paketo-buildpacks/packit/v2"
)

func decode(from io.Reader, to any) error {
	if _, err := toml.NewDecoder(from).Decode(to); err != nil {
		return packit.Fail.WithMessage("failed to decode buildpack.toml: %s", err)
	}
	return nil
}
