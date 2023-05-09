package tests

import (
	"net/http"
	"net/http/httptest"
)

type Route struct {
	Method string
	Path   string
}

type Response struct {
	Status int
	Body   string
}

func HTTPTestServer(routes map[Route]Response) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedRoute := Route{
			Method: r.Method,
			Path:   r.URL.Path,
		}

		response, exists := routes[requestedRoute]

		if !exists {
			http.NotFound(w, r)
			return
		}

		w.WriteHeader(response.Status)
		_, _ = w.Write([]byte(response.Body))
	}))

	return srv
}
