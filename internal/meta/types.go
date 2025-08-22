package meta

const BPLayerName = "runway-deno"

const DENO_CONFIG_FILE_JSON = "deno.json"
const DENO_CONFIG_FILE_JSONC = "deno.jsonc"

const DENO_BP_VERSION_FILE = "BP_RUNWAY_DENO_FILE_VERSION"
const DENO_BP_DVMRC_FILE = ".dvmrc"

type BuildpackConfig struct {
	Metadata struct {
		Configurations []struct {
			Default     string `toml:"default,omitempty"`
			Description string `toml:"description"`
			Name        string `toml:"name"`
		} `toml:"configurations"`
	} `toml:"metadata"`
}
