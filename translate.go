package aud

import "context"

type Translator interface {
	// Translate translate text from source language to target language
	Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error)
}
