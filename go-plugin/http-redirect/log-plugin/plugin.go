package main

import (
	"log/slog"
	"net/http"
	"net/http/httputil"

	"aoisoft.net/http-redirect/protocol"
)

func logRequest(req *http.Request) {
	result, err := httputil.DumpRequest(req, true)
	if err != nil {
		slog.Error("Failed to pring request", "err", err)
	}
	slog.Info("Request sent:", "req", result)
}

func logRequestLikeCUrl(req *http.Request) {
	panic("Unimplemented!")
}

type PluginStr struct{}

// PreRequestHook implements protocol.HttpRedirectPlugin.
func (p PluginStr) PreRequestHook(req *http.Request) {
	logRequest(req)
}

var _ protocol.HttpRedirectPlugin = PluginStr{}
var Plugin = PluginStr{}

func main() {}
