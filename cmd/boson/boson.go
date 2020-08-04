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
	var everyLine bool
	var filterRegexp string
	var webhookHTTPMethod string
	var webhookURL string
	var webhookHeadersOpt stringSlice
	var webhookBody string
	var webhookTimeoutSec int
	var urlEncodeBodyReplacement bool

	flag.BoolVar(&everyLine, "every-line", false, "Run with every-line mode")
	flag.StringVar(&filterRegexp, "filter-regexp", "", "Regexp for line filtering, run with regexp-filter mode")
	flag.StringVar(&webhookHTTPMethod, "http-method", "POST", "HTTP method for webhook request")
	flag.StringVar(&webhookURL, "url", "", "URL for webhook endpoint. It replaces \"{{ line }}\" token with got line string")
	flag.Var(&webhookHeadersOpt, "header", "HTTP header for webhook (example: `Content-Type: application/json`). It replaces \"{{ line }}\" token with got line string")
	flag.StringVar(&webhookBody, "body", "", "HTTP body for webhook request")
	flag.IntVar(&webhookTimeoutSec, "timeout-sec", 0, "HTTP timeout for webhook request (default 0, i.e. no-timeout)")
	flag.BoolVar(&urlEncodeBodyReplacement, "url-encode-body-replacement", false, "Encode the replacement of the body (i.e. the contents of \"{{ line }}\") of webhook request with url (percent) encoding; for \"application/x-www-form-urlencoded\"")

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
