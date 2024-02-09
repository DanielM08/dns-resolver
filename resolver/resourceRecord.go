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

func DecodeDomainNameWithPointer(buffer []byte, offset int) (string, int) {
	nameByteLength := 0
	var name string
	idx := offset
	var length int
	shouldIncreasePointer := false
	for idx < len(buffer) {
		length = int(buffer[idx])
		if length >= 192 { //val greater than 192 indicates a pointer
			idx = int(buffer[idx+1])
			shouldIncreasePointer = true
		} else {
			if length == 0 {
				break
			}
			name = name + string(buffer[idx+1:idx+length+1])
			name = name + "."
			idx = idx + 1 + length
		}
	}
	if name != "" {
		name = name[:len(name)-1]
	}
	if shouldIncreasePointer {
		nameByteLength = 2 //if it was a pointer its just 2
	} else {
		nameByteLength = len([]byte(name)) + 2 //2 for truncating 0
	}
	return name, offset + nameByteLength
}

func DecodeResourceRecord(buffer []byte, offset int) (*ResourceRecord, int, error) {
	name, offset := DecodeDomainNameWithPointer(buffer, offset)

	var rrType = uint16(buffer[offset])<<8 + uint16(buffer[offset+1])
	var rrClass = uint16(buffer[offset+2])<<8 + uint16(buffer[offset+3])
	var rrTtl = uint32(buffer[offset+4])<<24 + uint32(buffer[offset+5])<<16 + uint32(buffer[offset+6])<<8 + uint32(buffer[offset+7])
	var rrLength = uint16(buffer[offset+8])<<8 + uint16(buffer[offset+9])

	// Check if the buffer contains enough data for rData.
	if offset+10+int(rrLength) > len(buffer) {
		return nil, 0, fmt.Errorf("buffer too short for rData")
	}

	var rrData = ""
	if rrType == 1 || rrType == 28 {
		labels := []string{}

		for i := 0; i < int(rrLength); i++ {
			labels = append(labels, fmt.Sprintf("%d", buffer[offset+10+i]))
		}
		rrData = strings.Join(labels, ".")
	} else if rrType == 2 || rrType == 6 {
		rrData, _ = DecodeDomainNameWithPointer(buffer, offset+10)
	} else {
		rrData = string(buffer[offset+10 : offset+10+int(rrLength)])
	}

	return &ResourceRecord{
		Name:   name,
		Type:   rrType,
		Class:  rrClass,
		TTL:    rrTtl,
		Length: rrLength,
		Data:   rrData,
	}, offset + 10 + int(rrLength), nil
}
