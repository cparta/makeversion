package makeversion

import (
	"os"
	"testing"

	"github.com/matryer/is"
)

func Test_NewDefaultGitter_SucceedsNormally(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	is.True(dg != nil)
}

func Test_CheckGitRepo_SuceedsForCurrent(t *testing.T) {
	is := is.New(t)
	repo, err := CheckGitRepo(".")
	is.NoErr(err)
	is.True(repo != "")
}

func Test_CheckGitRepo_SuceedsForCmdMkver(t *testing.T) {
	is := is.New(t)
	repo, err := CheckGitRepo("./cmd/mkver")
	is.NoErr(err)
	is.True(repo != "")
}

func Test_CheckGitRepo_FailsForRoot(t *testing.T) {
	is := is.New(t)
	repo, err := CheckGitRepo("/")
	is.True(repo == "")
	is.True(err != nil)
}

func Test_CheckGitRepo_IgnoresFileNamedGit(t *testing.T) {
	is := is.New(t)
	const fileNamedGit = "./cmd/mkver/.git"
	if _, err := os.Stat(fileNamedGit); err != nil {
		if f, err := os.Create(fileNamedGit); err == nil {
			defer f.Close()
			defer os.Remove(fileNamedGit)
			repo, err := CheckGitRepo("./cmd/mkver")
			is.NoErr(err)
			is.True(repo != "")
		}
	} else {
		t.Logf("warning: '%s' already exists\n", fileNamedGit)
	}
}

func Test_DefaultGitter_GetTag(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	is.True(dg != nil)
	is.True(dg.GetTag(".") != "")
	is.Equal(dg.GetTag("/"), "")
}

func Test_DefaultGitter_GetBranch(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	is.True(dg != nil)
	is.True(dg.GetBranch(".") != "")
	is.Equal(dg.GetBranch("/"), "")
}

func Test_DefaultGitter_GetBranchFromTag(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	is.True(dg != nil)
	is.True(dg.GetBranchFromTag(".", "refs/tags/v1.0.0") == "main")
	is.Equal(dg.GetBranchFromTag("/", "refs/tags/v1.0.0"), "")
}

func Test_DefaultGitter_GetBuild(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	is.True(dg != nil)
	is.True(dg.GetBuild(".") != "")
	is.Equal(dg.GetBuild("/"), "")
}
