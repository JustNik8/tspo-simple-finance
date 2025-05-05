package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"simple-finance/internal/db"
	"simple-finance/internal/handler"
	appmiddleware "simple-finance/internal/handler/middleware"
	"simple-finance/internal/tokens"
)

const (
	postgresConnStr = "postgres://user:user@postgres:5432/finance_db"
	signingKey      = "J9&#YAVu+gRY7S0V(j)M@8fbr}?$8t"
)

func Run() {
	conn, err := pgx.Connect(context.Background(), postgresConnStr)
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	validate := validator.New(validator.WithRequiredStructEnabled())
	tokenManager, err := tokens.NewTokenManager(signingKey)
	if err != nil {
		panic(err)
	}

	financeDB := db.NewFinanceDB(conn)
	transactionHandler := handler.NewTransactionHandler(financeDB, validate)
	authHandler := handler.NewAuthHandler(validate, tokenManager, financeDB)

	authMiddleware := appmiddleware.NewAuthMiddleware(tokenManager)

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Use(authMiddleware.MakeAuth)

		r.Post("/transaction", transactionHandler.InsertTransaction)
		r.Get("/transaction", transactionHandler.GetTransactions)
		r.Get("/transaction/{transaction_uuid}", transactionHandler.GetTransactionByID)
		r.Delete("/transaction/{transaction_uuid}", transactionHandler.DeleteTransactionByID)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/sign_in", authHandler.SignIn)
	})

	addr := fmt.Sprintf(":%s", "8000")
	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
