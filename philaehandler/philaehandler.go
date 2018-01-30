package philaehandler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Scalingo/go-philae/prober"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type PhilaeHandler struct {
	prober  *prober.Prober
	verbose bool
	logger  logrus.FieldLogger
}

func (handler PhilaeHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	start := time.Now()
	result := handler.prober.Check()
	json.NewEncoder(response).Encode(result)
	duration := time.Now().Sub(start)

	if (handler.verbose && result.Healthy) || !result.Healthy {
		handler.logger.WithFields(logrus.Fields{
			"duration": duration.String(),
			"healthy":  result.Healthy,
		}).Info()
	}
}

type HandlerOpts struct {
	Verbose bool
	Logger  logrus.FieldLogger
}

func NewHandler(prober *prober.Prober, opts HandlerOpts) http.Handler {
	h := PhilaeHandler{
		prober:  prober,
		logger:  opts.Logger,
		verbose: opts.Verbose,
	}
	if h.logger != nil {
		h.logger = logrus.New()
	}
	return h
}

func NewPhilaeRouter(router http.Handler, prober *prober.Prober, opts HandlerOpts) *mux.Router {
	globalRouter := mux.NewRouter()
	globalRouter.Handle("/_health", NewHandler(prober, opts))
	globalRouter.Handle("/{any:.+}", router)
	return globalRouter
}
