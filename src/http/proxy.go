package http

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	libproxy "golang.org/x/net/proxy"
)

type proxy struct {
	client *http.Client
}

func newProxy(connectionString string) (*proxy, error) {
	result := &proxy{client: &http.Client{}}
	u, err := url.Parse(connectionString)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case `socks5`:
		// auth
		var auth *libproxy.Auth
		if u.User != nil {
			auth = &libproxy.Auth{User: u.User.Username()}
			if password, ok := u.User.Password(); ok {
				auth.Password = password
			}
		}
		socksProxy, err := libproxy.SOCKS5("tcp", u.Host, auth, libproxy.Direct)
		if err != nil {
			return nil, err
		}
		transport := &http.Transport{
			Dial:            socksProxy.Dial,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		result.client.Transport = transport
		log.Printf("[INFO] created socks5 proxy")
	default:
		return nil, fmt.Errorf("unknown proxy type: %s", u.Scheme)
	}
	return result, nil
}

func (s *proxy) get(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	url, err := getUrl(req)
	if err != nil {
		log.Printf("[ERROR] GET %s: %s\n", url, err.Error())
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	resp, err := s.client.Get(url)
	if err != nil {
		log.Printf("[ERROR] GET %s: %s\n", url, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("[ERROR] GET %s: %s\n", url, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("[INFO] GET %s [%.2f]\n", url, time.Now().Sub(start).Seconds())
}

func (s *proxy) post(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	url, err := getUrl(req)
	if err != nil {
		log.Printf("[ERROR] GET %s: %s\n", url, err.Error())
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	resp, err := s.client.Post(url, req.Header.Get("Content-type"), req.Body)
	if err != nil {
		log.Printf("[ERROR] POST %s: %s\n", url, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("[ERROR] POST %s: %s\n", url, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("[INFO] POST %s [%.2f]\n", url, time.Now().Sub(start).Seconds())
}
