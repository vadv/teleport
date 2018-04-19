package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	libhttp "http"
)

var (
	proxy       = flag.String("proxy", "socks5://user:password@ip:port", "proxy")
	listenHTTP  = flag.String("listen-http", "127.0.0.127:80", "listen ip")
	listenHTTPS = flag.String("listen-https", "", "listen ssl")
	keyServer   = flag.String("ssl-key", "server.key", "path to server.key")
	certServer  = flag.String("ssl-crt", "server.crt", "path to server.crt")
)

func main() {

	if !flag.Parsed() {
		flag.Parse()
	}

	s, err := libhttp.NewServer(*proxy)
	if err != nil {
		log.Printf("[FATAL] build proxy: %s\n", err.Error())
		os.Exit(2)
	}

	errChan := make(chan error, 1)

	go func() {
		if addr := *listenHTTP; addr != "" {
			l, err := net.Listen("tcp", addr)
			if err != nil {
				if err != nil {
					errChan <- err
					return
				}
			}
			log.Printf("[INFO] start http server at %s\n", addr)
			errChan <- http.Serve(l, s)
		}
	}()

	go func() {
		if addr := *listenHTTPS; addr != "" {
			l, err := net.Listen("tcp", addr)
			if err != nil {
				if err != nil {
					errChan <- err
					return
				}
			}
			log.Printf("[INFO] start https server at %s\n", addr)
			errChan <- http.ServeTLS(l, s, *certServer, *keyServer)
		}
	}()

	if err := <-errChan; err != nil {
		log.Printf("%s\n", err.Error())
		os.Exit(3)
	}

}
