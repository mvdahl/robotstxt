package robotstxt

import (
	"bytes"
	"net/url"
	"strings"
)

// https://www.rfc-editor.org/rfc/rfc9309.html

type parseToken byte

const (
	parseTokenName    parseToken = 0
	parseTokenValue   parseToken = 1
	parseTokenComment parseToken = 2
)

var (
	userAgentName = []byte("user-agent")
	allowName     = []byte("allow")
	disallowName  = []byte("disallow")
)

type stringReader struct {
	buff       []byte
	StartIndex int
	LastIndex  int
	Result     []byte
}

func (r *stringReader) Reading() bool {
	return r.StartIndex != -1
}

func (r *stringReader) Finalized() bool {
	return r.Result != nil
}

// expects either whitespace or regular line characters
func (r *stringReader) Consider(i int) {
	if isWhitespace(r.buff[i]) {
		return
	}

	if r.StartIndex == -1 {
		r.StartIndex = i
		r.LastIndex = i
	} else {
		r.LastIndex = i
	}
}

func (r *stringReader) Finalize() {
	if !r.Reading() {
		r.Result = []byte{}
		return
	}

	r.Result = r.buff[r.StartIndex : r.LastIndex+1]
	r.StartIndex = -1
	r.LastIndex = -1
}

func (r *stringReader) Reset() {
	r.Result = nil
	r.StartIndex = -1
	r.LastIndex = -1
}

func newStringReader(buff []byte) *stringReader {
	r := &stringReader{
		buff: buff,
	}

	r.Reset()

	return r
}

func isWhitespace(char byte) bool {
	return char == 0x20 || char == 0x9
}

func isWhitespaceRune(char rune) bool {
	return char == 0x20 || char == 0x9
}

func isNewline(char byte) bool {
	return char == '\r' || char == '\n'
}

// b (the value checked against) is expected to be all-lowercase
func equalCaseInsensitive(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i, bA := range a {
		lowercaseVal := bA
		if bA >= 0x41 && bA <= 0x5A {
			lowercaseVal += 0x20
		}

		if b[i] != lowercaseVal {
			return false
		}
	}

	return true
}

func Parse(raw []byte, productToken string) Ruleset {
	productToken = strings.ToLower(productToken)

	var ruleset Ruleset

	lastLineWasUserAgent := false
	currentlyAcceptingRules := false

	expectedToken := parseTokenName

	nameReader := newStringReader(raw)
	valueReader := newStringReader(raw)

	i := 0

	// special case handling for Byte Order Mark
	if len(raw) >= 3 && raw[0] == 0xEF && raw[1] == 0xBB && raw[2] == 0xBF {
		i = 3
	}

	for ; i < len(raw); i++ {
		switch expectedToken {
		case parseTokenName:
			if isNewline(raw[i]) {
				nameReader.Reset()
			} else if raw[i] == '#' {
				if nameReader.Reading() || nameReader.Finalized() {
					nameReader.Reset()
				}

				expectedToken = parseTokenComment
			} else if raw[i] == ':' {
				nameReader.Finalize()
				expectedToken = parseTokenValue
			} else {
				nameReader.Consider(i)
			}
		case parseTokenValue:
			isValueEndingSymbol := isNewline(raw[i]) || raw[i] == '#'

			if isValueEndingSymbol || i == len(raw)-1 {
				if !isValueEndingSymbol {
					valueReader.Consider(i)
				}

				valueReader.Finalize()

				if equalCaseInsensitive(nameReader.Result, userAgentName) {
					if !lastLineWasUserAgent {
						currentlyAcceptingRules = false
					}

					lastLineWasUserAgent = true

					if (len(valueReader.Result) == 1 && valueReader.Result[0] == '*') || strings.Contains(productToken, strings.ToLower(string(valueReader.Result))) {
						currentlyAcceptingRules = true
					}
				} else if currentlyAcceptingRules {
					lastLineWasUserAgent = false

					if len(valueReader.Result) > 0 && (equalCaseInsensitive(nameReader.Result, allowName) || equalCaseInsensitive(nameReader.Result, disallowName)) {
						// note: over time it may be faster to do our own impl of PathUnescape that accepts []byte,
						// but for now i'm fine with just doing this, considering the vast majority of paths will not use unicode characters
						decodedResult := valueReader.Result
						valid := true
						if bytes.IndexByte(valueReader.Result, '%') != -1 {
							decodedVal, err := url.PathUnescape(string(valueReader.Result))
							if err == nil {
								decodedResult = []byte(decodedVal)
							} else {
								valid = false
							}
						}

						if valid {
							pattern := parsePattern(decodedResult)

							if equalCaseInsensitive(nameReader.Result, allowName) {
								ruleset.allow = append(ruleset.allow, pattern)
							} else {
								ruleset.disallow = append(ruleset.disallow, pattern)
							}
						}
					}
				}

				nameReader.Reset()
				valueReader.Reset()

				if raw[i] == '#' {
					expectedToken = parseTokenComment
				} else {
					expectedToken = parseTokenName
				}
			} else {
				valueReader.Consider(i)
			}
		case parseTokenComment:
			if isNewline(raw[i]) {
				expectedToken = parseTokenName
			}
		}
	}

	return ruleset
}
