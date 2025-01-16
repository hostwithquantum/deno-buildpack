package meta

const BPLayerName = "runway-deno"

const DENO_CONFIG_FILE_JSON = "deno.json"
const DENO_CONFIG_FILE_JSONC = "deno.jsonc"

const DENO_BP_VERSION_FILE = "BP_RUNWAY_DENO_FILE_VERSION"
const DENO_BP_DVMRC_FILE = ".dvmrc"

type AppEnv struct {
	AllowEnv    string `env:"BP_RUNWAY_DENO_PERM_ENV" envDefault:"PORT"`
	AllowHRTime bool   `env:"BP_RUNWAY_DENO_PERM_HRTIME"`
	AllowNet    string `env:"BP_RUNWAY_DENO_PERM_NET" envDefault:"true"`
	AllowFFI    bool   `env:"BP_RUNWAY_DENO_PERM_FFI" envDefault:"false"`
	AllowRead   string `env:"BP_RUNWAY_DENO_PERM_READ" envDefault:"true"`
	AllowRun    string `env:"BP_RUNWAY_DENO_PERM_RUN" envDefault:"false"`
	AllowWrite  string `env:"BP_RUNWAY_DENO_PERM_WRITE" envDefault:"false"`
	AllowAll    bool   `env:"BP_RUNWAY_DENO_PERM_ALL" envDefault:"false"`

	DenoVersion     string `env:"BP_RUNWAY_DENO_VERSION" envDefault:"v1.46.3"`
	DenoFileVersion string `env:"BP_RUNWAY_DENO_FILE_VERSION" envDefault:"runtime.txt"`
	DenoMain        string `env:"BP_RUNWAY_DENO_MAIN" envDefault:"server.ts"`
}
