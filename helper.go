package aud

import (
	"regexp"

	"github.com/mattn/go-runewidth"
)

const (
	patternMustBreak = `[。！？.!?] ?$`
	patternCanBreak  = `[。！？，：；]|[.!?,:;] |——|—`
)

var (
	mustBreak = regexp.MustCompile(patternMustBreak)
	canBreak  = regexp.MustCompile(patternCanBreak)
)

// CJK no space between words
func wordSpace(lang string) string {
	switch lang {
	case "ja":
		fallthrough
	case "ko":
		fallthrough
	case "zh":
		return ""
	default:
		return " "
	}
}

func breakLineByPunctuation(text string, limit int) []string {
	ret := []string{}
	width := runewidth.StringWidth(text)
	if width <= limit {
		return []string{text}
	}
	indexs := canBreak.FindAllStringIndex(text, -1)
	for i := len(indexs) - 1; i >= 0; i-- {
		s := text[:indexs[i][1]]
		if runewidth.StringWidth(s) <= limit {
			if runewidth.StringWidth(s) < 10 {
				if i+1 < len(indexs)-1 {
					ret = append(ret, text[:indexs[i+1][1]])
					ret = append(ret, breakLineByPunctuation(text[indexs[i+1][1]:], limit)...)
					return ret
				}
				break
			}
			ret = append(ret, s)
			ret = append(ret, breakLineByPunctuation(text[indexs[i][1]:], limit)...)
			return ret
		}
	}
	return []string{text}
}
