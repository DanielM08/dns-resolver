package resolver

import (
	"fmt"
)

type ResourceRecord struct {
	Name   string
	Type   uint16
	Class  uint16
	TTL    uint32
	Length uint16
	Data   string
}

func DecodeResourceRecord(buffer []byte, offset int) (*ResourceRecord, int, error) {
	name, offset, err := DecodeDomainName(buffer, offset)

	if err != nil {
		return nil, 0, err
	}

	var rrType = uint16(buffer[offset])<<8 + uint16(buffer[offset+1])
	var rrClass = uint16(buffer[offset+2])<<8 + uint16(buffer[offset+3])
	var rrTtl = uint32(buffer[offset+4])<<24 + uint32(buffer[offset+5])<<16 + uint32(buffer[offset+6])<<8 + uint32(buffer[offset+7])
	var rrLength = uint16(buffer[offset+8])<<8 + uint16(buffer[offset+9])

	// Check if the buffer contains enough data for rData.
	if offset+10+int(rrLength) > len(buffer) {
		return nil, 0, fmt.Errorf("buffer too short for rData")
	}

	return &ResourceRecord{
		Name:   name,
		Type:   rrType,
		Class:  rrClass,
		TTL:    rrTtl,
		Length: rrLength,
		Data:   string(buffer[offset+10 : offset+10+int(rrLength)]),
	}, offset + 10 + int(rrLength), nil
}

func DecodeDomainName(buffer []byte, offset int) (string, int, error) {
	const mask = 0xC0                // 192 in decimal
	if buffer[offset]&mask == mask { // A Pointer
		pointerOffset := int(buffer[offset]&^mask)<<8 + int(buffer[offset+1])

		domainName, _, err := DecodeDomainName(buffer, pointerOffset)

		if err != nil {
			return "", 0, err
		}

		return domainName, offset + 2, nil
	}
	return DecodeLabelNames(buffer, offset)
}
