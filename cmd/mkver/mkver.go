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
	_, err = f.WriteString(content)
	return
}

var (
	flagName  = flag.String("name", "", "write Go source with given package name")
	flagRepo  = flag.String("repo", "", "repository to examine")
	flagOut   = flag.String("out", "", "file path relative to repo to write to (defaults to stdout)")
	flagGit   = flag.String("git", "git", "name of Git executable")
	flagFetch = flag.Bool("fetch", false, "fetch remote tags")
)

func main() {
	flag.Parse()

	var err error
	var repoDir string
	var vs *makeversion.VersionStringer
	var vi makeversion.VersionInfo
	var content string

	if repoDir = os.ExpandEnv(*flagRepo); repoDir == "" {
		if repoDir = flag.Arg(0); repoDir == "" {
			repoDir = "."
		}
	}

	if vs, err = makeversion.NewVersionStringer(*flagGit); err == nil {
		if repoDir, err = vs.Git.CheckGitRepo(repoDir); err == nil {
			if *flagFetch {
				err = vs.Git.FetchTags(repoDir)
			}
			if err == nil {
				if vi, err = vs.GetVersion(repoDir); err == nil {
					if content, err = vi.Render(*flagName); err == nil {
						outpath := os.ExpandEnv(*flagOut)
						if outpath != "" {
							outpath = path.Join(repoDir, outpath)
						}
						err = writeOutput(outpath, content)
					}
				}
			}
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%q: %v\n", repoDir, err.Error())
		os.Exit(1)
	}
}
