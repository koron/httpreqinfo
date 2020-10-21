package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"syscall"

	"github.com/koron-go/ctxsrv"
	"github.com/koron-go/sigctx"
)

var (
	addr   string
	silent bool
	dump   bool
	out    = log.New(os.Stderr, "", log.LstdFlags)
)

func main() {
	flag.StringVar(&addr, "addr", ":8000", `listen address`)
	flag.BoolVar(&silent, "silent", false, `suppress any outputs`)
	flag.BoolVar(&dump, "dump", false, `dump received requests`)
	flag.Parse()
	if silent {
		out = log.New(ioutil.Discard, "", 0)
	}
	out.Printf("listen on %s", addr)
	srv := &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(handle),
	}
	cfg := ctxsrv.HTTP(srv).
		WithDoneServer(func() {
			log.Println("done server")
		}).
		WithDoneContext(func() {
			log.Println("done context")
		})
	ctx, cancel := sigctx.WithCancelSignal(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer cancel()
	err := ctxsrv.Serve(ctx, *cfg)
	if err != nil {
		log.Fatal(err)
	}
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
	if dump {
		b, err := httputil.DumpRequest(r, false)
		if err != nil {
			log.Printf("failed to dump: %s", err)
		} else {
			log.Print("dump:{{{\n", string(b), "}}}:dump")
		}
	}
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
	out.Print(string(b))
}
