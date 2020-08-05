package boson

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/moznion/boson/webhook"
	"github.com/stretchr/testify/assert"
)

var httpHeaders = http.Header{
	"Content-Type": []string{"application/json"},
	"X-Boson-Test": []string{"for testing"},
}

func launchHTTPListenerWithLineTesters(t *testing.T, expects []string) (int, chan struct{}) {
	finished := make(chan struct{}, 1)

	tcpListener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	port := tcpListener.Addr().(*net.TCPAddr).Port

	cnt := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		for key, value := range httpHeaders {
			assert.Equal(t, r.Header[key], value)
		}

		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)

		assert.Equal(t, expects[cnt], string(body))

		cnt++
		if cnt == len(expects) {
			finished <- struct{}{}
		}
	})

	go func() {
		err = http.Serve(tcpListener, mux)
		assert.NoError(t, err)
	}()

	return port, finished
}

func TestEveryLineMode(t *testing.T) {
	expects := []string{
		`{"message":"foo"}`,
		`{"message":"bar"}`,
		`{"message":"buz"}`,
	}

	port, finished := launchHTTPListenerWithLineTesters(t, expects)

	webhookSender, err := webhook.NewHTTPSender(
		"POST",
		fmt.Sprintf("http://localhost:%d", port),
		httpHeaders,
		`{"message":"{{ line }}"}`,
		1*time.Second,
		false,
		false,
	)
	if err != nil {
		t.Fatal(err)
	}

	mockedStdin := new(bytes.Buffer)
	mockedStdin.Write([]byte("foo\nbar\nbuz"))

	err = Run(mockedStdin, &Opt{
		EveryLine:     true,
		FilterRegexp:  "",
		WebhookSender: webhookSender,
	}, true)
	assert.NoError(t, err)

	<-finished
}

func TestRegexpFilterMode(t *testing.T) {
	expects := []string{
		`{"message":"barbar"}`,
		`{"message":"buzbuz"}`,
	}

	port, finished := launchHTTPListenerWithLineTesters(t, expects)

	webhookSender, err := webhook.NewHTTPSender(
		"POST",
		fmt.Sprintf("http://localhost:%d", port),
		httpHeaders,
		`{"message":"{{ line }}{{ $1 }}{{ $2 }}"}`,
		1*time.Second,
		false,
		false,
	)
	if err != nil {
		t.Fatal(err)
	}

	mockedStdin := new(bytes.Buffer)
	mockedStdin.Write([]byte("foo\nbar\nbuz"))

	err = Run(mockedStdin, &Opt{
		EveryLine:     false,
		FilterRegexp:  "(b)(.+)",
		WebhookSender: webhookSender,
	}, true)
	assert.NoError(t, err)

	<-finished
}

func TestURLEncodingForBodyReplacement(t *testing.T) {
	expects := []string{
		`message=%5Berror%5D+foo`,
	}

	port, finished := launchHTTPListenerWithLineTesters(t, expects)

	webhookSender, err := webhook.NewHTTPSender(
		"POST",
		fmt.Sprintf("http://localhost:%d", port),
		httpHeaders,
		`message={{ line }}`,
		1*time.Second,
		true,
		false,
	)
	if err != nil {
		t.Fatal(err)
	}

	mockedStdin := new(bytes.Buffer)
	mockedStdin.Write([]byte("[error] foo\n"))

	err = Run(mockedStdin, &Opt{
		EveryLine:     true,
		FilterRegexp:  "",
		WebhookSender: webhookSender,
	}, true)
	assert.NoError(t, err)

	<-finished
}

func TestDryRunMode(t *testing.T) {
	webhookSender, err := webhook.NewHTTPSender(
		"POST",
		"http://localhost",
		httpHeaders,
		`{"message":"{{ line }}"}`,
		1*time.Second,
		false,
		true,
	)
	if err != nil {
		t.Fatal(err)
	}

	mockedStdin := new(bytes.Buffer)
	mockedStdin.Write([]byte("foo"))

	err = Run(mockedStdin, &Opt{
		EveryLine:     true,
		FilterRegexp:  "",
		WebhookSender: webhookSender,
	}, true)
	assert.NoError(t, err)
}
