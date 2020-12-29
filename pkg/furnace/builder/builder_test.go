package builder

import (
	"testing"

	"github.com/go-git/go-git/v5"
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
