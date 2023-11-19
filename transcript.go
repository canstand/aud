package aud

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/asticode/go-astisub"
)

const (
	DefaultMaxLineLength    = 72
	DefaultCJKMaxLineLength = 42
)

// Transcript is for speech-to-text transcription
type Transcript struct {
	sync.RWMutex
	Segments     []SingleSegment `json:"segments"`
	WordSegments []SingleWord    `json:"word_segments"`
	Language     string          `json:"language"`
}

// AvailableLangs returns available languages, the last one is the original.
func (t *Transcript) AvailableLangs() []string {
	t.RLock()
	defer t.RUnlock()
	langs := make([]string, 0)
	if len(t.Segments) > 0 {
		for lang := range t.Segments[0].Translations {
			if !slices.Contains(langs, lang) {
				langs = append(langs, lang)
			}
		}
	}
	langs = append(langs, t.Language)
	return langs
}

// GenSubtitle generates subtitle
func (t *Transcript) GenSubtitle(title string, langs []SubtitleOption, optimize bool) (*astisub.Subtitles, error) {
	t.RLock()
	defer t.RUnlock()
	s := astisub.NewSubtitles()
	s.Metadata = ssaDefaultMetadata
	s.Metadata.Title = title

	for _, lang := range langs {
		if !slices.Contains(t.AvailableLangs(), lang.LangCode) {
			return nil, fmt.Errorf("%w: %s", ErrLangNotAvailable, lang.LangCode)
		}
		if lang.Style == nil {
			lang.Style = SSADefaultStyle
		}
		s.Styles[lang.Style.ID] = lang.Style
		if lang.LineBreak && lang.MaxLineLength == 0 {
			switch lang.LangCode {
			case "zh":
				fallthrough
			case "ko":
				fallthrough
			case "ja":
				lang.MaxLineLength = DefaultCJKMaxLineLength
			default:
				lang.MaxLineLength = DefaultMaxLineLength
			}
		}
		for _, segment := range t.Segments {
			s.Items = append(s.Items, segment.genItems(&lang)...)
		}
	}
	if optimize {
		OptimizeGaps(s)
	}
	return s, nil
}

// SplitSegment split a segment into multiple segments based on \n
//
//	Only if the text and translation have the same number of lines
func (t *Transcript) SplitSegment(index int) error {
	t.Lock()
	defer t.Unlock()
	if index < 0 || index >= len(t.Segments) {
		return fmt.Errorf("%w: %d", ErrOutOfRange, index)
	}

	texts := strings.Split(t.Segments[index].Text, "\n")
	translatedTexts := make(map[string][]string)
	for k, v := range t.Segments[index].Translations {
		translatedTexts[k] = strings.Split(v, "\n")
		if len(translatedTexts[k]) != len(texts) {
			return fmt.Errorf("%w: lang %s requires %d lines, got %d", ErrLinesNotEqual, k, len(texts), len(translatedTexts[k]))
		}
	}

	newSegments := make([]SingleSegment, 0)
	words := slices.Clone(t.Segments[index].Words)
	wordIndex := 0
	for i, text := range texts {
		segment := SingleSegment{
			Start:        words[wordIndex].Start,
			Text:         text,
			Translations: make(map[string]string),
		}
		line := ""
		for j := wordIndex; j < len(words); j++ {
			line += " " + words[j].Word
			segment.Words = append(segment.Words, words[j])
			if len(strings.TrimSpace(line)) > len(strings.TrimSpace(text))-1 {
				wordIndex = j + 1
				break
			}
		}
		segment.End = segment.Words[len(segment.Words)-1].End
		for k, v := range translatedTexts {
			segment.Translations[k] = v[i]
		}
		newSegments = append(newSegments, segment)
	}

	t.Segments = slices.Insert(slices.Delete[[]SingleSegment](t.Segments, index, index+1), index, newSegments...)
	return nil
}

// resegmentByWords, use before translate only, otherwise lose translations
func (t *Transcript) resegmentByWords(opt ResegmentOption) error {
	opt.LangCode = t.Language // always use source language
	if opt.MaxInterval == 0 {
		opt.MaxInterval = 0.5
	}
	if opt.MaxLineLength == 0 {
		opt.MaxLineLength = 200
	}
	var (
		segment  = &SingleSegment{}
		segments = make([]SingleSegment, 0)
		space    = wordSpace(t.Language)
	)
	for _, word := range t.WordSegments {
		if word.End == 0 && word.Start == 0 {
			return fmt.Errorf("word lost start and end: %v", word.Word)
		}
		segment.Words = append(segment.Words, word)
		if len(segment.Words) == 1 {
			segment.Start = word.Start
			segment.Text += strings.TrimSpace(word.Word)
		} else {
			segment.Text += space + strings.TrimSpace(word.Word)
		}
		segment.End = word.End
		if !word.mustBreak() { // not sentence end
			continue
		}

		// if interval less than maxInterval and total length less than maxLineLength
		if len(segments) > 0 && (segment.Start-segments[len(segments)-1].End) < opt.MaxInterval && len(segment.Text)+len(segments[len(segments)-1].Text) < opt.MaxLineLength {
			// merge to previous segment
			segments[len(segments)-1].mergeWith(segment)
			segment = &SingleSegment{}
			continue
		}
		segments = append(segments, *segment)
		segment = &SingleSegment{}
	}
	// may miss last segment
	if segment.Text != "" && segments[len(segments)-1].Text != segment.Text {
		segments = append(segments, *segment)
	}
	t.Segments = segments
	return nil
}

// Resegment, merge those short segments
func (t *Transcript) Resegment(option *ResegmentOption) []SingleSegment {
	t.Lock()
	defer t.Unlock()
	if option.LangCode == "" {
		option.LangCode = t.Language
	}
	if option.MaxInterval == 0 {
		option.MaxInterval = 0.5
	}
	if option.MaxLineLength == 0 {
		switch option.LangCode {
		case "zh":
			fallthrough
		case "ko":
			fallthrough
		case "ja":
			option.MaxLineLength = DefaultCJKMaxLineLength
		default:
			option.MaxLineLength = DefaultMaxLineLength
		}
	}
	var (
		segments = make([]SingleSegment, 0)
	)
	for _, segment := range t.Segments {
		if len(segments) > 0 && (segment.Start-segments[len(segments)-1].End) < option.MaxInterval && allowMerge(&segment, &(segments[len(segments)-1]), t.Language, option) {
			segments[len(segments)-1].mergeWith(&segment)
		} else {
			segments = append(segments, segment)
		}
	}
	t.Segments = segments
	return segments
}

// Translate each segments into target language
func (t *Transcript) Translate(ctx context.Context, svc Translator, targetLang string) error {
	return t.TranslateWithOption(ctx, svc, targetLang, TranslateOption{})
}

// TranslateWithOption translate each segments into target language
func (t *Transcript) TranslateWithOption(ctx context.Context, svc Translator, targetLang string, option TranslateOption) error {
	t.Lock()
	defer t.Unlock()
	if targetLang == t.Language {
		return nil // nothing to do
	}
	if option.ResegmentByWords {
		if err := t.resegmentByWords(ResegmentOption{}); err != nil {
			return err
		}
	}
	total := len(t.Segments)
	for index, segment := range t.Segments {
		if !option.Override && t.Segments[index].Translations != nil {
			_, ok := t.Segments[index].Translations[targetLang]
			if ok {
				continue
			}
		}
		translation, err := svc.Translate(ctx, segment.Text, t.Language, targetLang)
		if err != nil {
			return err
		}
		if t.Segments[index].Translations == nil {
			t.Segments[index].Translations = make(map[string]string)
		}
		t.Segments[index].Translations[targetLang] = translation
		if option.Callback != nil {
			option.Callback(index+1, total)
		}
	}
	return nil
}

type TranslateOption struct {
	Override         bool                      // override existing translation, or else continue
	ResegmentByWords bool                      // Warning: possible to get better results, but all previous translations for all languages will be purged first
	Callback         func(finished, total int) // for progress report
}
