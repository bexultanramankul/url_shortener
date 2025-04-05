package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"url_shortener/internal/config"
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

	logger.Log.Infof("Received request to shorten URL: %s", r.URL.Path)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Log.Errorf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hash, err := h.usecase.ShortenUrl(request.Url)
	if err != nil {
		logger.Log.Errorf("Error shortening URL %s: %v", request.Url, err)
		http.Error(w, "Error shortening URL", http.StatusInternalServerError)
		return
	}

	shortenedUrl := config.AppConfig.Server.BaseURL + "/redirect/" + hash

	logger.Log.Infof("Successfully shortened URL %s to %s", request.Url, shortenedUrl)

	response := map[string]string{"shortened_url": shortenedUrl}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func (h *UrlHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	logger.Log.Infof("Received redirect request for hash: %s", hash)

	url, err := h.usecase.GetUrl(hash)
	if err != nil {
		logger.Log.Errorf("Error retrieving URL for hash %s: %v", hash, err)
		http.Error(w, "Error retrieving URL", http.StatusInternalServerError)
		return
	}

	if url == "" {
		logger.Log.Warnf("URL not found for hash %s", hash)
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
