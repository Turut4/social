package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Password  password `json:"-"`
	Email     string   `json:"email"`
	RoleID    int64    `json:"role_id"`
	Role      Role
	CreatedAt string `json:"created_at"`
	IsActive  bool   `json:"is_active"`
}

type password struct {
	text *string
	hash []byte
}

type UserStore struct {
	db *sql.DB
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

func (p *password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(text))
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (username, password, email, role_id)
		VALUES ($1, $2, $3, (SELECT id FROM roles WHERE name = $4))
		RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	role := user.Role.Name
	if role == "" {
		role = "user"
	}

	err := tx.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password.hash,
		user.Email,
		role,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		switch err.Error() {
		case `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case `pq: duplicate key value violates unique constraint "users_username_key""`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) Delete(ctx context.Context, userID int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, userID); err != nil {
			return err
		}

		if err := s.deleteUserInvitation(ctx, tx, userID); err != nil {
			return err
		}

		return nil

	})
}

func (s *UserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT users.id, username, email, password, created_at, roles.*
		FROM users
		JOIN roles ON (users.role_id = roles.id)
		WHERE users.id = $1 AND is_active = true
	`

	user := &User{}
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.Level,
		&user.Role.Description,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStore) GetByUsernameOrEmail(ctx context.Context, emailOrUsername string) (*User, error) {
	query := `
		SELECT id, username, password, email FROM users
		WHERE username = $1 OR email = $1 AND is_active = true
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	if err := s.db.QueryRowContext(ctx, query, emailOrUsername).Scan(
		&user.ID,
		&user.Username,
		&user.Password.hash,
		&user.Email,
	); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}
func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := s.createUserInvitation(ctx, tx, token, exp, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		user.IsActive = true

		if err := s.update(ctx, tx, user); err != nil {
			return err
		}

		if err := s.deleteUserInvitation(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})

}

func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
			SELECT u.id, u.username, u.email, u.created_at, u.is_active
			FROM users u
			JOIN user_invitations ui ON u.id = ui.user_id
			WHERE ui.token = $1 AND ui.expiry > $2
		`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil

}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userID int64) error {
	query := `
		INSERT INTO user_invitations (user_id, token, expiry)
		VALUES ($1, $2, $3)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID, token, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		UPDATE users 
		SET username = $1, email = $2, is_active = $3
		WHERE id = $4
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) delete(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM user_invitations WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}
