package common

import (
	"strings"

	"golang.org/x/text/unicode/norm"
)

func KanaToCommonHiragana(input string) string {
	input = string(norm.NFC.Bytes(([]byte)(input)))

	output := strings.Builder{}
	for _, chr := range input {
		output.WriteString(mapToHiragana(chr))
	}
	return output.String()
}

func mapToHiragana(chr rune) string {
	switch chr {
	case 'ン':
		return "ん"
	case 'ア':
		return "あ"
	case 'イ':
		return "い"
	case 'ウ':
		return "う"
	case 'エ':
		return "え"
	case 'オ':
		return "お"
	case 'カ':
		return "か"
	case 'キ':
		return "き"
	case 'ク':
		return "く"
	case 'ケ':
		return "け"
	case 'コ':
		return "こ"
	case 'ガ':
		return "が"
	case 'ギ':
		return "ぎ"
	case 'グ':
		return "ぐ"
	case 'ゲ':
		return "げ"
	case 'ゴ':
		return "ご"
	case 'サ':
		return "さ"
	case 'シ':
		return "し"
	case 'ス':
		return "す"
	case 'セ':
		return "せ"
	case 'ソ':
		return "そ"
	case 'ザ':
		return "ざ"
	case 'ジ':
		return "じ"
	case 'ズ':
		return "ず"
	case 'ゼ':
		return "ぜ"
	case 'ゾ':
		return "ぞ"
	case 'タ':
		return "た"
	case 'チ':
		return "ち"
	case 'ツ':
		return "つ"
	case 'テ':
		return "て"
	case 'ト':
		return "と"
	case 'ダ':
		return "だ"
	case 'ヂ':
		return "ぢ"
	case 'ヅ':
		return "づ"
	case 'デ':
		return "で"
	case 'ド':
		return "ど"
	case 'ナ':
		return "な"
	case 'ニ':
		return "に"
	case 'ヌ':
		return "ぬ"
	case 'ネ':
		return "ね"
	case 'ノ':
		return "の"
	case 'ハ':
		return "は"
	case 'ヒ':
		return "ひ"
	case 'フ':
		return "ふ"
	case 'ヘ':
		return "へ"
	case 'ホ':
		return "ほ"
	case 'バ':
		return "ば"
	case 'ビ':
		return "び"
	case 'ブ':
		return "ぶ"
	case 'ベ':
		return "べ"
	case 'ボ':
		return "ぼ"
	case 'パ':
		return "ぱ"
	case 'ピ':
		return "ぴ"
	case 'プ':
		return "ぷ"
	case 'ペ':
		return "ぺ"
	case 'ポ':
		return "ぽ"
	case 'マ':
		return "ま"
	case 'ミ':
		return "み"
	case 'ム':
		return "む"
	case 'メ':
		return "め"
	case 'モ':
		return "も"
	case 'ヤ':
		return "や"
	case 'ユ':
		return "ゆ"
	case 'ヨ':
		return "よ"
	case 'ラ':
		return "ら"
	case 'リ':
		return "り"
	case 'ル':
		return "る"
	case 'レ':
		return "れ"
	case 'ロ':
		return "ろ"
	case 'ワ':
		return "わ"
	case 'ヰ':
		return "い"
	case 'ヱ':
		return "え"
	case 'ヲ':
		return "を"
	case 'ヴ':
		return "ゔ"

	// rare chars
	case '〼':
		return "ます"

	// rare hiragana
	case 'ゐ':
		return "い"
	case 'ゑ':
		return "え"
	case 'ゟ':
		return "より"
	case '𛀁':
		return "え"
	case 'ゕ':
		return "か"
	case 'ゖ':
		return "け"
	case 'ゎ':
		return "わ"

	// small chars
	case 'ァ':
		return "ぁ"
	case 'ィ':
		return "ぃ"
	case 'ゥ':
		return "ぅ"
	case 'ェ':
		return "ぇ"
	case 'ォ':
		return "ぉ"
	case 'ッ':
		return "っ"
	case 'ャ':
		return "ゃ"
	case 'ュ':
		return "ゅ"
	case 'ョ':
		return "ょ"

	// rare katakana
	case 'ヵ':
		return "か"
	case 'ヶ':
		return "か"
	case 'ヮ':
		return "わ"
	case 'ヷ':
		return "ゔぁ"
	case 'ヸ':
		return "ゔぃ"
	case 'ヹ':
		return "ゔぇ"
	case 'ヺ':
		return "ゔぉ"
	case '𛀀':
		return "え"
	case 'ヿ':
		return "こと"

	// small extensions
	case 'ㇰ':
		return "く"
	case 'ㇱ':
		return "し"
	case 'ㇲ':
		return "す"
	case 'ㇳ':
		return "と"
	case 'ㇴ':
		return "ぬ"
	case 'ㇵ':
		return "は"
	case 'ㇶ':
		return "ひ"
	case 'ㇷ':
		return "ふ"
	case 'ㇸ':
		return "へ"
	case 'ㇹ':
		return "ほ"
	case 'ㇺ':
		return "む"
	case 'ㇻ':
		return "ら"
	case 'ㇼ':
		return "り"
	case 'ㇽ':
		return "る"
	case 'ㇾ':
		return "れ"
	case 'ㇿ':
		return "ろ"

	default:
		return string(chr)
	}
}
