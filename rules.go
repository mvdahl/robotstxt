package robotstxt

type rulePattern struct {
	Raw       []byte
	Sequences []string
	HasEnd    bool
}

type Ruleset struct {
	allow    []rulePattern
	disallow []rulePattern
}

func (r Ruleset) Test(path string) bool {
	mostSpecificDisallow, disallowed := findMostSpecificMatch(r.disallow, path, 0)

	if !disallowed {
		return true
	}

	mostSpecificAllow, allowed := findMostSpecificMatch(r.allow, path, len(mostSpecificDisallow.Raw))

	return allowed && len(mostSpecificAllow.Raw) >= len(mostSpecificDisallow.Raw)
}

func (r Ruleset) Encode() ([]byte, error) {
	return EncodeRuleset(r)
}
