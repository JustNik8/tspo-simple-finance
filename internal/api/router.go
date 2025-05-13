package api

import (
	"github.com/go-chi/chi/v5"
	m "github.com/go-chi/chi/v5/middleware"
	"net/http"
	"simple-finance/internal/handler"
	"simple-finance/internal/handler/middleware"
)

type Router struct {
	transactionHandler *handler.TransactionHandler
	authHandler        *handler.AuthHandler
	authMiddleware     *middleware.AuthMiddleware
	router             *chi.Mux
}

func NewRouter(h *handler.AuthHandler, t *handler.TransactionHandler, m *middleware.AuthMiddleware) *Router {
	r := &Router{
		transactionHandler: t,
		authHandler:        h,
		authMiddleware:     m,
		router:             chi.NewRouter(),
	}

	r.setupMiddleware()
	r.setupRoutes()

	return r
}

func (r *Router) setupMiddleware() {
	// Базовые middleware
	r.router.Use(m.Logger)
	r.router.Use(m.Recoverer)
	r.router.Use(m.RealIP)
	r.router.Use(m.RequestID)
}

func (r *Router) setupRoutes() {

	r.router.Route("/api", func(router chi.Router) {
		router.Use(r.authMiddleware.MakeAuth)

		router.Post("/transaction", r.transactionHandler.InsertTransaction)
		router.Get("/transaction", r.transactionHandler.GetTransactions)
		router.Get("/transaction/{transaction_uuid}", r.transactionHandler.GetTransactionByID)
		router.Delete("/transaction/{transaction_uuid}", r.transactionHandler.DeleteTransactionByID)
	})

	r.router.Route("/auth", func(router chi.Router) {
		router.Post("/sign_in", r.authHandler.SignIn)
		router.Post("/sign_up", r.authHandler.SignUp)
		router.Post("/refresh/tokens", r.authHandler.RefreshTokens)
	})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
