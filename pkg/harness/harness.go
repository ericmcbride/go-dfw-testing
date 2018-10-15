package harness

import (
	"github.com/ericmcbride/go-dfw-testing/pkg/clients"
	"github.com/ericmcbride/go-dfw-testing/pkg/logging"
	"os"
	"testing"
)

func Run(m *testing.M) {
	logging.ConfigureLogger("ERROR")
	SetEnvironmentals()
	CreateDatabase()
	code := m.Run()
	DropDatabase()
	os.Exit(code)
}

func SetEnvironmentals() {
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_USER", "postgres")
	os.Setenv("POSTGRES_PASSWORD", "password")
	os.Setenv("POSTGRES_NAME", "go_dfw_test")
	os.Setenv("X-CAR-TOKEN", "1234")
}

func CreateDatabase() {
	client, err := clients.NewDbConn()
	if err != nil {
		panic(err)
	}

	defer clients.Close(&client)

	query := `
	CREATE TABLE IF NOT EXISTS cars (
		id uuid PRIMARY KEY,
		model character(128),
		make character(128),
		color character(128),
		year integer
	)`

	_, err = client.Db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func DropDatabase() {
	client, err := clients.NewDbConn()
	if err != nil {
		panic(err)
	}

	defer clients.Close(&client)

	query := `DROP TABLE IF EXISTS cars;`

	_, err = client.Db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func Truncate() {
	client, err := clients.NewDbConn()
	if err != nil {
		panic(err)
	}

	defer clients.Close(&client)

	query := `TRUNCATE ONLY cars;`

	_, err = client.Db.Exec(query)
	if err != nil {
		panic(err)
	}
}
