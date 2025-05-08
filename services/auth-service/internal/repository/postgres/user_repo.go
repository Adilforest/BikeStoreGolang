package postgres

import (
	"BikeStoreGolang/services/auth-service/internal/domain"
	"context"
	"database/sql"
	"errors"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create добавляет нового пользователя в базу данных
func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	query := `
        INSERT INTO users (id, name, email, password_hash, role, is_active)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash, user.Role, user.IsActive)
	return err
}

// GetByID возвращает пользователя по ID
func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
        SELECT id, name, email, password_hash, role, is_active
        FROM users
        WHERE id = $1
    `
	row := r.db.QueryRowContext(ctx, query, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, role, is_active
		FROM users
		WHERE email = $1
	`
	row := r.db.QueryRowContext(ctx, query, email)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetAll возвращает всех пользователей
func (r *UserRepo) GetAll(ctx context.Context) ([]*domain.User, error) {
	query := `
        SELECT id, name, email, password_hash, role, is_active
        FROM users
    `
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// Update обновляет данные пользователя
func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	query := `
        UPDATE users
        SET name = $1, email = $2, password_hash = $3, role = $4, is_active = $5
        WHERE id = $6
    `
	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.PasswordHash, user.Role, user.IsActive, user.ID)
	return err
}

// DeleteByID удаляет пользователя по ID
func (r *UserRepo) DeleteByID(ctx context.Context, id string) error {
	query := `
        DELETE FROM users
        WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// DeleteAll удаляет всех пользователей
func (r *UserRepo) DeleteAll(ctx context.Context) error {
	query := `
        DELETE FROM users
    `
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *UserRepo) Count(ctx context.Context) (int, error) {
    var count int
    query := "SELECT COUNT(*) FROM users"
    err := r.db.QueryRowContext(ctx, query).Scan(&count)
    if err != nil {
        return 0, err
    }
    return count, nil
}

func (r *UserRepo) GetAllWithPagination(ctx context.Context, limit, offset int) ([]*domain.User, error) {
    query := "SELECT id, name, email FROM users LIMIT $1 OFFSET $2"
    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []*domain.User
    for rows.Next() {
        user := &domain.User{} // Create a pointer to domain.User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
            return nil, err
        }
        users = append(users, user) // Append the pointer to the slice
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return users, nil
}