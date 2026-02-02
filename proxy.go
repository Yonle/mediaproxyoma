package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var hc = http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		DisableCompression: true,
	},
}

func proxy(ctx context.Context, r *http.Request, origin_url string) (resp *http.Response, err error) {
	clientHeader := r.Header.Clone()

	clientHeader.Set("User-Agent", "mediaproxyoma - v0.1-dev")

	// anonymize & workarounds
	clientHeader.Del("x-forwarded-for")

	req, err := http.NewRequestWithContext(ctx, r.Method, origin_url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = clientHeader

	return hc.Do(req)
}

func buildUrl(upstr string) string {
	return fmt.Sprintf("%s?url=%s&bw=0&l=20", proxyhost, url.QueryEscape(upstr))
}
