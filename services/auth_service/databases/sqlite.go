package databases

import (
	"context"
	"database/sql"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/piyushsharma67/movie_booking/services/auth_service/models"
)

type SqliteDb struct {
	db *sql.DB
}

func NewSqliteDB(db *sql.DB) Database {
	return &SqliteDb{db: db}
}


var (
	sharedDB *sql.DB

	// 3. Insert i
	once     sync.Once
)

func InitSharedSqliteTestDB() *sql.DB {
	once.Do(func() {
		db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
		if err != nil {
			panic(err)
		}

		schema := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL
		);`

		if _, err := db.Exec(schema); err != nil {
			panic(err)
		}

		sharedDB = db
	})

	return sharedDB
}

func (s *SqliteDb) InsertUser(ctx context.Context, user *models.User) error {
	// Timer to simulate slow DB
	simulateSlow := time.NewTimer(0 * time.Second)
	defer simulateSlow.Stop()

	select {
	case <-simulateSlow.C:
		// Actual DB insert with context
		_, err := s.db.ExecContext(ctx,
			`INSERT INTO users (name, email, password_hash, role)
			 VALUES (?, ?, ?, ?)`,
			user.Name, user.Email, user.PasswordHash, user.Role,
		)
		return err
	case <-ctx.Done():
		// Context canceled or timed out
		return ctx.Err()
	}
}

func (s *SqliteDb) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	err := s.db.QueryRowContext(ctx,
		`SELECT id, name, email, password_hash, role
		 FROM users WHERE email = ?`,
		email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role)

	return u, err
}
