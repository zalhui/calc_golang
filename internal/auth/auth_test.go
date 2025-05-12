package auth

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestRegisterAndLoginUser(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE users (
		id TEXT PRIMARY KEY,
		login TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	)`)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Successful registration and login", func(t *testing.T) {
		userID, err := RegisterUser(db, "testuser", "password123")
		if err != nil {
			t.Errorf("Registration failed: %v", err)
		}
		if userID == "" {
			t.Error("Empty user ID returned")
		}

		token, err := LoginUser(db, "testuser", "password123")
		if err != nil {
			t.Errorf("Login failed: %v", err)
		}
		if token == "" {
			t.Error("Empty token returned")
		}
	})

	t.Run("Duplicate registration", func(t *testing.T) {
		_, err := RegisterUser(db, "testuser", "password123")
		if err == nil {
			t.Error("Expected error for duplicate user")
		}
	})

	t.Run("Invalid login", func(t *testing.T) {
		_, err := LoginUser(db, "wronguser", "wrongpass")
		if err == nil {
			t.Error("Expected error for invalid credentials")
		}
	})
}
