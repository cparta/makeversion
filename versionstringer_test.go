// @author jli@cparta.se

package makeversion

import (
	"testing"

	"github.com/matryer/is"
)

func Test_NewVersionStringer_SucceedsNormally(t *testing.T) {
	is := is.New(t)
	vs, err := NewVersionStringer("git")
	is.NoErr(err)
	is.True(vs != nil)
}

func Test_NewVersionStringer_FailsWithBadBinary(t *testing.T) {
	is := is.New(t)
	vs, err := NewVersionStringer("./versionstringer.go")
	is.True(err != nil)
	is.Equal(vs, nil)
}

func Test_VersionStringer_IsEnvTrue(t *testing.T) {
	is := is.New(t)
	vs := VersionStringer{
		Env: MockEnvironment{
			"TEST_EMPTY":      "",
			"TEST_FALSE":      "false",
			"TEST_TRUE_LOWER": "true",
			"TEST_TRUE_UPPER": "TRUE",
		},
	}
	is.Equal(vs.IsEnvTrue("TEST_MISSING"), false)
	is.Equal(vs.IsEnvTrue("TEST_MISSING"), false)
	is.Equal(vs.IsEnvTrue("TEST_EMPTY"), false)
	is.Equal(vs.IsEnvTrue("TEST_FALSE"), false)
	is.Equal(vs.IsEnvTrue("TEST_TRUE_LOWER"), true)
	is.Equal(vs.IsEnvTrue("TEST_TRUE_UPPER"), true)
}

func Test_VersionStringer_IsReleaseBranch(t *testing.T) {
	is := is.New(t)
	const branchName = "testbranch"
	env := MockEnvironment{}
	vs := VersionStringer{Env: env}

	is.True(vs.IsReleaseBranch("default"))
	is.True(vs.IsReleaseBranch("main"))
	is.True(vs.IsReleaseBranch("master"))
	is.True(!vs.IsReleaseBranch(branchName))

	env["CI_DEFAULT_BRANCH"] = branchName
	is.True(vs.IsReleaseBranch(branchName))
	delete(env, "CI_DEFAULT_BRANCH")

	env["CI_COMMIT_REF_PROTECTED"] = "true"
	is.True(vs.IsReleaseBranch(branchName))
	delete(env, "CI_COMMIT_REF_PROTECTED")

	env["GITHUB_REF_PROTECTED"] = "true"
	is.True(vs.IsReleaseBranch(branchName))
	delete(env, "GITHUB_REF_PROTECTED")

	is.True(!vs.IsReleaseBranch(branchName))
}

func Test_VersionStringer_GetTag(t *testing.T) {
	is := is.New(t)
	env := MockEnvironment{}
	git := &MockGitter{}
	vs := VersionStringer{Git: git, Env: env}

	var tag string
	var sametree bool

	tag, sametree = vs.GetTag("/")
	is.Equal("v0.0.0", tag)
	is.Equal(false, sametree)

	tag, sametree = vs.GetTag(".")
	is.Equal("v6.0.0", tag)
	is.Equal(false, sametree)

	git.treehash = "tree-4"
	tag, sametree = vs.GetTag(".")
	is.Equal("v4.0.0", tag)
	is.Equal(true, sametree)

	git.treehash = ""
	env["CI_COMMIT_TAG"] = "v3"
	tag, sametree = vs.GetTag(".")
	is.Equal("v3", tag)
	is.Equal(true, sametree)
}

func Test_VersionStringer_GetBranch(t *testing.T) {
	is := is.New(t)
	env := MockEnvironment{}
	git := &MockGitter{}
	vs := VersionStringer{Git: git, Env: env}

	git.branch = "zomg"
	text, name := vs.GetBranch(".")
	is.Equal("zomg", text)
	is.Equal("zomg", name)

	git.branch = "branch.with..dots"
	text, name = vs.GetBranch(".")
	is.Equal("branch-with-dots", text)
	is.Equal("branch.with..dots", name)

	env["CI_COMMIT_REF_NAME"] = "gitlab---branch"
	text, name = vs.GetBranch(".")
	is.Equal("gitlab-branch", text)
	is.Equal("gitlab---branch", name)
	delete(env, "CI_COMMIT_REF_NAME")

	env["GITHUB_REF_NAME"] = "github.branch"
	text, name = vs.GetBranch(".")
	is.Equal("github.branch", name)
	is.Equal("github-branch", text)
	delete(env, "GITHUB_REF_NAME")
}

func Test_VersionStringer_GetBranchFromTag_GitLab(t *testing.T) {
	is := is.New(t)
	env := MockEnvironment{}
	git := &MockGitter{}
	vs := VersionStringer{Git: git, Env: env}

	env["CI_COMMIT_TAG"] = "v1.0.0"
	env["CI_COMMIT_REF_NAME"] = "v1.0.0"
	text, name := vs.GetBranch(".")
	is.Equal("main", name)
	is.Equal("main", text)
}

func Test_VersionStringer_GetBranchFromTag_GitHub(t *testing.T) {
	is := is.New(t)
	env := MockEnvironment{}
	git := &MockGitter{}
	vs := VersionStringer{Git: git, Env: env}

	env["GITHUB_REF_TYPE"] = "tag"
	env["GITHUB_REF_NAME"] = "v1.0.0"
	text, name := vs.GetBranch(".")
	is.Equal("main", name)
	is.Equal("main", text)

	env["GITHUB_REF_NAME"] = "v1"
	text, name = vs.GetBranch(".")
	is.Equal("onepointoh", name)
	is.Equal("onepointoh", text)
}

func Test_VersionStringer_GetBuild(t *testing.T) {
	is := is.New(t)
	env := MockEnvironment{}
	git := &MockGitter{}
	vs := VersionStringer{Git: git, Env: env}

	build := vs.GetBuild(".")
	is.Equal("build", build)

	env["CI_PIPELINE_IID"] = "456"
	build = vs.GetBuild(".")
	is.Equal("456", build)
	delete(env, "CI_PIPELINE_IID")

	env["GITHUB_RUN_NUMBER"] = "789"
	build = vs.GetBuild(".")
	is.Equal("789", build)
	delete(env, "CI_PIPELINE_IID")
}

func Test_VersionStringer_GetVersion(t *testing.T) {
	is := is.New(t)
	env := MockEnvironment{}
	git := &MockGitter{}
	vs := VersionStringer{Git: git, Env: env}

	vi, err := vs.GetVersion("/") // invalid repo
	is.NoErr(err)
	is.Equal("v0.0.0", vi.Version)

	vi, err = vs.GetVersion(".")
	is.NoErr(err)
	is.Equal("v6.0.0-main.build", vi.Version)

	git.treehash = "tree-6"
	vi, err = vs.GetVersion(".")
	is.NoErr(err)
	is.Equal("v6.0.0", vi.Version)

	git.treehash = ""
	env["CI_COMMIT_REF_NAME"] = "HEAD"
	vi, err = vs.GetVersion(".")
	is.NoErr(err)
	is.Equal("v6.0.0-head.build", vi.Version)

	delete(env, "CI_COMMIT_REF_NAME")
	env["GITHUB_RUN_NUMBER"] = "789"
	vi, err = vs.GetVersion(".")
	is.NoErr(err)
	is.Equal("v6.0.0-main.789", vi.Version)

	env["CI_COMMIT_REF_NAME"] = "*Branch--.--ONE*-*"
	env["GITHUB_RUN_NUMBER"] = "789"
	vi, err = vs.GetVersion(".")
	is.NoErr(err)
	is.Equal("v6.0.0-branch-one.789", vi.Version)

	env["CI_COMMIT_REF_NAME"] = "main"
	vi, err = vs.GetVersion(".")
	is.NoErr(err)
	is.Equal("v6.0.0-main.789", vi.Version)
}
