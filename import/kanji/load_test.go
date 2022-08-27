package kanji_test

import (
	"bufio"
	"testing"

	"github.com/ax-lab/tango/import/kanji"
	"github.com/stretchr/testify/require"
)

func TestLoadsImportArchive(t *testing.T) {
	input, err := kanji.Load(kanji.DefaultFilePath)

	test := require.New(t)
	test.NoError(err)
	test.NotNil(input)

	reader := bufio.NewReader(input)
	line, err := reader.ReadString('\n')
	test.NoError(err)
	test.Contains(line, "<?xml")
}

func TestLoadReturnsErrorIfNotFound(t *testing.T) {
	input, err := kanji.Load("vendor/data/entries/Kanji-invalid-name.gz")

	test := require.New(t)
	test.Error(err)
	test.Nil(input)
}
