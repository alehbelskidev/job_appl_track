package auth

import (
	"database/sql"
	"errors"
)

type repository struct {
	db *sql.DB
}

func newRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) getUser(email string) (*User, error) {
	userRow := &User{}
	err := r.db.QueryRow("SELECT id, email, password FROM users WHERE email=?", email).Scan(&userRow.ID, &userRow.Email, &userRow.Password)

	if err == sql.ErrNoRows {
		return nil, errors.New("not found")
	}
	if err != nil {
		return nil, err
	}

	return userRow, nil
}

func (r *repository) createUser(dto *RegisterDTO) error {
	_, err := r.db.Exec(
		"INSERT INTO users(email, password) VALUES(?, ?)",
		dto.Email, dto.Password,
	)

	return err
}

func (r *repository) createRefreshToken(rt string) error {
	_, err := r.db.Exec(
		"INSERT INTO refresh_tokens(token) VALUES(?)",
		rt,
	)

	return err
}

func (r *repository) getRefreshToken(rt string) (*string, error) {
	var token *string
	err := r.db.QueryRow("SELECT token FROM refresh_tokens WHERE token=?", rt).Scan(&token)

	if err == sql.ErrNoRows {
		return nil, errors.New("not found")
	}
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (r *repository) deleteRefreshToken(rt string) error {
	_, err := r.db.Exec("DELETE FROM refresh_tokens WHERE token=?", rt)
	return err
}
