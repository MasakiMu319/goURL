package utils

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"sort"
	"strings"
	"time"
)

type headers []string

func (h headers) String() string {
	var o []string
	for _, v := range h {
		o = append(o, "-H "+v)
	}
	return strings.Join(o, " ")
}

func (h *headers) Set(v string) error {
	*h = append(*h, v)
	return nil
}

func (h headers) Len() int      { return len(h) }
func (h headers) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h headers) Less(i, j int) bool {
	a, b := h[i], h[j]

	// server always sorts at the top
	if a == "Server" {
		return true
	}
	if b == "Server" {
		return false
	}

	endtoend := func(n string) bool {
		// https://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html#sec13.5.1
		switch n {
		case "Connection",
			"Keep-Alive",
			"Proxy-Authenticate",
			"Proxy-Authorization",
			"TE",
			"Trailers",
			"Transfer-Encoding",
			"Upgrade":
			return false
		default:
			return true
		}
	}

	x, y := endtoend(a), endtoend(b)
	if x == y {
		// both are of the same class
		return a < b
	}
	return x
}

var (
	// Command line flags
	HttpMethod       string // http method
	HttpResponseHead bool   // response head
	HttpConnectInfo  bool   // connect information

	ShowVersion bool	// show program version

	Version = "Dev"
)

func printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(color.Output, format, a...)
}

func grayscale(code color.Attribute) func(string, ...interface{}) string {
	return color.New(code + 232).SprintfFunc()
}

func VisitURL(url *url.URL) error {
	// TODO: data body have not set flag
	req, err := newRequest(HttpMethod, url, "")
	if err != nil {
		return err
	}

	// TODO: count time cost

	trace := &httptrace.ClientTrace{
		ConnectDone: func(net, addr string, err error) {
			if err != nil {
				log.Fatalf("unable to connect to host %v: %v", addr, err)
			}

			printf("\n%s%s\n", color.GreenString("Connected to "), color.CyanString(addr))
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(context.Background(), trace))

	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		MaxIdleConns: 100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
	}

	// TODO: choose IPv4 or IPv6

	switch url.Scheme {
	case "https":
		host, _, err := net.SplitHostPort(req.Host)
		if err != nil {
			host = req.Host
		}

		tr.TLSClientConfig = &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
		}
	}

	client := &http.Client{
		Transport: tr,
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.New(color.HiRedString("failed to read response:", err))
	}
	// Print SSL/TLS version which is used for connection
	connectedVia := "plaintext"
	if resp.TLS != nil {
		switch resp.TLS.Version {
		case tls.VersionTLS12:
			connectedVia = "TLSv1.2"
		case tls.VersionTLS13:
			connectedVia = "TLSv1.3"
		}
	}
	printf("\n%s %s\n", color.GreenString("Connected via"), color.CyanString("%s", connectedVia))

	resp.Body.Close()

	names := make([]string, 0, len(resp.Header))
	for k := range resp.Header {
		names = append(names, k)
	}
	sort.Sort(headers(names))
	for _, k := range names {
		printf("%s %s\n", grayscale(14)(k+":"), color.CyanString(strings.Join(resp.Header[k], ",")))
	}

	return nil
}

func newRequest(method string, url *url.URL, body string) (*http.Request, error) {
	req, err := http.NewRequest(method, url.String(), createBody(body))
	if err != nil {
		return nil, errors.New(color.HiRedString("Unable to create request:", err))
	}
	// TODO: add headers for request
	return req, nil
}

func createBody(body string) io.Reader {
	return strings.NewReader(body)
}