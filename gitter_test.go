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

func Test_CheckGitRepo_SucceedsForCurrent(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	repo, err := dg.CheckGitRepo(".")
	is.NoErr(err)
	is.True(repo != "")
}

func Test_CheckGitRepo_SuceedsForCmdMkver(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	repo, err := dg.CheckGitRepo("./cmd/mkver")
	is.NoErr(err)
	is.True(repo != "")
}

func Test_CheckGitRepo_FailsForRoot(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	repo, err := dg.CheckGitRepo("/")
	is.True(repo == "/")
	is.True(err != nil)
}

func Test_CheckGitRepo_IgnoresFileNamedGit(t *testing.T) {
	is := is.New(t)
	const fileNamedGit = "./cmd/mkver/.git"
	if _, err := os.Stat(fileNamedGit); err != nil {
		if f, err := os.Create(fileNamedGit); err == nil {
			defer f.Close()
			defer os.Remove(fileNamedGit)
			dg, err := NewDefaultGitter("git")
			is.NoErr(err)
			repo, err := dg.CheckGitRepo("./cmd/mkver")
			is.NoErr(err)
			is.True(repo != "")
		}
	} else {
		t.Logf("warning: '%s' already exists\n", fileNamedGit)
	}
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

func Test_DefaultGitter_GetTags(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	is.True(dg != nil)
	is.Equal(dg.GetTags("/"), nil)
	alltags := dg.GetTags(".")
	if len(alltags) == 0 {
		t.Error("no tags")
	}
}

func Test_DefaultGitter_GetCurrentTreeHash(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	is.True(dg != nil)
	is.Equal(dg.GetCurrentTreeHash("/"), "")
	s := dg.GetCurrentTreeHash(".")
	if len(s) == 0 {
		t.Error("no tree hash")
	}
}

func Test_DefaultGitter_GetTreeHash(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	is.True(dg != nil)
	is.Equal(dg.GetTreeHash("/", "v1.0.0"), "")
	is.Equal(dg.GetTreeHash(".", "v1.0.0"), "0efbb9e3dce88d590a0bfa4b67e0d5341d2d8cb8")
}

func Test_DefaultGitter_GetCommits(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	is.True(dg != nil)
	is.Equal(dg.GetCommits("/"), nil)
	commits := dg.GetCommits(".")
	for _, s := range commits {
		if s == "40dadd20bd3bf243e1597e06ebec5ba61b7099af" {
			return
		}
	}
	t.Error(commits)
}

func Test_DefaultGitter_GetClosestTag(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	is.True(dg != nil)
	is.Equal(dg.GetClosestTag("/", ""), "")
	tag := dg.GetClosestTag(".", "2e4ae09e864e47f9f0505c14206d7438d811e1ea")
	if tag != "v1.8.0" {
		t.Error(tag)
	}
}

func Test_DefaultGitter_GetBranchFromTag(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	is.True(dg != nil)
	is.Equal(dg.GetBranchesFromTag("/", "refs/tags/v1.0.0"), nil)
	is.Equal(dg.GetBranchesFromTag(".", "refs/tags/v1.0.0"), []string{"main"})
}

func Test_DefaultGitter_GetBuild(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	is.True(dg != nil)
	is.True(dg.GetBuild(".") != "")
	is.Equal(dg.GetBuild("/"), "")
}

func Test_DefaultGitter_FetchTags(t *testing.T) {
	is := is.New(t)
	dg, err := NewDefaultGitter("git")
	is.NoErr(err)
	is.True(dg != nil)
	dg.FetchTags(".")
}
