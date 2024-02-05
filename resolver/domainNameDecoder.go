package resolver

import (
	"fmt"
	"strings"
)

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
	return decodeLabelNames(buffer, offset)
}

func decodeLabelNames(data []byte, offset int) (string, int, error) {
	labels := []string{}

	for offset < len(data) {
		labelLength := int(data[offset])
		offset += 1
		if labelLength == 0 {
			break
		}

		if labelLength > 63 || offset+labelLength > len(data) {
			return "", 0, fmt.Errorf(
				"Invalid label in question. "+
					"Expected at most 63 bytes or %d bytes, got %d",
				len(data),
				labelLength,
			)
		}

		labels = append(labels, string(data[offset:offset+labelLength]))
		offset += labelLength
	}

	if offset+4 > len(data) {
		return "", 0, fmt.Errorf(
			"Invalid question. Expected at least 4 bytes for Type and Class, got %d",
			len(data)-offset,
		)
	}

	return strings.Join(labels, "."), offset, nil
}
