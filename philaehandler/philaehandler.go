package philaehandler

import (
	"encoding/json"
	"net/http"

	handlers "github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-philae/prober"
)

type PhilaeHandler struct {
	prober *prober.Prober
}

func (handler PhilaeHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	result := handler.prober.Check()
	json.NewEncoder(response).Encode(result)

}

func NewHandler(prober *prober.Prober) http.Handler {
	return PhilaeHandler{
		prober: prober,
	}
}

func NewScalingoHandler(prober *prober.Prober) handlers.HandlerFunc {
	return func(response http.ResponseWriter, _ *http.Request, _ map[string]string) error {
		json.NewEncoder(response).Encode(prober.Check())
		return nil
	}
}
