package frequency_test

import (
	"testing"

	"github.com/ax-lab/tango/import/frequency"
	"github.com/stretchr/testify/require"
)

func TestParsePair(t *testing.T) {
	test := require.New(t)
	pair, err := frequency.ParsePair("21086758\tの")
	test.NoError(err)
	test.Equal(int64(21086758), pair.Count)
	test.Equal("の", pair.Entry)

	pair, err = frequency.ParsePair("\t1234\ttest\t")
	test.NoError(err)
	test.Equal(int64(1234), pair.Count)
	test.Equal("test", pair.Entry)

	checkEmpty := func(line string) {
		pair, err := frequency.ParsePair(line)
		test.NoError(err)
		test.Nil(pair)
	}

	checkError := func(line string) {
		pair, err := frequency.ParsePair(line)
		test.ErrorContains(err, "parsing pair frequency")
		test.Nil(pair)
	}

	checkEmpty("")
	checkEmpty("   \t   ")
	checkError("abc")
	checkError("123")
	checkError("123\t")
	checkError("a pair")
	checkError("123\tpair\textra")
}
