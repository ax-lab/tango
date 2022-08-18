package jmdict

type Entry struct {
	Sequence int `xml:"ent_seq"`

	Kanji   []EntryKanji   `xml:"k_ele"`
	Reading []EntryReading `xml:"r_ele"`
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
