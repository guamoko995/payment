// В данном файле описана логика взаимодействий с базой данных.

package database

import (
	"database/sql"
	"fmt"
	"sync"

	_ "modernc.org/sqlite"
)

// schemaSQL - заголовок схемы базы данных.
var schemaSQL = `
CREATE TABLE IF NOT EXISTS accounts(
    id  INTEGER PRIMARY KEY,
    balance MONEY
);

CREATE INDEX IF NOT EXISTS accounts_id ON accounts(id);
`

// wownBalance - Заголовок запроса списания.
var downBalance = `
UPDATE accounts
SET balance = balance - ?
WHERE id==?;
`

// upBalance - Заголовок запроса пополнения.
var upBalance = `
UPDATE accounts
SET balance = balance + ?
WHERE id==?;
`

// getBalance - Заголовок запроса баланса.
var getBalance = ` 
SELECT balance FROM accounts
WHERE id == ?
`

// addAccount - Заголовок запроса добавления платежного аккаунта.
var addAccount = `
INSERT INTO accounts (
    balance
) VALUES (
    0
)
`

// Account - запись в базе данных
type Account struct {
	ID      uint64
	Balance int64
}

// DB - база данных платежных аккаунтов.
type DB struct {
	mu  sync.RWMutex
	sql *sql.DB

	// Предварительно скомпилированные запросы
	downBalance *sql.Stmt
	upBalance   *sql.Stmt
	getBalance  *sql.Stmt
	addAccount  *sql.Stmt
}

// New создает/открывает базу данных.
func NewDB(dbFile string) (*DB, error) {
	db := DB{}
	var err error

	// создает/открывает файл базы данных
	db.sql, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	// создает/проверяет наличие таблицы accounts,
	// соответствующей schemaSQL
	if _, err = db.sql.Exec(schemaSQL); err != nil {
		return nil, err
	}

	// предварительная компеляция заголовков
	db.addAccount, err = db.sql.Prepare(addAccount)
	if err != nil {
		return nil, err
	}

	db.getBalance, err = db.sql.Prepare(getBalance)
	if err != nil {
		return nil, err
	}

	db.downBalance, err = db.sql.Prepare(downBalance)
	if err != nil {
		return nil, err
	}

	db.upBalance, err = db.sql.Prepare(upBalance)
	if err != nil {
		return nil, err
	}

	return &db, nil
}

// AddAccount - добавляет платежный аккаунт в базу данных.
func (db *DB) AddAccount() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	tx, err := db.sql.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(db.addAccount).Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// GetBalance - возвращает баланс по ID.
func (db *DB) GetBalance(ID uint64) (int64, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	tx, err := db.sql.Begin()
	if err != nil {
		return 0, err
	}

	row := tx.Stmt(db.getBalance).QueryRow(ID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var balance int64

	err = row.Scan(&balance)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return balance, tx.Commit()
}

// UpBalance - пополняет баланс по ID на указанную сумму.
func (db *DB) UpBalance(ID uint64, sum int64) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	tx, err := db.sql.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(db.upBalance).Exec(sum, ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// DownBalance - списывает с баланса по ID указанную сумму.
func (db *DB) DownBalance(ID uint64, sum int64) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	tx, err := db.sql.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(db.downBalance).Exec(sum, ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// AmountTransfer - осуществляет перевод суммы с одного ID на другой.
func (db *DB) AmountTransfer(senderID, geterID uint64, sum int64) error {
	// Получение баланса отправителя.
	senderBalance, err := db.GetBalance(senderID)
	if err != nil {
		return err
	}

	// Проверка наличия достаточного количества средств
	if senderBalance < sum {
		return ErrorInsufficientFunds{
			ID:            senderID,
			RequestAmount: sum,
		}
	}

	db.mu.Lock()
	defer db.mu.Unlock()
	// Начало операции перевода
	tx, err := db.sql.Begin()
	if err != nil {
		return err
	}

	// Списание средств у отправителя
	_, err = tx.Stmt(db.downBalance).Exec(sum, senderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Начисление средств получателю
	_, err = tx.Stmt(db.upBalance).Exec(sum, geterID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Завершение операции перевода
	return tx.Commit()
}

type ErrorInsufficientFunds struct {
	ID            uint64
	RequestAmount int64
}

func (err ErrorInsufficientFunds) Error() string {
	return fmt.Sprintf("На счете %v недостаточно средств", err.ID)
}
