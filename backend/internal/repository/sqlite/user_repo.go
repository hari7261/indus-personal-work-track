package sqlite

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"indus-task-manager/internal/domain"
)

type UserRepository struct {
	db *sqlx.DB
}

type userRow struct {
	ID        uuid.UUID   `db:"id"`
	Username  string      `db:"username"`
	Role      domain.Role `db:"role"`
	CreatedAt string      `db:"created_at"`
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	user.ID = uuid.New()
	user.CreatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	_, err := r.db.Exec(
		"INSERT INTO users (id, username, role, created_at) VALUES (?, ?, ?, ?)",
		user.ID, user.Username, user.Role, user.CreatedAt,
	)
	return err
}

func (r *UserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	var row userRow
	err := r.db.Get(&row, "SELECT id, username, role, created_at FROM users WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	createdAt, err := parseSQLiteTime(row.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        row.ID,
		Username:  row.Username,
		Role:      row.Role,
		CreatedAt: createdAt,
	}, nil
}

func (r *UserRepository) GetByUsername(username string) (*domain.User, error) {
	var row userRow
	err := r.db.Get(&row, "SELECT id, username, role, created_at FROM users WHERE username = ?", username)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	createdAt, err := parseSQLiteTime(row.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        row.ID,
		Username:  row.Username,
		Role:      row.Role,
		CreatedAt: createdAt,
	}, nil
}

func (r *UserRepository) List() ([]domain.User, error) {
	var rows []userRow
	if err := r.db.Select(&rows, "SELECT id, username, role, created_at FROM users ORDER BY username"); err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, len(rows))
	for _, row := range rows {
		createdAt, err := parseSQLiteTime(row.CreatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, domain.User{
			ID:        row.ID,
			Username:  row.Username,
			Role:      row.Role,
			CreatedAt: createdAt,
		})
	}

	return users, nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}
