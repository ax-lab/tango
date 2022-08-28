package frequency_test

import (
	"testing"

	"github.com/ax-lab/tango/import/frequency"
	"github.com/stretchr/testify/require"
)

func TestParseInfo(t *testing.T) {
	test := require.New(t)

	check := func(line string, expected frequency.Info) {
		info, err := frequency.ParseInfo(line)
		test.NoError(err)
		test.Equal(&expected, info)

		info, err = frequency.ParseInfo(line + "\t")
		test.NoError(err)
		test.Equal(&expected, info)
	}

	checkEmpty := func(line string) {
		info, err := frequency.ParseInfo(line)
		test.NoError(err)
		test.Nil(info)
	}

	checkError := func(line string) {
		info, err := frequency.ParseInfo(line)
		test.ErrorContains(err, "parsing frequency info")
		test.Nil(info)
	}

	check(
		"の\t783900\t52752.36\t346504\t52.1601\t588537\t49415.37\t315749\t47.3302\t942403\t66742.42\t271217\t86.6741",
		frequency.Info{
			Entry: "の",
			Blog: frequency.InfoData{
				Freq:   783900,
				FreqPM: "52752.36",
				CD:     346504,
				CDPc:   "52.1601",
			},
			Twitter: frequency.InfoData{
				Freq:   588537,
				FreqPM: "49415.37",
				CD:     315749,
				CDPc:   "47.3302",
			},
			News: frequency.InfoData{
				Freq:   942403,
				FreqPM: "66742.42",
				CD:     271217,
				CDPc:   "86.6741",
			},
		},
	)

	check(
		"い\t862956\t33525.87\t363144\t54.6649\t943164\t52368.91\t431207\t64.6372\t562980\t21719.91\t212494\t67.9077",
		frequency.Info{
			Entry: "い",
			Blog: frequency.InfoData{
				Freq:   862956,
				FreqPM: "33525.87",
				CD:     363144,
				CDPc:   "54.6649",
			},
			Twitter: frequency.InfoData{
				Freq:   943164,
				FreqPM: "52368.91",
				CD:     431207,
				CDPc:   "64.6372",
			},
			News: frequency.InfoData{
				Freq:   562980,
				FreqPM: "21719.91",
				CD:     212494,
				CDPc:   "67.9077",
			},
		},
	)

	checkEmpty("")
	checkEmpty("   \t   ")

	checkError("abc")
	checkError("123")
	checkError("123\t")
	checkError("a word")
	checkError("123\tword\textra")

	// spaces and empty entries occur in the source file
	checkEmpty("\t0\t0\t0\t0\t0\t0\t0\t0\t0\t0\t0\t0")
	checkEmpty(" \t0\t0\t0\t0\t0\t0\t0\t0\t0\t0\t0\t0")

	checkError("い\t_\t0\t0\t0\t0\t0\t0\t0\t0\t0\t0\t0")
	checkError("い\t0\t_\t0\t0\t0\t0\t0\t0\t0\t0\t0\t0")
	checkError("い\t0\t0\t_\t0\t0\t0\t0\t0\t0\t0\t0\t0")
	checkError("い\t0\t0\t0\t_\t0\t0\t0\t0\t0\t0\t0\t0")
	checkError("い\t0\t0\t0\t0\t_\t0\t0\t0\t0\t0\t0\t0")
	checkError("い\t0\t0\t0\t0\t0\t_\t0\t0\t0\t0\t0\t0")
	checkError("い\t0\t0\t0\t0\t0\t0\t_\t0\t0\t0\t0\t0")
	checkError("い\t0\t0\t0\t0\t0\t0\t0\t_\t0\t0\t0\t0")
	checkError("い\t0\t0\t0\t0\t0\t0\t0\t0\t_\t0\t0\t0")
	checkError("い\t0\t0\t0\t0\t0\t0\t0\t0\t0\t_\t0\t0")
	checkError("い\t0\t0\t0\t0\t0\t0\t0\t0\t0\t0\t_\t0")
	checkError("い\t0\t0\t0\t0\t0\t0\t0\t0\t0\t0\t0\t_")
	checkError("い\t0\t0\t0\t0\t0\t0\t0\t0\t0\t0\t0\t_")
}
