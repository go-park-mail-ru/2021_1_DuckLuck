package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func NewSessionPostgresqlRepository(db *sql.DB) user.Repository {
	return &PostgresqlRepository{
		db: db,
	}
}

func (pr *PostgresqlRepository) AddProfile(user *models.SignupUser) (uint64, error) {
	row := pr.db.QueryRow(
		"INSERT INTO users(email, password) VALUES ($1, $2) RETURNING id",
		user.Email,
		user.Password,
	)

	var userId uint64
	err := row.Scan(&userId)

	if err != nil {
		return 0, errors.ErrDBInternalError
	}

	return userId, nil
}

func (pr *PostgresqlRepository) SelectProfileByEmail(email string) (*models.ProfileUser, error) {
	row := pr.db.QueryRow(
		"SELECT id, firstName, lastName, email, password, avatar "+
		"FROM users WHERE email = $1",
		email,
	)

	userByEmail := models.ProfileUser{}

	err := row.Scan(
		&userByEmail.Id,
		&userByEmail.FirstName,
		&userByEmail.LastName,
		&userByEmail.Email,
		&userByEmail.Password,
		&userByEmail.Avatar.Url,
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

func (pr *PostgresqlRepository) SelectProfileById(userId uint64) (*models.ProfileUser, error) {
	row := pr.db.QueryRow(
		"SELECT id, firstName, lastName, email, password, avatar "+
		"FROM users WHERE id = $1",
		userId,
	)

	userByEmail := models.ProfileUser{}

	err := row.Scan(
		&userByEmail.Id,
		&userByEmail.FirstName,
		&userByEmail.LastName,
		&userByEmail.Email,
		&userByEmail.Password,
		&userByEmail.Avatar.Url,
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

func (pr *PostgresqlRepository) UpdateProfile(userId uint64, user *models.UpdateUser) error {
	_, err := pr.db.Exec(
		"UPDATE users SET "+
		"firstName = $1, "+
		"lastName = $2 "+
		"WHERE id = $3",
		user.FirstName,
		user.LastName,
		userId,
	)
	if err != nil {
		return errors.ErrDBInternalError
	}

	return nil
}

func (pr *PostgresqlRepository) UpdateAvatar(userId uint64, avatarUrl string) error {
	_, err := pr.db.Exec(
		"UPDATE users SET "+
		"avatar = $1 "+
		"WHERE id = $2",
		avatarUrl,
		userId,
	)
	if err != nil {
		return errors.ErrDBInternalError
	}

	return nil
}
