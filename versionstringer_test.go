// @author jli@cparta.se

package makeversion

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

type MockEnvironment map[string]string

func (me MockEnvironment) Getenv(key string) string {
	return me[key]
}

func (me MockEnvironment) LookupEnv(key string) (val string, ok bool) {
	val, ok = me[key]
	return
}

type MockGitter struct {
	TopTag string
}

func (mg MockGitter) GetTag(repo string) string {
	return repo
}

func (mg MockGitter) GetBranch(repo string) string {
	return repo
}

func (mg MockGitter) GetBranchesFromTag(repo, tag string) (branches []string) {
	if strings.HasPrefix(tag, "v1.0") {
		branches = append(branches, "main")
	}
	if strings.HasPrefix(tag, "v1") {
		branches = append(branches, "onepointoh")
	}
	return
}

func (mg MockGitter) GetBuild(repo string) string {
	return repo
}

func (mg MockGitter) GetTagForHEAD(repo string) string {
	return mg.TopTag
}

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
	git := MockGitter{}
	vs := VersionStringer{Git: git, Env: env}

	tag, err := vs.GetTag("")
	is.NoErr(err)
	is.Equal("v0.0.0", tag)

	tag, err = vs.GetTag("v1.2.3")
	is.NoErr(err)
	is.Equal("v1.2.3", tag)

	tag, err = vs.GetTag("foo")
	is.True(err != nil)
	is.Equal("foo", tag)

	env["CI_COMMIT_TAG"] = "v3"
	tag, err = vs.GetTag("")
	is.NoErr(err)
	is.Equal("v3", tag)
}

func Test_VersionStringer_GetBranch(t *testing.T) {
	is := is.New(t)
	env := MockEnvironment{}
	git := MockGitter{}
	vs := VersionStringer{Git: git, Env: env}

	text, name := vs.GetBranch("branch")
	is.Equal("branch", name)
	is.Equal("branch", text)

	text, name = vs.GetBranch("branch.with..dots")
	is.Equal("branch.with..dots", name)
	is.Equal("branch-with-dots", text)

	env["CI_COMMIT_REF_NAME"] = "gitlab-branch"
	text, name = vs.GetBranch("")
	is.Equal("gitlab-branch", name)
	is.Equal("gitlab-branch", text)
	delete(env, "CI_COMMIT_REF_NAME")

	env["GITHUB_REF_NAME"] = "github.branch"
	text, name = vs.GetBranch("")
	is.Equal("github.branch", name)
	is.Equal("github-branch", text)
	delete(env, "GITHUB_REF_NAME")
}

func Test_VersionStringer_GetBranchFromTag_GitLab(t *testing.T) {
	is := is.New(t)
	env := MockEnvironment{}
	git := MockGitter{}
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
	git := MockGitter{}
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
	git := MockGitter{}
	vs := VersionStringer{Git: git, Env: env}

	build := vs.GetBuild("123")
	is.Equal("123", build)

	env["CI_PIPELINE_IID"] = "456"
	build = vs.GetBuild("")
	is.Equal("456", build)
	delete(env, "CI_PIPELINE_IID")

	env["GITHUB_RUN_NUMBER"] = "789"
	build = vs.GetBuild("345")
	is.Equal("789", build)
	delete(env, "CI_PIPELINE_IID")
}

func Test_VersionStringer_GetVersion(t *testing.T) {
	is := is.New(t)
	env := MockEnvironment{}
	git := &MockGitter{}
	vs := VersionStringer{Git: git, Env: env}

	vi, err := vs.GetVersion("v1", false)
	is.NoErr(err)
	is.Equal("v1-v1.v1", vi.Version)

	vi, err = vs.GetVersion("", false)
	is.NoErr(err)
	is.Equal("v0.0.0", vi.Version)

	env["CI_COMMIT_REF_NAME"] = "HEAD"
	vi, err = vs.GetVersion("", false)
	is.NoErr(err)
	is.Equal("v0.0.0-HEAD", vi.Version)

	delete(env, "CI_COMMIT_REF_NAME")
	env["GITHUB_RUN_NUMBER"] = "789"
	vi, err = vs.GetVersion("", false)
	is.NoErr(err)
	is.Equal("v0.0.0-789", vi.Version)

	env["CI_COMMIT_REF_NAME"] = "*Branch--.--ONE*-*"
	env["GITHUB_RUN_NUMBER"] = "789"
	vi, err = vs.GetVersion("v2.0", false)
	is.NoErr(err)
	is.Equal("v2.0-branch-one.789", vi.Version)

	vi, err = vs.GetVersion("v3.4.5", true)
	is.NoErr(err)
	is.Equal("v3.4.5-branch-one.789", vi.Version)

	env["CI_COMMIT_REF_NAME"] = "main"
	git.TopTag = "v3.4.5"
	vi, err = vs.GetVersion("v3.4.5", true)
	is.NoErr(err)
	is.Equal("v3.4.5", vi.Version)
}
