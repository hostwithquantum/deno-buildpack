package meta_test

import (
	"testing"

	"github.com/caarlos0/env"
	"github.com/hostwithquantum/deno-buildpack/internal/meta"
	"github.com/stretchr/testify/assert"
)

func TestTypes(t *testing.T) {
	var allTheVars meta.AppEnv
	err := env.Parse(&allTheVars)
	assert.NoError(t, err)
	assert.Len(t, allTheVars.DenoMain, 2)
	assert.Equal(t, []string{"main.ts", "server.ts"}, allTheVars.DenoMain)
}
