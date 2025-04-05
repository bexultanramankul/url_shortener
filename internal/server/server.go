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

	logger.Log.Info("Repositories initialized: URL, Hash, UniqueId, UrlCache")

	hashGenerator := generator.NewHashGenerator(uniqueIdRepo, hashRepo)
	hashCache := cache.NewHashCache(hashRepo, uniqueIdRepo, hashGenerator)
	logger.Log.Info("Hash generator and hash cache initialized.")

	urlUsecase := usecase.NewUrlUsecase(urlRepo, urlCacheRepo, hashCache)
	logger.Log.Info("URL usecase initialized.")

	urlHandler := handler.NewUrlHandler(urlUsecase)
	logger.Log.Info("URL handler initialized.")

	router := httpdelivery.NewRouter(urlHandler)
	logger.Log.Info("Router initialized.")

	return &Server{
		Router: router,
	}
}

func (s *Server) Run() {
	port := config.AppConfig.Server.Port
	addr := fmt.Sprintf(":%s", port)

	logger.Log.Infof("Starting server on port %s...", port)

	if err := http.ListenAndServe(addr, s.Router); err != nil {
		logger.Log.Fatal("Server error: ", err)
	}
}
