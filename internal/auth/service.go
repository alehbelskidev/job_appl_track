package auth

import (
	"context"
	"errors"
	"time"

	"github.com/alehbelskidev/job_appl_track/internal/repo"
	"github.com/alehbelskidev/job_appl_track/internal/shared"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	config *shared.Config
	q      repo.Querier
}

func newService(q repo.Querier, config *shared.Config) *service {
	return &service{config: config, q: q}
}

func (s *service) register(ctx context.Context, dto *RegisterDTO) (*TokensDTO, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	dto.Password = string(hash)

	user, err := s.q.CreateUser(ctx, repo.CreateUserParams{Email: dto.Email, Password: dto.Password})
	if err != nil {
		return nil, err
	}

	tokens, err := s.createTokens(ctx, user.ID.String())
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *service) login(ctx context.Context, dto *LoginDTO) (*TokensDTO, error) {
	user, err := s.q.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		return nil, err
	}

	tokens, err := s.createTokens(ctx, user.ID.String())
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *service) refreshToken(ctx context.Context, token string) (*TokensDTO, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return s.config.JwtSecret, nil
	})
	if err != nil || !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	_, err = s.q.GetRefreshToken(ctx, token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims := parsed.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(string)

	if err := s.q.DeleteRefreshToken(ctx, token); err != nil {
		return nil, err
	}

	return s.createTokens(ctx, userId)
}

func (s *service) createTokens(ctx context.Context, userId string) (*TokensDTO, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(5 * time.Minute).Unix(),
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"jti":     uuid.New().String(),
	})

	accessTokenStr, err := accessToken.SignedString(s.config.JwtSecret)
	if err != nil {
		return nil, err
	}

	refreshTokenStr, err := refreshToken.SignedString(s.config.JwtSecret)
	if err != nil {
		return nil, err
	}

	err = s.q.AddRefreshToken(ctx, refreshTokenStr)
	if err != nil {
		return nil, err
	}

	return &TokensDTO{AccessToken: accessTokenStr, RefreshToken: refreshTokenStr}, nil
}
