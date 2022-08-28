package frequency_test

import (
	"testing"

	"github.com/ax-lab/tango/import/frequency"
	"github.com/stretchr/testify/require"
)

func TestLoadInfo(t *testing.T) {
	test := require.New(t)

	checkError := func(file, msg string) {
		word, character, err := frequency.LoadInfo(file)
		test.Empty(word)
		test.Empty(character)
		test.ErrorContains(err, msg)
	}

	word, character, err := frequency.LoadInfo("testdata/frequency/valid-freq.zip")
	test.NoError(err)
	test.NotEmpty(word)
	test.NotEmpty(character)

	checkError("testdata/frequency/some-non-existing-file.zip", "not found")
	checkError("testdata/frequency/invalid.txt", "zip file")
	checkError("testdata/frequency/invalid-freq-empty.zip", "loading")
	checkError("testdata/frequency/invalid-freq-missing-char.zip", "loading char info")
	checkError("testdata/frequency/invalid-freq-missing-word.zip", "loading word info")
	checkError("testdata/frequency/invalid-freq-invalid-char.zip", "loading char info: parsing")
	checkError("testdata/frequency/invalid-freq-invalid-word.zip", "loading word info: parsing")

	word, character, err = frequency.LoadInfo(frequency.DefaultInfoFile)
	test.NoError(err)
	test.NotEmpty(word)
	test.NotEmpty(character)
	test.Greater(len(word), 50000)
	test.Greater(len(character), 5000)
}

func TestLoadPairs(t *testing.T) {
	test := require.New(t)

	checkError := func(file, msg string) {
		jparser, mecab, kanji, err := frequency.LoadPairs(file)
		test.Empty(jparser)
		test.Empty(mecab)
		test.Empty(kanji)
		test.ErrorContains(err, msg)
	}

	jparser, mecab, kanji, err := frequency.LoadPairs("testdata/frequency/valid-novel.zip")
	test.NoError(err)
	test.NotEmpty(jparser)
	test.NotEmpty(mecab)
	test.NotEmpty(kanji)

	checkError("testdata/frequency/some-non-existing-file.zip", "not found")
	checkError("testdata/frequency/invalid.txt", "zip file")
	checkError("testdata/frequency/invalid-novel-empty.zip", "loading")
	checkError("testdata/frequency/invalid-novel-missing-jparser.zip", "loading jparser pairs")
	checkError("testdata/frequency/invalid-novel-missing-mecab.zip", "loading mecab pairs")
	checkError("testdata/frequency/invalid-novel-missing-kanji.zip", "loading kanji pairs")
	checkError("testdata/frequency/invalid-novel-invalid-jparser.zip", "loading jparser pairs: parsing")
	checkError("testdata/frequency/invalid-novel-invalid-mecab.zip", "loading mecab pairs: parsing")
	checkError("testdata/frequency/invalid-novel-invalid-kanji.zip", "loading kanji pairs: parsing")

	jparser, mecab, kanji, err = frequency.LoadPairs(frequency.DefaultPairsFile)
	test.NoError(err)
	test.NotEmpty(jparser)
	test.NotEmpty(mecab)
	test.NotEmpty(kanji)
	test.Greater(len(jparser), 50000)
	test.Greater(len(mecab), 50000)
	test.Greater(len(kanji), 5000)
}
