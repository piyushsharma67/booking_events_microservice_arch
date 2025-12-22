package databases

import (
	"bufio"
	"context"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/piyushsharma67/movie_booking/services/auth_service/models"
	"github.com/piyushsharma67/movie_booking/services/auth_service/postgresdb"
)

type SqlDb struct {
	queries *postgresdb.Queries
}

func NewPostgresDB(q *postgresdb.Queries) Database {
	return &SqlDb{
		queries: q,
	}
}

func InitPostgres() (*pgxpool.Pool, *postgresdb.Queries) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode
	slog.Info("Connecting to Postgres", "dsn", dsn)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		panic(err)
	}

	// Wait until Postgres is ready
	if err := pool.Ping(context.Background()); err != nil {
		slog.Error("failed to connect to postgres:", "err", err)
		os.Exit(1)
	}

	// Run schema.sql
	content, err := os.ReadFile("postgresdb/schema.sql")
	if err != nil {
		panic(err)
	}

	// Split statements by semicolon and execute individually
	statements := splitSQLStatements(string(content))
	for _, stmt := range statements {
		if strings.TrimSpace(stmt) == "" {
			continue
		}
		if _, err := pool.Exec(context.Background(), stmt); err != nil {
			slog.Error("failed to execute statement:", "err", err, "stmt", stmt)
			panic(err)
		}
	}

	queries := postgresdb.New(pool)
	return pool, queries
}

// splitSQLStatements splits SQL by semicolon and ignores semicolons inside dollar-quoted blocks
func splitSQLStatements(sql string) []string {
	var stmts []string
	scanner := bufio.NewScanner(strings.NewReader(sql))
	scanner.Split(bufio.ScanLines)
	var sb strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "--") {
			continue // skip comments
		}
		sb.WriteString(line + "\n")
		if strings.HasSuffix(strings.TrimSpace(line), ";") {
			stmts = append(stmts, sb.String())
			sb.Reset()
		}
	}
	if sb.Len() > 0 {
		stmts = append(stmts, sb.String())
	}
	return stmts
}
func (s *SqlDb) InsertUser(ctx context.Context, user *models.User) error {
	// Timer to simulate slow DB
	simulateSlow := time.NewTimer(4 * time.Second)
	defer simulateSlow.Stop()

	select {
	case <-simulateSlow.C:
		// Actual DB insert
		_, err := s.queries.CreateUser(ctx, postgresdb.CreateUserParams{
			Name:         user.Name,
			Email:        user.Email,
			PasswordHash: user.PasswordHash,
			Role:         user.Role,
		})
		return err
	case <-ctx.Done():
		// Context timeout or cancellation
		return context.DeadlineExceeded
	}
}

func (s *SqlDb) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {

	time.Sleep(2 * time.Second)
	u, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:           strconv.Itoa(int(u.ID)),
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         u.Role,
	}, nil
}
