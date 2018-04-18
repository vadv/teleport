package http

import (
	"net/http"
)

type server struct {
	p *proxy
}

func NewServer(connectionString string) (*server, error) {
	p, err := newProxy(connectionString)
	if err != nil {
		return nil, err
	}
	return &server{p: p}, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		s.p.get(w, req)
	case http.MethodPost:
		s.p.post(w, req)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}
