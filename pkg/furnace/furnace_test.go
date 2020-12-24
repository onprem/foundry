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
				Packages: []*api.Package{{Name: "yay"}},
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
					{Name: "timeshift"},
					{Name: "yay"},
					{Name: "octopi"},
					{Name: "polybar"},
					{Name: "godot"},
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
