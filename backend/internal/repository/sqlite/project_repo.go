package sqlite

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"indus-task-manager/internal/domain"
)

type ProjectRepository struct {
	db *sqlx.DB
}

type projectRow struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   string    `db:"created_at"`
	UpdatedAt   string    `db:"updated_at"`
}

func NewProjectRepository(db *sqlx.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *domain.Project) error {
	project.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339Nano)
	project.CreatedAt = now
	project.UpdatedAt = now
	_, err := r.db.Exec(
		"INSERT INTO projects (id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		project.ID, project.Name, project.Description, project.CreatedAt, project.UpdatedAt,
	)
	return err
}

func (r *ProjectRepository) GetByID(id uuid.UUID) (*domain.Project, error) {
	var row projectRow
	err := r.db.Get(&row, "SELECT id, name, description, created_at, updated_at FROM projects WHERE id = ?", id)
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
	updatedAt, err := parseSQLiteTime(row.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &domain.Project{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func (r *ProjectRepository) List() ([]domain.Project, error) {
	var rows []projectRow
	if err := r.db.Select(&rows, "SELECT id, name, description, created_at, updated_at FROM projects ORDER BY name"); err != nil {
		return nil, err
	}

	projects := make([]domain.Project, 0, len(rows))
	for _, row := range rows {
		createdAt, err := parseSQLiteTime(row.CreatedAt)
		if err != nil {
			return nil, err
		}
		updatedAt, err := parseSQLiteTime(row.UpdatedAt)
		if err != nil {
			return nil, err
		}

		projects = append(projects, domain.Project{
			ID:          row.ID,
			Name:        row.Name,
			Description: row.Description,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
	}

	return projects, nil
}

func (r *ProjectRepository) Update(project *domain.Project) error {
	project.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	_, err := r.db.Exec(
		"UPDATE projects SET name = ?, description = ?, updated_at = ? WHERE id = ?",
		project.Name, project.Description, project.UpdatedAt, project.ID,
	)
	return err
}

func (r *ProjectRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM projects WHERE id = ?", id)
	return err
}
