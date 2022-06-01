package makeversion

import (
	"os"
	"testing"

	"github.com/matryer/is"
)

func Test_OsEnvironment_Getenv(t *testing.T) {
	is := is.New(t)
	const VarName = "MKENV_TEST3141592654"
	env := OsEnvironment{}
	_, expectOk := os.LookupEnv(VarName)
	_, actualOk := env.LookupEnv(VarName)
	is.Equal(expectOk, actualOk)
	expect := os.Getenv(VarName)
	actual := env.Getenv(VarName)
	is.Equal(expect, actual)
}
