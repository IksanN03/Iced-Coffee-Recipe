package route

import (
	"be-test/database"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Set up test database connection
	os.Setenv("DATABASE_URL", "postgres://postgres:12345678@localhost:5432/test?sslmode=disable")
	database.Init()

	// Run tests
	code := m.Run()

	os.Exit(code)
}
