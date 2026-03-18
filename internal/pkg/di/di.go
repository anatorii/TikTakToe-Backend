package di

import (
	"context"
	"log"
	"net/http"
	"tiktaktoe/internal/pkg/datasource"
	"tiktaktoe/internal/pkg/datasource/repository"
	drepository "tiktaktoe/internal/pkg/domain/repository"
	"tiktaktoe/internal/pkg/domain/service"
	"tiktaktoe/internal/pkg/web"
	wservice "tiktaktoe/internal/pkg/web/service"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

// зависимости приложения
var Module = fx.Module("app",
	fx.Provide(
		// Инициализация зависимостей
		datasource.NewDbPool,
		repository.NewGameRepository,
		repository.NewUserRepository,
		drepository.NewGameRepo,
		drepository.NewUserRepo,
		service.NewGameService,
		service.NewUserService,
		wservice.NewGameServ,
		wservice.NewUserServ,
		wservice.NewAuthenticator,
		wservice.NewAuthenticationService,
		wservice.NewJwtProvider,
		web.NewAuthHandler,
		web.NewGameHandler,
		web.NewRouter,
		// HTTP сервер
		NewHTTPServer,
	),
	fx.Invoke(
		web.RegisterAuthRoutes,
		web.RegisterGameRoutes,
		RegisterHooks),
)

// конфигурация сервера
type HTTPServerConfig struct {
	Addr string
}

// HTTP сервер с mux.Router
func NewHTTPServer(lc fx.Lifecycle, router *mux.Router) *http.Server {
	config := &HTTPServerConfig{
		Addr: ":8080",
	}
	server := &http.Server{
		Addr:    config.Addr,
		Handler: router,
	}
	return server
}

// жизненный цикл сервера
func RegisterHooks(lc fx.Lifecycle, server *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Printf("Server starting on %s", server.Addr)
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("Failed to start server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down server...")
			return server.Shutdown(ctx)
		},
	})
}
