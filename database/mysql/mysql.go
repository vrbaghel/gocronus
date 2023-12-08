package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	Client *sql.DB
}

/** creates and returns a new sql client instance */
func NewMySQL(connString string) (*MySQL, error) {
	client, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, fmt.Errorf("mysql : unable to establish connection with mysql server %s", err.Error())
	}
	return &MySQL{
		Client: client,
	}, nil
}

/** verifies if connection to the database client is alive */
func (m *MySQL) Ping() error {
	if err := m.Client.Ping(); err != nil {
		return fmt.Errorf("mysql : connection status failed %s", err.Error())
	}
	return nil
}

/** closes the database client */
func (m *MySQL) CloseConnection() error {
	if err := m.Client.Close(); err != nil {
		return fmt.Errorf("mysql : unable to close connection with mysql server %s", err.Error())
	}
	return nil
}
