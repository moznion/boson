package filter

// AllPassFilter is a filter that passes all of the lines.
type AllPassFilter struct {
}

// Find always returns `[]string{line}`.
func (f *AllPassFilter) Find(line string) []string {
	return []string{line}
}
