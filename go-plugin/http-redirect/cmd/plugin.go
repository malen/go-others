package main

import (
	"encoding/json"
	"net/http"
	"os"
	"plugin"

	"aoisoft.net/http-redirect/protocol"
)

var pluginPathList []string
var pluginList []*protocol.HttpRedirectPlugin

func LoadConfig() {
	f, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal(f, &pluginPathList)
}

func LoadPlugins() {
	pluginList = make([]*protocol.HttpRedirectPlugin, 0, len(pluginPathList))
	for _, p := range pluginPathList {
		plg, err := plugin.Open(p)
		if err != nil {
			panic(err)
		}
		// pluginList = append(pluginList, plg)
		v, err := plg.Lookup("Plugin")
		if err != nil {
			panic(err)
		}

		castV, ok := v.(protocol.HttpRedirectPlugin)
		if !ok {
			panic("Could not cast plugin")
		}
		pluginList = append(pluginList, &castV)
	}
}

func init() {
	LoadConfig()
	LoadPlugins()
}

func PreRequestHook(req *http.Request) {
	for _, plg := range pluginList {
		(*plg).PreRequestHook((req))
	}
}
