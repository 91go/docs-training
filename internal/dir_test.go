package internal

import (
	"testing"
)

func TestExtractQuestion(t *testing.T) {
	filePath := "../testdata/test.md"

	questions := ExtractQuestion(filePath)

	expectedQuestions := []Question{
		{
			Text: "There is no royal road to learning？.",
			URL:  "../testdata/test#there-is-no-royal-road-to-learning.",
		},
		{
			Text: "There is no royal road to learning?.",
			URL:  "../testdata/test#there-is-no-royal-road-to-learning.",
		},
		{
			Text: "【There is no royal road to learning.】",
			URL:  "<BaseURL>/testdata/test.md#Question%203",
		},
		{
			Text: "“Variety is the spice of life.”",
			URL:  "../testdata/test#variety-is-the-spice-of-life.",
		},
		{
			Text: "\"Variety is the spice of life.\"",
			URL:  "../testdata/test#variety-is-the-spice-of-life.",
		},
		{
			Text: "Variety is the spice of life.+-",
			URL:  "../testdata/test#variety-is-the-spice-of-life.",
		},
		{
			Text: "Doubt is the key to knowledge.",
			URL:  "../testdata/test#doubt-is-the-key-to-knowledge.",
		},
		{
			Text: "Doubt is the key to knowledge.",
			URL:  "../testdata/test#doubt-is-the-key-to-knowledge.",
		},
		{
			Text: "Doubt is the key to knowledge.",
			URL:  "../testdata/test#doubt-is-the-key-to-knowledge.",
		},
	}

	// 验证问题数量是否一致
	if len(questions) != len(expectedQuestions) {
		t.Errorf("Unexpected number of questions. Expected: %d, Got: %d", len(expectedQuestions), len(questions))
	}

	// 逐个验证每个问题的文本和URL
	for i, question := range questions {
		if question.Text != expectedQuestions[i].Text {
			t.Errorf("Unexpected question text. Expected: %s, Got: %s", expectedQuestions[i].Text, question.Text)
		}

		// if question.URL != expectedQuestions[i].URL {
		// 	t.Errorf("Unexpected question URL. Expected: %s, Got: %s", expectedQuestions[i].URL, question.URL)
		// }
	}
}
