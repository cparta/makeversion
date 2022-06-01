package makeversion

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func Test_VersionInfo_Render(t *testing.T) {
	is := is.New(t)
	const VersionText = "v1.2.3-mybranch.456"
	vi := &VersionInfo{Version: VersionText}

	txt, err := vi.Render("")
	is.NoErr(err)
	is.Equal(VersionText+"\n", txt)

	txt, err = vi.Render("FooBar")
	is.NoErr(err)
	is.True(txt != "")
	is.True(strings.Contains(txt, "package foobar"))
	is.True(strings.Contains(txt, "const PkgName = \"FooBar\""))

	txt, err = vi.Render("123")
	is.True(err != nil)
	is.Equal(txt, "")
}
