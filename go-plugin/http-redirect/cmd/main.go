package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

var from int
var to string

func init() {
	flag.IntVar(&from, "from", 5555, "Local port to get requests")
	flag.StringVar(&to, "to", "", "Target server to redirect request to")
}

func main() {
	flag.Parse()
	listen()
}

type proxy struct{}

func listen() {
	p := &proxy{}
	srvr := http.Server{
		Addr:    fmt.Sprintf(":%d", from),
		Handler: p,
	}
	if err := srvr.ListenAndServe(); err != nil {
		slog.Error("Server is down", "Error", err)
	}
}

func (p *proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	PreRequestHook(req)
	req.RequestURI = ""
	req.URL.Host = to
	if req.TLS == nil {
		req.URL.Scheme = "http"
	} else {
		req.URL.Scheme = "https"
	}

	DropHoopHeader(&req.Header)

	SetProxyHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(rw, "Server Error: Redirect failed", http.StatusInternalServerError)
	}
	defer resp.Body.Close()
	DropHoopHeader(&req.Header)

	CopyHeader(rw.Header(), &resp.Header)
	rw.WriteHeader(resp.StatusCode)
	if _, err = io.Copy(rw, resp.Body); err != nil {
		slog.Error("Error writing response", "error", err)
	}
}

func CopyHeader(src http.Header, dst *http.Header) {
	for headingName, headingValues := range src {
		for _, value := range headingValues {
			dst.Add(headingName, value)
		}
	}
}

var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func DropHoopHeader(head *http.Header) {
	for _, header := range hopHeaders {
		head.Del(header)
	}
}

func SetProxyHeader(req *http.Request) {
	headerName := "X-Forwarded-for"
	target := to
	if prior, ok := req.Header[headerName]; ok {
		target = strings.Join(prior, ",") + "," + target
	}
	req.Header.Set(headerName, target)
}
