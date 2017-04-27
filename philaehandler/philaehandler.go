package philaehandler

import (
	"encoding/json"
	"net/http"

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
