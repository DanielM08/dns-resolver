package resolver

import (
	"math/rand"
)

type DNSMessage struct {
	Header     Header
	Question   Question
	Answers    []ResourceRecord
	Authority  []ResourceRecord
	Additional []ResourceRecord
}

func GenMessage(domainName string) ([]byte, error) {
	header := Header{
		ID:                    uint16(rand.Intn(1 << 16)),
		Flags:                 0,
		QuestionCount:         1,
		AnswerRecordCount:     0,
		AuthorityRecordCount:  0,
		AdditionalRecordCount: 0,
	}
	encodedHeader := header.Encode()

	question := Question{
		Name:  domainName,
		Type:  1,
		Class: 1,
	}

	questionEncoded, err := question.Encode()

	if err != nil {
		return nil, err
	}

	message := append(encodedHeader, questionEncoded...)
	message = append(message, 0, 0, 0)
	return message, nil
}

func DecodeMessage(dnsResponse []byte) (DNSMessage, error) {
	header, offset, err := DecodeHeader(dnsResponse)

	if err != nil {
		return DNSMessage{}, err
	}

	question, offset, err := DecodeQuestion(dnsResponse, offset)

	if err != nil {
		return DNSMessage{}, err
	}

	answerRecords, offset, err := getResourceRecords(dnsResponse, int(header.AnswerRecordCount), offset)

	if err != nil {
		return DNSMessage{}, err
	}

	authorityRecords, offset, err := getResourceRecords(dnsResponse, int(header.AuthorityRecordCount), offset)

	if err != nil {
		return DNSMessage{}, err
	}

	additionalRecords, _, err := getResourceRecords(dnsResponse, int(header.AdditionalRecordCount), offset)

	if err != nil {
		return DNSMessage{}, err
	}

	return DNSMessage{
		Header:     *header,
		Question:   *question,
		Answers:    answerRecords,
		Authority:  authorityRecords,
		Additional: additionalRecords,
	}, nil
}

func getResourceRecords(buffer []byte, count, offset int) ([]ResourceRecord, int, error) {
	var resourceRecords = make([]ResourceRecord, count)
	for i := 0; i < int(count); i++ {
		answerRecord, newOffset, err := DecodeResourceRecord(buffer, offset)

		if err != nil {
			return resourceRecords, 0, err
		}

		offset = newOffset
		resourceRecords[i] = *answerRecord
	}

	return resourceRecords, offset, nil
}
