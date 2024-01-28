package resolver

import "math/rand"

type DNSMessage struct {
	Header     Header
	Questions  []Question
	Answers    []ResourceRecord
	Authority  []ResourceRecord
	Additional []ResourceRecord
}

func GenMessage() ([]byte, error) {
	header := Header{
		ID:                    uint16(rand.Intn(1 << 16)),
		Flags:                 128,
		QuestionCount:         1,
		AnswerRecordCount:     0,
		AuthorityRecordCount:  0,
		AdditionalRecordCount: 0,
	}
	encodedHeader := header.Encode()

	question := Question{
		Name:  "www.google.com",
		Type:  1,
		Class: 1,
	}
	questionEncoded, err := question.Encode()

	if err != nil {
		return nil, err
	}

	message := append(encodedHeader, questionEncoded...)

	return message, nil
}
