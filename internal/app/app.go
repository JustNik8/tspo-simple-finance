package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"simple-finance/internal/db"
	"simple-finance/internal/handler"
	appmiddleware "simple-finance/internal/handler/middleware"
	"simple-finance/internal/tokens"
)

const (
	signingKey    = "J9&#YAVu+gRY7S0V(j)M@8fbr}?$8t"
	serverPortKey = "SERVER_PORT"
)

func Run() {
	logger := setupLogger()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("No .env file found, using system environment variables")
	}

	serverPort := os.Getenv(serverPortKey)
	if serverPort == "" {
		logger.Fatal("env var SERVER_PORT is empty")
	}

	conn, err := pgx.Connect(context.Background(), getPostgresConn(logger))
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		err := conn.Close(context.Background())
		if err != nil {
			logger.Warn(err)
		}
	}()

	validate := validator.New(validator.WithRequiredStructEnabled())
	tokenManager, err := tokens.NewTokenManager(signingKey)
	if err != nil {
		logger.Fatal(err)
	}

	financeDB := db.NewFinanceDB(conn)
	transactionHandler := handler.NewTransactionHandler(financeDB, validate, logger)
	authHandler := handler.NewAuthHandler(validate, tokenManager, financeDB, logger)
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

	addr := fmt.Sprintf(":%s", serverPort)
	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	err = server.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}
}

func getPostgresConn(logger *logrus.Logger) string {
	dbUser, found := os.LookupEnv("DB_USER")
	if !found {
		logger.Fatal("DB_USER not found")
	}

	dbPass, found := os.LookupEnv("DB_PASS")
	if !found {
		logger.Fatal("DB_PASS not found")
	}

	dbHost, found := os.LookupEnv("DB_HOST")
	if !found {
		logger.Fatal("DB_HOST not found")
	}

	dbPort, found := os.LookupEnv("DB_PORT")
	if !found {
		logger.Fatal("DB_PORT not found")
	}

	dbName, found := os.LookupEnv("DB_NAME")
	if !found {
		logger.Fatal("DB_NAME not found")
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
}

func setupLogger() *logrus.Logger {
	logger := logrus.New()

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutput(file)
	} else {
		logger.Warn("Не удалось открыть файл логов, логирование в консоль")
	}

	return logger
}
