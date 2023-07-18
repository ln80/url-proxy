package urlproxy

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type ProxyConfig struct {
	DenyHosts []string
	UserAgent string
}

type Proxy struct {
	cfg    ProxyConfig
	client *http.Client
}

var _ http.Handler = &Proxy{}

// NewProxy creates a new instance of Proxy http.Handler
func NewProxy(client *http.Client, opts ...func(*ProxyConfig)) *Proxy {
	cfg := ProxyConfig{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&cfg)
	}
	if cfg.DenyHosts == nil {
		cfg.DenyHosts = make([]string, 0)
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = "ln80/url-proxy - github.com/ln80/url-proxy"
	}

	if client == nil {
		// In some cases, using a cookie jar can prevent infinite redirect issues.
		jar, _ := cookiejar.New(nil)
		client = &http.Client{
			Timeout: 10 * time.Second,
			Jar:     jar,
			CheckRedirect: func(newreq *http.Request, via []*http.Request) error {
				if len(via) > 10 {
					return errors.New("too many redirects")
				}
				if hostMatches(newreq.URL, cfg.DenyHosts) {
					return errors.New("redirect to a denied host")
				}
				return nil
			},
		}
	}

	return &Proxy{
		client: client,
		cfg:    cfg,
	}
}

// ServeHTTP implements http.handler interface
func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	u, err := url.Parse(req.URL.Query().Get("url"))
	if err != nil || hostMatches(u, p.cfg.DenyHosts) {
		http.Error(w, "invalid source URL", http.StatusBadRequest)
		return
	}

	actualReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		http.Error(w, "invalid source URL", http.StatusBadRequest)
		return
	}

	// Privacy: Do not send referer to source
	// https://developer.mozilla.org/en-US/docs/Web/Security/Referer_header:_privacy_and_security_concerns
	copyHeader(actualReq.Header, req.Header, "Accept", "Accept-Encoding", "User-Agent")
	if agent := p.cfg.UserAgent; agent != "" {
		actualReq.Header.Add("User-Agent", agent)
	}

	resp, err := p.client.Do(actualReq)
	if err != nil {
		http.Error(w, "error fetching source URL", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Privacy: Do not copy Set-Cookie header to resp
	copyHeader(w.Header(), resp.Header, "Content-Type", "Content-Length")

	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("error copying response: %v", err)
	}
}

func copyHeader(dst, src http.Header, headerNames ...string) {
	for _, name := range headerNames {
		k := http.CanonicalHeaderKey(name)
		for _, v := range src[k] {
			dst.Add(k, v)
		}
	}
}

func hostMatches(u *url.URL, hosts []string) bool {
	if u == nil {
		return false
	}
	for _, h := range hosts {
		if h == u.Hostname() {
			return true
		}
	}
	return false
}
