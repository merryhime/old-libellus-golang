package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/vharitonsky/iniflags"
	"golang.org/x/crypto/acme/autocert"
)

var (
	canonicalDomain = flag.String("domain", "", "canonical domain")
	altDomains      = flag.String("alt_domains", "", "alternative domains, to redirect to the canonical domain, seperated by semicolons")
	httpsEndpoint   = flag.String("https_endpoint", "127.0.0.1:8081", "HTTPS endpoint")
	httpEndpoint    = flag.String("http_endpoint", "127.0.0.1:8080", "HTTP endpoint")
	privateDir      = flag.String("private_dir", "./libellus_private/", "private data directory")
)

func parseConfig() {
	iniflags.Parse()

	if *canonicalDomain == "" {
		panic("canonical domain required")
	}
}

func isAltDomain(host string) bool {
	if *altDomains == "" || host == "" {
		return false
	}

	alts := strings.Split(*altDomains, ";")
	for _, v := range alts {
		if v == host {
			return true
		}
	}
	return false
}

func makeAutocertManager() *autocert.Manager {
	invalidHostError := errors.New("Invalid Host")

	hostPolicy := func(ctx context.Context, host string) error {
		if host == "" || strings.ContainsRune(host, ';') {
			return invalidHostError
		}
		if host == *canonicalDomain {
			return nil
		}
		if isAltDomain(host) {
			return nil
		}
		return invalidHostError
	}

	certCacheDir := *privateDir + "/cert_cache/"

	return &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocert.DirCache(certCacheDir),
	}
}

func runRedirectServer(acm *autocert.Manager) {
	srv := http.Server{
		Addr:         *httpEndpoint,
		Handler:      acm.HTTPHandler(nil),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	go func() {
		log.Printf("Starting HTTP redirector")
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("HTTP redirector failed with %s", err)
		}
	}()
}

func makeMux() *http.ServeMux {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		newURI := "https://" + *canonicalDomain + r.URL.String()
		http.Redirect(w, r, newURI, http.StatusFound)
	})
	mux.HandleFunc(*canonicalDomain+"/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "It works!")
	})
	return mux
}

func main() {
	parseConfig()

	acm := makeAutocertManager()
	mux := makeMux()
	srv := &http.Server{
		Addr:         *httpsEndpoint,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig:    &tls.Config{GetCertificate: acm.GetCertificate},
		Handler:      mux,
	}

	runRedirectServer(acm)

	log.Printf("Starting HTTPS server")
	err := srv.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("srv.ListendAndServeTLS() failed with %s", err)
	}
}
