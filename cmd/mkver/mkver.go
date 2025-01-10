package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/cparta/makeversion/v2"
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
	fmt.Println(fileName, content)
	_, err = f.WriteString(content)
	return
}

var (
	flagName = flag.String("name", "", "write Go source with given package name")
	flagRepo = flag.String("repo", ".", "repository to examine")
	flagOut  = flag.String("out", "", "file path relative to repo to write to (defaults to stdout)")
	flagGit  = flag.String("git", "git", "name of Git executable")
)

func main() {
	flag.Parse()

	var err error
	var repoDir string
	var vs *makeversion.VersionStringer
	var vi makeversion.VersionInfo
	var content string

	repoDir = os.ExpandEnv(*flagRepo)
	if repoDir, err = makeversion.CheckGitRepo(repoDir); err != nil {
		repoDir = *flagRepo
		fmt.Fprintf(os.Stderr, "warning: '%s' is not a git repository\n", repoDir)
	}

	if vs, err = makeversion.NewVersionStringer(*flagGit); err == nil {
		if vi, err = vs.GetVersion(repoDir); err == nil {
			if content, err = vi.Render(*flagName); err == nil {
				outpath := path.Join(repoDir, os.ExpandEnv(*flagOut))
				err = writeOutput(outpath, content)
			}
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
