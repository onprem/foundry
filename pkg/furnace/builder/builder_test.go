package builder

import (
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/mikkeloscar/aur"
	"github.com/prometheus/client_golang/prometheus"
	promtest "github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

// Very basic test that just checks for no error.
func TestGitClone(t *testing.T) {
	dir := t.TempDir()
	url := "https://github.com/prmsrswt/foundry"

	err := gitClone(dir, url)
	assert.NoError(t, err)
	// verify that it has been cloned
	_, err = git.PlainOpen(dir)
	assert.NoError(t, err)
	// clone again at same path
	err = gitClone(dir, url)
	assert.NoError(t, err)
	// verify again that there is a repo there
	_, err = git.PlainOpen(dir)
	assert.NoError(t, err)
}

func TestInstrumentedBuilder(t *testing.T) {
	cases := []struct {
		shouldFail []bool
	}{
		// Single successful build.
		{shouldFail: []bool{false}},
		// Single failing build.
		{shouldFail: []bool{true}},
		// No builds.
		{shouldFail: []bool{}},
		// A mix of passing and failing builds.
		{shouldFail: []bool{false, true, false, false, true, false}},
	}
	for _, v := range cases {
		tc := v
		t.Run("builds", func(t *testing.T) {
			t.Parallel()
			reg := prometheus.NewRegistry()
			noop := &NoOpBuilder{}
			b := BuilderWithMetrics(noop, reg)
			var numFail int
			for _, v := range tc.shouldFail {
				if v {
					numFail++
				}
				noop.ShouldFail = v
				err := b.Build(aur.Pkg{Name: "example"})
				if v {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}
			assert.Equal(t, float64(len(tc.shouldFail)), promtest.ToFloat64(b.buildOps))
			assert.Equal(t, float64(len(tc.shouldFail)), promtest.ToFloat64(b.buildOpsCompleted))
			assert.Equal(t, float64(numFail), promtest.ToFloat64(b.buildOpsFailure))
		})
	}
}
