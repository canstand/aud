package aud_test

import (
	"fmt"
	"log"
	"os"

	"github.com/canstand/aud"
)

func ExampleTranscript_GenSubtitle() {
	f, err := os.OpenFile("testdata/libri.json", os.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	transcript, err := aud.ReadTranscript(f)
	if err != nil {
		log.Fatal(err)
	}
	s, err := transcript.GenSubtitle("from libri", []aud.SubtitleOption{
		{
			LangCode: "en",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(s.Items[len(s.Items)-1].String())
	// Output:
	// They were induced to confess that one summer's night, the same she had mentioned, themselves and another friend being out on a stroll with Rodolfo, they had been concerned in the adduction of a girl whom Rodolfo carried off, whilst the rest of them detained her family, who made a great outcry and would have defended her if they could.
}

func ExampleTranscript_GenSubtitle_breakLines() {
	f, err := os.OpenFile("testdata/libri.json", os.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	transcript, err := aud.ReadTranscript(f)
	if err != nil {
		log.Fatal(err)
	}
	s, err := transcript.GenSubtitle("from libri", []aud.SubtitleOption{
		{
			LangCode:  "en",
			LineBreak: true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(s.Items[len(s.Items)-1].String())
	// Output:
	// who made a great outcry and would have defended her if they could.
}
