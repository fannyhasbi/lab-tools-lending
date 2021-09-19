package postgres

import (
	"database/sql"

	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type UserRepositoryPostgres struct {
	DB *sql.DB
}

func NewUserRepositoryPostgres(DB *sql.DB) repository.UserRepository {
	return &UserRepositoryPostgres{
		DB: DB,
	}
}

func (ur *UserRepositoryPostgres) Save(user *types.User) (types.User, error) {
	row := ur.DB.QueryRow(`INSERT INTO users (id, name, nim, batch, address, user_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, nim, batch, address, created_at, user_type`, user.ID, user.Name, user.NIM, user.Batch, user.Address, user.UserType)

	u := types.User{}
	err := row.Scan(
		&u.ID,
		&u.Name,
		&u.NIM,
		&u.Batch,
		&u.Address,
		&u.CreatedAt,
		&u.UserType,
	)
	if err != nil {
		return types.User{}, err
	}

	return u, nil
}

func (ur *UserRepositoryPostgres) Update(user *types.User) (types.User, error) {
	row := ur.DB.QueryRow(`UPDATE users SET name = $1, nim = $2, batch = $3, address = $4
		WHERE id = $5
		RETURNING id, name, nim, batch, address, created_at`, user.Name, user.NIM, user.Batch, user.Address, user.ID)

	u := types.User{}
	err := row.Scan(
		&u.ID,
		&u.Name,
		&u.NIM,
		&u.Batch,
		&u.Address,
		&u.CreatedAt,
	)
	if err != nil {
		return types.User{}, err
	}

	return u, nil
}

func (ur *UserRepositoryPostgres) Delete(id int64) error {
	_, err := ur.DB.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}

func (ur *UserRepositoryPostgres) UpdateUserType(id int64, userType types.UserType) error {
	_, err := ur.DB.Exec(`UPDATE users SET user_type = $1 WHERE id = $2`, userType, id)
	return err
}
