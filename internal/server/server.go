package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"url_shortener/internal/config"
	"url_shortener/internal/delivery/httpdelivery"
	"url_shortener/internal/delivery/httpdelivery/handler"
	"url_shortener/internal/repository"
	"url_shortener/internal/storage"
	"url_shortener/internal/usecase"
	"url_shortener/internal/usecase/cache"
	"url_shortener/internal/usecase/generator"
	"url_shortener/pkg/logger"
)

type Server struct {
	Router *mux.Router
}

func NewServer() *Server {
	urlRepo := repository.NewUrlRepository(storage.DB)
	hashRepo := repository.NewHashRepository(storage.DB)
	uniqueIdRepo := repository.NewUniqueIdRepository(storage.DB)
	urlCacheRepo := repository.NewUrlCacheRepository()

	hashGenerator := generator.NewHashGenerator(uniqueIdRepo, hashRepo)
	hashCache := cache.NewHashCache(hashRepo, uniqueIdRepo, hashGenerator)
	urlUsecase := usecase.NewUrlUsecase(urlRepo, urlCacheRepo, hashCache)
	urlHandler := handler.NewUrlHandler(urlUsecase)

	router := httpdelivery.NewRouter(urlHandler)

	return &Server{
		Router: router,
	}
}

func (s *Server) Run() {
	port := config.AppConfig.Server.Port
	addr := fmt.Sprintf(":%s", port)

	logger.Log.Info("Starting server on port ", port)

	if err := http.ListenAndServe(addr, s.Router); err != nil {
		logger.Log.Fatal("Server error: ", err)
	}
}
