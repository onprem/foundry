package furnace

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"

	api "github.com/prmsrswt/foundry/pkg/furnace/furnacepb"
)

func TestBuild(t *testing.T) {
	cases := []struct {
		name    string
		req     *api.BuildRequest
		want    *api.BuildResponse
		wantErr bool
	}{
		{
			name: "req with single package",
			req: &api.BuildRequest{
				Packages: []*api.Package{{Name: "yay", Version: "10.1.2-1"}},
			},
			want:    &api.BuildResponse{},
			wantErr: false,
		},
		{
			name: "req with empty packages array",
			req: &api.BuildRequest{
				Packages: []*api.Package{},
			},
			want:    &api.BuildResponse{},
			wantErr: false,
		},
		{
			name: "req with multiple packages",
			req: &api.BuildRequest{
				Packages: []*api.Package{
					{Name: "timeshift", Version: "20.11.1+3+g08d0e59-2"},
					{Name: "yay", Version: "10.1.2-1"},
					{Name: "octopi", Version: "0.10.0-2"},
					{Name: "polybar", Version: "3.5.2-1"},
					{Name: "godot", Version: "3.2.3-1"},
				},
			},
			want:    &api.BuildResponse{},
			wantErr: false,
		},
	}
	for _, v := range cases {
		tc := v
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := NewFurnace(1, log.NewNopLogger())
			res, err := f.Build(context.Background(), tc.req)
			if tc.wantErr {
				assert.NotNil(t, err, "error should not be nil")
				assert.Nil(t, res, "response should be nil")
			} else {
				assert.Equal(t, res, tc.want)
			}
		})
	}
}
