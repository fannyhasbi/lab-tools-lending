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
	row := ur.DB.QueryRow(`INSERT INTO users (id, name, nim, batch, address)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, nim, batch, address, created_at`, user.ID, user.Name, user.NIM, user.Batch, user.Address)

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
