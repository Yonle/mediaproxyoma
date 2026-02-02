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
	req, err := http.NewRequestWithContext(ctx, r.Method, origin_url, nil)
	if err != nil {
		return nil, err
	}

	copyClientHeaders(req.Header, r.Header)
	req.Header.Set("User-Agent", "mediaproxyoma [https://github.com/Yonle/mediaproxyoma] - v0.2-dev")

	return hc.Do(req)
}

func buildUrl(upstr string) string {
	return fmt.Sprintf("%s?url=%s&bw=0&l=20", proxyhost, url.QueryEscape(upstr))
}
