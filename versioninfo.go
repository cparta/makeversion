package makeversion

import (
	"fmt"
	"go/token"
	"os"
	"path"
	"strconv"
	"time"
)

type VersionInfo struct {
	Tag     string // git tag, e.g. "v1.2.3"
	Branch  string // git branch, e.g. "mybranch"
	Build   string // git or CI build number, e.g. "456"
	Version string // composite version, e.g. "v1.2.3-mybranch.456"
}

// Render returns either the Version string followed by a newline,
// or, if the pkgName is not an empty string, a small piece of
// Go code defining a global variable named "Version" with
// the contents of Version.
// If the pkgName is given but isn't a valid Go identifier,
// an error is returned.
func (vi *VersionInfo) Render(pkgName string) (string, error) {
	if pkgName == "" {
		return vi.Version + "\n", nil
	}
	if !token.IsIdentifier(pkgName) {
		return "", fmt.Errorf("'%s' is not a valid Go identifier", pkgName)
	}
	generatedBy := ""
	if executable, err := os.Executable(); err == nil {
		generatedBy = " by " + path.Base(executable)
	}
	return fmt.Sprintf(`// Code generated%s at %s DO NOT EDIT.
// branch %s, build %s
package %s

const Version = %s
`,
		generatedBy, time.Now().Format(time.ANSIC),
		strconv.Quote(vi.Branch), vi.Build,
		pkgName,
		strconv.QuoteToASCII(vi.Version)), nil
}