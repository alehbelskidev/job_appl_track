package repo

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) *sql.DB {
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:16",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.WithInitScripts("../../sql/schema.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err)

	t.Cleanup(func() { _ = container.Terminate(ctx) })

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sql.Open("pgx", connStr)
	require.NoError(t, err)

	return db
}

func TestCreateUser_Integration(t *testing.T) {
	db := setupTestDB(t)
	q := New(db)

	user, err := q.CreateUser(context.Background(), CreateUserParams{
		Email:    "test@test.com",
		Password: "hashedpassword",
	})

	require.NoError(t, err)
	assert.Equal(t, "test@test.com", user.Email)
	assert.NotEmpty(t, user.ID)
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	q := New(db)

	_, err := q.GetUserByEmail(context.Background(), "nobody@test.com")

	assert.Error(t, err)
}

func TestCreateJobApplication_Integration(t *testing.T) {
	db := setupTestDB(t)
	q := New(db)

	user, err := q.CreateUser(context.Background(), CreateUserParams{
		Email:    "test@test.com",
		Password: "hash",
	})
	require.NoError(t, err)

	app, err := q.CreateJobApplication(context.Background(), CreateJobApplicationParams{
		Company:     "Google",
		Role:        "SWE",
		DateApplied: time.Now(),
		Status:      0,
		OwnerID:     user.ID,
	})

	require.NoError(t, err)
	assert.Equal(t, "Google", app.Company)
	assert.Equal(t, user.ID, app.OwnerID)
}
