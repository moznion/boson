package webhook

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// HTTPSender is a webhook client based on HTTP.
type HTTPSender struct {
	HTTPMethod string
	URL        string
	Headers    http.Header
	Body       string
	Client     *http.Client
}

// NewHTTPSender returns a HTTPSender instance.
func NewHTTPSender(httpMethod string, url string, headers http.Header, body string, timeout time.Duration) (*HTTPSender, error) {
	if httpMethod == "" {
		return nil, errors.New("http method has to have some value, but that is empty")
	}
	if url == "" {
		return nil, errors.New("url has to have some value, but that is empty")
	}

	return &HTTPSender{
		HTTPMethod: httpMethod,
		URL:        url,
		Headers:    headers,
		Body:       body,
		Client:     &http.Client{
			Timeout: timeout,
		},
	}, nil
}

// Send sends the request to webhook endpoint over HTTP.
func (s *HTTPSender) Send(line string) error {
	url := replacePlaceholder(s.URL, line)

	headers := make(http.Header)
	for key, headerValues := range s.Headers {
		headers[key] = func() []string {
			vs := make([]string, len(headerValues))
			for i, v := range headerValues {
				vs[i] = replacePlaceholder(v, line)
			}
			return vs
		}()
	}

	body := replacePlaceholder(s.Body, line)

	req, err := http.NewRequest(s.HTTPMethod, url, strings.NewReader(body))
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
	return nil
}
