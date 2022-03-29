package makeversion

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OsEnvironment_Getenv(t *testing.T) {
	const VarName = "MKENV_TEST3141592654"
	env := OsEnvironment{}
	_, expectOk := os.LookupEnv(VarName)
	_, actualOk := env.LookupEnv(VarName)
	assert.Equal(t, expectOk, actualOk)
	expect := os.Getenv(VarName)
	actual := env.Getenv(VarName)
	assert.Equal(t, expect, actual)
}
