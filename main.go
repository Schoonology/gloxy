package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

/**
 * Empty transport type.
 */
type GloxyTransport struct{}

/**
 * The crux of the logging implementation: a Transport.RoundTrips that logs all
 * requests made of it.
 */
func (self *GloxyTransport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	reqDump, _ := httputil.DumpRequest(req, true)

	res, err = http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return
	}

	resDump, _ := httputil.DumpResponse(res, true)
	fmt.Printf("\n---\n\n")
	fmt.Printf("%s > %s:\n%s", req.RemoteAddr, flag.Arg(0), reqDump)
	fmt.Printf("%s < %s:\n%s", flag.Arg(0), req.RemoteAddr, resDump)

	return
}

/**
 * Creates a new reverse proxy associated with our logging Transport.
 */
func NewGloxy(rawurl string) *httputil.ReverseProxy {
	target, err := url.Parse(rawurl)
	if err != nil {
		fmt.Printf("Bad target URL: %v", rawurl)
		os.Exit(1)
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &GloxyTransport{}
	return proxy
}

/**
 * Parse command-line flags and start the proxy server.
 */
func main() {
	port := flag.Int("port", 8080, "The port to listen on.")
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Listening on :%d and proxying to %s...", *port, flag.Arg(0))
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), NewGloxy(flag.Arg(0)))
	if err != nil {
		fmt.Printf("Failed to Listen with %v", err)
		os.Exit(1)
	}
}
