package main

import (
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

var addr string

func main() {
	flag.StringVar(&addr, "addr", ":8000", `listen address`)
	flag.Parse()
	log.Printf("listen on %s", addr)
	log.Fatal(http.ListenAndServe(addr, http.HandlerFunc(handle)))
}

type Request struct {
	Method           string                 `json:"method"`
	URL              string                 `json:"url"`
	Proto            string                 `json:"proto"`
	Header           map[string]interface{} `json:"header"`
	ContentLength    int64                  `json:"content-length"`
	BodySize         *int64                 `json:"body-size,omitempty"`
	BodyHash         *string                `json:"body-hash,omitempty"`
	TransferEncoding []string               `json:"transfer-encoding,omitempty"`
	Host             string                 `json:"host,omitempty"`
}

func procHeader(src http.Header) map[string]interface{} {
	dst := map[string]interface{}{}
	for k, v := range src {
		switch len(v) {
		case 0:
			dst[k] = nil
		case 1:
			dst[k] = v[0]
		default:
			dst[k] = v
		}
	}
	return dst
}

func procBody(src io.ReadCloser) (hash string, size int64, err error) {
	h := md5.New()
	size, err = io.Copy(h, src)
	if err != nil {
		return "", 0, err
	}
	hash = fmt.Sprintf("%x", h.Sum(nil))
	return hash, size, nil
}

func handle(w http.ResponseWriter, r *http.Request) {
	v := &Request{
		Method:           r.Method,
		URL:              r.URL.String(),
		Proto:            r.Proto,
		Header:           procHeader(r.Header),
		ContentLength:    r.ContentLength,
		TransferEncoding: r.TransferEncoding,
		Host:             r.Host,
	}
	hash, size, err := procBody(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if size != 0 || r.ContentLength != 0 {
		v.BodyHash, v.BodySize = &hash, &size
	}

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	log.Print(string(b))
}
