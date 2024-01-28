package resolver

import (
	"fmt"
	"strings"
)

type Question struct {
	Name  string
	Type  uint16
	Class uint16
}

func (q *Question) Encode() ([]byte, error) {
	var buffer []byte

	dnsLabels := strings.Split(q.Name, ".")

	for _, label := range dnsLabels {
		if len(label) > 63 {
			return nil, fmt.Errorf("Invalid label. Label %s exceeds 63 characters", label)
		}

		buffer = append(buffer, byte(len(label))) // 1 byte
		buffer = append(buffer, []byte(label)...) // n bytes, where maximum value of n is 63
	}
	buffer = append(buffer, 0) // 1 byte to indicate the end of the domain name

	buffer = append(buffer, byte(q.Type>>8), byte(q.Type&0xFF))   // 2 bytes
	buffer = append(buffer, byte(q.Class>>8), byte(q.Class&0xFF)) // 2 bytes

	// Maximum buffer size -> (4 * 63 + 1 + 2 + 2) = 257 bytes
	return buffer, nil
}

func DecodeQuestion(data []byte) (q *Question, err error) {
	q = &Question{}

	labels := []string{}

	var i = 0
	for i < len(data) {
		labelLength := int(data[i])
		i += 1
		if labelLength == 0 {
			break
		}

		if labelLength > 63 || i+labelLength > len(data) {
			return nil, fmt.Errorf(
				"Invalid label in question. "+
					"Expected at most 63 bytes or %d bytes, got %d",
				len(data),
				labelLength,
			)
		}

		labels = append(labels, string(data[i:i+labelLength]))
		i += labelLength
	}

	if i+4 > len(data) {
		return nil, fmt.Errorf(
			"Invalid question. Expected at least 4 bytes for Type and Class, got %d",
			len(data),
		)
	}

	q.Name = strings.Join(labels, ".")
	q.Type = uint16(data[len(data)-4])<<8 | uint16(data[len(data)-3])
	q.Class = uint16(data[len(data)-2])<<8 | uint16(data[len(data)-1])

	return q, nil
}
