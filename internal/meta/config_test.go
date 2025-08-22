package meta_test

import (
	"os"
	"testing"

	"github.com/hostwithquantum/deno-buildpack/internal/meta"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	var (
		ctxDeno = packit.BuildContext{
			CNBPath:    "../../",
			WorkingDir: "../../samples/deno",
		}
		logger = scribe.NewEmitter(os.Stdout)
	)

	t.Run("samples/deno", func(t *testing.T) {
		runArgs, err := meta.Config(ctxDeno, logger)
		assert.NoError(t, err)
		assert.Equal(t, []string{
			"run",
			"--allow-env=PORT",
			"--allow-net",
			"--allow-read",
			"server.ts",
		}, runArgs)
	})

	t.Run("samples/deno: env", func(t *testing.T) {
		t.Setenv("BP_RUNWAY_DENO_PERM_ENV", "PORT,DATABASE_URL")
		runArgs, err := meta.Config(ctxDeno, logger)
		assert.NoError(t, err)
		assert.Equal(t, []string{
			"run",
			"--allow-env=PORT,DATABASE_URL",
			"--allow-net",
			"--allow-read",
			"server.ts",
		}, runArgs)
	})

	t.Run("samples/deno: all env", func(t *testing.T) {
		t.Setenv("BP_RUNWAY_DENO_PERM_ENV", "true")
		runArgs, err := meta.Config(ctxDeno, logger)
		assert.NoError(t, err)
		assert.Equal(t, []string{
			"run",
			"--allow-env",
			"--allow-net",
			"--allow-read",
			"server.ts",
		}, runArgs)
	})

	t.Run("samples/deno: ffi", func(t *testing.T) {
		t.Setenv("BP_RUNWAY_DENO_PERM_ENV", "PORT,DATABASE_URL,API_KEY")
		t.Setenv("BP_RUNWAY_DENO_PERM_FFI", "true")
		runArgs, err := meta.Config(ctxDeno, logger)
		assert.NoError(t, err)
		assert.Equal(t, []string{
			"run",
			"--allow-env=PORT,DATABASE_URL,API_KEY",
			"--allow-net",
			"--allow-ffi",
			"--allow-read",
			"server.ts",
		}, runArgs)
	})

	t.Run("samples/deno: --allow-all", func(t *testing.T) {
		t.Setenv("BP_RUNWAY_DENO_PERM_ALL", "true")
		runArgs, err := meta.Config(ctxDeno, logger)
		assert.NoError(t, err)
		assert.Equal(t, []string{
			"run",
			"--allow-all",
			"server.ts",
		}, runArgs)
	})

	t.Run("samples/no-deno", func(t *testing.T) {
		_, err := meta.Config(packit.BuildContext{
			CNBPath:    "../../",
			WorkingDir: "../../samples/no-deno",
		}, logger)
		assert.Error(t, err)
	})

	t.Run("samples/deno2-serve", func(t *testing.T) {
		t.Setenv("BP_RUNWAY_DENO_PERM_ALL", "true")
		t.Setenv("BP_RUNWAY_DENO_SERVE", "true")
		runArgs, err := meta.Config(packit.BuildContext{
			CNBPath:    "../../",
			WorkingDir: "../../samples/deno2-serve",
		}, logger)
		assert.NoError(t, err)
		assert.Equal(t, []string{"serve", "--allow-all", "main.ts"}, runArgs)
	})
}
