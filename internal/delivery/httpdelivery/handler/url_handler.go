package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"url_shortener/internal/config" // Импортируем конфиг
	"url_shortener/internal/usecase"
	"url_shortener/pkg/logger"
)

type UrlHandler struct {
	usecase *usecase.UrlUsecase
}

func NewUrlHandler(usecase *usecase.UrlUsecase) *UrlHandler {
	return &UrlHandler{usecase: usecase}
}

func (h *UrlHandler) ShortenUrl(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Url string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hash, err := h.usecase.ShortenUrl(request.Url)
	if err != nil {
		logger.Log.Error(err)
		http.Error(w, "Error shortening URL", http.StatusInternalServerError)
		return
	}

	shortenedUrl := config.AppConfig.Server.BaseURL + "/redirect/" + hash

	response := map[string]string{"shortened_url": shortenedUrl}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UrlHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	url, err := h.usecase.GetUrl(hash)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
