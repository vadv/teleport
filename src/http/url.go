package http

import (
	"fmt"
	"net/http"
)

// криво но не понятно что с этим делать

func getUrl(req *http.Request) (string, error) {

	proto := `http`
	if req.TLS != nil {
		proto = `https`
	}

	host := req.Host
	if host == `` {
		host = req.Header.Get(`Host`)
	}

	result := fmt.Sprintf("%s://%s%s", proto, host, req.RequestURI)
	return result, nil
}
