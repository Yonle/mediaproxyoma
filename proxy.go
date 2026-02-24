package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

var ua = "mediaproxyoma [https://github.com/Yonle/mediaproxyoma] - v0.3"

func init() {
	if ua_n, e := os.LookupEnv("USER_AGENT"); e {
		ua = ua_n
	}
}

var hc = http.Client{
	Timeout: 30 * time.Second,
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
	req.Header.Set("User-Agent", ua)

	return hc.Do(req)
}

func buildUrl(upstr string, isStatic bool) string {
	var anim int
	if isStatic {
		anim = 0
	} else {
		anim = 1
	}
	return fmt.Sprintf("%s?url=%s&bw=0&l=40&nr=1&a=%d", proxyhost, url.QueryEscape(upstr), anim)
}
