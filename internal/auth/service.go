package auth

import (
	"errors"
	"time"

	"github.com/alehbelskidev/job_appl_track/internal/shared"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo   *repository
	config *shared.Config
}

func newService(repo *repository, config *shared.Config) *service {
	return &service{repo: repo, config: config}
}

func (s *service) register(dto *RegisterDTO) (*TokensDTO, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	dto.Password = string(hash)

	err = s.repo.createUser(dto)
	if err != nil {
		return nil, err
	}

	tokens, err := s.createTokens(dto.Email)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *service) login(dto *LoginDTO) (*TokensDTO, error) {
	user, err := s.repo.getUser(dto.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		return nil, err
	}

	tokens, err := s.createTokens(dto.Email)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *service) refreshToken(token string) (*TokensDTO, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return s.config.JwtSecret, nil
	})
	if err != nil || !parsed.Valid {
		return nil, errors.New("invalid token")
	}

	claims := parsed.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	if err := s.repo.deleteRefreshToken(token); err != nil {
		return nil, err
	}

	return s.createTokens(email)
}

func (s *service) createTokens(email string) (*TokensDTO, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(5 * time.Minute).Unix(),
	})
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	accessTokenStr, err := accessToken.SignedString(s.config.JwtSecret)
	if err != nil {
		return nil, err
	}

	refreshTokenStr, err := refreshToken.SignedString(s.config.JwtSecret)
	if err != nil {
		return nil, err
	}

	err = s.repo.createRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, err
	}

	return &TokensDTO{AccessToken: accessTokenStr, RefreshToken: refreshTokenStr}, nil
}
