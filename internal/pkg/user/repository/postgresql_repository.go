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

// Add new user profile in db
func (r *PostgresqlRepository) AddProfile(user *models.ProfileUser) (uint64, error) {
	row := r.db.QueryRow(
		"INSERT INTO users(email, password) VALUES ($1, $2) RETURNING id",
		user.Email,
		user.Password,
	)

	var userId uint64
	if err := row.Scan(&userId); err != nil {
		return 0, errors.ErrDBInternalError
	}

	return userId, nil
}

// Select one profile by email
func (r *PostgresqlRepository) SelectProfileByEmail(email string) (*models.ProfileUser, error) {
	row := r.db.QueryRow(
		"SELECT id, first_name, last_name, email, password, avatar "+
			"FROM users WHERE email = $1",
		email,
	)

	userByEmail := models.ProfileUser{}

	firstName := sql.NullString{}
	lastName := sql.NullString{}
	avatarUrl := sql.NullString{}
	err := row.Scan(
		&userByEmail.Id,
		&firstName,
		&lastName,
		&userByEmail.Email,
		&userByEmail.Password,
		&avatarUrl,
	)
	userByEmail.FirstName = firstName.String
	userByEmail.LastName = lastName.String
	userByEmail.Avatar.Url = avatarUrl.String

	switch err {
	case sql.ErrNoRows:
		return nil, errors.ErrUserNotFound
	case nil:
		return &userByEmail, nil
	default:
		return nil, errors.ErrDBInternalError
	}
}

// Select one profile by id
func (r *PostgresqlRepository) SelectProfileById(userId uint64) (*models.ProfileUser, error) {
	row := r.db.QueryRow(
		"SELECT id, first_name, last_name, email, password, avatar "+
			"FROM users WHERE id = $1",
		userId,
	)

	userByEmail := models.ProfileUser{}

	firstName := sql.NullString{}
	lastName := sql.NullString{}
	avatarUrl := sql.NullString{}
	err := row.Scan(
		&userByEmail.Id,
		&firstName,
		&lastName,
		&userByEmail.Email,
		&userByEmail.Password,
		&avatarUrl,
	)
	userByEmail.FirstName = firstName.String
	userByEmail.LastName = lastName.String
	userByEmail.Avatar.Url = avatarUrl.String

	switch err {
	case sql.ErrNoRows:
		return nil, errors.ErrUserNotFound
	case nil:
		return &userByEmail, nil
	default:
		return nil, errors.ErrDBInternalError
	}
}

// Update info in user profile
func (r *PostgresqlRepository) UpdateProfile(userId uint64, user *models.UpdateUser) error {
	_, err := r.db.Exec(
		"UPDATE users SET "+
			"first_name = $1, "+
			"last_name = $2 "+
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

// Update user avatar
func (r *PostgresqlRepository) UpdateAvatar(userId uint64, avatarUrl string) error {
	_, err := r.db.Exec(
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
