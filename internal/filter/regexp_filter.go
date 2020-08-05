package filter

import "regexp"

// RegexpFilter is a filter to decide whether should send a line to the webhook endpoint based on a regular expression.
type RegexpFilter struct {
	Regexp *regexp.Regexp
}

// Find returns the sub-matched string slice for line according to given regular expression.
func (f *RegexpFilter) Find(line string) []string {
	return f.Regexp.FindStringSubmatch(line)
}
