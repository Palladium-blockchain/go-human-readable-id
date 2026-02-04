package words

import (
	_ "embed"
	"strings"
)

//go:embed dict/adj.txt
var adjData string

//go:embed dict/noun.txt
var nounData string

//go:embed dict/verb.txt
var verbData string

var (
	AdjWords  []string
	NounWords []string
	VerbWords []string
)

func init() {
	AdjWords = parseLines(adjData)
	NounWords = parseLines(nounData)
	VerbWords = parseLines(verbData)
}

func parseLines(s string) []string {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		out = append(out, l)
	}
	return out
}
