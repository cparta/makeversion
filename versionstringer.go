// @author jli@cparta.se

package makeversion

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	reCheckTag  = regexp.MustCompile(`^v\d+(\.\d+(\.\d+)?)?$`)
	reOnlyWords = regexp.MustCompile(`[^\w]`)
)

type VersionStringer struct {
	Git Gitter      // Git
	Env Environment // environment
}

// NewVersionStringer returns a VersionStringer ready to examine
// the git repositories using the given Git binary.
func NewVersionStringer(gitBin string) (vs *VersionStringer, err error) {
	var git Gitter
	if git, err = NewDefaultGitter(gitBin); err == nil {
		vs = &VersionStringer{
			Git: git,
			Env: OsEnvironment{},
		}
	}
	return
}

// IsEnvTrue returns true if the given environment variable
// exists and is set to the string "true" (not case sensitive).
func (vs *VersionStringer) IsEnvTrue(envvar string) bool {
	return "true" == strings.ToLower(strings.TrimSpace(vs.Env.Getenv(envvar)))
}

// IsReleaseBranch returns true if the given branch name should
// be allowed to use 'release mode', where the version string
// doesn't contains build information suffix.
func (vs *VersionStringer) IsReleaseBranch(branchName string) bool {
	// A GitLab or GitHub protected branch allows release mode.
	if vs.IsEnvTrue("CI_COMMIT_REF_PROTECTED") || vs.IsEnvTrue("GITHUB_REF_PROTECTED") {
		return true
	}

	// If the branch isn't protected, we only allow release
	// mode for the 'default' branch.

	// GitLab gives us the default branch name directly.
	if defBranch, ok := vs.Env.LookupEnv("CI_DEFAULT_BRANCH"); ok {
		defBranch = strings.TrimSpace(defBranch)
		return branchName == defBranch
	}

	// Fallback to common default branch names.
	switch branchName {
	case "default":
		return true
	case "master":
		return true
	case "main":
		return true
	}

	return false
}

// GetTag returns the git version tag. Returns an error if the tag is not in the form "vX.Y.Z".
func (vs *VersionStringer) GetTag(repo string) (tag string, err error) {
	if tag = strings.TrimSpace(vs.Env.Getenv("CI_COMMIT_TAG")); tag == "" {
		if tag = vs.Git.GetTag(repo); tag == "" {
			tag = "v0.0.0"
		}
	}
	if !reCheckTag.MatchString(tag) {
		err = fmt.Errorf("tag doesn't match 'vN(.N(.N))': '%s'", tag)
	}
	return
}

// GetBranch returns the current branch as a string suitable
// for inclusion in the semver text as well as the actual
// branch name in the build system or Git.
func (vs *VersionStringer) GetBranch(repo string) (branchText, branchName string) {
	if branchName = strings.TrimSpace(vs.Env.Getenv("CI_COMMIT_REF_NAME")); branchName == "" {
		if branchName = strings.TrimSpace(vs.Env.Getenv("GITHUB_REF_NAME")); branchName == "" {
			branchName = vs.Git.GetBranch(repo)
		}
	}

	branchText = branchName
	if branchName != "HEAD" {
		branchText = strings.ReplaceAll(strings.ToLower(reOnlyWords.ReplaceAllString(branchText, "-")), "--", "-")
	}

	return
}

// GetBuild returns the build counter. This is taken from the CI system if available,
// otherwise the Git commit count is used.
func (vs *VersionStringer) GetBuild(repo string) (build string) {
	if build = strings.TrimSpace(vs.Env.Getenv("CI_PIPELINE_IID")); build == "" {
		if build = strings.TrimSpace(vs.Env.Getenv("GITHUB_RUN_NUMBER")); build == "" {
			build = vs.Git.GetBuild(repo)
		}
	}
	return
}

// GetVersion returns a version string for the source code in the Git repository.
func (vs *VersionStringer) GetVersion(repo string, forRelease bool) (vi VersionInfo, err error) {
	if vi.Tag, err = vs.GetTag(repo); err == nil {
		vi.Version = vi.Tag
		vi.Build = vs.GetBuild(repo)
		branchText, branchName := vs.GetBranch(repo)
		vi.Branch = branchName
		if forRelease {
			if !vs.IsReleaseBranch(branchName) {
				err = fmt.Errorf("release version must be on default branch, not '%s'\n", branchName)
			}
		} else {
			vi.Version += "-"
			if branchText != "" {
				vi.Version += branchText + "."
			}
			vi.Version += vi.Build
		}
	}
	return
}
