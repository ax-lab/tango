package data

func LoadNames(importDir string) ([]Data, error) {
	db := OpenDB(importDir, "names.db")
	names := db.LoadTable("SELECT * FROM name")

	bySequence := make(map[int64]Data)
	for _, row := range names {
		row["senses"] = make([]Data, 0)
		row["kanji"] = splitTabs(row["kanji"])
		row["reading"] = splitTabs(row["reading"])
		bySequence[row["sequence"].(int64)] = row
	}

	senses := db.LoadTable("SELECT * FROM name_sense")
	for _, row := range senses {
		name := bySequence[row["sequence"].(int64)]
		delete(row, "sequence")
		delete(row, "order")
		row["info"] = splitTabs(row["info"])
		row["xref"] = splitTabs(row["xref"])
		row["translation"] = splitTabs(row["translation"])
		name["senses"] = append(name["senses"].([]Data), row)
	}

	err := db.Done()
	return names, err
}
