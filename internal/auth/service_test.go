package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alehbelskidev/job_appl_track/internal/repo"
	"github.com/alehbelskidev/job_appl_track/internal/shared"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockQuerier struct {
	users         map[string]repo.User
	refreshTokens map[string]struct{}
}

func newMockQuerier() *mockQuerier {
	return &mockQuerier{
		users:         make(map[string]repo.User),
		refreshTokens: make(map[string]struct{}),
	}
}

func (m *mockQuerier) CreateUser(_ context.Context, arg repo.CreateUserParams) (repo.User, error) {
	user := repo.User{ID: uuid.New(), Email: arg.Email, Password: arg.Password}
	m.users[arg.Email] = user
	return user, nil
}

func (m *mockQuerier) GetUserByEmail(_ context.Context, email string) (repo.User, error) {
	user, ok := m.users[email]
	if !ok {
		return repo.User{}, errors.New("not found")
	}
	return user, nil
}

func (m *mockQuerier) AddRefreshToken(_ context.Context, token string) error {
	m.refreshTokens[token] = struct{}{}
	return nil
}

func (m *mockQuerier) GetRefreshToken(_ context.Context, token string) (string, error) {
	if _, ok := m.refreshTokens[token]; !ok {
		return "", errors.New("not found")
	}
	return token, nil
}

func (m *mockQuerier) DeleteRefreshToken(_ context.Context, token string) error {
	delete(m.refreshTokens, token)
	return nil
}

func (m *mockQuerier) CreateJobApplication(_ context.Context, arg repo.CreateJobApplicationParams) (repo.JobApplication, error) {
	return repo.JobApplication{}, nil
}

func (m *mockQuerier) GetJobApplications(_ context.Context, ownerID uuid.UUID) ([]repo.GetJobApplicationsRow, error) {
	return nil, nil
}

func (m *mockQuerier) UpdateJobApplication(_ context.Context, arg repo.UpdateJobApplicationParams) (repo.JobApplication, error) {
	return repo.JobApplication{}, nil
}

func newTestService() (*service, *mockQuerier) {
	q := newMockQuerier()
	config := &shared.Config{JwtSecret: []byte("testsecret")}
	return newService(q, config), q
}

func TestRegister_Success(t *testing.T) {
	s, _ := newTestService()

	tokens, err := s.register(context.Background(), &RegisterDTO{
		Email:    "test@test.com",
		Password: "password123",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
}

func TestRegister_HashesPassword(t *testing.T) {
	s, q := newTestService()

	dto := &RegisterDTO{Email: "test@test.com", Password: "password123"}
	_, err := s.register(context.Background(), dto)
	require.NoError(t, err)

	user := q.users["test@test.com"]
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
	assert.NoError(t, err)
}

func TestLogin_Success(t *testing.T) {
	s, q := newTestService()

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	q.users["test@test.com"] = repo.User{
		ID:       uuid.New(),
		Email:    "test@test.com",
		Password: string(hash),
	}

	tokens, err := s.login(context.Background(), &LoginDTO{
		Email:    "test@test.com",
		Password: "password123",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
}

func TestLogin_WrongPassword(t *testing.T) {
	s, q := newTestService()

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	q.users["test@test.com"] = repo.User{
		ID:       uuid.New(),
		Email:    "test@test.com",
		Password: string(hash),
	}

	_, err := s.login(context.Background(), &LoginDTO{
		Email:    "test@test.com",
		Password: "wrongpassword",
	})

	assert.Error(t, err)
}

func TestLogin_UserNotFound(t *testing.T) {
	s, _ := newTestService()

	_, err := s.login(context.Background(), &LoginDTO{
		Email:    "nobody@test.com",
		Password: "password123",
	})

	assert.Error(t, err)
}

func TestRefreshToken_Success(t *testing.T) {
	s, _ := newTestService()

	tokens, err := s.createTokens(context.Background(), uuid.New().String())
	require.NoError(t, err)

	time.Sleep(time.Second)
	newTokens, err := s.refreshToken(context.Background(), tokens.RefreshToken)

	require.NoError(t, err)
	assert.NotEmpty(t, newTokens.AccessToken)
	assert.NotEmpty(t, newTokens.RefreshToken)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	s, _ := newTestService()

	_, err := s.refreshToken(context.Background(), "invalidtoken")

	assert.Error(t, err)
}

func TestRefreshToken_TokenNotInDB(t *testing.T) {
	s, _ := newTestService()

	config := &shared.Config{JwtSecret: []byte("testsecret")}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": uuid.New().String(),
	})
	tokenStr, _ := token.SignedString(config.JwtSecret)

	_, err := s.refreshToken(context.Background(), tokenStr)

	assert.Error(t, err)
}
