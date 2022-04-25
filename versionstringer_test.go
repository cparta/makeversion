// @author jli@cparta.se

package makeversion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockEnvironment map[string]string

func (me MockEnvironment) Getenv(key string) string {
	return me[key]
}

func (me MockEnvironment) LookupEnv(key string) (val string, ok bool) {
	val, ok = me[key]
	return
}

type MockGitter struct{}

func (mg MockGitter) GetTag(repo string) string {
	return repo
}

func (mg MockGitter) GetBranch(repo string) string {
	return repo
}

func (mg MockGitter) GetBuild(repo string) string {
	return repo
}

func Test_NewVersionStringer_SucceedsNormally(t *testing.T) {
	vs, err := NewVersionStringer("git")
	assert.NoError(t, err)
	assert.NotNil(t, vs)
}

func Test_NewVersionStringer_FailsWithBadBinary(t *testing.T) {
	vs, err := NewVersionStringer("./versionstringer.go")
	assert.Error(t, err)
	assert.Nil(t, vs)
}

func Test_VersionStringer_IsEnvTrue(t *testing.T) {
	vs := VersionStringer{
		Env: MockEnvironment{
			"TEST_EMPTY":      "",
			"TEST_FALSE":      "false",
			"TEST_TRUE_LOWER": "true",
			"TEST_TRUE_UPPER": "TRUE",
		},
	}
	assert.False(t, vs.IsEnvTrue("TEST_MISSING"))
	assert.False(t, vs.IsEnvTrue("TEST_EMPTY"))
	assert.False(t, vs.IsEnvTrue("TEST_FALSE"))
	assert.True(t, vs.IsEnvTrue("TEST_TRUE_LOWER"))
	assert.True(t, vs.IsEnvTrue("TEST_TRUE_UPPER"))
}

func Test_VersionStringer_IsReleaseBranch(t *testing.T) {
	const branchName = "testbranch"
	env := MockEnvironment{}
	vs := VersionStringer{Env: env}

	assert.True(t, vs.IsReleaseBranch("default"))
	assert.True(t, vs.IsReleaseBranch("main"))
	assert.True(t, vs.IsReleaseBranch("master"))
	assert.False(t, vs.IsReleaseBranch(branchName))

	env["CI_DEFAULT_BRANCH"] = branchName
	assert.True(t, vs.IsReleaseBranch(branchName))
	delete(env, "CI_DEFAULT_BRANCH")

	env["CI_COMMIT_REF_PROTECTED"] = "true"
	assert.True(t, vs.IsReleaseBranch(branchName))
	delete(env, "CI_COMMIT_REF_PROTECTED")

	env["GITHUB_REF_PROTECTED"] = "true"
	assert.True(t, vs.IsReleaseBranch(branchName))
	delete(env, "GITHUB_REF_PROTECTED")

	assert.False(t, vs.IsReleaseBranch(branchName))
}

func Test_VersionStringer_GetTag(t *testing.T) {
	env := MockEnvironment{}
	git := MockGitter{}
	vs := VersionStringer{Git: git, Env: env}

	tag, err := vs.GetTag("")
	assert.NoError(t, err)
	assert.Equal(t, "v0.0.0", tag)

	tag, err = vs.GetTag("v1.2.3")
	assert.NoError(t, err)
	assert.Equal(t, "v1.2.3", tag)

	tag, err = vs.GetTag("foo")
	assert.Error(t, err)
	assert.Equal(t, "foo", tag)

	env["CI_COMMIT_TAG"] = "v3"
	tag, err = vs.GetTag("")
	assert.NoError(t, err)
	assert.Equal(t, "v3", tag)
}

func Test_VersionStringer_GetBranch(t *testing.T) {
	env := MockEnvironment{}
	git := MockGitter{}
	vs := VersionStringer{Git: git, Env: env}

	text, name := vs.GetBranch("branch")
	assert.Equal(t, "branch", name)
	assert.Equal(t, "branch", text)

	text, name = vs.GetBranch("branch.with..dots")
	assert.Equal(t, "branch.with..dots", name)
	assert.Equal(t, "branch-with-dots", text)

	env["CI_COMMIT_REF_NAME"] = "gitlab-branch"
	text, name = vs.GetBranch("")
	assert.Equal(t, "gitlab-branch", name)
	assert.Equal(t, "gitlab-branch", text)
	delete(env, "CI_COMMIT_REF_NAME")

	env["GITHUB_REF_NAME"] = "github.branch"
	text, name = vs.GetBranch("")
	assert.Equal(t, "github.branch", name)
	assert.Equal(t, "github-branch", text)
	delete(env, "GITHUB_REF_NAME")
}

func Test_VersionStringer_GetBuild(t *testing.T) {
	env := MockEnvironment{}
	git := MockGitter{}
	vs := VersionStringer{Git: git, Env: env}

	build := vs.GetBuild("123")
	assert.Equal(t, "123", build)

	env["CI_PIPELINE_IID"] = "456"
	build = vs.GetBuild("")
	assert.Equal(t, "456", build)
	delete(env, "CI_PIPELINE_IID")

	env["GITHUB_RUN_NUMBER"] = "789"
	build = vs.GetBuild("345")
	assert.Equal(t, "789", build)
	delete(env, "CI_PIPELINE_IID")
}

func Test_VersionStringer_GetVersion(t *testing.T) {
	env := MockEnvironment{}
	git := MockGitter{}
	vs := VersionStringer{Git: git, Env: env}

	vi, err := vs.GetVersion("v1", false)
	assert.NoError(t, err)
	assert.Equal(t, "v1-v1.v1", vi.Version)

	vi, err = vs.GetVersion("", false)
	assert.NoError(t, err)
	assert.Equal(t, "v0.0.0", vi.Version)

	env["CI_COMMIT_REF_NAME"] = "HEAD"
	vi, err = vs.GetVersion("", false)
	assert.NoError(t, err)
	assert.Equal(t, "v0.0.0-HEAD", vi.Version)

	delete(env, "CI_COMMIT_REF_NAME")
	env["GITHUB_RUN_NUMBER"] = "789"
	vi, err = vs.GetVersion("", false)
	assert.NoError(t, err)
	assert.Equal(t, "v0.0.0-789", vi.Version)

	env["CI_COMMIT_REF_NAME"] = "*Branch--.--ONE*-*"
	env["GITHUB_RUN_NUMBER"] = "789"
	vi, err = vs.GetVersion("v2.0", false)
	assert.NoError(t, err)
	assert.Equal(t, "v2.0-branch-one.789", vi.Version)

	vi, err = vs.GetVersion("v3.4.5", true)
	assert.Error(t, err)

	env["CI_COMMIT_REF_NAME"] = "main"
	vi, err = vs.GetVersion("v3.4.5", true)
	assert.NoError(t, err)
	assert.Equal(t, "v3.4.5", vi.Version)
}
