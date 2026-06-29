package robotstxt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

var encodingEndianness = binary.LittleEndian

func computeListBufferLen(list []rulePattern) (buffLen int, validLengths bool) {
	totalPayloadLen := len(list) * 2 // u16 headers
	for _, elem := range list {
		if len(elem.Raw) > math.MaxUint16 {
			return 0, false
		}

		totalPayloadLen += len(elem.Raw)
	}

	return 4 + totalPayloadLen, true
}

// buf is expected to be large enough in advance
func encodeList(list []rulePattern, buf []byte) {
	encodingEndianness.PutUint32(buf, uint32(len(buf)-4))

	currOffset := 4

	for _, elem := range list {
		encodingEndianness.PutUint16(buf[currOffset:], uint16(len(elem.Raw)))

		copy(buf[currOffset+2:], []byte(elem.Raw))

		currOffset += 2 + len(elem.Raw)
	}
}

func EncodeRuleset(ruleset Ruleset) ([]byte, error) {
	if len(ruleset.allow) == 0 && len(ruleset.disallow) == 0 {
		return []byte{}, nil
	}

	allowLen, allowsValid := computeListBufferLen(ruleset.allow)
	if !allowsValid {
		return nil, errors.New("allow list contained paths too long for encoding")
	}

	disallowLen, disallowsValid := computeListBufferLen(ruleset.disallow)
	if !disallowsValid {
		return nil, errors.New("disallow list contained paths too long for encoding")
	}

	buf := make([]byte, allowLen+disallowLen)
	encodeList(ruleset.allow, buf[:allowLen])
	encodeList(ruleset.disallow, buf[allowLen:])

	return buf, nil
}

func consumeList(data []byte) (list []rulePattern, remaningData []byte, err error) {
	if len(data) < 4 {
		return nil, nil, fmt.Errorf("array cannot fit length header")
	}

	listLen := int(encodingEndianness.Uint32(data[0:4]))

	data = data[4:]

	if listLen > len(data) {
		return nil, nil, fmt.Errorf("got invalid list length header at offset 0")
	}

	offset := 0
	for offset < listLen {
		strStartsAt := offset + 2

		strLen := encodingEndianness.Uint16(data[offset:strStartsAt])

		nextItemStartsAt := strStartsAt + int(strLen)

		if nextItemStartsAt > listLen {
			return nil, nil, fmt.Errorf("got invalid string length header at offset %d", offset)
		}

		list = append(list, parsePattern(data[strStartsAt:nextItemStartsAt]))

		offset = nextItemStartsAt
	}

	return list, data[offset:], nil
}

func DecodeRuleset(encodedRuleset []byte) (Ruleset, error) {
	if len(encodedRuleset) == 0 {
		return Ruleset{
			allow:    []rulePattern{},
			disallow: []rulePattern{},
		}, nil
	}

	var decoded Ruleset

	var err error

	decoded.allow, encodedRuleset, err = consumeList(encodedRuleset)
	if err != nil {
		return Ruleset{}, fmt.Errorf("failed to consume allow list: %s", err.Error())
	}

	decoded.disallow, encodedRuleset, err = consumeList(encodedRuleset)
	if err != nil {
		return Ruleset{}, fmt.Errorf("failed to consume disallow list: %s", err.Error())
	}

	if len(encodedRuleset) > 0 {
		return Ruleset{}, fmt.Errorf("encoded ruleset has more data than expected after lists")
	}

	return decoded, nil
}
