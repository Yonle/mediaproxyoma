package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var proxypath = "/proxy/"
var proxypreviewpath = "/proxy/preview/"
var proxyhost = os.Getenv("BWHERO_HOST")
var listen = os.Getenv("LISTEN")
var allowOrigin = "*"

func init() {
	if ao, e := os.LookupEnv("ALLOW_ORIGIN"); e {
		allowOrigin = ao
	}
}

func main() {
	if proxyhost == "" || listen == "" {
		log.Println("please set LISTEN and BWHERO_HOST environment variable before running.")
		log.Println("  LISTEN value could be 0.0.0.0:2000")
		log.Println("  BWHERO_HOST could be http://localhost:8080/")
		return
	}

	http.HandleFunc(proxypath, func(w http.ResponseWriter, r *http.Request) {
		url, err := getUrl(r.URL.Path[len(proxypath):])

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filename := r.URL.Query().Get("name")

		log.Println("proxying", url)
		resp, err := proxy(r.Context(), r, url)
		if err != nil {
			log.Println("give up on fetching to upstream:", err)
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		sex(w, resp.Body, resp, filename)
	})

	http.HandleFunc(proxypreviewpath, func(w http.ResponseWriter, r *http.Request) {
		url, err := getUrl(r.URL.Path[len(proxypreviewpath):])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		q := r.URL.Query()
		isStatic := q.Get("static") == "true"
		filename := q.Get("name")

		log.Println("previewing", url)
		resp, err := proxy(r.Context(), r, buildUrl(url, isStatic))
		if err != nil {
			log.Println("can't fetch from bwhero:", err)

			resp, err := proxy(r.Context(), r, url)
			if err != nil {
				log.Println("give up on fetching to upstream:", err)
				http.Error(w, "clusterbad", http.StatusBadGateway)
				return
			}

			sex(w, resp.Body, resp, filename)
			return
		}

		sex(w, resp.Body, resp, filename)
	})

	log.Println("Listening at", listen)
	if err := http.ListenAndServe(listen, nil); err != nil {
		panic(err)
	}
}

func sex(hole http.ResponseWriter, dih io.ReadCloser, sperm *http.Response, childname string) {
	defer dih.Close()

	h := hole.Header()
	h.Set("Access-Control-Allow-Credentials", "true")
	h.Set("Access-Control-Allow-Origin", allowOrigin)

	h.Set("Cache-Control", "public, max-age=604800, immutable")

	if sperm.ContentLength > 0 {
		h.Set("Content-Length", strconv.FormatInt(sperm.ContentLength, 10))
	}

	copyUpstreamHeaders(h, sperm.Header)

	if len(childname) > 0 {
		h.Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", childname))
	}

	hole.WriteHeader(sperm.StatusCode)

	io.Copy(hole, dih)
}
