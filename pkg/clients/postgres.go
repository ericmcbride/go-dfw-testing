package clients

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

type DBClient struct {
	Db *sql.DB
}

func NewDbConn() (DBClient, error) {
	var client DBClient
	host := os.Getenv("POSTGRES_HOST")
	name := os.Getenv("POSTGRES_NAME")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	sslmode := "disable"

	connection := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=%s", user, password, name, host, sslmode)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		return DBClient{}, err
	}

	client.Db = db
	return client, nil
}

func Close(db *DBClient) error {
	if db.Db == nil {
		return nil
	}

	err := db.Db.Close()
	if err != nil {
		return err
	}
	return nil
}
