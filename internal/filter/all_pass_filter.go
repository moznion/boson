package filter

// AllPassFilter is a filter that passes all of the lines.
type AllPassFilter struct {
}

// Match always returns true.
func (f *AllPassFilter) Match(line string) bool {
	return true
}
