package repository

import (
	"log"
	"tiktaktoe/internal/pkg/datasource"
	"tiktaktoe/internal/pkg/datasource/mapper"
	"tiktaktoe/internal/pkg/domain"

	"github.com/google/uuid"
)

type userRepository struct {
	Postgres datasource.Postgres
}

func NewUserRepository() userRepository {
	return userRepository{
		Postgres: datasource.NewPostgres(),
	}
}

func (r userRepository) CreateUser(user domain.User) error {
	pool := r.Postgres.Pool
	pool.Connect()
	defer pool.Close()

	query := `
        INSERT INTO users (id, login, password)
        VALUES ($1, $2, $3)
        ON CONFLICT (id) 
        DO UPDATE SET 
            login = EXCLUDED.login,
            password = EXCLUDED.password
    `
	err := pool.Exec(query,
		user.Id,
		user.Login,
		user.Password,
	)
	if err != nil {
		log.Fatalf("Query error: %v", err)
		return err
	}
	return nil
}

func (r userRepository) GetUser(id uuid.UUID) (domain.User, error) {
	pool := r.Postgres.Pool
	pool.Connect()
	defer pool.Close()

	query := `
		SELECT users.id, users.login, users.password
		FROM users
		WHERE users.id = $1
	`
	var user datasource.User
	row, err := pool.QueryRow(query, id.String())
	err = row.Scan(&user.Id, &user.Login, &user.Password)
	if err != nil {
		return domain.User{}, err
	}
	return *mapper.DatasourceToUser(&user), nil
}

func (r userRepository) GetUserByLogin(login string) (domain.User, error) {
	pool := r.Postgres.Pool
	pool.Connect()
	defer pool.Close()

	query := `
		SELECT users.id, users.login, users.password
		FROM users
		WHERE users.login = $1
	`
	var user datasource.User
	row, err := pool.QueryRow(query, login)
	err = row.Scan(&user.Id, &user.Login, &user.Password)
	if err != nil {
		return domain.User{}, err
	}
	return *mapper.DatasourceToUser(&user), nil
}

func (r userRepository) UpdateUser(user domain.User) (domain.User, error) {
	pool := r.Postgres.Pool
	pool.Connect()
	defer pool.Close()

	query := `
        INSERT INTO users (id, login, password)
        VALUES ($1, $2, $3)
        ON CONFLICT (id) 
        DO UPDATE SET 
            login = EXCLUDED.login,
            password = EXCLUDED.password
    `
	err := pool.Exec(query,
		user.Id,
		user.Login,
		user.Password,
	)
	if err != nil {
		log.Fatalf("Query error: %v", err)
		return user, err
	}
	return user, nil
}

func (r userRepository) DeleteUser(id uuid.UUID) error {
	pool := r.Postgres.Pool
	pool.Connect()
	defer pool.Close()

	query := `
        DELETE FROM users WHERE id = $1
    `
	err := pool.Exec(query, id)
	if err != nil {
		log.Fatalf("Query error: %v", err)
		return err
	}
	return nil
}
