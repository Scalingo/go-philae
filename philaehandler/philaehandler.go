package philaehandler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Scalingo/go-philae/prober"
	"github.com/gorilla/mux"
)

type PhilaeHandler struct {
	prober *prober.Prober
}

func (handler PhilaeHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	start := time.Now()
	result := handler.prober.Check()
	json.NewEncoder(response).Encode(result)
	duration := time.Now().Sub(start)
	log.Printf("[PHILAE] Probe check done. Duration: %s, Healthy: %t", duration.String(), result.Healthy)

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
