package boson

import (
	"bufio"
	"errors"
	"io"
	"log"
	"regexp"

	"github.com/moznion/boson/internal/filter"
	"github.com/moznion/boson/webhook"
)

// Opt is the struct that represents the option for this application.
type Opt struct {
	EveryLine     bool
	FilterRegexp  string
	WebhookSender webhook.Sender
}

// NewOpt returns the new instance of Opt struct.
func NewOpt(everyLine bool, filterRegexp string, webhookSender webhook.Sender) (*Opt, error) {
	if everyLine && filterRegexp != "" {
		return nil, errors.New("every-line mode and regexp-filter mode are exclusive, please specify either one")
	}
	return &Opt{
		EveryLine:     everyLine,
		FilterRegexp:  filterRegexp,
		WebhookSender: webhookSender,
	}, nil
}

// GetFilter returns the filter instance from the option.
func (o *Opt) GetFilter() filter.Filter {
	if o.EveryLine {
		return &filter.AllPassFilter{}
	}

	return &filter.RegexpFilter{
		Regexp: regexp.MustCompile(o.FilterRegexp),
	}
}

// Run is the entry point of this application.
func Run(reader io.Reader, opt *Opt, dieOnError bool) error {
	var stdinScanner = bufio.NewScanner(reader)

	lineFilter := opt.GetFilter()

	for stdinScanner.Scan() {
		line := stdinScanner.Text()
		if parts := lineFilter.Find(line); len(parts) > 0 {
			if err := opt.WebhookSender.Send(parts); err != nil {
				log.Printf("[error] %s", err)
				if dieOnError {
					return err
				}
			}
		}
	}

	return nil
}
