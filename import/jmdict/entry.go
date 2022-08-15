package jmdict

type Entry struct {
	Sequence int `xml:"ent_seq"`

	Kanji []EntryKanji `xml:"k_ele"`
}

type EntryKanji struct {
	Text     string   `xml:"keb"`
	Info     []string `xml:"ke_inf"`
	Priority []string `xml:"ke_pri"`
}
