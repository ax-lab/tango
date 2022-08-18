package jmdict

type Entry struct {
	Sequence int `xml:"ent_seq"`

	Kanji   []EntryKanji   `xml:"k_ele"`
	Reading []EntryReading `xml:"r_ele"`
	Sense   []EntrySense   `xml:"sense"`
}

type EntryKanji struct {
	Text     string   `xml:"keb"`
	Info     []string `xml:"ke_inf"`
	Priority []string `xml:"ke_pri"`
}

type EntryReading struct {
	Text        string   `xml:"reb"`
	Info        []string `xml:"re_inf"`
	Priority    []string `xml:"re_pri"`
	Restriction []string `xml:"re_restr"`

	NoKanji    bool
	NoKanjiRaw *string `xml:"re_nokanji"`
}

type EntrySense struct {
	Glossary     []EntrySenseGlossary `xml:"gloss"`
	Info         []string             `xml:"s_inf"`
	PartOfSpeech []string             `xml:"pos"`
	StagKanji    []string             `xml:"stagk"`
	StagReading  []string             `xml:"stagr"`
	Field        []string             `xml:"field"`
	Misc         []string             `xml:"misc"`
	Dialect      []string             `xml:"dial"`
	Antonym      []string             `xml:"ant"`
	XRef         []string             `xml:"xref"`
	Source       []EntrySenseSource   `xml:"lsource"`
}

type EntrySenseGlossary struct {
	Text string `xml:",chardata"`
	Lang string `xml:"lang,attr"`
	Type string `xml:"g_type,attr"`
}

type EntrySenseSource struct {
	Text  string `xml:",chardata"`
	Lang  string `xml:"lang,attr"`
	Type  string `xml:"ls_type,attr"`
	Wasei string `xml:"ls_wasei,attr"`
}
