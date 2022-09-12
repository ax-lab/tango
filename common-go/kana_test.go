package common_test

import (
	"testing"

	"github.com/ax-lab/tango/common"
	"github.com/stretchr/testify/require"
)

func TestKanaToCommonHiragana(t *testing.T) {
	test := require.New(t)

	check := func(input, expected string) {
		test.Equal(expected, common.KanaToCommonHiragana(input))
	}

	checkIdentity := func(input string) {
		check(input, input)
	}

	checkIdentity("")
	checkIdentity("あいうえおかがきぎくぐけげこごさざしじすずせぜそぞただちぢつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもやゆよらりるれろわをん")
	checkIdentity("あいうえお－ぁぃぅぇぉっゃゅょ")

	check(
		"アイウエオカガキギクグケゲコゴサザシジスズセゼソゾタダチヂツヅテデトドナニヌネノハバパヒビピフブプヘベペホボポマミムメモヤユヨラリルレロワヲン",
		"あいうえおかがきぎくぐけげこごさざしじすずせぜそぞただちぢつづてでとどなにぬねのはばぱひびぴふぶぷへべぺほぼぽまみむめもやゆよらりるれろわをん")

	check(
		"ァィゥェォッャュョ",
		"ぁぃぅぇぉっゃゅょ")

	check("ヮヰヱヵヶ", "わいえかか")
	check(
		"ㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇺㇻㇼㇽㇾㇿ",
		"くしすとぬはひふへほむらりるれろ")
	check("ヷヸヹヺヴ𛀀", "ゔぁゔぃゔぇゔぉゔえ")

	// ゑ・𛀁・ゐ -> e/e/i
	check("ゑゎゕゖ𛀁ゐ", "えわかけえい")
	check("ゟヿ〼", "よりことます")

	checkVoiced := func(input string, voiced, semiVoiced string) {
		check(input+"\u3099", voiced)
		if semiVoiced != "" {
			check(input+"\u309A", semiVoiced)
		}
	}

	checkVoiced("か", "が", "")
	checkVoiced("き", "ぎ", "")
	checkVoiced("く", "ぐ", "")
	checkVoiced("け", "げ", "")
	checkVoiced("こ", "ご", "")
	checkVoiced("さ", "ざ", "")
	checkVoiced("し", "じ", "")
	checkVoiced("す", "ず", "")
	checkVoiced("せ", "ぜ", "")
	checkVoiced("そ", "ぞ", "")
	checkVoiced("た", "だ", "")
	checkVoiced("ち", "ぢ", "")
	checkVoiced("つ", "づ", "")
	checkVoiced("て", "で", "")
	checkVoiced("と", "ど", "")
	checkVoiced("は", "ば", "ぱ")
	checkVoiced("ひ", "び", "ぴ")
	checkVoiced("ふ", "ぶ", "ぷ")
	checkVoiced("へ", "べ", "ぺ")
	checkVoiced("ほ", "ぼ", "ぽ")
	checkVoiced("う", "ゔ", "")

	checkVoiced("カ", "が", "")
	checkVoiced("キ", "ぎ", "")
	checkVoiced("ク", "ぐ", "")
	checkVoiced("ケ", "げ", "")
	checkVoiced("コ", "ご", "")
	checkVoiced("サ", "ざ", "")
	checkVoiced("シ", "じ", "")
	checkVoiced("ス", "ず", "")
	checkVoiced("セ", "ぜ", "")
	checkVoiced("ソ", "ぞ", "")
	checkVoiced("タ", "だ", "")
	checkVoiced("チ", "ぢ", "")
	checkVoiced("ツ", "づ", "")
	checkVoiced("テ", "で", "")
	checkVoiced("ト", "ど", "")
	checkVoiced("ハ", "ば", "ぱ")
	checkVoiced("ヒ", "び", "ぴ")
	checkVoiced("フ", "ぶ", "ぷ")
	checkVoiced("ヘ", "べ", "ぺ")
	checkVoiced("ホ", "ぼ", "ぽ")
	checkVoiced("ウ", "ゔ", "")

	checkVoiced("ワ", "ゔぁ", "")
	checkVoiced("ヰ", "ゔぃ", "")
	checkVoiced("ヱ", "ゔぇ", "")
	checkVoiced("ヲ", "ゔぉ", "")
}
