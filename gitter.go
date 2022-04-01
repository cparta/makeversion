package makeversion

import (
	"os/exec"
	"strconv"
	"strings"
)

// Gitter is an interface exposing the required Git functionality
type Gitter interface {
	GetTag(repo string) string
	GetBranch(repo string) string
	GetBuild(repo string) string
}

type DefaultGitter string

func NewDefaultGitter(gitBin string) (gitter Gitter, err error) {
	if gitBin, err = exec.LookPath(gitBin); err == nil {
		gitter = DefaultGitter(gitBin)
	}
	return
}

// GetTag returns the latest Git tag that starts with a lowercase 'v' followed by a number, otherwise an empty string.
func (dg DefaultGitter) GetTag(repo string) string {
	if b, _ := exec.Command(string(dg), "-C", repo, "describe", "--tags", "--match", "v[0-9]*", "--abbrev=0").Output(); len(b) > 0 {
		return strings.TrimSpace(string(b))
	}
	return ""
}

// GetBranch returns the current branch in the repository, or the string "HEAD"
func (dg DefaultGitter) GetBranch(repo string) string {
	if b, _ := exec.Command(string(dg), "-C", repo, "rev-parse", "--abbrev-ref", "HEAD").Output(); len(b) > 0 {
		return strings.TrimSpace(string(b))
	}
	return "HEAD"
}

// GetBuild returns the number of commits in the currently checked out branch as a string, or an empty string
func (dg DefaultGitter) GetBuild(repo string) string {
	if b, _ := exec.Command(string(dg), "-C", repo, "rev-list", "HEAD", "--count").Output(); len(b) > 0 {
		str := strings.TrimSpace(string(b))
		if num, err := strconv.Atoi(str); err == nil && num > 0 {
			return str
		}
	}
	return ""
}
