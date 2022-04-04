package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/cparta/makeversion"
)

func checkDir(p string) error {
	fi, err := os.Stat(p)
	if err == nil {
		if fi.IsDir() {
			return nil
		}
		err = fmt.Errorf("'%s' is not a directory", p)
	}
	return err
}

func dirOrParentHasGitSubdir(s string) bool {
	if err := checkDir(path.Join(s, ".git")); err == nil {
		return true
	}
	if s = path.Dir(s); s == "/" {
		return false
	}
	return dirOrParentHasGitSubdir(s)
}

func checkRepoDir(s string) (repo string, err error) {
	if repo, err = filepath.Abs(s); err == nil {
		if err = checkDir(repo); err == nil {
			if !dirOrParentHasGitSubdir(repo) {
				err = errors.New("can't find .git directory")
			}
		}
	}
	return
}

func writeOutput(fileName, content string) (err error) {
	f := os.Stdout
	if len(fileName) > 0 {
		fileName = path.Clean(fileName)
		if f, err = os.Create(fileName); err != nil {
			return
		}
		defer f.Close()
	}
	_, err = f.WriteString(content)
	return
}

var (
	flagName    = flag.String("name", "", "write Go source with given package name")
	flagOut     = flag.String("out", "", "file to write to (defaults to stdout)")
	flagGit     = flag.String("git", "git", "name of Git executable")
	flagRepo    = flag.String("repo", ".", "repository to examine")
	flagRelease = flag.Bool("release", false, "output release version without build info suffix")
)

func main() {
	flag.Parse()

	var err error
	var repoDir string
	var vs *makeversion.VersionStringer
	var vi makeversion.VersionInfo
	var content string

	if repoDir, err = checkRepoDir(*flagRepo); err == nil {
		if vs, err = makeversion.NewVersionStringer(*flagGit); err == nil {
			if vi, err = vs.GetVersion(repoDir, *flagRelease); err == nil {
				if content, err = vi.Render(*flagName); err == nil {
					err = writeOutput(*flagOut, content)
				}
			}
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
