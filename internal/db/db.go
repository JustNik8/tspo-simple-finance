package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"simple-finance/internal/models"
)

type FinanceDB struct {
	conn *pgx.Conn
}

func NewFinanceDB(conn *pgx.Conn) *FinanceDB {
	return &FinanceDB{
		conn: conn,
	}
}

func (db *FinanceDB) InsertTransaction(ctx context.Context, transaction models.Transaction) (string, error) {
	const query = `
	INSERT INTO transactions (id, user_id, amount, category_id, comment, date, created_at)
	VALUES($1, $2, $3, $4, $5, $6, NOW())
	RETURNING id
	`

	row := db.conn.QueryRow(ctx, query,
		transaction.ID,
		transaction.UserID,
		transaction.Amount,
		transaction.CategoryID,
		transaction.Comment,
		transaction.Date,
	)

	var transactionID string
	err := row.Scan(&transactionID)

	return transactionID, err
}

func (db *FinanceDB) GetTransactions(ctx context.Context, userID string) ([]models.Transaction, error) {
	const query = "SELECT * FROM transactions WHERE user_id = $1"

	rows, err := db.conn.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	transactions := make([]models.Transaction, 0)

	for rows.Next() {
		var transaction models.Transaction

		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Amount,
			&transaction.CategoryID,
			&transaction.Comment,
			&transaction.Date,
			&transaction.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (db *FinanceDB) GetTransactionByID(ctx context.Context, userID string, transactionID string) (models.Transaction, error) {
	const query = "SELECT * FROM transactions WHERE user_id = $1 AND id = $2 LIMIT 1"

	row := db.conn.QueryRow(ctx, query, userID, transactionID)
	var transaction models.Transaction

	err := row.Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Amount,
		&transaction.CategoryID,
		&transaction.Comment,
		&transaction.Date,
		&transaction.CreatedAt,
	)

	return transaction, err
}

func (db *FinanceDB) DeleteTransactionByID(ctx context.Context, userID string, transactionID string) error {
	const query = "DELETE FROM transactions WHERE user_id = $1 AND id = $2"

	_, err := db.conn.Exec(ctx, query, userID, transactionID)
	return err
}

func (db *FinanceDB) GetUserID(ctx context.Context, username string) (string, error) {
	const query = `SELECT id FROM users WHERE username = $1 LIMIT 1`

	row := db.conn.QueryRow(ctx, query, username)
	var userID string

	err := row.Scan(&userID)

	return userID, err
}

func (db *FinanceDB) InsertUser(ctx context.Context, userInfo models.UserInfo) (models.UserInfo, error) {
	const query = `
		INSERT INTO users(id, email, username, hash_pass, created_at)
		VALUES($1, $2, $3, $4, NOW())
		RETURNING created_at
	`

	row := db.conn.QueryRow(ctx, query,
		userInfo.ID,
		userInfo.Email,
		userInfo.UserName,
		userInfo.Password,
	)

	var createdAt time.Time
	err := row.Scan(&createdAt)
	if err != nil {
		return models.UserInfo{}, err
	}

	return models.UserInfo{
		ID:        userInfo.ID,
		Email:     userInfo.Email,
		UserName:  userInfo.UserName,
		CreatedAt: createdAt,
	}, nil
}

func (db *FinanceDB) GetUserInfo(ctx context.Context, userName string) (models.UserInfo, error) {
	const query = `
		SELECT id, email, username, hash_pass, created_at
		FROM users
		WHERE username = $1
		LIMIT 1
	`
	row := db.conn.QueryRow(ctx, query, userName)
	var userInfo models.UserInfo
	err := row.Scan(
		&userInfo.ID,
		&userInfo.Email,
		&userInfo.UserName,
		&userInfo.Password,
		&userInfo.CreatedAt,
	)

	return userInfo, err
}
