package resolver

import (
	"fmt"
	"strings"
)

type ResourceRecord struct {
	Name   string
	Type   uint16
	Class  uint16
	TTL    uint32
	Length uint16
	Data   string
}

func DecodeResource(buffer []byte, startoffset int) (*ResourceRecord, int, error) {

	return &ResourceRecord{}, 0, nil
}

func DecodeDomainName(buffer []byte, offset int) (string, int, error) {
	if offset >= len(buffer) {
		return "", 0, fmt.Errorf("Invalid domain name. Offset %d exceeds buffer size %d", offset, len(buffer))
	}

	const mask = 0xC0                // 192 in decimal
	if buffer[offset]&mask == mask { // A Pointer
		pointerOffset := int(buffer[offset]&^mask)<<8 + int(buffer[offset+1])

		domainName, _, err := DecodeDomainName(buffer, pointerOffset)

		if err != nil {
			return "", 0, err
		}

		return domainName, offset + 2, nil
	}

	var labels []string
	for {
		labelLength := int(buffer[offset])

		offset += 1
		if labelLength == 0 {
			break
		}

		if labelLength > 63 || offset+labelLength > len(buffer) {
			return "", 0, fmt.Errorf(
				"Invalid label in question. "+
					"Expected at most 63 bytes or %d bytes, got %d",
				len(buffer),
				labelLength,
			)
		}

		labels = append(labels, string(buffer[offset:offset+labelLength]))
		offset += labelLength

		if buffer[offset] == 0 {
			offset += 1
			break
		}
	}
	return strings.Join(labels, "."), offset, nil
}
