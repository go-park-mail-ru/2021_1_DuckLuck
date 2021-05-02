package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/pkg/user"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/server/errors"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func NewSessionPostgresqlRepository(db *sql.DB) user.Repository {
	return &PostgresqlRepository{
		db: db,
	}
}

func (r *PostgresqlRepository) AddProfile(newUser *models.AuthUser) (uint64, error) {
	row := r.db.QueryRow(
		"INSERT INTO auth_users(email, password) "+
			"VALUES ($1, $2) RETURNING id",
		newUser.Email,
		newUser.Password,
	)

	var userId uint64
	if err := row.Scan(&userId); err != nil {
		return 0, errors.ErrDBInternalError
	}

	return userId, nil
}

func (r *PostgresqlRepository) SelectUserByEmail(email string) (*models.AuthUser, error) {
	row := r.db.QueryRow(
		"SELECT id, email, password "+
			"FROM auth_users WHERE email = $1",
		email,
	)

	userByEmail := models.AuthUser{}
	err := row.Scan(
		&userByEmail.Id,
		&userByEmail.Email,
		&userByEmail.Password,
	)

	switch err {
	case sql.ErrNoRows:
		return nil, errors.ErrUserNotFound
	case nil:
		return &userByEmail, nil
	default:
		return nil, errors.ErrDBInternalError
	}
}
