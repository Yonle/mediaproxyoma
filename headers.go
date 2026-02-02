package main

import "net/http"

// client -> upstream
var clientHeaders = []string{
	"Range",
	"If-None-Match",
	"If-Modified-Since",
	"Accept",
	"Accept-Encoding",
}

// upstream -> client
var upstreamHeaders = []string{
	"Content-Type",
	//"Content-Length", <-- already dealed via http.Response.ContentLength
	"Content-Encoding",
	"Content-Range",
	"Accept-Ranges",

	"ETag",
	"Last-Modified",
	"Cache-Control",
	"Expires",
	"Vary",

	"Content-Disposition",
}

func copyUpstreamHeaders(dst http.Header, src http.Header) {
	for _, k := range upstreamHeaders {
		if vv, ok := src[k]; ok {
			for _, v := range vv {
				dst.Add(k, v)
			}
		}
	}
}

func copyClientHeaders(dst http.Header, src http.Header) {
	for _, k := range clientHeaders {
		if vv, ok := src[k]; ok {
			for _, v := range vv {
				dst.Add(k, v)
			}
		}
	}
}
