package resolver

import "fmt"

type Header struct {
	ID                    uint16
	Flags                 uint16
	QuestionCount         uint16
	AnswerRecordCount     uint16
	AuthorityRecordCount  uint16
	AdditionalRecordCount uint16
}

func (h *Header) Encode() []byte {

	buffer := make([]byte, 12)

	// Ex: 1010 0111 0100 1111
	buffer[0] = byte(h.ID >> 8) // 1º shifting the bits 8 places to the right -> 0000 0000 1010 0111
	/*
		2º

		h.ID:     1010 0111 0100 1111
		0xFF:     0000 0000 1111 1111
		------------------------------
		AND:      0000 0000 0100 1111
	*/
	buffer[1] = byte(h.ID & 0xFF)
	buffer[2] = byte(h.Flags >> 8)
	buffer[3] = byte(h.Flags & 0xFF)
	buffer[4] = byte(h.QuestionCount >> 8)
	buffer[5] = byte(h.QuestionCount & 0xFF)
	buffer[6] = byte(h.AnswerRecordCount >> 8)
	buffer[7] = byte(h.AnswerRecordCount & 0xFF)
	buffer[8] = byte(h.AuthorityRecordCount >> 8)
	buffer[9] = byte(h.AuthorityRecordCount & 0xFF)
	buffer[10] = byte(h.AdditionalRecordCount >> 8)
	buffer[11] = byte(h.AdditionalRecordCount & 0xFF)

	return buffer
}

func DecodeHeader(data []byte) (h *Header, offset int, err error) {

	if len(data) < 12 {
		return nil, 0, fmt.Errorf("Invalid header. Expected 12 bytes, got %d", len(data))
	}

	h = &Header{}

	/*
		uint16(data[0])<<8: This takes the first byte (which is the high byte),
		converts it to uint16 to prevent data loss during the shift operation,
		and then shifts it 8 places to the left.
		This effectively moves the high byte to its original position in the uint16.

		uint16(data[1]): This takes the second byte (which is the low byte) and converts it to uint16.

		|: This is the bitwise OR operator. It combines the high byte and the low byte to form the original uint16.
	*/
	h.ID = uint16(data[0])<<8 | uint16(data[1])
	h.Flags = uint16(data[2])<<8 | uint16(data[3])
	h.QuestionCount = uint16(data[4])<<8 | uint16(data[5])
	h.AnswerRecordCount = uint16(data[6])<<8 | uint16(data[7])
	h.AuthorityRecordCount = uint16(data[8])<<8 | uint16(data[9])
	h.AdditionalRecordCount = uint16(data[10])<<8 | uint16(data[11])

	return h, 12, nil
}
