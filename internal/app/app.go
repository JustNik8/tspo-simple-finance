package app

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	_ "net/http"
	"os"
	_ "simple-finance/docs"
	"simple-finance/internal/api"
	"simple-finance/internal/closer"
)

const (
	signingKey    = "J9&#YAVu+gRY7S0V(j)M@8fbr}?$8t"
	serverPortKey = "SERVER_PORT"
	salt          = "dxetkyhvxkhpndxbfnmwkctqqekanrmq"
)

type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runHttpServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHttpServer}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initHttpServer(ctx context.Context) error {
	router := api.NewRouter(a.serviceProvider.GetAuthHandler(), a.serviceProvider.GetTransactionHandler(), a.serviceProvider.GetAuthMiddleware())

	serverPort := os.Getenv(serverPortKey)
	if serverPort == "" {
		log.Fatal("env var SERVER_PORT is empty")
	}
	addr := fmt.Sprintf(":%s", serverPort)
	a.httpServer = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return nil
}

func (a *App) runHttpServer() error {
	log.Println("starting http server on port", a.httpServer.Addr)
	return a.httpServer.ListenAndServe()
}
