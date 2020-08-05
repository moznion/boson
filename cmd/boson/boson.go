package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/moznion/boson"
	"github.com/moznion/boson/webhook"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	var (
		everyLine                bool
		filterRegexp             string
		webhookHTTPMethod        string
		webhookURL               string
		webhookHeadersOpt        stringSlice
		webhookBody              string
		webhookTimeoutSec        int
		urlEncodeBodyReplacement bool
		dryRun                   bool
	)

	flag.BoolVar(&everyLine, "every-line", false, "Run with every-line mode. This mode sends all of lines to the webhook endpoint")
	flag.StringVar(&filterRegexp, "filter-regexp", "", `Regexp for line filtering, run with regexp-filter mode; it allows using regexp group and it can use matched groups as "{{ $1 }}", "{{ $2 }}"...`)
	flag.StringVar(&webhookHTTPMethod, "http-method", "POST", "HTTP method for webhook request")
	flag.StringVar(&webhookURL, "url", "", `URL for webhook endpoint. It replaces "{{ line }}" token with got line string and "{{ $N }}" token with correspond group`)
	flag.Var(&webhookHeadersOpt, "header", `HTTP header for webhook (example: "Content-Type: application/json"). It replaces "{{ line }}" token with got line string and "{{ $N }}" token with correspond group`)
	flag.StringVar(&webhookBody, "body", "", "HTTP body for webhook request")
	flag.IntVar(&webhookTimeoutSec, "timeout-sec", 0, "HTTP timeout for webhook request (default 0, i.e. no-timeout)")
	flag.BoolVar(&urlEncodeBodyReplacement, "url-encode-body-replacement", false, `Encode the replacement of the body (i.e. the contents of "{{ line }}" and "{{ $N }}") of webhook request with url (percent) encoding; for "application/x-www-form-urlencoded"`)
	flag.BoolVar(&dryRun, "dry-run", false, "Run this application with dry-run (i.e. it doesn't send any requests to the webhook endpoint)")

	flag.Parse()

	webhookHeaders := make(http.Header)
	for _, headerOpt := range webhookHeadersOpt {
		kv := strings.SplitN(headerOpt, ":", 2)
		if len(kv) < 2 {
			log.Fatal("[error] invalid header has come; it has to be colon-separated notation like `Content-Type: application/json`")
		}

		headerKey := strings.TrimSpace(kv[0])
		webhookHeaders[headerKey] = append(webhookHeaders[headerKey], strings.TrimSpace(kv[1]))
	}

	httpSender, err := webhook.NewHTTPSender(
		webhookHTTPMethod,
		webhookURL,
		webhookHeaders,
		webhookBody,
		time.Duration(webhookTimeoutSec)*time.Second,
		urlEncodeBodyReplacement,
		dryRun,
	)
	if err != nil {
		log.Fatalf("[error] %s", err)
	}

	opt, err := boson.NewOpt(everyLine, filterRegexp, httpSender)
	if err != nil {
		log.Fatalf("[error] %s", err)
	}
	boson.Run(opt)
}
