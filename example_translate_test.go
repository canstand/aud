package aud_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/canstand/aud"
)

var (
	toZh = map[string]string{
		"The Honorable Charles Smith, Miss Sarah's brother, was walking swiftly uptown from Mr. Easterly's Wall Street office and his face was pale.": "查尔斯-史密斯阁下，莎拉小姐的哥哥，正从伊斯特里先生的华尔街办公室快步向市区走来，脸色苍白。",
		"At last the cotton combine was to all appearances an assured fact and he was slated for the Senate.":                                         "终于，棉花联合的事情看起来已经稳妥了，他也被提名为参议员。",
		"Why should he not be as other men?": "他为什么不能像其他人一样呢？",
		"She was not herself a notably intelligent woman, she greatly admired intelligence or whatever looked to her like intelligence in others.":                                                                                                                                                                                                         "她自己并不是一个特别聪明的女人，但她非常欣赏别人的聪明才智，或者在她看来像聪明才智的东西。",
		"So persuasive were her entreaties, and so strong her assurances that no harm whatever could result to them, from the information she sought.":                                                                                                                                                                                                     "她的恳求如此有说服力，她的保证如此有力，她所寻求的信息不会对他们造成任何伤害。",
		"They were induced to confess that one summer's night, the same she had mentioned, themselves and another friend being out on a stroll with Rodolfo, they had been concerned in the adduction of a girl whom Rodolfo carried off, whilst the rest of them detained her family, who made a great outcry and would have defended her if they could.": "在她的诱导下，他们承认，就在她提到的那个夏夜，他们和另一个朋友与鲁道夫一起出去散步，他们参与了绑架一个女孩，鲁道夫把她带走了，而其他人则拦住了她的家人，她的家人大吵大闹，如果可以的话，他们会保护她的。",
	}
)

type mockTranslator struct{}

func (s *mockTranslator) Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	if targetLang != "zh" {
		return "", fmt.Errorf("language not supported: %s", targetLang)
	}
	result, ok := toZh[strings.TrimSpace(text)]
	if !ok {
		return "", fmt.Errorf("text not found: %s", text)
	}
	return result, nil
}

func ExampleTranscript_Translate() {
	f, err := os.OpenFile("testdata/libri.json", os.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	transcript, err := aud.ReadTranscript(f)
	if err != nil {
		log.Fatal(err)
	}
	err = transcript.Translate(context.Background(), &mockTranslator{}, "zh")
	if err != nil {
		log.Fatal(err)
	}
	s, err := transcript.GenSubtitle("from libri", []aud.SubtitleOption{
		{
			LangCode:  "zh",
			LineBreak: true,
		},
	}, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s.Items[len(s.Items)-2].String())
	fmt.Print(s.Items[len(s.Items)-1].String())
	// Output:
	// 她的家人大吵大闹，如果可以的话，
	// 他们会保护她的。
}
