package main

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var proxypath = "/proxy/"
var proxypreviewpath = "/proxy/preview/"
var proxyhost = os.Getenv("BWHERO_HOST")
var listen = os.Getenv("LISTEN")

func main() {
	if proxyhost == "" || listen == "" {
		log.Println("please set LISTEN and BWHERO_HOST environment variable before running.")
		log.Println("  LISTEN value could be 0.0.0.0:2000")
		log.Println("  BWHERO_HOST could be http://localhost:8080/")
		return
	}

	http.HandleFunc(proxypath, func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path[len(proxypath):]
		sp := strings.SplitN(url, "/", 3)

		if len(sp) < 3 {
			http.Error(w, "bad", http.StatusBadRequest)
			return
		}

		decb, err := base64.RawURLEncoding.DecodeString(sp[1])
		if err != nil {
			http.Error(w, "bad", http.StatusBadRequest)
			return
		}

		url = string(decb)

		log.Println("proxying", url)
		resp, err := proxy(r.Context(), r, url)
		if err != nil {
			log.Println("give up on fetching to upstream:", err)
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		sex(w, resp.Body, resp)
	})

	http.HandleFunc(proxypreviewpath, func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path[len(proxypreviewpath):]
		sp := strings.SplitN(url, "/", 3)

		if len(sp) < 3 {
			http.Error(w, "bad", http.StatusBadRequest)
			return
		}

		decb, err := base64.RawURLEncoding.DecodeString(sp[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		url = string(decb)

		log.Println("previewing", url)
		resp, err := proxy(r.Context(), r, buildUrl(url))
		if err != nil {
			log.Println("can't fetch from bwhero:", err)
			resp, err := proxy(r.Context(), r, url)
			if err != nil {
				log.Println("give up on fetching to upstream:", err)
				http.Error(w, "clusterbad", http.StatusBadGateway)
				return
			}

			sex(w, resp.Body, resp)
			return
		}

		sex(w, resp.Body, resp)
	})

	log.Println("Listening at", listen)
	if err := http.ListenAndServe(listen, nil); err != nil {
		panic(err)
	}
}

func sex(hole http.ResponseWriter, dih io.ReadCloser, sperm *http.Response) {
	defer dih.Close()

	h := hole.Header()
	h.Set("Access-Control-Allow-Credentials", "true")
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Cache-Control", "public, max-age=604800, stale-while-revalidate=86400")

	rh := sperm.Header
	h.Set("Content-Type", rh.Get("Content-Type"))

	// workaround for videos
	if cr := rh.Get("Content-Range"); cr != "" {
		h.Set("Content-Range", cr)
	}
	if ar := rh.Get("Accept-Ranges"); ar != "" {
		h.Set("Accept-Ranges", ar)
	}
	if ce := rh.Get("Content-Encoding"); ce != "" {
		h.Set("Content-Encoding", ce)
	}
	if sperm.ContentLength > 0 {
		h.Set("Content-Length", strconv.FormatInt(sperm.ContentLength, 10))
	}

	hole.WriteHeader(sperm.StatusCode)

	io.Copy(hole, dih)
}
