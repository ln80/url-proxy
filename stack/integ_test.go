//go:build integ

package stack_test

import (
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestProxy(t *testing.T) {
	proxyBaseURL := os.Getenv("PROXY_FUNCTION_URL")
	if proxyBaseURL == "" {
		t.Fatal("invalid proxy base URL")
	}

	src := "https://go.dev/images/go-logo-white.svg"

	url := proxyBaseURL + "/proxy?url=" + url.QueryEscape(src)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if status := resp.StatusCode; status != http.StatusOK {
		t.Fatal("invalid status code", status)
	}
}
