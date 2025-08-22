package meta_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/hostwithquantum/deno-buildpack/internal/meta"
	"github.com/stretchr/testify/assert"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
}

func TestFinderFind(t *testing.T) {
	testCases := []struct {
		Desc  string
		Path  string
		Match bool
	}{
		{
			Desc:  "deno",
			Path:  "../../samples/deno",
			Match: true,
		},
		{
			Desc:  "not deno",
			Path:  "../../samples/not-deno",
			Match: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Desc, func(t *testing.T) {
			finder := meta.Factory()
			assert.NoError(t, finder.Find(tc.Path))
			assert.Equal(t, tc.Match, finder.HasMatch())

			if tc.Match {
				assert.GreaterOrEqual(t, len(finder.GetMatches()), 1)
			} else {
				assert.Len(t, finder.GetMatches(), 0)
			}
		})
	}
}
