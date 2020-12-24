package furnace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultPackages = []Package{{Name: "default"}, {Name: "test"}}

func TestEnqueue(t *testing.T) {
	cases := []struct {
		name  string
		input []Package
	}{
		{
			name:  "enqueueing a few packages",
			input: []Package{{Name: "yay"}, {Name: "polybar"}, {Name: "timeshift"}, {Name: "godot"}},
		},
		{
			name:  "enqueueing zero packages",
			input: []Package(nil),
		},
	}

	for _, v := range cases {
		tc := v
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			q := &queue{}
			q.enqueue(tc.input)

			assert.Equal(t, tc.input, q.items)
		})
	}
}

func TestDequeue(t *testing.T) {
	q := &queue{
		items: defaultPackages,
	}

	for _, v := range defaultPackages {
		pkg := v
		t.Run("dequeueing concurrently", func(t *testing.T) {
			got, ok := q.dequeue()
			assert.True(t, ok)
			assert.Equal(t, got, pkg)
		})
	}
	_, ok := q.dequeue()
	assert.False(t, ok)
}

func TestDequeueConcurrent(t *testing.T) {
	q := &queue{
		items: defaultPackages,
	}

	t.Run("group", func(t *testing.T) {
		for i := 0; i < len(defaultPackages); i++ {
			t.Run("dequeueing concurrently", func(t *testing.T) {
				t.Parallel()
				_, ok := q.dequeue()
				assert.True(t, ok)
			})
		}
	})

	// This should be false as we already dequeued everything before this.
	_, ok := q.dequeue()
	assert.False(t, ok)
}

func TestGetItems(t *testing.T) {
	cases := []struct {
		name  string
		input []Package
		// want []Package
	}{
		{
			name:  "a few items in the queue",
			input: defaultPackages,
			// want: defaultPackages,
		},
		{
			name: "zero items in the queue",
		},
	}
	for _, v := range cases {
		tc := v
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			q := &queue{}
			q.enqueue(tc.input)
			assert.Equal(t, tc.input, q.getItems())
		})
	}
}

func TestIsQueued(t *testing.T) {
	q := &queue{items: defaultPackages}
	cases := []struct {
		name  string
		input Package
		want  bool
	}{
		{
			name:  "package exists in queue",
			input: defaultPackages[0],
			want:  true,
		},
		{
			name:  "package doesn't exist in queue",
			input: Package{Name: "this-does-n0t-exist"},
			want:  false,
		},
		{
			name:  "empty package",
			input: Package{},
			want:  false,
		},
	}
	for _, v := range cases {
		tc := v
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, q.isQueued(tc.input))
		})
	}
}
