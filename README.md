# boson

A command line tool to send lines from through STDIN to the webhook endpoint.

## Usage

```
Usage of boson:
  -body string
        HTTP body for webhook request
  -every-line
        Run with every-line mode
  -filter-regexp string
        Regexp for line filtering, run with regexp-filter mode
  -header Content-Type: application/json
        HTTP header for webhook (example: Content-Type: application/json). It replaces "{{ line }}" token with got line string
  -http-method string
        HTTP method for webhook request (default "POST")
  -timeout-sec int
        HTTP timeout for webhook request (default 0, i.e. no-timeout)
  -url string
        URL for webhook endpoint. It replaces "{{ line }}" token with got line string
  -url-encode-body-replacement
        Encode the replacement of the body (i.e. the contents of "{{ line }}") of webhook request with url (percent) encoding; for "application/x-www-form-urlencoded"
```

### Example: Send every line to the webhook endpoint

```
tail -F something.log | boson --every-line --header "Content-Type: application/json" --header "X-Test: Heigh-Ho" --url http://localhost:8080 --body '{"message":"{{ line }}"}'
```

### Example: Send only matched line with a regular expression to the webhook endpoint

```
tail -F something.log | boson --filter-regexp="^\[error].+" --header "Content-Type: application/json" --header "X-Test: Heigh-Ho" --url http://localhost:8080 --body '{"message":"{{ line }}"}'
```

## Author

moznion (<moznion@gmail.com>)

