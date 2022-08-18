package jmdict_test

import (
	"strings"
	"testing"

	"github.com/ax-lab/tango/import/jmdict"
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
		<JMdict>
		<entry
		</JMdict>
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
		<JMdict>
		<entry></entry>
		</JMdict>
	`)

	checkError("decoding entry", `
		<?xml version="1.0" encoding="UTF-8"?>
		<JMdict>
		<entry><ent_seq>abc</ent_seq></entry>
		</JMdict>
	`)

	checkError("decoding entry", `
		<?xml version="1.0" encoding="UTF-8"?>
		<JMdict>
		<entry>
		</JMdict>
	`)
}

func TestDecoderIgnoresUnknownEntries(t *testing.T) {
	test := require.New(t)

	input := openXML(`
		<?xml version="1.0" encoding="UTF-8"?>
		<JMdict>
		<some>123</some>
		<some>456</some>
		</JMdict>
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
		<JMdict>
		<entry><ent_seq>1000</ent_seq></entry>
		<entry><ent_seq>1001</ent_seq></entry>
		</JMdict>
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
		<?xml version="1.0" encoding="UTF-8"?>
		<!DOCTYPE JMdict [
			<!ELEMENT JMdict (entry*)>
			<!--  -->
			<!ELEMENT entry (ent_seq, k_ele*, r_ele+, sense+)>
			<!--
			some comment
			-->
			<!ELEMENT ent_seq (#PCDATA)>
			<!ELEMENT k_ele (keb, ke_inf*, ke_pri*)>
			<!-- entities -->
			<!-- <dial> (dialect) entities -->
			<!ENTITY bra "Brazilian">
			<!ENTITY hob "Hokkaido-ben">
			<!ENTITY ksb "Kansai-ben">
			<!-- <field> entities -->
			<!ENTITY agric "agriculture">
			<!ENTITY anat "anatomy">
			<!-- test a separate quote -->
			<!ENTITY uK 'word usually written using kanji alone'>

			<!-- external entity, not included -->
			<!ENTITY entityname [PUBLIC "public-identifier"] SYSTEM "system-identifier">

			<!-- this is valid, even if highly not recommended -->
			<!ENTITY weird ">x<">

			<!-- invalid entities are not included -->
			<!ENTITY test1>
			<!ENTITY "test2">
		]>

		<JMdict>
			<custom>&bra;</custom>
			<custom>&uK;</custom>
		</JMdict>
	`)

	entry, err := input.ReadEntry()
	test.NoError(err)
	test.Nil(entry)

	tags := map[string]string{
		"bra":   "Brazilian",
		"hob":   "Hokkaido-ben",
		"ksb":   "Kansai-ben",
		"agric": "agriculture",
		"anat":  "anatomy",
		"uK":    "word usually written using kanji alone",
		"weird": ">x<",
	}
	test.Equal(tags, input.Tags)
}

func TestDecoderParsesLoadedFile(t *testing.T) {
	test := require.New(t)

	input, err := jmdict.Load(jmdict.DefaultFilePath)
	test.NoError(err)

	decoder := jmdict.NewDecoder(input)
	entry, err := decoder.ReadEntry()
	test.NoError(err)
	test.NotNil(entry)
	test.True(entry.Sequence >= 1000000)
	test.NotEmpty(decoder.Tags)
}

func TestDecoderParsesFullEntryData(t *testing.T) {
	const entriesToLoad = 5000

	test := require.New(t)

	input, err := jmdict.Load(jmdict.DefaultFilePath)
	test.NoError(err)

	var entries []*jmdict.Entry
	decoder := jmdict.NewDecoder(input)
	for i := 0; i < entriesToLoad; i++ {
		entry, err := decoder.ReadEntry()
		test.NoError(err)
		test.NotNil(entry)
		entries = append(entries, entry)
	}

	checkAny := func(name string, check func(entry *jmdict.Entry) bool) {
		for _, entry := range entries {
			if check(entry) {
				return
			}
		}
		test.FailNow("check any failed", "assertion `%s` did not match any of the %d entries", name, len(entries))
	}

	checkAny("parses kanji", func(entry *jmdict.Entry) bool {
		if len(entry.Kanji) == 0 {
			return false
		}
		for _, it := range entry.Kanji {
			test.NotEmpty(it.Text)
		}
		return true
	})

	checkAny("parses kanji info", func(entry *jmdict.Entry) bool {
		for _, it := range entry.Kanji {
			if len(it.Info) > 0 && it.Info[0] != "" {
				return true
			}
		}
		return false
	})

	checkAny("parses kanji priority", func(entry *jmdict.Entry) bool {
		for _, it := range entry.Kanji {
			if len(it.Priority) > 0 && it.Priority[0] != "" {
				return true
			}
		}
		return false
	})

	checkAny("parses reading", func(entry *jmdict.Entry) bool {
		return len(entry.Reading) > 0 && entry.Reading[0].Text != ""
	})

	checkAny("parses reading info", func(entry *jmdict.Entry) bool {
		for _, it := range entry.Reading {
			if len(it.Info) > 0 && it.Info[0] != "" {
				return true
			}
		}
		return false
	})

	checkAny("parses reading priority", func(entry *jmdict.Entry) bool {
		for _, it := range entry.Reading {
			if len(it.Priority) > 0 && it.Priority[0] != "" {
				return true
			}
		}
		return false
	})

	checkAny("parses reading restriction", func(entry *jmdict.Entry) bool {
		for _, it := range entry.Reading {
			if len(it.Restriction) > 0 && it.Restriction[0] != "" {
				return true
			}
		}
		return false
	})

	checkAny("parses reading no-kanji", func(entry *jmdict.Entry) bool {
		for _, it := range entry.Reading {
			if it.NoKanji {
				return true
			}
		}
		return false
	})

	checkAny("parses sense", func(entry *jmdict.Entry) bool {
		return len(entry.Sense) > 0 && len(entry.Sense[0].Glossary) > 0 &&
			entry.Sense[0].Glossary[0].Text != ""
	})

	checkAny("parses sense info", func(entry *jmdict.Entry) bool {
		for _, it := range entry.Sense {
			if len(it.Info) > 0 && it.Info[0] != "" {
				return true
			}
		}
		return false
	})

	checkAny("parses sense pos", func(entry *jmdict.Entry) bool {
		for _, it := range entry.Sense {
			if len(it.PartOfSpeech) > 0 && it.PartOfSpeech[0] != "" {
				return true
			}
		}
		return false
	})

	checkAny("parses sense stag kanji", func(entry *jmdict.Entry) bool {
		for _, it := range entry.Sense {
			if len(it.StagKanji) > 0 && it.StagKanji[0] != "" {
				return true
			}
		}
		return false
	})

	checkAny("parses sense stag reading", func(entry *jmdict.Entry) bool {
		for _, it := range entry.Sense {
			if len(it.StagReading) > 0 && it.StagReading[0] != "" {
				return true
			}
		}
		return false
	})

	checkAny("parses sense field", func(entry *jmdict.Entry) bool {
		for _, it := range entry.Sense {
			if len(it.Field) > 0 && it.Field[0] != "" {
				return true
			}
		}
		return false
	})

	checkAny("parses sense misc", func(entry *jmdict.Entry) bool {
		for _, it := range entry.Sense {
			if len(it.Misc) > 0 && it.Misc[0] != "" {
				return true
			}
		}
		return false
	})

	checkAny("parses sense dialect", func(entry *jmdict.Entry) bool {
		for _, it := range entry.Sense {
			if len(it.Dialect) > 0 && it.Dialect[0] != "" {
				return true
			}
		}
		return false
	})

	checkAny("parses sense antonym", func(entry *jmdict.Entry) bool {
		for _, it := range entry.Sense {
			if len(it.Antonym) > 0 && it.Antonym[0] != "" {
				return true
			}
		}
		return false
	})

	checkAny("parses sense xref", func(entry *jmdict.Entry) bool {
		for _, it := range entry.Sense {
			if len(it.XRef) > 0 && it.XRef[0] != "" {
				return true
			}
		}
		return false
	})

	checkAny("parses sense glossary lang", func(entry *jmdict.Entry) bool {
		for _, sense := range entry.Sense {
			for _, it := range sense.Glossary {
				if it.Lang != "" {
					return true
				}
			}
		}
		return false
	})

	checkAny("parses sense glossary type", func(entry *jmdict.Entry) bool {
		for _, sense := range entry.Sense {
			for _, it := range sense.Glossary {
				if it.Type != "" {
					return true
				}
			}
		}
		return false
	})

	checkAny("parses sense source", func(entry *jmdict.Entry) bool {
		for _, sense := range entry.Sense {
			if len(sense.Source) > 0 && sense.Source[0].Text != "" {
				return true
			}
		}
		return false
	})

	checkAny("parses sense source lang", func(entry *jmdict.Entry) bool {
		for _, sense := range entry.Sense {
			for _, it := range sense.Source {
				if it.Lang != "" {
					return true
				}
			}
		}
		return false
	})

	checkAny("parses sense source type", func(entry *jmdict.Entry) bool {
		for _, sense := range entry.Sense {
			for _, it := range sense.Source {
				if it.Type == "part" {
					return true
				}
			}
		}
		return false
	})

	checkAny("parses sense source wasei", func(entry *jmdict.Entry) bool {
		for _, sense := range entry.Sense {
			for _, it := range sense.Source {
				if it.Wasei == "y" {
					return true
				}
			}
		}
		return false
	})
}

func openXML(input string) *jmdict.Decoder {
	reader := strings.NewReader(input)
	return jmdict.NewDecoder(reader)
}
