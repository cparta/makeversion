package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/cparta/makeversion"
)

func writeOutput(fileName, content string) (err error) {
	f := os.Stdout
	if len(fileName) > 0 {
		fileName = path.Clean(fileName)
		if f, err = os.Create(fileName); err != nil /* #nosec G304 */ {
			return
		}
		defer f.Close()
	}
	_, err = f.WriteString(content)
	return
}

var (
	flagName = flag.String("name", "", "write Go source with given package name")
	flagOut  = flag.String("out", "", "file to write to (defaults to stdout)")
	flagGit  = flag.String("git", "git", "name of Git executable")
	flagRepo = flag.String("repo", ".", "repository to examine")
)

func main() {
	flag.Parse()

	var err error
	var repoDir string
	var vs *makeversion.VersionStringer
	var vi makeversion.VersionInfo
	var content string

	if repoDir, err = makeversion.CheckGitRepo(*flagRepo); err != nil {
		repoDir = *flagRepo
		fmt.Fprintf(os.Stderr, "warning: '%s' is not a git repository\n", repoDir)
	}

	if vs, err = makeversion.NewVersionStringer(*flagGit); err == nil {
		if vi, err = vs.GetVersion(repoDir); err == nil {
			if content, err = vi.Render(*flagName); err == nil {
				err = writeOutput(*flagOut, content)
			}
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
