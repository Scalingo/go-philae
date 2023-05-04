package tests

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/jarcoal/httpmock"
)

type Route struct {
	Method string
	Path   string
}

type MatchResponder struct {
	Matcher   httpmock.Matcher
	Responder httpmock.Responder
}

func HTTPTestServer(routes map[Route]MatchResponder) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedRoute := Route{
			Method: r.Method,
			Path:   r.URL.Path,
		}

		responder, exists := routes[requestedRoute]

		if !exists {
			http.NotFound(w, r)
			return
		}

		if !responder.Matcher.Check(r) {
			http.Error(w, "headers not matching", 400)
			return
		}

		response, _ := responder.Responder(r)
		defer response.Body.Close()

		w.WriteHeader(response.StatusCode)
		_, _ = io.Copy(w, response.Body)
	}))

	return srv
}
