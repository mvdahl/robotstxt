package robotstxt

import (
	"testing"

	_ "embed"

	"github.com/temoto/robotstxt"
)

func TestParseExample(t *testing.T) {
	expectedRuleset := Ruleset{
		allow:    parsePatterns([]string{"/publications/", "/this/has/utf8/é/é", "/example/page.html", "/example/allowed.gif"}),
		disallow: parsePatterns([]string{"*.gif$", "/example/", "/"}),
	}

	ruleset := Parse(exampleRobotsTxt, exampleRobotsProduct)

	if !equivelentRuleLists(expectedRuleset.allow, ruleset.allow) {
		t.Fatal("incorrect allow list")
	}

	if !equivelentRuleLists(expectedRuleset.disallow, ruleset.disallow) {
		t.Fatal("incorrect disallow list")
	}
}

func TestParseGoogle(t *testing.T) {
	ruleset := Parse(googleRobotsTxt, googleRobotsProduct)

	if len(ruleset.allow) != 64 {
		t.Fatal("unexpected amount of allows", len(ruleset.allow))
	}

	if len(ruleset.disallow) != 170 {
		t.Fatal("unexpected amount of disallows", len(ruleset.disallow))
	}
}

func TestParseGitHub(t *testing.T) {
	ruleset := Parse(githubRobotsTxt, githubRobotsProduct)

	if len(ruleset.allow) != 1 {
		t.Fatal("unexpected amount of allows", len(ruleset.allow))
	}

	if len(ruleset.disallow) != 66 {
		t.Fatal("unexpected amount of disallows", len(ruleset.disallow))
	}
}

func BenchmarkOwnParse(b *testing.B) {
	for b.Loop() {
		Parse(googleRobotsTxt, googleRobotsProduct)
	}
}

func BenchmarkTemotoParse(b *testing.B) {
	for b.Loop() {
		thing, _ := robotstxt.FromBytes(googleRobotsTxt)

		thing.FindGroup(googleRobotsProduct)
	}
}
