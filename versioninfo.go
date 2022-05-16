package makeversion

import (
	"fmt"
	"go/token"
	"os"
	"path"
	"strconv"
	"strings"
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
// Go code defining global variables named "PkgName" and "PkgVersion"
// with the given pkgName and the contents of Version.
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

const PkgName = %s
const PkgVersion = %s
`,
		generatedBy, time.Now().Format(time.ANSIC),
		strconv.Quote(vi.Branch), vi.Build,
		strings.ToLower(pkgName),
		strconv.Quote(pkgName),
		strconv.QuoteToASCII(vi.Version)), nil
}
