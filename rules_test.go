package robotstxt

import (
	"testing"

	"github.com/temoto/robotstxt"
)

func TestRules(t *testing.T) {
	ruleset := Ruleset{
		disallow: parsePatterns([]string{"/", "/app/secrets"}),
		allow:    parsePatterns([]string{"/app", "/app/secrets/public"}),
	}

	paths := map[string]bool{
		"/":                        false,
		"/help":                    false,
		"/app":                     true,
		"/app/dashboard":           true,
		"/app/secrets":             false,
		"/app/secrets/public/info": true,
	}

	for path, shouldMatch := range paths {
		if ruleset.Test(path) != shouldMatch {
			t.Fatal("path", path, "should have result", shouldMatch, "but got the opposite")
		}
	}
}

var benchMatchPath = "/maps?input=bogus&output=classic"

func BenchmarkOwnPathTests(b *testing.B) {
	ruleset := Parse(googleRobotsTxt, googleRobotsProduct)

	b.ResetTimer()

	for b.Loop() {
		ruleset.Test(benchMatchPath)
	}
}

func BenchmarkTemotoPathTests(b *testing.B) {
	thing, err := robotstxt.FromBytes(googleRobotsTxt)
	if err != nil {
		b.Fatal(err.Error())
	}

	gr := thing.FindGroup(googleRobotsProduct)

	b.ResetTimer()

	for b.Loop() {
		gr.Test(benchMatchPath)
	}
}
