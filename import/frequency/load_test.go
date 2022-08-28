package frequency_test

import (
	"testing"

	"github.com/ax-lab/tango/import/frequency"
	"github.com/stretchr/testify/require"
)

func TestLoadFrequencyEntries(t *testing.T) {
	test := require.New(t)

	checkError := func(file, msg string) {
		word, character, err := frequency.LoadFrequencyEntries(file)
		test.Empty(word)
		test.Empty(character)
		test.ErrorContains(err, msg)
	}

	word, character, err := frequency.LoadFrequencyEntries("testdata/valid-frequency.zip")
	test.NoError(err)
	test.NotEmpty(word)
	test.NotEmpty(character)

	checkError("testdata/some-non-existing-file.zip", "not found")
	checkError("testdata/invalid-frequency.txt", "zip file")
	checkError("testdata/invalid-frequency-empty.zip", "loading")
	checkError("testdata/invalid-frequency-missing-char.zip", "loading char entries")
	checkError("testdata/invalid-frequency-missing-word.zip", "loading word entries")
	checkError("testdata/invalid-frequency-invalid-char.zip", "loading char entries: parsing")
	checkError("testdata/invalid-frequency-invalid-word.zip", "loading word entries: parsing")

	word, character, err = frequency.LoadFrequencyEntries(frequency.DefaultFrequencyEntries)
	test.NoError(err)
	test.NotEmpty(word)
	test.NotEmpty(character)
	test.Greater(len(word), 50000)
	test.Greater(len(character), 5000)
}
