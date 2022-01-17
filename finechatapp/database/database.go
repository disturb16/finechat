package database

import (
	"context"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// DatabaseConfiguration is the configuration for a database.
type DatabaseConfiguration struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

// CreateMysqlConnection returns a mysql connection.
func CreateMysqlConnection(ctx context.Context, config DatabaseConfiguration) (*sqlx.DB, error) {
	connString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)

	log.Println("Connecting to database")
	return sqlx.ConnectContext(ctx, "mysql", connString)
}
