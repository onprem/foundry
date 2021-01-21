package furnace

import (
	"testing"

	"github.com/mikkeloscar/aur"
	"github.com/stretchr/testify/assert"
)

func TestPkgMapGetSet(t *testing.T) {
	cases := []struct {
		name  string
		pkg   aur.Pkg
		input status
		want  status
	}{
		{
			name: "default status",
			pkg:  aur.Pkg{Name: "yay-bin"},
			want: StatusDefault,
		},
		{
			name:  "building status",
			pkg:   aur.Pkg{Name: "yay-bin"},
			input: StatusBuilding,
			want:  StatusBuilding,
		},
		{
			name:  "error status",
			pkg:   aur.Pkg{Name: "yay-bin"},
			input: StatusError,
			want:  StatusError,
		},
		{
			name: "empty package and default status",
			pkg:  aur.Pkg{},
			want: StatusDefault,
		},
		{
			name:  "empty package with non-default status",
			pkg:   aur.Pkg{},
			input: StatusBuilding,
			want:  StatusBuilding,
		},
	}
	for _, v := range cases {
		tc := v
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			pm := pkgMap{
				pkgs: make(map[string]status),
			}
			if tc.input != "" {
				pm.set(tc.pkg, tc.input)
			}
			assert.Equal(t, tc.want, pm.get(tc.pkg))
		})
	}
}

func TestMapSetQueued(t *testing.T) {
	n := 600
	pm := pkgMap{
		pkgs: make(map[string]status),
	}
	pkg := aur.Pkg{
		Name: "package",
	}
	errchan := make(chan error, 600)
	for i := 0; i < n; i++ {
		go func() {
			err := pm.setQueued(pkg)
			errchan <- err
		}()
	}
	var errCount int
	for i := 0; i < n; i++ {
		err := <-errchan
		if err != nil {
			errCount++
		}
	}
	assert.Equal(t, n-1, errCount, "we should get ErrAlreadyBuilding for all except first")
}
