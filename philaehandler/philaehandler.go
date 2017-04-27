package philaehandler

import (
	"encoding/json"
	"net/http"

	"github.com/Scalingo/go-philae/prober"
	"github.com/gorilla/mux"
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

func NewPhilaeRouter(router http.Handler, prober *prober.Prober) *mux.Router {
	globalRouter := mux.NewRouter()
	globalRouter.Handle("/_health", NewHandler(prober))
	globalRouter.Handle("/{any:.+}", router)
	return globalRouter
}
