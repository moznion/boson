package filter

// Filter is an interface that has the responsibility to filter the input to decide whether should send the line to the webhook endpoint.
type Filter interface {
	// Find returns the string slice to send to the webhook.
	// If the length of this slice is greater than 0, it has to send the line to webhook endpoint.
	//
	// The first element of the slice is matched entire line, and trailing elements are grouped part.
	// e.g. [{{ line }}, {{ $1 }}, {{ $2 }}, ...]
	Find(line string) []string
}
