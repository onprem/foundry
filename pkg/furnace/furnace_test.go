package furnace

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"

	api "github.com/prmsrswt/foundry/pkg/furnace/furnacepb"
)

func TestFurnaceBuild(t *testing.T) {
	cases := []struct {
		name    string
		req     *api.BuildRequest
		wantErr bool
	}{
		{
			name: "req with single package",
			req: &api.BuildRequest{
				Packages: []*api.Package{{Name: "yay"}},
			},
			wantErr: false,
		},
		{
			name: "req with empty packages array",
			req: &api.BuildRequest{
				Packages: []*api.Package{},
			},
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
				assert.Nil(t, err)
			}
		})
	}
}

func TestFurnaceIsQueued(t *testing.T) {
	f := NewFurnace(1, log.NewNopLogger())
	_, err := f.Build(context.Background(), &api.BuildRequest{
		Packages: []*api.Package{
			{Name: "timeshift"},
			{Name: "yay"},
			{Name: "octopi"},
			{Name: "polybar"},
			{Name: "godot"},
		},
	})
	assert.NoError(t, err)

	cases := []struct {
		name    string
		req     *api.IsQueuedRequest
		res     *api.IsQueuedResponse
		wantErr bool
	}{
		{
			name: "req with in-queue package",
			req: &api.IsQueuedRequest{
				Package: &api.Package{Name: "yay"},
			},
			res: &api.IsQueuedResponse{
				Status: true,
			},
			wantErr: false,
		},
		{
			name: "req with a package that's not queued",
			req: &api.IsQueuedRequest{
				Package: &api.Package{Name: "this-does-n0t-exist"},
			},
			res: &api.IsQueuedResponse{
				Status: false,
			},
			wantErr: false,
		},
		{
			name: "req with nil request",
			res: &api.IsQueuedResponse{
				Status: false,
			},
			wantErr: false,
		},
	}
	for _, v := range cases {
		tc := v
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res, err := f.IsQueued(context.Background(), tc.req)
			if tc.wantErr {
				assert.Error(t, err, "error should not be nil")
				assert.Nil(t, res, "response should be nil")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.res, res)
			}
		})
	}
}
