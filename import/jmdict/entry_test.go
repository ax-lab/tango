package jmdict_test

import (
	"strings"
	"testing"

	"github.com/ax-lab/tango/import/jmdict"
	"github.com/stretchr/testify/require"
)

func TestEntryParsesKanji(t *testing.T) {
	test := require.New(t)
	entry := parseEntry(test, `
		<entry>
			<ent_seq>123</ent_seq>
			<k_ele>
				<keb>kanji 1</keb>
				<ke_inf>info1</ke_inf>
				<ke_pri>news1</ke_pri>
			</k_ele>
			<k_ele>
				<keb>kanji 2</keb>
				<ke_inf>info2a</ke_inf>
				<ke_inf>info2b</ke_inf>
				<ke_pri>news2a</ke_pri>
				<ke_pri>news2b</ke_pri>
			</k_ele>
		</entry>`)

	expected := &jmdict.Entry{
		Sequence: 123,
		Kanji: []jmdict.EntryKanji{
			{
				Text:     "kanji 1",
				Info:     []string{"info1"},
				Priority: []string{"news1"},
			},
			{
				Text:     "kanji 2",
				Info:     []string{"info2a", "info2b"},
				Priority: []string{"news2a", "news2b"},
			},
		},
	}
	test.Equal(expected, entry)
}

func TestEntryParsesReading(t *testing.T) {
	test := require.New(t)
	entry := parseEntry(test, `
		<entry>
			<ent_seq>123</ent_seq>
			<r_ele>
				<reb>reading 1</reb>
			</r_ele>
			<r_ele>
				<reb>reading 2</reb>
				<re_restr>kanji 2</re_restr>
				<re_inf>info2</re_inf>
				<re_pri>news2</re_pri>
			</r_ele>
			<r_ele>
				<reb>reading 3</reb>
				<re_nokanji/>
				<re_restr>kanji 3a</re_restr>
				<re_restr>kanji 3b</re_restr>
				<re_inf>info3a</re_inf>
				<re_inf>info3b</re_inf>
				<re_pri>news3a</re_pri>
				<re_pri>news3b</re_pri>
			</r_ele>
		</entry>`)

	expected := &jmdict.Entry{
		Sequence: 123,
		Reading: []jmdict.EntryReading{
			{
				Text: "reading 1",
			},
			{
				Text:        "reading 2",
				Info:        []string{"info2"},
				Priority:    []string{"news2"},
				Restriction: []string{"kanji 2"},
			},
			{
				Text:        "reading 3",
				NoKanji:     true,
				Info:        []string{"info3a", "info3b"},
				Priority:    []string{"news3a", "news3b"},
				Restriction: []string{"kanji 3a", "kanji 3b"},
			},
		},
	}
	test.Equal(expected, entry)
}

func TestEntryParsesSense(t *testing.T) {
	test := require.New(t)
	entry := parseEntry(test, `
		<entry>
			<ent_seq>123</ent_seq>
			<sense>
				<gloss>sense A</gloss>
			</sense>
			<sense>
				<gloss>sense B1</gloss>
				<gloss>sense B2</gloss>
				<stagk>stagk B1</stagk>
				<stagk>stagk B2</stagk>
				<stagr>stagr B1</stagr>
				<stagr>stagr B2</stagr>
				<pos>pos B1</pos>
				<pos>pos B2</pos>
				<xref>xref B1</xref>
				<xref>xref B2</xref>
				<ant>ant B1</ant>
				<ant>ant B2</ant>
				<field>field B1</field>
				<field>field B2</field>
				<misc>misc B1</misc>
				<misc>misc B2</misc>
				<s_inf>s_inf B1</s_inf>
				<s_inf>s_inf B2</s_inf>
				<dial>dial B1</dial>
				<dial>dial B2</dial>
			</sense>
		</entry>`)

	expected := &jmdict.Entry{
		Sequence: 123,
		Sense: []jmdict.EntrySense{
			{
				Glossary: []jmdict.EntrySenseGlossary{
					{Text: "sense A"},
				},
			},
			{
				Glossary: []jmdict.EntrySenseGlossary{
					{Text: "sense B1"},
					{Text: "sense B2"},
				},
				StagKanji:    []string{"stagk B1", "stagk B2"},
				StagReading:  []string{"stagr B1", "stagr B2"},
				PartOfSpeech: []string{"pos B1", "pos B2"},
				XRef:         []string{"xref B1", "xref B2"},
				Antonym:      []string{"ant B1", "ant B2"},
				Field:        []string{"field B1", "field B2"},
				Misc:         []string{"misc B1", "misc B2"},
				Info:         []string{"s_inf B1", "s_inf B2"},
				Dialect:      []string{"dial B1", "dial B2"},
			},
		},
	}
	test.Equal(expected, entry)
}

func TestEntryParsesSenseGlossaryAttributes(t *testing.T) {
	test := require.New(t)
	entry := parseEntry(test, `
		<entry>
			<ent_seq>123</ent_seq>
			<sense>
				<gloss xml:lang="por">portuguese</gloss>
				<gloss g_type="fig">figurative</gloss>
			</sense>
		</entry>`)

	expected := &jmdict.Entry{
		Sequence: 123,
		Sense: []jmdict.EntrySense{
			{
				Glossary: []jmdict.EntrySenseGlossary{
					{Text: "portuguese", Lang: "por"},
					{Text: "figurative", Type: "fig"},
				},
			},
		},
	}
	test.Equal(expected, entry)
}

func TestEntryParsesSenseSource(t *testing.T) {
	test := require.New(t)
	entry := parseEntry(test, `
		<entry>
			<ent_seq>123</ent_seq>
			<sense>
				<lsource>source A</lsource>
				<lsource xml:lang="ger">source B</lsource>
				<lsource ls_type="part">source C</lsource>
				<lsource ls_wasei="y">source D</lsource>
			</sense>
		</entry>`)

	expected := &jmdict.Entry{
		Sequence: 123,
		Sense: []jmdict.EntrySense{
			{
				Source: []jmdict.EntrySenseSource{
					{Text: "source A"},
					{Text: "source B", Lang: "ger"},
					{Text: "source C", Type: "part"},
					{Text: "source D", Wasei: "y"},
				},
			},
		},
	}
	test.Equal(expected, entry)
}

func parseEntry(test *require.Assertions, entryXML string) *jmdict.Entry {
	entries := parseEntries(test, entryXML)
	return entries[0]
}

func parseEntries(test *require.Assertions, entriesXML ...string) (out []*jmdict.Entry) {
	xml := []string{
		`<?xml version="1.0" encoding="UTF-8"?>`,
		`<JMdict>`,
	}
	xml = append(xml, entriesXML...)
	xml = append(xml, `</JMdict>`)

	input := openXML(strings.Join(xml, "\n"))
	for {
		entry, err := input.ReadEntry()
		test.NoError(err)
		if entry != nil {
			out = append(out, entry)
		} else {
			break
		}
	}
	return out
}
