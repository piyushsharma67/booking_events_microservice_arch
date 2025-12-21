package databases

import (
	"context"
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteDb struct {
	db *sql.DB
}

var (
	once   sync.Once
	testDB Database
	err    error
)

func NewSqliteDB(db *sql.DB) Database {
	return &SqliteDb{db: db}
}

func InitSqliteTestDB() (Database, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL
	);`

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return NewSqliteDB(db), nil
}

func GetTestDB() Database {
	once.Do(func() {
		testDB, err = InitSqliteTestDB()
		if err != nil {
			panic(err)
		}
	})
	return testDB
}

func (s *SqliteDb) InsertUser(ctx context.Context, user *User) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO users (name, email, password_hash, role)
		 VALUES (?, ?, ?, ?)`,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Role,
	)
	return err
}

func (s *SqliteDb) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := s.db.QueryRowContext(ctx,
		`SELECT id, name, email, password_hash, role
		 FROM users WHERE email = ?`,
		email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role)

	return u, err
}
