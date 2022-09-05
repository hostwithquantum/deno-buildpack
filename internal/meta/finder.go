package meta

import (
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2/fs"
)

type Finder struct {
	matched bool
	matches map[string]string

	Files []string
}

func Factory() *Finder {
	finder := &Finder{
		matched: false,
		matches: make(map[string]string),
	}

	finder.Files = []string{DENO_CONFIG_FILE_JSON, DENO_CONFIG_FILE_JSONC, DENO_BP_DVMRC_FILE}

	versionFile, ok := os.LookupEnv(DENO_BP_VERSION_FILE)
	if !ok {
		versionFile = "runtime.txt"
	}

	finder.Files = append(finder.Files, versionFile)

	return finder
}

func (f *Finder) Find(workingDir string) error {
	for _, metaFile := range f.Files {

		l := filepath.Join(workingDir, metaFile)

		exist, err := fs.Exists(l)
		if err != nil {
			return err
		}

		if exist {
			f.matched = true
			f.matches[metaFile] = l
		}
	}

	return nil
}

func (f *Finder) GetMatches() map[string]string {
	return f.matches
}

func (f *Finder) HasMatch() bool {
	return f.matched
}
