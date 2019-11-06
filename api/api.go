package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"github.com/m1/go-dependency-check/cache"
	"github.com/m1/go-dependency-check/client"

	"github.com/go-chi/chi"
)

type API struct {
	Config

	client client.Client
	router *chi.Mux
	cache  cache.Cache
	logger *zap.Logger
}

type Config struct {
	Redis string
	Port  string
}

func New(config Config) *API {
	r := chi.NewRouter()
	return &API{router: r, Config: config}
}

func (a *API) Run() error {
	var err error
	a.cache, err = cache.NewRedis(a.Config.Redis)
	if err != nil {
		return err
	}

	a.logger, err = zap.NewProduction()
	if err != nil {
		return err
	}

	packagesHandler := NewPackagesHandler(a)
	a.router.Use(middleware.StripSlashes)
	a.router.Mount("/packages", packagesHandler.GetRoutes())

	ctx := context.Background()
	srv := http.Server{
		Addr:    fmt.Sprintf(":%v", a.Config.Port),
		Handler: chi.ServerBaseContext(ctx, a.router),
	}
	a.logger.Info(fmt.Sprintf("listening on :%v", a.Config.Port))

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		a.logger.Error(fmt.Sprintf("HTTP server ListenAndServe: %v", err))
	}

	return nil
}
