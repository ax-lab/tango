package frequency_test

import (
	"testing"

	"github.com/ax-lab/tango/import/frequency"
	"github.com/stretchr/testify/require"
)

func TestParseWord(t *testing.T) {
	test := require.New(t)
	word, err := frequency.ParseWord("21086758\tの")
	test.NoError(err)
	test.Equal(int64(21086758), word.Count)
	test.Equal("の", word.Entry)

	word, err = frequency.ParseWord("\t1234\ttest\t")
	test.NoError(err)
	test.Equal(int64(1234), word.Count)
	test.Equal("test", word.Entry)

	checkEmpty := func(line string) {
		word, err := frequency.ParseWord(line)
		test.NoError(err)
		test.Nil(word)
	}

	checkError := func(line string) {
		word, err := frequency.ParseWord(line)
		test.ErrorContains(err, "parsing word frequency")
		test.Nil(word)
	}

	checkEmpty("")
	checkEmpty("   \t   ")
	checkError("abc")
	checkError("123")
	checkError("123\t")
	checkError("a word")
	checkError("123\tword\textra")
}
