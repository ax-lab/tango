package kanji_test

import (
	"strings"
	"testing"

	"github.com/ax-lab/tango/import/kanji"
	"github.com/stretchr/testify/require"
)

func TestCharacterParsesLiteral(t *testing.T) {
	test := require.New(t)
	character := parseName(test, `
		<character>
			<literal>日</literal>
		</character>`)

	expected := &kanji.Character{
		Literal: "日",
	}
	test.Equal(expected, character)
}

func TestCharacterParsesCodepoint(t *testing.T) {
	test := require.New(t)
	character := parseName(test, `
		<character>
			<literal>X</literal>
			<codepoint>
				<cp_value cp_type="ucs">1234</cp_value>
				<cp_value cp_type="jis208">1-16-01</cp_value>
			</codepoint>
		</character>`)

	expected := &kanji.Character{
		Literal: "X",
		Codepoint: []kanji.CharacterCodepoint{
			{Type: "ucs", Text: "1234"},
			{Type: "jis208", Text: "1-16-01"},
		},
	}
	test.Equal(expected, character)
}

func TestCharacterParsesRadical(t *testing.T) {
	test := require.New(t)
	character := parseName(test, `
		<character>
			<literal>X</literal>
			<radical>
				<rad_value rad_type="classical">7</rad_value>
				<rad_value rad_type="nelson">1</rad_value>
			</radical>
		</character>`)

	expected := &kanji.Character{
		Literal: "X",
		Radical: []kanji.CharacterRadical{
			{Type: "classical", Text: "7"},
			{Type: "nelson", Text: "1"},
		},
	}
	test.Equal(expected, character)
}

func TestCharacterParsesMiscInfo(t *testing.T) {
	test := require.New(t)
	character := parseName(test, `
		<character>
			<literal>X</literal>
			<misc>
				<grade>8</grade>
				<stroke_count>10</stroke_count>
				<stroke_count>11</stroke_count>
				<freq>1234</freq>
				<jlpt>5</jlpt>
				<variant var_type="jis208">1-16-01</variant>
				<variant var_type="ucs">ABCD</variant>
				<rad_name>のぎ</rad_name>
				<rad_name>のぎへん</rad_name>
			</misc>
		</character>`)

	expected := &kanji.Character{
		Literal:     "X",
		Grade:       8,
		Strokes:     []int{10, 11},
		Frequency:   1234,
		JLPT:        5,
		RadicalName: []string{"のぎ", "のぎへん"},
		Variant: []kanji.CharacterVariant{
			{Type: "jis208", Text: "1-16-01"},
			{Type: "ucs", Text: "ABCD"},
		},
	}
	test.Equal(expected, character)
}

func TestCharacterParsesReference(t *testing.T) {
	test := require.New(t)
	character := parseName(test, `
		<character>
			<literal>X</literal>
			<dic_number>
				<dic_ref dr_type="nelson_n">4121</dic_ref>
				<dic_ref dr_type="oneill_names">220</dic_ref>
				<dic_ref dr_type="moro" m_vol="8" m_page="0522">24906</dic_ref>
			</dic_number>
		</character>`)

	expected := &kanji.Character{
		Literal: "X",
		Reference: []kanji.CharacterReference{
			{Type: "nelson_n", Text: "4121"},
			{Type: "oneill_names", Text: "220"},
			{Type: "moro", Text: "24906", Volume: "8", Page: "0522"},
		},
	}
	test.Equal(expected, character)
}

func TestCharacterParsesQueryCode(t *testing.T) {
	test := require.New(t)
	character := parseName(test, `
		<character>
			<literal>X</literal>
			<query_code>
				<q_code qc_type="skip">4-5-3</q_code>
				<q_code qc_type="four_corner">2090.4</q_code>
				<q_code qc_type="skip" skip_misclass="stroke_count">2-1-4</q_code>
			</query_code>
		</character>`)

	expected := &kanji.Character{
		Literal: "X",
		QueryCode: []kanji.CharacterQueryCode{
			{Type: "skip", Text: "4-5-3"},
			{Type: "four_corner", Text: "2090.4"},
			{Type: "skip", Text: "2-1-4", SkipMisclass: "stroke_count"},
		},
	}
	test.Equal(expected, character)
}

func TestCharacterParsesReadingMeaning(t *testing.T) {
	test := require.New(t)
	character := parseName(test, `
		<character>
			<literal>X</literal>
			<reading_meaning>
				<rmgroup>
					<reading r_type="ja_on">カ</reading>
					<reading r_type="ja_kun">かせ.ぐ</reading>
					<meaning>earnings</meaning>
					<meaning>work</meaning>
					<meaning>earn money</meaning>
					<meaning m_lang="fr">gains</meaning>
				</rmgroup>
				<rmgroup>
					<reading r_type="ja_on">カ</reading>
					<reading r_type="ja_kun">コ</reading>
					<meaning>counter for articles</meaning>
				</rmgroup>
				<nanori>な</nanori>
				<nanori>なす</nanori>
			</reading_meaning>
		</character>`)

	expected := &kanji.Character{
		Literal: "X",
		ReadingMeanings: []kanji.CharacterReadingMeaningGroup{
			{
				Reading: []kanji.CharacterReading{
					{Type: "ja_on", Text: "カ"},
					{Type: "ja_kun", Text: "かせ.ぐ"},
				},
				Meaning: []kanji.CharacterMeaning{
					{Text: "earnings", Lang: ""},
					{Text: "work", Lang: ""},
					{Text: "earn money", Lang: ""},
					{Text: "gains", Lang: "fr"},
				},
			},
			{
				Reading: []kanji.CharacterReading{
					{Type: "ja_on", Text: "カ"},
					{Type: "ja_kun", Text: "コ"},
				},
				Meaning: []kanji.CharacterMeaning{
					{Text: "counter for articles", Lang: ""},
				},
			},
		},
		Nanori: []string{"な", "なす"},
	}
	test.Equal(expected, character)
}

func parseName(test *require.Assertions, characterXML string) *kanji.Character {
	entries := parseCharacters(test, characterXML)
	return entries[0]
}

func parseCharacters(test *require.Assertions, entriesXML ...string) (out []*kanji.Character) {
	xml := []string{
		`<?xml version="1.0" encoding="UTF-8"?>`,
		`<kanjidic2>`,
	}
	xml = append(xml, entriesXML...)
	xml = append(xml, `</kanjidic2>`)

	input := openXML(strings.Join(xml, "\n"))
	for {
		character, err := input.ReadCharacter()
		test.NoError(err)
		if character != nil {
			out = append(out, character)
		} else {
			break
		}
	}
	return out
}
