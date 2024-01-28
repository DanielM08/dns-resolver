package resolver

import (
	"fmt"
	"math/rand"
)

type DNSMessage struct {
	Header     Header
	Question   Question
	Answers    []ResourceRecord
	Authority  []ResourceRecord
	Additional []ResourceRecord
}

func GenMessage() ([]byte, error) {
	header := Header{
		ID:                    uint16(rand.Intn(1 << 16)),
		Flags:                 256,
		QuestionCount:         1,
		AnswerRecordCount:     0,
		AuthorityRecordCount:  0,
		AdditionalRecordCount: 0,
	}
	encodedHeader := header.Encode()

	question := Question{
		Name:  "dns.google.com",
		Type:  1,
		Class: 1,
	}

	fmt.Printf("Header: %+v\n", header)
	fmt.Printf("Question: %+v\n", question)

	questionEncoded, err := question.Encode()

	if err != nil {
		return nil, err
	}

	message := append(encodedHeader, questionEncoded...)

	return message, nil
}

func DecodeMessage(dnsResponse []byte) (DNSMessage, error) {
	fmt.Println("Decoding message..., message size: ", len(dnsResponse))

	header, err := DecodeHeader(dnsResponse[:12])

	if err != nil {
		return DNSMessage{}, err
	}

	question, bytesConsumed, err := DecodeQuestion(dnsResponse[12:])

	if err != nil {
		return DNSMessage{}, err
	}

	// answers, err := DecodeResourceRecords(dnsResponse, header.AnswerRecordCount)

	// if err != nil {
	// 	return DNSMessage{}, err
	// }

	// authority, err := DecodeResourceRecords(dnsResponse, header.AuthorityRecordCount)

	// if err != nil {
	// 	return DNSMessage{}, err
	// }

	// additional, err := DecodeResourceRecords(dnsResponse, header.AdditionalRecordCount)

	// if err != nil {
	// 	return DNSMessage{}, err
	// }

	fmt.Printf("Response Header: %+v\n", header)
	fmt.Printf("Response Question: %+v\n", question)
	fmt.Printf("Other bytes: %x\n", dnsResponse[12+bytesConsumed:])

	return DNSMessage{
		Header:     *header,
		Question:   *question,
		Answers:    nil,
		Authority:  nil,
		Additional: nil,
	}, nil

}
