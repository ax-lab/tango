package jmnedict_test

import (
	"strings"
	"testing"

	"github.com/ax-lab/tango/import/jmnedict"
	"github.com/stretchr/testify/require"
)

func TestDecoderReturnsErrorOnInvalidXML(t *testing.T) {
	test := require.New(t)

	checkError := func(expectedMessage string, xml string) {
		input := openXML(xml)
		entry, err := input.ReadEntry()
		test.Nil(entry)
		test.ErrorContains(err, expectedMessage)
	}

	checkError("XML syntax error", "<not valid")
	checkError("XML syntax error", `
		<?xml version="1.0" encoding="UTF-8"?>
		<JMnedict>
		<entry
		</JMnedict>
	`)
}

func TestDecoderReturnsErrorOnInvalidSchema(t *testing.T) {
	test := require.New(t)

	input := openXML(`
		<?xml version="1.0" encoding="UTF-8"?>
		<invalid>
		<entry><ent_seq>1000</ent_seq></entry>
		<entry><ent_seq>1001</ent_seq></entry>
		</invalid>
	`)
	entry, err := input.ReadEntry()

	test.Nil(entry)
	test.ErrorContains(err, "invalid schema")
}

func TestDecoderReturnsErrorOnInvalidEntry(t *testing.T) {
	test := require.New(t)

	checkError := func(expectedMessage string, xml string) {
		input := openXML(xml)
		entry, err := input.ReadEntry()
		test.Nil(entry)
		test.ErrorContains(err, expectedMessage)
	}

	checkError("invalid entry: missing sequence", `
		<?xml version="1.0" encoding="UTF-8"?>
		<JMnedict>
		<entry></entry>
		</JMnedict>
	`)

	checkError("decoding entry", `
		<?xml version="1.0" encoding="UTF-8"?>
		<JMnedict>
		<entry><ent_seq>abc</ent_seq></entry>
		</JMnedict>
	`)

	checkError("decoding entry", `
		<?xml version="1.0" encoding="UTF-8"?>
		<JMnedict>
		<entry>
		</JMnedict>
	`)
}

func TestDecoderIgnoresUnknownEntries(t *testing.T) {
	test := require.New(t)

	input := openXML(`
		<?xml version="1.0" encoding="UTF-8"?>
		<JMnedict>
		<some>123</some>
		<some>456</some>
		</JMnedict>
	`)

	entry, err := input.ReadEntry()
	test.NoError(err)
	test.Nil(entry)
}

func TestDecoderReadsEntries(t *testing.T) {
	test := require.New(t)

	input := openXML(`
		<?xml version="1.0" encoding="UTF-8"?>
		<!-- some comment -->
		<JMnedict>
		<entry><ent_seq>1000</ent_seq></entry>
		<entry><ent_seq>1001</ent_seq></entry>
		</JMnedict>
	`)

	check := func(expectedSequence int) {
		entry, err := input.ReadEntry()
		test.NoError(err)
		test.NotNil(entry)
		test.Equal(expectedSequence, entry.Sequence)
	}

	checkEOF := func() {
		entry, err := input.ReadEntry()
		test.NoError(err)
		test.Nil(entry)
	}

	check(1000)
	check(1001)
	checkEOF()
}

func TestDecoderParsesCustomEntities(t *testing.T) {
	test := require.New(t)

	input := openXML(`
		<?xml version="1.0"?>
		<!DOCTYPE JMnedict [
		<!-- some comment -->
		<!ELEMENT JMnedict (entry*)>
		<!-- -->
		<!ELEMENT entry (ent_seq, k_ele*, r_ele+, trans+)*>
			<!-- some element -->
		<!ELEMENT ent_seq (#PCDATA)>
			<!-- another element -->
		<!ATTLIST trans_det xml:lang CDATA #IMPLIED>
			<!-- The xml:lang attribute defines the target language of the
			translated name. -->
		<!-- <name_type> entities -->
		<!ENTITY char "character">
		<!ENTITY company "company name">
		<!ENTITY doc "document">
		<!ENTITY ev "event">
		<!ENTITY fem "female given name or forename">
		<!ENTITY surname "family or surname">

		<!-- external entity, not included -->
		<!ENTITY entityname [PUBLIC "public-identifier"] SYSTEM "system-identifier">

		<!-- this is valid, even if highly not recommended -->
		<!ENTITY weird ">x<">

		<!-- invalid entities are not included -->
		<!ENTITY test1>
		<!ENTITY "test2">

		]>
		<!-- JMnedict created: 2022-08-21 -->
		<JMnedict>
		<entry>
		<ent_seq>5000000</ent_seq>
		</entry>
		</JMnedict>
	`)

	entry, err := input.ReadEntry()
	test.NoError(err)
	test.NotNil(entry)
	test.Equal(5000000, entry.Sequence)

	tags := map[string]string{
		"char":    "character",
		"company": "company name",
		"doc":     "document",
		"ev":      "event",
		"fem":     "female given name or forename",
		"surname": "family or surname",
		"weird":   ">x<",
	}
	test.Equal(tags, input.Tags)
}

func TestDecoderParsesLoadedFile(t *testing.T) {
	test := require.New(t)

	input, err := jmnedict.Load(jmnedict.DefaultFilePath)
	test.NoError(err)

	decoder := jmnedict.NewDecoder(input)
	entry, err := decoder.ReadEntry()
	test.NoError(err)
	test.NotNil(entry)
	test.True(entry.Sequence >= 5000000)
	test.NotEmpty(decoder.Tags)
}

func TestDecoderParsesFullEntryData(t *testing.T) {
	const entriesToLoad = 20000

	test := require.New(t)

	input, err := jmnedict.Load(jmnedict.DefaultFilePath)
	test.NoError(err)

	var entries []*jmnedict.Name
	decoder := jmnedict.NewDecoder(input)
	for i := 0; i < entriesToLoad; i++ {
		entry, err := decoder.ReadEntry()
		test.NoError(err)
		test.NotNil(entry)
		entries = append(entries, entry)
	}

	checkAny := func(name string, check func(entry *jmnedict.Name) bool) {
		for _, entry := range entries {
			if check(entry) {
				return
			}
		}
		test.FailNow("check any failed", "assertion `%s` did not match any of the %d entries", name, len(entries))
	}

	checkAny("parses kanji", func(entry *jmnedict.Name) bool {
		return len(entry.Kanji) > 0 && entry.Kanji[0] != ""
	})

	checkAny("parses reading", func(entry *jmnedict.Name) bool {
		return len(entry.Reading) > 0 && entry.Reading[0] != ""
	})

	checkAny("parses sense", func(entry *jmnedict.Name) bool {
		return len(entry.Sense) > 0 && len(entry.Sense[0].Translation) > 0 &&
			entry.Sense[0].Translation[0] != ""
	})

	checkAny("parses sense type", func(entry *jmnedict.Name) bool {
		for _, it := range entry.Sense {
			if len(it.Type) > 0 && it.Type[0] != "" {
				return true
			}
		}
		return false
	})

	checkAny("parses sense type as tag", func(entry *jmnedict.Name) bool {
		for _, it := range entry.Sense {
			if it.Type[0] != "fem" {
				return true
			}
		}
		return false
	})

	checkAny("parses sense xref", func(entry *jmnedict.Name) bool {
		for _, it := range entry.Sense {
			if len(it.XRef) > 0 && it.XRef[0] != "" {
				return true
			}
		}
		return false
	})
}

func openXML(input string) *jmnedict.Decoder {
	reader := strings.NewReader(input)
	return jmnedict.NewDecoder(reader)
}
