# boson

A command line tool to send lines from through STDIN to the webhook endpoint.

## Usage

```
Usage of boson:
  -body string
        HTTP body for webhook request
  -dry-run
        Run this application with dry-run (i.e. it doesn't send any requests to the webhook endpoint)
  -every-line
        Run with every-line mode. This mode sends all of lines to the webhook endpoint
  -filter-regexp string
        Regexp for line filtering, run with regexp-filter mode; it allows using regexp group and it can use matched groups as "{{ $1 }}", "{{ $2 }}"...
  -header value
        HTTP header for webhook (example: "Content-Type: application/json"). It replaces "{{ line }}" token with got line string and "{{ $N }}" token with correspond group
  -http-method string
        HTTP method for webhook request (default "POST")
  -timeout-sec int
        HTTP timeout for webhook request (default 0, i.e. no-timeout)
  -url string
        URL for webhook endpoint. It replaces "{{ line }}" token with got line string and "{{ $N }}" token with correspond group
  -url-encode-body-replacement
        Encode the replacement of the body (i.e. the contents of "{{ line }}" and "{{ $N }}") of webhook request with url (percent) encoding; for "application/x-www-form-urlencoded"
```

### Example: Send every line to the webhook endpoint

```
tail -F something.log | boson --every-line --header "Content-Type: application/json" --header "X-Test: Heigh-Ho" --url http://localhost:8080 --body '{"message":"{{ line }}"}'
```

### Example: Send only matched line with a regular expression to the webhook endpoint

```
tail -F something.log | boson --filter-regexp="^\[error].+" --header "Content-Type: application/json" --header "X-Test: Heigh-Ho" --url http://localhost:8080 --body '{"message":"{{ line }}"}'
```

### Example: Send only matched line with a regular expression with grouping to the webhook endpoint

```
echo "2020/08/05 18:11:25 [error] hello" | boson --filter-regexp="^(\d+/\d+/\d+ \d+:\d+:\d+) \[error] (.+)" --header "Content-Type: application/json" --header "X-Test: Heigh-Ho" --url http://localhost:8080 --body '{"timestamp":"{{ $1 }}","message":"{{ $2 }}"}'
```

## Author

moznion (<moznion@gmail.com>)

