package app

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"log"
	"simple-finance/internal/auth"
	"simple-finance/internal/closer"
	"simple-finance/internal/config"
	"simple-finance/internal/db"
	"simple-finance/internal/handler"
	"simple-finance/internal/handler/middleware"
	"simple-finance/internal/tokens"
	"simple-finance/pkg/hash"
)

type serviceProvider struct {
	pgConfig config.PGConfig

	conn *pgx.Conn
	db   *db.FinanceDB

	logger *logrus.Logger

	validate *validator.Validate
	hasher   *hash.SHA1Hasher

	tokenManager *tokens.TokenManager

	auth *auth.Manager

	//redisConfig redisConfig
	redisClient *redis.Client

	authHandler *handler.AuthHandler

	transactionHandler *handler.TransactionHandler

	authMiddleware *middleware.AuthMiddleware
}

type redisConfig interface {
	RedisURL() string
	RedisPassword() string
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GetPGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Panicln(nil, "Database host is not set. db.host should be configured.")
		}
		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GetLogger() *logrus.Logger {
	if s.logger == nil {
		logger := logrus.New()
		s.logger = logger
	}

	return s.logger
}

func (s *serviceProvider) GetConn() *pgx.Conn {

	if s.conn == nil {
		conn, err := pgx.Connect(context.Background(), s.GetPGConfig().DSN())
		if err != nil {
			log.Panicln(nil, "Database connection failed.", err)
		}
		closer.Add(func() error {
			err := conn.Close(context.Background())
			if err != nil {
				return err
			}
			return nil
		})
		s.conn = conn
	}

	return s.conn
}

func (s *serviceProvider) GetFinanceDb() *db.FinanceDB {
	if s.db == nil {
		s.db = db.NewFinanceDB(s.GetConn())
	}

	return s.db
}

//func (s *serviceProvider) GetRedisConfig() redisConfig {
//	if s.redisConfig == nil {
//		cfg, err := config.NewRedisConfig()
//		if err != nil {
//			log.Panicln("Redis config error:", err)
//		}
//		s.redisConfig = cfg
//	}
//	return s.redisConfig
//}

func (s *serviceProvider) GetRedisClient() *redis.Client {
	if s.redisClient == nil {
		// Жестко задаем параметры подключения
		opt := &redis.Options{
			Addr:     "redis:6379", // совпадает с именем сервиса в docker-compose
			Password: "",           // без пароля
			DB:       0,
		}

		client := redis.NewClient(opt)
		_, err := client.Ping(context.Background()).Result()
		if err != nil {
			log.Panicln("Failed to connect to Redis:", err)
		}

		closer.Add(func() error {
			return client.Close()
		})

		s.redisClient = client
	}
	return s.redisClient
}
func (s *serviceProvider) GetHasher() *hash.SHA1Hasher {
	if s.hasher == nil {
		s.hasher = hash.NewSHA1Hasher(salt)
	}

	return s.hasher
}

func (s *serviceProvider) GetValidator() *validator.Validate {
	if s.validate == nil {
		s.validate = validator.New(validator.WithRequiredStructEnabled())
	}

	return s.validate
}
func (s *serviceProvider) GetTokenManager() *tokens.TokenManager {
	if s.tokenManager == nil {
		tokenManager, err := tokens.NewTokenManager(signingKey)
		if err != nil {
			log.Panicln(nil, "Token manager failed.", err)
		}
		s.tokenManager = tokenManager
	}

	return s.tokenManager
}

func (s *serviceProvider) GetAuthManager() *auth.Manager {
	if s.auth == nil {
		s.auth = auth.NewManager(s.GetFinanceDb(), s.GetHasher(), s.GetTokenManager())
	}
	return s.auth
}

func (s *serviceProvider) GetAuthHandler() *handler.AuthHandler {
	if s.authHandler == nil {
		s.authHandler = handler.NewAuthHandler(s.GetValidator(), s.GetFinanceDb(), s.GetLogger(), s.GetHasher(), s.GetAuthManager())
	}
	return s.authHandler
}
func (s *serviceProvider) GetTransactionHandler() *handler.TransactionHandler {
	if s.transactionHandler == nil {
		s.transactionHandler = handler.NewTransactionHandler(s.GetFinanceDb(), s.GetValidator(), s.GetLogger(), s.GetRedisClient())
	}
	return s.transactionHandler
}

func (s *serviceProvider) GetAuthMiddleware() *middleware.AuthMiddleware {
	if s.authMiddleware == nil {
		s.authMiddleware = middleware.NewAuthMiddleware(s.GetTokenManager())
	}
	return s.authMiddleware
}
