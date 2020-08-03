package webhook

import "regexp"

// Sender is an interface that has the responsibility to send request to webhook endpoint.
type Sender interface {
	// Send sends the request to webhook endpoint.
	Send(line string) error
}

const placeholder = "{{ line }}"

var placeholderRegexp = regexp.MustCompile(placeholder)

func replacePlaceholder(src string, line string) string {
	return placeholderRegexp.ReplaceAllString(src, line)
}
