package robotstxt

import (
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	ruleset := Ruleset{
		allow:    parsePatterns([]string{"allow1", "allow2", "allow3", "allow-with-utf8-é"}),
		disallow: parsePatterns([]string{"disallow1", "disallow2"}),
	}

	encoded, err := ruleset.Encode()
	if err != nil {
		t.Fatal("failed to encode:", err.Error())
		return
	}

	decoded, err := DecodeRuleset(encoded)
	if err != nil {
		t.Fatal("failed to decode:", err.Error())
		return
	}

	if !equivelentRuleLists(ruleset.allow, decoded.allow) {
		t.Fatal("allow lists are not equivelent")
	}

	if !equivelentRuleLists(ruleset.disallow, decoded.disallow) {
		t.Fatal("disallow lists are not equivelent")
	}
}

func TestEncodingStorage(t *testing.T) {
	testRobotsFiles := [][]byte{
		exampleRobotsTxt,
		googleRobotsTxt,
		githubRobotsTxt,
	}

	testRobotsProducts := []string{
		exampleRobotsProduct,
		googleRobotsProduct,
		githubRobotsProduct,
	}

	for i, content := range testRobotsFiles {
		ruleset := Parse(content, testRobotsProducts[i])

		encoded, err := ruleset.Encode()
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("list %d encoded takes up %d bytes vs original size of %d", i+1, len(encoded), len(content))
	}
}

func BenchmarkEncode(b *testing.B) {
	ruleset := Parse(googleRobotsTxt, googleRobotsProduct)

	b.ResetTimer()

	for b.Loop() {
		ruleset.Encode()
	}
}

func BenchmarkDecode(b *testing.B) {
	ruleset := Parse(googleRobotsTxt, googleRobotsProduct)

	encoded, err := ruleset.Encode()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for b.Loop() {
		DecodeRuleset(encoded)
	}
}
