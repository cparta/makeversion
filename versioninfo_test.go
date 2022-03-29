package makeversion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_VersionInfo_Render(t *testing.T) {
	const VersionText = "v1.2.3-mybranch.456"
	vi := &VersionInfo{Version: VersionText}

	txt, err := vi.Render("")
	assert.NoError(t, err)
	assert.Equal(t, VersionText+"\n", txt)

	txt, err = vi.Render("foo")
	assert.NoError(t, err)
	assert.NotEmpty(t, txt)

	txt, err = vi.Render("123")
	assert.Error(t, err)
	assert.Empty(t, txt)
}
