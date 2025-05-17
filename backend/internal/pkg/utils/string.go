package utils

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/huydq/test/internal/pkg/utils/static"
)

var kanaReplacers = map[string][]*strings.Replacer{
	// k: 全角カタカナから半角カタカナ
	static.CONV_HALFWIDTHKATA_FROM_FULLWIDTHKATA: {strings.NewReplacer(static.HalfWidthKataFromFullWidthKata...)},
	// K: 半角カタカナから全角カタカナ
	static.CONV_FULLWIDTHKATA_FROM_HALFWIDTHKATA: {strings.NewReplacer(static.FullWidthKataFromHalfWidthKata...)},
	// c: 全角カタカナから全角ひらがな
	static.CONV_HALFWIDTHKATA_FROM_HALFWIDTHKANA: {strings.NewReplacer(static.FullWidthHiraFromFullWidthKata...)},
	// KV: 半角カタカナから全角カタカナ + 濁点を1文字
	static.CONV_FULLWIDTHKATA_FROM_HALFWIDTHKATA_AND_DAKUTEN: {strings.NewReplacer(static.FullWidthKataFromHalfWidthKata...), strings.NewReplacer(static.KanaDakuten...)},
	// h:  全角ひらがなから半角カタカナ
	static.CONV_HALFWIDTHKANA_FROM_FULLWIDTHHIRA: {strings.NewReplacer(static.HalfWidthKataFromFullWidthHira...)},
	// H: 半角カタカナから全角ひらがな
	static.CONV_FULLWIDTHHIRA_FROM_HALFWIDTHKANA: {strings.NewReplacer(static.FullWidthHiraFromHalfWidthKata...)},
	// HV: 半角カタカナから全角ひらがな + 濁点を1文字
	static.CONV_FULLWIDTHHIRA_FROM_HALFWIDTHKANA_AND_DAKUTEN: {strings.NewReplacer(static.FullWidthHiraFromHalfWidthKata...), strings.NewReplacer(static.KanaDakuten...)},
	// a: 全角英数字から半角英数字
	static.CONV_HALFWIDTHEISU_FROM_FULLWIDTHEISU: {strings.NewReplacer(static.HalfWidthAlphaNumFromFullWidthAlphaNum...)},
	// A: 半角英数字から全角英数字
	static.CONV_FULLWIDTHEISU_FROM_HALFWIDTHEISU: {strings.NewReplacer(static.FullWidthAlphaNumFromHalfWidthAlphaNum...)},
	// s: 全角スペースから半角スペース
	static.CONV_HALFWIDTHSPACE_FROM_FULLWIDTHSPACE: {strings.NewReplacer(static.HalfWidthSpaceFromFullWidthSpace...)},
	// S: 全角スペースから半角スペース
	static.CONV_FULLWIDTHSPACE_FROM_HALFWIDTHSPACE: {strings.NewReplacer(static.FullWidthSpaceFromHalfWidthSpace...)},
}

func ConvertKana(s string, mode []string) string {
	if s == "" {
		return s
	}

	for i := range mode {
		if _, ok := kanaReplacers[mode[i]]; ok {
			for j := range kanaReplacers[mode[i]] {
				s = kanaReplacers[mode[i]][j].Replace(s)
			}
		}
	}

	return s
}

// RuneWriter UTF8 →　ShiftJIS　変換で書き込みを失敗した場合、?に置き換えるWriter
type RuneWriter struct {
	W io.Writer
}

func (rw *RuneWriter) Write(b []byte) (int, error) {
	var err error
	l := 0

loop:
	for len(b) > 0 {
		_, n := utf8.DecodeRune(b)
		if n == 0 {
			break loop
		}
		_, err = rw.W.Write(b[:n])
		if err != nil {
			_, err = rw.W.Write([]byte{'?'})
			if err != nil {
				break loop
			}
		}
		l += n
		b = b[n:]
	}
	return l, err
}
