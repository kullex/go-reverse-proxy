package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

func main() {
	var target string
	var port int

	defPort, _ := strconv.Atoi(os.Getenv("GO_RP_PORT"))

	flag.StringVar(&target, "target", os.Getenv("GO_RP_TARGET"), "target address")
	flag.IntVar(&port, "port", defPort, "service port number")

	flag.Parse()

	remote, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.URL)
			r.Host = remote.Host
			p.ServeHTTP(w, r)
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	http.HandleFunc("/", handler(proxy))
	log.Printf("Server Proxy %s at %d\n", target, port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}
