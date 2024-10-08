package makeversion

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// Gitter is an interface exposing the required Git functionality
type Gitter interface {
	// GetTag returns the latest Git tag that starts with a lowercase 'v' followed by a number, otherwise an empty string.
	GetTag(repo string) string
	// GetBranch returns the current branch in the repository, or the string "HEAD" if we're in a Git repo, otherwise an empty string.
	GetBranch(repo string) string
	// GetBranchesFromTag returns the non-HEAD branches in the repository that have the tag, otherwise an empty string.
	GetBranchesFromTag(repo, tag string) []string
	// GetBuild returns the number of commits in the currently checked out branch as a string, or an empty string
	GetBuild(repo string) string
	// GetTagForHEAD returns the tag that exactly match the current HEAD, or an empty string
	GetTagForHEAD(repo string) string
}

type DefaultGitter string

func NewDefaultGitter(gitBin string) (gitter Gitter, err error) {
	if gitBin, err = exec.LookPath(gitBin); err == nil {
		gitter = DefaultGitter(gitBin)
	}
	return
}

// checkDir checks that the given path is accessible and is a directory.
// Returns nil if it is, else an error.
func checkDir(dir string) (err error) {
	var fi os.FileInfo
	if fi, err = os.Stat(dir); err == nil {
		if !fi.IsDir() {
			err = fmt.Errorf("'%s' is not a directory", dir)
		}
	}
	return
}

// dirOrParentHasGitSubdir returns the name of a directory containing
// a '.git' subdirectory or an empty string. It searches starting from
// the given directory and looks in that and it's parents.
func dirOrParentHasGitSubdir(dir string) string {
	if dir != "/" && dir != "." {
		if checkDir(path.Join(dir, ".git")) == nil {
			return dir
		}
		return dirOrParentHasGitSubdir(path.Dir(dir))
	}
	return ""
}

// CheckGitRepo checks that the given directory is part of a git repository,
// meaning that it or one of it's parent directories has a '.git' subdirectory.
// If it is, it returns the absolute path of the git repo and a nil error.
func CheckGitRepo(dir string) (repo string, err error) {
	if dir, err = filepath.Abs(dir); err == nil {
		if err = checkDir(dir); err == nil {
			if repo = dirOrParentHasGitSubdir(dir); repo == "" {
				err = errors.New("can't find .git directory")
			}
		}
	}
	return
}

func (dg DefaultGitter) GetTag(repo string) string {
	if repo, _ = CheckGitRepo(repo); repo != "" {
		if b, _ := exec.Command(string(dg), "-C", repo, "describe", "--tags", "--match", "v[0-9]*", "--abbrev=0").Output(); len(b) > 0 /* #nosec G204 */ {
			return strings.TrimSpace(string(b))
		}
	}
	return ""
}

func lastName(s string) string {
	if idx := strings.LastIndexByte(s, '/'); idx > -1 {
		s = s[idx+1:]
	}
	return s
}

func (dg DefaultGitter) GetBranchesFromTag(repo, tag string) (branches []string) {
	if repo, _ = CheckGitRepo(repo); repo != "" {
		tag = strings.TrimPrefix(tag, "refs/")
		tag = strings.TrimPrefix(tag, "tags/")
		if b, _ := exec.Command(string(dg), "-C", repo, "branch", "--all", "--no-color", "--contains", "tags/"+tag).Output(); len(b) > 0 /* #nosec G204 */ {
			for _, s := range strings.Split(string(b), "\n") {
				if s = strings.TrimSpace(s); len(s) > 1 {
					if !strings.Contains(s, "HEAD") {
						starred := s[0] == '*'
						s = strings.TrimSpace(strings.TrimPrefix(s, "*"))
						if len(s) > 0 && !strings.Contains(s, " ") {
							branches = append(branches, lastName(s))
							if starred {
								branches = branches[len(branches)-1:]
								break
							}
						}
					}
				}
			}
		}
	}
	return
}

func (dg DefaultGitter) GetTagForHEAD(repo string) (tag string) {
	if repo, _ = CheckGitRepo(repo); repo != "" {
		if b, _ := exec.Command(string(dg), "-C", repo, "describe", "--exact-match", "--tags", "HEAD").Output(); len(b) > 0 /* #nosec G204 */ {
			if s := strings.TrimSpace(string(b)); len(s) > 1 {
				tag = s
			}
		}
	}
	return
}

func (dg DefaultGitter) GetBranch(repo string) (branch string) {
	if repo, _ = CheckGitRepo(repo); repo != "" {
		branch = "HEAD"
		if b, _ := exec.Command(string(dg), "-C", repo, "branch", "--show-current").Output(); len(b) > 0 /* #nosec G204 */ {
			branch = strings.TrimSpace(string(b))
		}
	}
	return
}

func (dg DefaultGitter) GetBuild(repo string) string {
	if repo, _ = CheckGitRepo(repo); repo != "" {
		if b, _ := exec.Command(string(dg), "-C", repo, "rev-list", "HEAD", "--count").Output(); len(b) > 0 /* #nosec G204 */ {
			str := strings.TrimSpace(string(b))
			if num, err := strconv.Atoi(str); err == nil && num > 0 {
				return str
			}
		}
	}
	return ""
}
