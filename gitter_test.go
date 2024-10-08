package makeversion

import (
	"os"
	"os/exec"
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

func Test_lastName(t *testing.T) {
	is := is.New(t)
	is.Equal("foo", lastName("foo"))
	is.Equal("bar", lastName("foo/bar"))
}

func Test_DefaultGitter_GetBranchFromTag(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	is.True(dg != nil)
	is.Equal(dg.GetBranchesFromTag("/", "refs/tags/v1.0.0"), nil)
	is.Equal(dg.GetBranchesFromTag(".", "refs/tags/v1.0.0"), []string{"main"})
}

func Test_DefaultGitter_GetTagForHEAD(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	is.True(dg != nil)
	err = exec.Command("git", "tag", "test-tag").Run()
	is.NoErr(err)
	if err == nil {
		defer exec.Command("git", "tag", "-d", "test-tag").Run()
		is.Equal(dg.GetTagForHEAD("."), "test-tag")
	}
}

func Test_DefaultGitter_GetBuild(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	is.True(dg != nil)
	is.True(dg.GetBuild(".") != "")
	is.Equal(dg.GetBuild("/"), "")
}
