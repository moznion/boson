package webhook

import (
	"fmt"
	"net/url"
	"regexp"
)

// Sender is an interface that has the responsibility to send request to webhook endpoint.
type Sender interface {
	// Send sends the request to webhook endpoint.
	Send(parts []string) error
}

const linePlaceholder = "{{ line }}"

var linePlaceholderRegexp = regexp.MustCompile(linePlaceholder)

func replacePlaceholder(src string, parts []string, urlEncodeBodyReplacement bool) string {
	// NOTE: performance of this code is not effective because it scans full-text multiple times.

	for i, part := range parts {
		var placeholderRegexp *regexp.Regexp
		if i == 0 {
			placeholderRegexp = linePlaceholderRegexp
		} else {
			placeholderRegexp = regexp.MustCompile(fmt.Sprintf("{{ \\$%d }}", i))
		}

		if urlEncodeBodyReplacement {
			part = url.QueryEscape(part)
		}

		src = placeholderRegexp.ReplaceAllString(src, part)
	}
	return src
}
