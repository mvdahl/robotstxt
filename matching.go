package robotstxt

import (
	"bytes"
	"strings"
)

func parsePattern(raw []byte) rulePattern {
	mustEndAt := bytes.IndexByte(raw, '$')

	effectiveLen := len(raw)
	if mustEndAt != -1 {
		effectiveLen = mustEndAt
	}

	return rulePattern{
		Raw:       raw,
		Sequences: strings.Split(string(raw[:effectiveLen]), "*"), // todo: would be awesome to find a more performant way of doing this
		HasEnd:    mustEndAt != -1,
	}
}

func parsePatterns(list []string) []rulePattern {
	patterns := make([]rulePattern, len(list))

	for i, item := range list {
		patterns[i] = parsePattern([]byte(item))
	}

	return patterns
}

func isMatch(pattern rulePattern, str string) bool {
	for i, seq := range pattern.Sequences {
		wildcard := i > 0

		if i == len(pattern.Sequences)-1 && wildcard && seq == "" {
			str = ""
			break
		}

		if len(str) < len(seq) {
			return false
		}

		if !wildcard {
			if str[:len(seq)] != seq {
				return false
			}

			str = str[len(seq):]
		} else {
			startI := strings.Index(str, seq)
			if startI == -1 {
				return false
			}

			str = str[startI+len(seq):]
		}
	}

	if pattern.HasEnd && len(str) > 0 {
		return false
	}

	return true
}

func findMostSpecificMatch(patterns []rulePattern, str string, minLen int) (mostSpecific rulePattern, found bool) {
	for _, pattern := range patterns {
		if len(pattern.Raw) > len(mostSpecific.Raw) && len(pattern.Raw) >= minLen && isMatch(pattern, str) {
			mostSpecific = pattern
			found = true
		}
	}

	return mostSpecific, found
}
