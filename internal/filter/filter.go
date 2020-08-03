package filter

// Filter is an interface that has the responsibility to filter the input to decide whether should send the line to the webhook endpoint.
type Filter interface {
	// Match returns the line is matched with the condition.
	// If this returns true, it has to send the line to webhook endpoint.
	Match(line string) bool
}
