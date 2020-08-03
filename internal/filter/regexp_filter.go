package filter

import "regexp"

// RegexpFilter is a filter to decide whether should send a line to the webhook endpoint based on a regular expression.
type RegexpFilter struct {
	Regexp *regexp.Regexp
}

// Match returns the line is matched with given regular expression.
func (f *RegexpFilter) Match(line string) bool {
	return f.Regexp.MatchString(line)
}
