package frequency_test

import (
	"testing"

	"github.com/ax-lab/tango/import/frequency"
	"github.com/stretchr/testify/require"
)

func TestLoadEntries(t *testing.T) {
	test := require.New(t)

	checkError := func(file, msg string) {
		word, character, err := frequency.LoadEntries(file)
		test.Empty(word)
		test.Empty(character)
		test.ErrorContains(err, msg)
	}

	word, character, err := frequency.LoadEntries("testdata/frequency/valid-entries.zip")
	test.NoError(err)
	test.NotEmpty(word)
	test.NotEmpty(character)

	checkError("testdata/frequency/some-non-existing-file.zip", "not found")
	checkError("testdata/frequency/invalid.txt", "zip file")
	checkError("testdata/frequency/invalid-entries-empty.zip", "loading")
	checkError("testdata/frequency/invalid-entries-missing-char.zip", "loading char entries")
	checkError("testdata/frequency/invalid-entries-missing-word.zip", "loading word entries")
	checkError("testdata/frequency/invalid-entries-invalid-char.zip", "loading char entries: parsing")
	checkError("testdata/frequency/invalid-entries-invalid-word.zip", "loading word entries: parsing")

	word, character, err = frequency.LoadEntries(frequency.DefaultEntriesFile)
	test.NoError(err)
	test.NotEmpty(word)
	test.NotEmpty(character)
	test.Greater(len(word), 50000)
	test.Greater(len(character), 5000)
}

func TestLoadWords(t *testing.T) {
	test := require.New(t)

	checkError := func(file, msg string) {
		jparser, mecab, kanji, err := frequency.LoadWords(file)
		test.Empty(jparser)
		test.Empty(mecab)
		test.Empty(kanji)
		test.ErrorContains(err, msg)
	}

	jparser, mecab, kanji, err := frequency.LoadWords("testdata/frequency/valid-novel.zip")
	test.NoError(err)
	test.NotEmpty(jparser)
	test.NotEmpty(mecab)
	test.NotEmpty(kanji)

	checkError("testdata/frequency/some-non-existing-file.zip", "not found")
	checkError("testdata/frequency/invalid.txt", "zip file")
	checkError("testdata/frequency/invalid-novel-empty.zip", "loading")
	checkError("testdata/frequency/invalid-novel-missing-jparser.zip", "loading jparser entries")
	checkError("testdata/frequency/invalid-novel-missing-mecab.zip", "loading mecab entries")
	checkError("testdata/frequency/invalid-novel-missing-kanji.zip", "loading kanji entries")
	checkError("testdata/frequency/invalid-novel-invalid-jparser.zip", "loading jparser entries: parsing")
	checkError("testdata/frequency/invalid-novel-invalid-mecab.zip", "loading mecab entries: parsing")
	checkError("testdata/frequency/invalid-novel-invalid-kanji.zip", "loading kanji entries: parsing")

	jparser, mecab, kanji, err = frequency.LoadWords(frequency.DefaultWordsFile)
	test.NoError(err)
	test.NotEmpty(jparser)
	test.NotEmpty(mecab)
	test.NotEmpty(kanji)
	test.Greater(len(jparser), 50000)
	test.Greater(len(mecab), 50000)
	test.Greater(len(kanji), 5000)
}
