package jmnedict

type Name struct {
	Sequence int `xml:"ent_seq"`

	Kanji   []string    `xml:"k_ele>keb"`
	Reading []string    `xml:"r_ele>reb"`
	Sense   []NameSense `xml:"trans"`
}

type NameSense struct {
	Type        []string `xml:"name_type"`
	XRef        []string `xml:"xref"`
	Translation []string `xml:"trans_det"`
}
