package webhook

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// HTTPSender is a webhook client based on HTTP.
type HTTPSender struct {
	HTTPMethod               string
	URL                      string
	Headers                  http.Header
	Body                     string
	Client                   *http.Client
	URLEncodeBodyReplacement bool
	DryRun                   bool
}

// NewHTTPSender returns a HTTPSender instance.
func NewHTTPSender(httpMethod string, url string, headers http.Header, body string, timeout time.Duration, urlEncodeBodyReplacement bool, dryRun bool) (*HTTPSender, error) {
	if httpMethod == "" {
		return nil, errors.New("http method has to have some value, but that is empty")
	}
	if url == "" {
		return nil, errors.New("url has to have some value, but that is empty")
	}

	return &HTTPSender{
		HTTPMethod:               httpMethod,
		URL:                      url,
		Headers:                  headers,
		Body:                     body,
		URLEncodeBodyReplacement: urlEncodeBodyReplacement,
		DryRun:                   dryRun,
		Client: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

// Send sends the request to webhook endpoint over HTTP.
func (s *HTTPSender) Send(parts []string) error {
	webhookURL := replacePlaceholder(s.URL, parts, s.URLEncodeBodyReplacement)

	headers := make(http.Header)
	for key, headerValues := range s.Headers {
		headers[key] = func() []string {
			vs := make([]string, len(headerValues))
			for i, v := range headerValues {
				vs[i] = replacePlaceholder(v, parts, s.URLEncodeBodyReplacement)
			}
			return vs
		}()
	}

	body := replacePlaceholder(s.Body, parts, s.URLEncodeBodyReplacement)

	if s.DryRun {
		log.Printf(`[info] DRY-RUN MODE! %s request to %s; header = %v, body = "%s"`, s.HTTPMethod, webhookURL, headers, body)
		return nil
	}

	req, err := http.NewRequest(s.HTTPMethod, webhookURL, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header = headers

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}

	statusCode := resp.StatusCode
	if statusCode < 200 || statusCode >= 300 {
		return fmt.Errorf("failed to send data to the webhook destination; statusCode = %d", statusCode)
	}
	log.Printf(`[info] sent %s request to %s; statusCode = %d, body = "%s"`, s.HTTPMethod, webhookURL, statusCode, body)
	return nil
}
