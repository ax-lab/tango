package jmdict_test

import (
	"bufio"
	"testing"

	"github.com/ax-lab/tango/import/jmdict"
	"github.com/stretchr/testify/require"
)

func TestLoadsImportArchive(t *testing.T) {
	input, err := jmdict.Load(jmdict.DefaultFilePath)

	test := require.New(t)
	test.NoError(err)
	test.NotNil(input)

	reader := bufio.NewReader(input)
	line, err := reader.ReadString('\n')
	test.NoError(err)
	test.Contains(line, "<?xml")
}

func TestLoadReturnsErrorIfNotFound(t *testing.T) {
	input, err := jmdict.Load("vendor/data/entries/JMDict-invalid-name.gz")

	test := require.New(t)
	test.Error(err)
	test.Nil(input)
}
