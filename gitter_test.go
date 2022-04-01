package makeversion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewDefaultGitter_SuccedsNormally(t *testing.T) {
	dg, err := NewDefaultGitter("git")
	assert.NoError(t, err)
	assert.NotNil(t, dg)
}

func Test_DefaultGitter_GetTag(t *testing.T) {
	dg, err := NewDefaultGitter("git")
	if assert.NoError(t, err) && assert.NotNil(t, dg) {
		assert.NotEmpty(t, dg.GetTag("."))
		assert.Empty(t, dg.GetTag("/"))
	}
}

func Test_DefaultGitter_GetBranch(t *testing.T) {
	dg, err := NewDefaultGitter("git")
	if assert.NoError(t, err) && assert.NotNil(t, dg) {
		assert.NotEmpty(t, dg.GetBranch("."))
		assert.Equal(t, "HEAD", dg.GetBranch("/"))
	}
}

func Test_DefaultGitter_GetBuild(t *testing.T) {
	dg, err := NewDefaultGitter("git")
	if assert.NoError(t, err) && assert.NotNil(t, dg) {
		assert.NotEmpty(t, dg.GetBuild("."))
		assert.Empty(t, dg.GetBuild("/"))
	}
}
