package kanji

type Character struct {
	Literal   string `xml:"literal"`
	Grade     int    `xml:"misc>grade"`
	Strokes   []int  `xml:"misc>stroke_count"`
	Frequency int    `xml:"misc>freq"`
	JLPT      int    `xml:"misc>jlpt"`

	ReadingMeanings []CharacterReadingMeaningGroup `xml:"reading_meaning>rmgroup"`

	RadicalName []string `xml:"misc>rad_name"`

	Codepoint []CharacterCodepoint `xml:"codepoint>cp_value"`
	Radical   []CharacterRadical   `xml:"radical>rad_value"`
	Variant   []CharacterVariant   `xml:"misc>variant"`
	Reference []CharacterReference `xml:"dic_number>dic_ref"`
	QueryCode []CharacterQueryCode `xml:"query_code>q_code"`

	Nanori []string `xml:"reading_meaning>nanori"`
}

type CharacterCodepoint struct {
	Type string `xml:"cp_type,attr"`
	Text string `xml:",chardata"`
}

type CharacterRadical struct {
	Type string `xml:"rad_type,attr"`
	Text string `xml:",chardata"`
}

type CharacterVariant struct {
	Type string `xml:"var_type,attr"`
	Text string `xml:",chardata"`
}

type CharacterReference struct {
	Type   string `xml:"dr_type,attr"`
	Text   string `xml:",chardata"`
	Volume string `xml:"m_vol,attr"`
	Page   string `xml:"m_page,attr"`
}

type CharacterQueryCode struct {
	Type         string `xml:"qc_type,attr"`
	Text         string `xml:",chardata"`
	SkipMisclass string `xml:"skip_misclass,attr"`
}

type CharacterReadingMeaningGroup struct {
	Reading []CharacterReading `xml:"reading"`
	Meaning []CharacterMeaning `xml:"meaning"`
}

type CharacterReading struct {
	Type string `xml:"r_type,attr"`
	Text string `xml:",chardata"`
}

type CharacterMeaning struct {
	Lang string `xml:"m_lang,attr"`
	Text string `xml:",chardata"`
}
