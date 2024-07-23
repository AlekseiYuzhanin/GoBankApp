package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage interface{
	CreateAccount(*Account) error
	DeleteAccount(int) error
	GetAccounts() ([]*Account,error) 
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account,error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error){
	conn := "user=postgres dbname=postgres password=qweasd123 sslmode=disable"
	db,err := sql.Open("postgres", conn)
	if err != nil{
		return nil, err
	}

	if err := db.Ping(); err != nil{
		return nil,err
	}

	return &PostgresStore{
		db: db,
	},nil
}

func (s *PostgresStore) Init() error{
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error{
	query := `CREATE TABLE IF NOT EXISTS Account(
		id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		firstname CHARACTER VARYING(30) NOT NULL,
		lastname CHARACTER VARYING(30) NOT NULL,
		number serial,
		balance serial,
		created_at timestamp
	)`

	_,err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error{
	query:= (`INSERT INTO Account(firstname, lastname, number, balance, created_at)
	VALUES
	($1, $2, $3, $4, $5)
	
	`)
	_, err := s.db.Query(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt)
	if err != nil{
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error{
	return nil

}

func (s *PostgresStore) GetAccountById(id int) (*Account,error){
	return nil,nil

}

func (s *PostgresStore) UpdateAccount(*Account) error{
	return nil

}

func (s *PostgresStore) GetAccounts() ([]*Account,error){
	rows,err := s.db.Query("SELECT * FROM Account")
	if err != nil{
		return nil,err
	}
	accounts := []*Account{}
	for rows.Next(){
		account := new(Account)
		if err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt); err != nil{
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts,nil
}