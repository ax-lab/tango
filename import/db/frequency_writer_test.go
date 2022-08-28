package db_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/ax-lab/tango/import/db"
	"github.com/ax-lab/tango/import/frequency"
	"github.com/stretchr/testify/require"
)

func TestFrequencyWriter(t *testing.T) {
	dt := func(x int) frequency.InfoData {
		return frequency.InfoData{
			Freq:   int64(x * 100),
			CD:     int64(x*100 + 50),
			FreqPM: fmt.Sprintf("%.1f", float64(x)+0.5),
			CDPc:   fmt.Sprintf("%.1f", float64(x*10)+0.5),
		}
	}

	testFrequency(t,
		func(test *require.Assertions, db *db.FrequencyWriter) {
			db.AddCharInfo([]*frequency.Info{
				{Entry: "a", Blog: dt(1), Twitter: dt(2), News: dt(3)},
				{Entry: "b", Blog: dt(4), Twitter: dt(5), News: dt(6)},
				{Entry: "!", Blog: dt(7), Twitter: dt(8), News: dt(9)},
			})
			db.AddCharPairs([]*frequency.Pair{
				{Entry: "x", Count: 10},
				{Entry: "y", Count: 12},
				{Entry: "!", Count: 98},
			})
			db.AddWordInfo([]*frequency.Info{
				{Entry: "wa", Blog: dt(1), Twitter: dt(2), News: dt(3)},
				{Entry: "wb", Blog: dt(4), Twitter: dt(5), News: dt(6)},
				{Entry: "XX", Blog: dt(7), Twitter: dt(8), News: dt(9)},
			})
			db.AddWordPairs(
				[]*frequency.Pair{
					{Entry: "wx", Count: 10},
					{Entry: "wy", Count: 12},
					{Entry: "XX", Count: 98},
				},
				[]*frequency.Pair{
					{Entry: "wx", Count: 11},
					{Entry: "wy", Count: 13},
					{Entry: "XX", Count: 99},
				},
			)
			err := db.Write()
			test.NoError(err)
		},
		func(test *require.Assertions, db *sql.DB) {
			expected := func(word bool, a, b, c, d, e string) []Data {
				out := []Data{
					{
						"entry": a, "count": int64(98),
						"blog_freq": int64(700), "blog_pm": "7.5", "blog_cd": int64(750), "blog_cdp": "70.5",
						"twit_freq": int64(800), "twit_pm": "8.5", "twit_cd": int64(850), "twit_cdp": "80.5",
						"news_freq": int64(900), "news_pm": "9.5", "news_cd": int64(950), "news_cdp": "90.5",
					},
					{
						"entry": b, "count": int64(12),
						"blog_freq": nil, "blog_pm": nil, "blog_cd": nil, "blog_cdp": nil,
						"twit_freq": nil, "twit_pm": nil, "twit_cd": nil, "twit_cdp": nil,
						"news_freq": nil, "news_pm": nil, "news_cd": nil, "news_cdp": nil,
					},
					{
						"entry": c, "count": int64(10),
						"blog_freq": nil, "blog_pm": nil, "blog_cd": nil, "blog_cdp": nil,
						"twit_freq": nil, "twit_pm": nil, "twit_cd": nil, "twit_cdp": nil,
						"news_freq": nil, "news_pm": nil, "news_cd": nil, "news_cdp": nil,
					},
					{
						"entry": d, "count": nil,
						"blog_freq": int64(400), "blog_pm": "4.5", "blog_cd": int64(450), "blog_cdp": "40.5",
						"twit_freq": int64(500), "twit_pm": "5.5", "twit_cd": int64(550), "twit_cdp": "50.5",
						"news_freq": int64(600), "news_pm": "6.5", "news_cd": int64(650), "news_cdp": "60.5",
					},
					{
						"entry": e, "count": nil,
						"blog_freq": int64(100), "blog_pm": "1.5", "blog_cd": int64(150), "blog_cdp": "10.5",
						"twit_freq": int64(200), "twit_pm": "2.5", "twit_cd": int64(250), "twit_cdp": "20.5",
						"news_freq": int64(300), "news_pm": "3.5", "news_cd": int64(350), "news_cdp": "30.5",
					},
				}
				if word {
					for _, it := range out {
						it["count_m"] = nil
					}
					out[0]["count_m"] = int64(99)
					out[1]["count_m"] = int64(13)
					out[2]["count_m"] = int64(11)
				}
				return out
			}

			expectedChars := expected(false, "!", "y", "x", "b", "a")
			actualChars := testQuery(test, db, "SELECT * FROM freq_char")
			test.EqualValues(expectedChars, actualChars)

			expectedWords := expected(true, "XX", "wy", "wx", "wb", "wa")
			actualWords := testQuery(test, db, "SELECT * FROM freq_word")
			test.EqualValues(expectedWords, actualWords)
		})
}

func testFrequency(t *testing.T, prepare func(test *require.Assertions, db *db.FrequencyWriter), eval func(test *require.Assertions, db *sql.DB)) {
	testTempDB(t, func(test *require.Assertions, dbFile string) {
		func() {
			db, dbErr := db.NewFrequencyWriter(dbFile)
			if dbErr != nil {
				panic(dbErr)
			}

			defer db.Close()
			prepare(test, db)
		}()

		func() {
			db, err := sql.Open("sqlite3", dbFile)
			if err != nil {
				panic(err)
			}
			defer db.Close()
			eval(test, db)
		}()
	})
}
