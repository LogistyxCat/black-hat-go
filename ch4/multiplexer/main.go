package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

var (
	hostProxy = make(map[string]string)
	proxies   = make(map[string]*httputil.ReverseProxy)
)

func init() {
	hostProxy["notevil1.local"] = "http://10.0.1.20:10080"
	hostProxy["notevil2.local"] = "http://10.0.1.20:20080"

	for key, value := range hostProxy {
		remote, err := url.Parse(value)
		if err != nil {
			log.Fatal("Unable to parse proxy target")
		}
		proxies[key] = httputil.NewSingleHostReverseProxy(remote)
	}
}

func main() {
	r := mux.NewRouter()
	for host, proxy := range proxies {
		r.Host(host).Handler(proxy)
	}
	log.Fatal(http.ListenAndServe(":80", r))
}
