package resolver

import (
	"strings"
)

func DecodeDomainName(buffer []byte, offset int) (string, int) {
	idx := offset
	offsetIncrement := 0

	var labelLength int
	pointerFound := false

	labels := []string{}

	for idx < len(buffer) {
		labelLength = int(buffer[idx])
		idx += 1
		if labelLength >= 192 { //val greater than 192 indicates a pointer
			idx = int(buffer[idx])
			pointerFound = true
		} else {
			offsetIncrement += 1
			if labelLength == 0 {
				break
			}
			labels = append(labels, string(buffer[idx:idx+labelLength]))
			offsetIncrement += labelLength
			idx += labelLength
		}
	}

	if pointerFound {
		offsetIncrement = 2 //if it was a pointer its just 2
	}
	return strings.Join(labels, "."), offset + offsetIncrement
}
