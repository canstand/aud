package aud

import (
	"encoding/json"
	"fmt"
	"io"
)

var (
	ErrFormatNotSupported = fmt.Errorf("format not supported")
	ErrOutOfRange         = fmt.Errorf("out of range")
	ErrLinesNotEqual      = fmt.Errorf("lines not equal")
	ErrLangNotAvailable   = fmt.Errorf("language not available")
)

// ReadTranscript reads transcript from io.Reader
// only support whisperx output for now
func ReadTranscript(r io.Reader) (*Transcript, error) {
	var t Transcript
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	// whisperx output json should contain `segments` and `word_segments`
	if t.Language == "" || len(t.Segments) == 0 || len(t.WordSegments) == 0 {
		return nil, ErrFormatNotSupported
	}
	if err := json.NewDecoder(r).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}
