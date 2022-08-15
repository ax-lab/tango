package jmdict_test

import (
	"strings"
	"testing"

	"github.com/ax-lab/tango/import/jmdict"
	"github.com/stretchr/testify/require"
)

func TestEntryReadsKanji(t *testing.T) {
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
