package jmnedict_test

import (
	"strings"
	"testing"

	"github.com/ax-lab/tango/import/jmnedict"
	"github.com/stretchr/testify/require"
)

func TestNameParsesKanji(t *testing.T) {
	test := require.New(t)
	entry := parseName(test, `
		<entry>
			<ent_seq>123</ent_seq>
			<k_ele>
				<keb>kanji 1</keb>
			</k_ele>
			<k_ele>
				<keb>kanji 2</keb>
			</k_ele>
		</entry>`)

	expected := &jmnedict.Name{
		Sequence: 123,
		Kanji:    []string{"kanji 1", "kanji 2"},
	}
	test.Equal(expected, entry)
}

func TestNameParsesReading(t *testing.T) {
	test := require.New(t)
	entry := parseName(test, `
		<entry>
			<ent_seq>123</ent_seq>
			<r_ele>
				<reb>reading 1</reb>
			</r_ele>
			<r_ele>
				<reb>reading 2</reb>
			</r_ele>
		</entry>`)

	expected := &jmnedict.Name{
		Sequence: 123,
		Reading:  []string{"reading 1", "reading 2"},
	}
	test.Equal(expected, entry)
}

func TestNameParsesSense(t *testing.T) {
	test := require.New(t)
	entry := parseName(test, `
		<entry>
			<ent_seq>123</ent_seq>
			<trans>
				<trans_det>translation</trans_det>
			</trans>
			<trans>
				<name_type>info 1</name_type>
				<name_type>info 2</name_type>
				<xref>xref 1</xref>
				<xref>xref 2</xref>
				<trans_det>translation 1</trans_det>
				<trans_det>translation 2</trans_det>
			</trans>
		</entry>`)

	expected := &jmnedict.Name{
		Sequence: 123,
		Sense: []jmnedict.NameSense{
			{
				Translation: []string{"translation"},
			},
			{
				Info:        []string{"info 1", "info 2"},
				XRef:        []string{"xref 1", "xref 2"},
				Translation: []string{"translation 1", "translation 2"},
			},
		},
	}
	test.Equal(expected, entry)
}

func parseName(test *require.Assertions, entryXML string) *jmnedict.Name {
	entries := parseEntries(test, entryXML)
	return entries[0]
}

func parseEntries(test *require.Assertions, entriesXML ...string) (out []*jmnedict.Name) {
	xml := []string{
		`<?xml version="1.0" encoding="UTF-8"?>`,
		`<JMnedict>`,
	}
	xml = append(xml, entriesXML...)
	xml = append(xml, `</JMnedict>`)

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
