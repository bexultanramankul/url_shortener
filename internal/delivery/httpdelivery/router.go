package httpdelivery

import (
	"github.com/gorilla/mux"
	"net/http"
	"url_shortener/internal/delivery/httpdelivery/handler"
	"url_shortener/pkg/logger"
)

func NewRouter(urlHandler *handler.UrlHandler) *mux.Router {
	router := mux.NewRouter()

	logger.Log.Info("Configuring routes...")

	router.HandleFunc("/shorten", urlHandler.ShortenUrl).Methods(http.MethodPost)
	logger.Log.Info("Route registered: POST /shorten")

	router.HandleFunc("/redirect/{hash}", urlHandler.Redirect).Methods(http.MethodGet)
	logger.Log.Info("Route registered: GET /redirect/{hash}")

	logger.Log.Info("Routes configured successfully.")

	return router
}
