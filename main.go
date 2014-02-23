package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

/**
 * Empty transport type.
 */
type GloxyTransport struct{}

var PrintableTypes map[string]bool = map[string]bool{
	"application/atom+xml":              true,
	"application/ecmascript":            true,
	"application/json":                  true,
	"application/javascript":            true,
	"application/rdf+xml":               true,
	"application/rss+xml":               true,
	"application/soap+xml":              true,
	"application/xhtml+xml":             true,
	"application/xml":                   true,
	"application/xml-dtd":               true,
	"application/x-www-form-urlencoded": true,
	"text/css":                          true,
	"text/csv":                          true,
	"text/html":                         true,
	"text/javascript":                   true,
	"text/plain":                        true,
	"text/vcard":                        true,
	"text/xml":                          true,
}

/**
 * IsPrintable returns true if the Header indicates a printable Request or
 * Response.
 */
func IsPrintable(header http.Header) bool {
	mimeType := header.Get(http.CanonicalHeaderKey("content-type"))
	mimeType = strings.SplitN(mimeType, ";", 2)[0]
	mimeType = strings.TrimSpace(mimeType)

	if mimeType == "" {
		return true
	}

	return PrintableTypes[mimeType]
}

/**
 * The crux of the logging implementation: a Transport.RoundTrips that logs all
 * requests made of it.
 */
func (self *GloxyTransport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	res, err = http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return
	}

	body := IsPrintable(req.Header)
	reqDump, _ := httputil.DumpRequest(req, body)
	if !body {
		reqDump = append(reqDump, []byte("BINARY\n\n")...)
	}

	body = IsPrintable(res.Header)
	resDump, _ := httputil.DumpResponse(res, body)
	if !body {
		resDump = append(resDump, []byte("BINARY\n\n")...)
	}

	fmt.Printf("\n---\n\n")
	fmt.Printf("%s > %s:\n%s", req.RemoteAddr, flag.Arg(0), reqDump)
	fmt.Printf("%s < %s:\n%s", flag.Arg(0), req.RemoteAddr, resDump)
	fmt.Printf("\n")

	return
}

/**
 * Creates a new reverse proxy associated with our logging Transport.
 */
func NewGloxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &GloxyTransport{}
	return proxy
}

/**
 * Usage and help information
 */
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: gloxy [OPTIONS] URL\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.VisitAll(func(flag *flag.Flag) {
		fmt.Fprintf(os.Stderr, "  --%s (%v)  \t%s\n", flag.Name, flag.DefValue, flag.Usage)
	})
	fmt.Fprintf(os.Stderr, "\nFor more information, see https://github.com/Schoonology/gloxy.\n")
}

/**
 * Parse command-line flags and start the proxy server.
 */
func main() {
	port := flag.Int("port", 8080, "The port to listen on.")
	help := flag.Bool("help", false, "Show this help message, then exit.")
	flag.Parse()
	flag.Usage = Usage

	if *help {
		Usage()
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	target, err := url.Parse(flag.Arg(0))
	if err != nil {
		fmt.Printf("Bad target URL: %v", flag.Arg(0))
		os.Exit(1)
	}

	if target.Scheme == "" {
		target.Scheme = "http"

		port, err := strconv.ParseInt(target.Path, 0, 0)
		if err == nil {
			target.Host = fmt.Sprintf("127.0.0.1:%d", port)
			target.Path = ""
		} else {
			target.Host = target.Path
			target.Path = ""
		}
	}

	fmt.Printf("Listening on :%d and proxying to %v...\n", *port, target)
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), NewGloxy(target))
	if err != nil {
		fmt.Printf("Failed to Listen with %v", err)
		os.Exit(1)
	}
}
