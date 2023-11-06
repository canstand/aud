package aud

import (
	"strings"
	"time"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astisub"
	"github.com/mattn/go-runewidth"
)

type (
	SingleSegment struct {
		Start        float64           `json:"start"`
		End          float64           `json:"end"`
		Text         string            `json:"text"`
		Words        []SingleWord      `json:"words"`
		Speaker      string            `json:"speaker,omitempty"`
		Chars        []SingleChar      `json:"chars,omitempty"`
		Translations map[string]string `json:"translations,omitempty"`
	}

	SingleChar struct {
		Start float64 `json:"start"`
		End   float64 `json:"end"`
		Char  string  `json:"char"`
		Score float64 `json:"score"`
	}

	SingleWord struct {
		Word    string  `json:"word"`
		Start   float64 `json:"start"`
		End     float64 `json:"end"`
		Score   float64 `json:"score"`
		Speaker string  `json:"speaker,omitempty"`
	}

	// ResegmentOption is for segment merge
	ResegmentOption struct {
		MaxInterval   float64 // max time interval between words, in second
		MaxLineLength int     // max line length of each segment
		LangCode      string  // language used for calc length when multilanguage
	}

	// SubtitleOption
	SubtitleOption struct {
		LangCode      string
		LineBreak     bool // default false
		MaxLineLength int
		Layer         int
		Style         *astisub.Style
	}
)

func (s *SingleSegment) genItems(opt *SubtitleOption) []*astisub.Item {
	var items []*astisub.Item
	text := s.Text
	if s.Translations != nil {
		translated, ok := s.Translations[opt.LangCode]
		if ok {
			text = translated
		}
	}
	length := len(text)
	texts := []string{text}
	if opt.LineBreak {
		texts = breakLineByPunctuation(text, opt.MaxLineLength)
	}
	startAt := time.Duration(s.Start*1000) * time.Millisecond
	for _, line := range texts {
		i := strings.Index(text, line)

		item := &astisub.Item{
			Style: opt.Style,
			InlineStyle: &astisub.StyleAttributes{
				SSALayer: astikit.IntPtr(opt.Layer),
			},
			StartAt: startAt,
			Lines:   make([]astisub.Line, 0),
		}

		item.Lines = append(item.Lines, astisub.Line{
			Items: []astisub.LineItem{
				{
					Text: line,
				},
			},
		})
		if i+len(line) >= length {
			item.EndAt = time.Duration(s.End*1000) * time.Millisecond
		} else {
			startAt = time.Duration((s.Start+(s.End-s.Start)*float64(i+len(line))/float64(length))*1000) * time.Millisecond
			item.EndAt = startAt
		}

		items = append(items, item)
	}
	return items
}

func (s *SingleSegment) mergeWith(other *SingleSegment) {
	s.Words = append(s.Words, other.Words...)
	s.Text += " " + strings.TrimSpace(other.Text)
	s.End = other.Words[len(other.Words)-1].End
	if s.Translations != nil {
		for k, v := range s.Translations {
			s.Translations[k] = v + wordSpace(k) + strings.TrimSpace(other.Translations[k])
		}
	}
}

func (w *SingleWord) mustBreak() bool {
	return mustBreak.MatchString(w.Word)
	// return hasSuffixAny(strings.TrimSpace(w.Word), punctuationMustBreak...)
}

func allowMerge(segment1, segment2 *SingleSegment, srcLang string, option *ResegmentOption) bool {
	text1 := segment1.Text
	text2 := segment2.Text
	if option.LangCode != srcLang {
		t1, ok := segment1.Translations[option.LangCode]
		if ok {
			t2, ok := segment2.Translations[option.LangCode]
			if ok {
				text1 = t1
				text2 = t2
			}
		}
	}
	width := runewidth.StringWidth(text1) + runewidth.StringWidth(text2)
	return width <= option.MaxLineLength
}
