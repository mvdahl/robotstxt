package robotstxt

import (
	"testing"
)

type matchTestCase struct {
	Pattern     string
	ValidPath   string
	InvalidPath string
}

func TestValidMatches(t *testing.T) {
	testStr := "/one/two/three"

	validPatterns := parsePatterns([]string{"/one", "/one/two", "/one/two/three", "/*/two/three", "/one/*/three", "/one/*/thr", "/one/two/*", "/one/two/three$", "/one/*/three$", "/one/two/*$", "/one/two/*three", "*/one/two/three", "*/two", "*/two/thr*"})

	for i, patt := range validPatterns {
		if !isMatch(patt, testStr) {
			t.Fatal("pattern", i+1, "should have been a match")
		}
	}
}

func TestInvalidMatches(t *testing.T) {
	testStr := "/one/two/three"

	invalidPatterns := parsePatterns([]string{"/four", "/one/three", "/one/*/four", "/one$", "/one/two$", "/one/two/thr$", "/$", "$", "one"})

	for i, patt := range invalidPatterns {
		if isMatch(patt, testStr) {
			t.Fatal("pattern", i+1, "should not have been a match")
		}
	}
}

func TestRealWorldMatches(t *testing.T) {
	cases := []matchTestCase{
		{
			Pattern:     "/?gws_rd=ssl$",
			ValidPath:   "/?gws_rd=ssl",
			InvalidPath: "/?gws_rd=ssl&thing=that",
		},
		{
			Pattern:     "/?hl=*&gws_rd=ssl$",
			ValidPath:   "/?hl=hi&gws_rd=ssl",
			InvalidPath: "/?gws_rd=ssl",
		},
		{
			Pattern:     "*/tarball/",
			ValidPath:   "/tarball/yeah",
			InvalidPath: "/tarball",
		},
		{
			Pattern:     "/*setup_organization=*",
			ValidPath:   "/somesetup_organization=",
			InvalidPath: "/setup_organization",
		},
	}

	for i, c := range cases {
		parsed := parsePattern([]byte(c.Pattern))

		if !isMatch(parsed, c.ValidPath) {
			t.Fatal("valid path in case", i+1, "did not match")
		}

		if isMatch(parsed, c.InvalidPath) {
			t.Fatal("invalid path in case", i+1, "matched")
		}
	}
}
