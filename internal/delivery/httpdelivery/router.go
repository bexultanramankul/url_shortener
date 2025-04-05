package httpdelivery

import (
	"github.com/gorilla/mux"
	"net/http"
	"url_shortener/internal/delivery/httpdelivery/handler"
)

// NewRouter создает и настраивает маршруты
func NewRouter(urlHandler *handler.UrlHandler) *mux.Router {
	router := mux.NewRouter()

	// Определяем маршруты
	router.HandleFunc("/shorten", urlHandler.ShortenUrl).Methods(http.MethodPost)
	router.HandleFunc("/redirect/{hash}", urlHandler.Redirect).Methods(http.MethodGet)

	return router
}
