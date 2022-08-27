package kanji_test

import (
	"strings"
	"testing"

	"github.com/ax-lab/tango/import/kanji"
	"github.com/stretchr/testify/require"
)

func TestDecoderReturnsErrorOnInvalidXML(t *testing.T) {
	test := require.New(t)

	checkError := func(expectedMessage string, xml string) {
		input := openXML(xml)
		entry, err := input.ReadCharacter()
		test.Nil(entry)
		test.ErrorContains(err, expectedMessage)
	}

	checkError("XML syntax error", "<not valid")
	checkError("XML syntax error", `
		<?xml version="1.0" encoding="UTF-8"?>
		<kanjidic2>
		<header
		</kanjidic2>
	`)
}

func TestDecoderReturnsErrorOnInvalidSchema(t *testing.T) {
	test := require.New(t)

	input := openXML(`
		<?xml version="1.0" encoding="UTF-8"?>
		<invalid>
		<character><literal>日</literal></character>
		<character><literal>本</literal></character>
		</invalid>
	`)
	entry, err := input.ReadCharacter()

	test.Nil(entry)
	test.ErrorContains(err, "invalid schema")
}

func TestDecoderReturnsErrorOnInvalidEntry(t *testing.T) {
	test := require.New(t)

	checkError := func(expectedMessage string, xml string) {
		input := openXML(xml)
		entry, err := input.ReadCharacter()
		test.Nil(entry)
		test.ErrorContains(err, expectedMessage)
	}

	checkError("invalid character: missing literal", `
		<?xml version="1.0" encoding="UTF-8"?>
		<kanjidic2>
		<character></character>
		</kanjidic2>
	`)

	checkError("invalid character: missing literal", `
		<?xml version="1.0" encoding="UTF-8"?>
		<kanjidic2>
		<character><literal></literal></character>
		</kanjidic2>
	`)

	checkError("decoding character", `
		<?xml version="1.0" encoding="UTF-8"?>
		<kanjidic2>
		<character>
		</kanjidic2>
	`)

	checkError("decoding header", `
		<?xml version="1.0" encoding="UTF-8"?>
		<kanjidic2>
		<header>
		</kanjidic2>
	`)
}

func TestDecoderIgnoresUnknownEntries(t *testing.T) {
	test := require.New(t)

	input := openXML(`
		<?xml version="1.0" encoding="UTF-8"?>
		<kanjidic2>
		<some>123</some>
		<some>456</some>
		</kanjidic2>
	`)

	entry, err := input.ReadCharacter()
	test.NoError(err)
	test.Nil(entry)
}

func TestDecoderReadsHeader(t *testing.T) {
	test := require.New(t)

	input := openXML(`
		<?xml version="1.0" encoding="UTF-8"?>
		<!-- some comment -->
		<kanjidic2>
		<header>
		<database_version>2022-001</database_version>
		<date_of_creation>2022-01-01</date_of_creation>
		</header>
		</kanjidic2>
	`)

	character, err := input.ReadCharacter()
	test.NoError(err)
	test.Nil(character)

	test.Equal("2022-001", input.Info.Version)
	test.Equal("2022-01-01", input.Info.Created)
}

func TestDecoderReadsCharacters(t *testing.T) {
	test := require.New(t)

	input := openXML(`
		<?xml version="1.0" encoding="UTF-8"?>
		<!-- some comment -->
		<kanjidic2>
		<character><literal>日</literal></character>
		<character><literal>本</literal></character>
		</kanjidic2>
	`)

	check := func(expected string) {
		character, err := input.ReadCharacter()
		test.NoError(err)
		test.NotNil(character)
		test.Equal(expected, character.Literal)
	}

	checkEOF := func() {
		character, err := input.ReadCharacter()
		test.NoError(err)
		test.Nil(character)
	}

	check("日")
	check("本")
	checkEOF()
}

func TestDecoderParsesLoadedFile(t *testing.T) {
	test := require.New(t)

	input, err := kanji.Load(kanji.DefaultFilePath)
	test.NoError(err)

	decoder := kanji.NewDecoder(input)
	entry, err := decoder.ReadCharacter()
	test.NoError(err)
	test.NotNil(entry)
	test.True(entry.Literal != "")
	test.NotEmpty(decoder.Info.Version)
	test.NotEmpty(decoder.Info.Created)
}

func TestDecoderParsesFullEntryData(t *testing.T) {
	const entriesToLoad = 250

	test := require.New(t)

	input, err := kanji.Load(kanji.DefaultFilePath)
	test.NoError(err)

	var entries []*kanji.Character
	decoder := kanji.NewDecoder(input)
	for i := 0; i < entriesToLoad; i++ {
		character, err := decoder.ReadCharacter()
		test.NoError(err)
		test.NotNil(character)
		entries = append(entries, character)
	}

	checkAny := func(name string, check func(character *kanji.Character) bool) {
		for _, character := range entries {
			if check(character) {
				return
			}
		}
		test.FailNow("check any failed", "assertion `%s` did not match any of the %d entries", name, len(entries))
	}

	checkAny("parses literal", func(character *kanji.Character) bool {
		return character.Literal != ""
	})

	checkAny("parses codepoint", func(character *kanji.Character) bool {
		if len(character.Codepoint) > 0 {
			return character.Codepoint[0].Type != "" && character.Codepoint[0].Text != ""
		}
		return false
	})

	checkAny("parses radical", func(character *kanji.Character) bool {
		if len(character.Radical) > 0 {
			return character.Radical[0].Type != "" && character.Radical[0].Text != ""
		}
		return false
	})

	checkAny("parses misc grade", func(character *kanji.Character) bool {
		return character.Grade > 0
	})

	checkAny("parses misc strokes", func(character *kanji.Character) bool {
		return len(character.Strokes) > 1 && character.Strokes[0] > 0 && character.Strokes[1] > 0
	})

	checkAny("parses misc frequency", func(character *kanji.Character) bool {
		return character.Frequency > 0
	})

	checkAny("parses misc jlpt", func(character *kanji.Character) bool {
		return character.JLPT > 0
	})

	checkAny("parses misc radical name", func(character *kanji.Character) bool {
		return len(character.RadicalName) > 0 && character.RadicalName[0] != ""
	})

	checkAny("parses misc variant", func(character *kanji.Character) bool {
		if len(character.Variant) > 0 {
			return character.Variant[0].Type != "" && character.Variant[0].Text != ""
		}
		return false
	})

	checkAny("parses reference", func(character *kanji.Character) bool {
		if len(character.Reference) > 0 {
			hasPage := false
			for _, it := range character.Reference {
				if it.Type == "" || it.Text == "" {
					return false
				}
				if it.Page != "" && it.Volume != "" {
					hasPage = true
				}
			}
			return hasPage
		}
		return false
	})

	checkAny("parses query code", func(character *kanji.Character) bool {
		if len(character.QueryCode) > 0 {
			hasSkipMisclass := false
			for _, it := range character.QueryCode {
				if it.Type == "" || it.Text == "" {
					return false
				}
				if it.SkipMisclass != "" {
					hasSkipMisclass = true
				}
			}
			return hasSkipMisclass
		}
		return false
	})

	checkAny("parses reading meaning", func(character *kanji.Character) bool {
		if len(character.ReadingMeanings) > 0 {
			for _, it := range character.ReadingMeanings {
				if len(it.Meaning) == 0 || it.Meaning[0].Text == "" {
					return false
				}
				if len(it.Reading) == 0 || it.Reading[0].Text == "" {
					return false
				}
			}
			return true
		}
		return false
	})

	checkAny("parses reading meaning lang", func(character *kanji.Character) bool {
		if len(character.ReadingMeanings) > 0 {
			for _, it := range character.ReadingMeanings {
				for _, m := range it.Meaning {
					if m.Lang != "" {
						return true
					}
				}
			}
		}
		return false
	})

	checkAny("parses nanori", func(character *kanji.Character) bool {
		return len(character.Nanori) > 0 && character.Nanori[0] != ""
	})
}

func openXML(input string) *kanji.Decoder {
	reader := strings.NewReader(input)
	return kanji.NewDecoder(reader)
}
