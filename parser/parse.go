package parser

import (
	"errors"
	"github.com/fatih/color"
	"net/url"
	"strings"
)

func ParseURL(uri string) (*url.URL, error) {
	// we can only solve "*://url" and "//url".
	// if user types incorrect url unfortunately, like "///url",
	// it's a bad news for our program.
	// TODO: think about a solution for this.
	// one way to solve this is force check weather url is correct

	// If we choose to check all kinds of urls, it's too complex
	// and hard to think about all situations. Anyway it's a bad
	// practice.
	// TODO: Solution One
	// if !strings.Contains(uri, "://") && !strings.HasPrefix(uri, "//") {
	//	 uri = "//" + uri
	// }

	// we use defend coding to avoid complex situations.
	// TODO: Solution Two
	uri, err := checkURL(uri)
	if err != nil {
		return nil, err
	}

	url, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if url.Scheme == "" {
		url.Scheme = "http"
		if !strings.HasSuffix(url.Host, ":80") {
			url.Scheme += "s"
		}
	}
	return url, nil
}

// check input url format.
func checkURL(url string) (string, error) {
	// The correct url types is:
	// example.com
	// http://example.com
	// https://example.com
	// TODO: we can't check domain perfect
	if !strings.HasPrefix(url, "/") &&
		!strings.Contains(url, "//") &&
		!strings.Contains(url, "http") {
		return url, nil
	}
	if strings.Contains(url, "//") &&
		(strings.HasPrefix(url, "http:") ||
			strings.HasPrefix(url, "https:")) {
		return url, nil
	}
	return "", errors.New(color.HiRedString("URL using bad/illegal format or missing URL"))
}
