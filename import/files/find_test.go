package files_test

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/ax-lab/tango/import/files"
	"github.com/stretchr/testify/require"
)

func TestFindExistingFile(t *testing.T) {
	input, err := files.Find("vendor/data/README.md")

	test := require.New(t)
	test.NoError(err)
	test.NotNil(input)
}

func TestReturnsErrorIfNotFound(t *testing.T) {
	input, err := files.Find("vendor/data/this-does-not-exist.md")

	test := require.New(t)
	test.Nil(input)
	test.ErrorContains(err, "this-does-not-exist.md")
	test.ErrorContains(err, "not found")
	test.True(errors.Is(err, fs.ErrNotExist))
}
