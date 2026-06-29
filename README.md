# robotstxt
A library implementing the Robot Exclusion Protocol (REP), as defined by [RFC 9309](https://www.rfc-editor.org/rfc/rfc9309.html). Its primary purpose is to serve as a parser and path tester for clients (e.g. crawlers), meaning groups not relevant for your [product token](https://www.rfc-editor.org/rfc/rfc9309.html#name-the-user-agent-line) are silently dropped. Additionally, there is "encoding" functionality available for storing the relevant ruleset data in a much more efficient manner than storing the raw robots.txt file (sizes will often be reduced by roughly 50%). Throughout the library there is a focus on simplicity and efficiency - benchmarks are included below.

## Features not yet implemented
- Crawl-Duration rules
- Sitemap rules
- Parsing from io.Readers

## Usage
### Parsing data and testing against paths
The following encompasses the most important operations the library provides: parsing robots.txt files and testing different paths against the output ruleset. A return value of `true` means your client is free to visit the page.

```go
package main

import "github.com/mvdahl/robotstxt"

func main() {
	fileData := []byte(`
# This one does not apply to us
User-Agent: some-other-crawler
Disallow: /

# This one is relevant for us
User-Agent: my-crawler
Disallow: /users
Allow: /users/admin

# Rules applying to everyone
User-Agent: *
Disallow: /humans-only
`)

	productToken := "my-crawler"

	ruleset := robotstxt.Parse(fileData, productToken)

	println(ruleset.Test("/"))            // true
	println(ruleset.Test("/users/joe"))   // false
	println(ruleset.Test("/users/admin")) // true
	println(ruleset.Test("/humans-only")) // false
}
```

Note: Since the parser only concerns itself with the rules you will actually be using, only a single `Ruleset` object is returned, which is comprised of a simple allow and disallow list.

### Encoding and decoding
Another important thing you might need to do is store fetched robots.txt files for later use. Here is an example of how you can do so:

```go
package main

import "github.com/mvdahl/robotstxt"

func main() {
	fileData := []byte(`
# Rules applying to everyone
User-Agent: *
Disallow: /humans-only
`)

	productToken := "my-crawler"

	ruleset := robotstxt.Parse(fileData, productToken)

	encodedRules, err := ruleset.Encode() // []byte containing encoded rules
	if err != nil {
		panic(err)
	}

    // save the rules to a file, DB, etc.

	sameRuleset, err := robotstxt.DecodeRuleset(encodedRules)
	if err != nil {
		panic(err)
	}

	println(sameRuleset.Test("/humans-only")) // false
}
```

## Benchmarks
To give a sort of reference point for the benchmarks a similar library is used, which you can find [here](https://github.com/temoto/robotstxt). This library turns out to be much faster, perhaps owing to a simpler or more efficient design. The benchmarks were made using `go test -bench=.` - results from running the command on my own system are as follows, using the google.com robots.txt file as input data:

```
goos: linux
goarch: amd64
pkg: robotstxt
cpu: AMD Ryzen 9 9950X 16-Core Processor            
BenchmarkEncode-32               1336891               884.7 ns/op
BenchmarkDecode-32                113852             10629 ns/op
BenchmarkOwnParse-32               62268             19031 ns/op
BenchmarkTemotoParse-32             9888            114070 ns/op
BenchmarkOwnPathTests-32          1644067               725.8 ns/op
BenchmarkTemotoPathTests-32        795858              1521 ns/op
```

The parsing this library provides turns out to be around 5-6x faster on average, while path testing is around 2x faster. Additionally, decoding previously parsed data is around twice as fast as parsing the original file.
