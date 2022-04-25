package makeversion

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewDefaultGitter_SucceedsNormally(t *testing.T) {
	dg, err := NewDefaultGitter("git")
	assert.NoError(t, err)
	assert.NotNil(t, dg)
}

func Test_CheckGitRepo_SuceedsForCurrent(t *testing.T) {
	repo, err := CheckGitRepo(".")
	assert.NoError(t, err)
	assert.NotEmpty(t, repo)
}

func Test_CheckGitRepo_SuceedsForCmdMkver(t *testing.T) {
	repo, err := CheckGitRepo("./cmd/mkver")
	assert.NoError(t, err)
	assert.NotEmpty(t, repo)
}

func Test_CheckGitRepo_FailsForRoot(t *testing.T) {
	repo, err := CheckGitRepo("/")
	assert.Error(t, err)
	assert.Empty(t, repo)
}

func Test_CheckGitRepo_IgnoresFileNamedGit(t *testing.T) {
	const fileNamedGit = "./cmd/mkver/.git"
	if _, err := os.Stat(fileNamedGit); err != nil {
		if f, err := os.Create(fileNamedGit); err == nil {
			defer f.Close()
			defer os.Remove(fileNamedGit)
			repo, err := CheckGitRepo("./cmd/mkver")
			assert.NoError(t, err)
			assert.NotEmpty(t, repo)
		}
	} else {
		t.Logf("warning: '%s' already exists\n", fileNamedGit)
	}
}

func Test_DefaultGitter_GetTag(t *testing.T) {
	dg, err := NewDefaultGitter("git")
	if assert.NoError(t, err) && assert.NotNil(t, dg) {
		assert.NotEmpty(t, dg.GetTag("."))
		assert.Empty(t, dg.GetTag("/"))
	}
}

func Test_DefaultGitter_GetBranch(t *testing.T) {
	dg, err := NewDefaultGitter("git")
	if assert.NoError(t, err) && assert.NotNil(t, dg) {
		assert.NotEmpty(t, dg.GetBranch("."))
		assert.Equal(t, "", dg.GetBranch("/"))
	}
}

func Test_DefaultGitter_GetBuild(t *testing.T) {
	dg, err := NewDefaultGitter("git")
	if assert.NoError(t, err) && assert.NotNil(t, dg) {
		assert.NotEmpty(t, dg.GetBuild("."))
		assert.Empty(t, dg.GetBuild("/"))
	}
}
