package main

import (
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambdaurl"
	urlproxy "github.com/ln80/url-proxy"
)

func main() {
	proxy := urlproxy.NewProxy(
		nil,
		func(pc *urlproxy.ProxyConfig) {
			pc.DenyHosts = []string{"localhost"}
		},
	)
	mux := http.NewServeMux()
	mux.Handle("/proxy", proxy)

	if os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
		lambdaurl.Start(mux)
		return
	}

	server := &http.Server{
		Addr:         "localhost:9001",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	server.ListenAndServe()
}
