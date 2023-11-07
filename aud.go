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
	if err := decoder.Decode(&t); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFormatNotSupported, err)
	}
	// whisperx output json should contain `segments` and `word_segments`
	if t.Language == "" || len(t.Segments) == 0 || len(t.WordSegments) == 0 {
		return nil, ErrFormatNotSupported
	}
	return &t, nil
}
