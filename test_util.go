package robotstxt

import (
	"bytes"
	_ "embed"
)

//go:embed testdata/example.txt
var exampleRobotsTxt []byte

const exampleRobotsProduct = "foobot"

//go:embed testdata/google.txt
var googleRobotsTxt []byte

const googleRobotsProduct = "yandex"

//go:embed testdata/github.txt
var githubRobotsTxt []byte

const githubRobotsProduct = "bingbot"

func equivelentRuleLists(a []rulePattern, b []rulePattern) bool {
	if len(a) != len(b) {
		return false
	}

	for i, itemA := range a {
		if !bytes.Equal(b[i].Raw, itemA.Raw) {
			return false
		}
	}

	return true
}
