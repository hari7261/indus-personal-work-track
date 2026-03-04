package sqlite

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"indus-task-manager/internal/domain"
)

type IssueRepository struct {
	db *sqlx.DB
}

type issueRow struct {
	ID          uuid.UUID        `db:"id"`
	ProjectID   uuid.UUID        `db:"project_id"`
	Title       string           `db:"title"`
	Description string           `db:"description"`
	Status      string           `db:"status"`
	AssigneeID  *uuid.UUID       `db:"assignee_id"`
	Priority    domain.Priority  `db:"priority"`
	IsIncident  bool             `db:"is_incident"`
	Severity    *domain.Severity `db:"severity"`
	CreatedBy   uuid.UUID        `db:"created_by"`
	CreatedAt   string           `db:"created_at"`
	UpdatedAt   string           `db:"updated_at"`
}

type issueListRow struct {
	ID            uuid.UUID        `db:"id"`
	ProjectID     uuid.UUID        `db:"project_id"`
	Title         string           `db:"title"`
	Status        string           `db:"status"`
	Priority      domain.Priority  `db:"priority"`
	IsIncident    bool             `db:"is_incident"`
	Severity      *domain.Severity `db:"severity"`
	AssigneeID    *uuid.UUID       `db:"assignee_id"`
	Assignee      sql.NullString   `db:"assignee"`
	CreatedBy     uuid.UUID        `db:"created_by"`
	CreatedByName sql.NullString   `db:"created_by_name"`
	CreatedAt     string           `db:"created_at"`
	UpdatedAt     string           `db:"updated_at"`
}

func NewIssueRepository(db *sqlx.DB) *IssueRepository {
	return &IssueRepository{db: db}
}

func (r *IssueRepository) Create(issue *domain.Issue) error {
	issue.ID = uuid.New()
	now := time.Now().UTC().Format(time.RFC3339Nano)
	issue.CreatedAt = now
	issue.UpdatedAt = now
	_, err := r.db.Exec(
		`INSERT INTO issues (id, project_id, title, description, status, assignee_id, priority, is_incident, severity, created_by, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		issue.ID, issue.ProjectID, issue.Title, issue.Description, issue.Status, issue.AssigneeID,
		issue.Priority, issue.IsIncident, issue.Severity, issue.CreatedBy, issue.CreatedAt, issue.UpdatedAt,
	)
	return err
}

func (r *IssueRepository) GetByID(id uuid.UUID) (*domain.Issue, error) {
	var row issueRow
	err := r.db.Get(&row, `SELECT id, project_id, title, description, status, assignee_id, priority, 
		is_incident, severity, created_by, created_at, updated_at FROM issues WHERE id = ?`, id)
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

	return &domain.Issue{
		ID:          row.ID,
		ProjectID:   row.ProjectID,
		Title:       row.Title,
		Description: row.Description,
		Status:      row.Status,
		AssigneeID:  row.AssigneeID,
		Priority:    row.Priority,
		IsIncident:  row.IsIncident,
		Severity:    row.Severity,
		CreatedBy:   row.CreatedBy,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func (r *IssueRepository) List(filter domain.IssueFilter) ([]domain.IssueListItem, int, error) {
	var issueRows []issueListRow
	var total int

	query := `SELECT i.id, i.project_id, i.title, i.status, i.priority, i.is_incident, i.severity, 
			  i.assignee_id, u.username as assignee, i.created_by, cr.username as created_by_name, i.created_at, i.updated_at
			  FROM issues i
			  LEFT JOIN users u ON i.assignee_id = u.id
			  LEFT JOIN users cr ON i.created_by = cr.id
			  WHERE 1=1`

	countQuery := `SELECT COUNT(*) FROM issues i WHERE 1=1`

	args := []interface{}{}

	if filter.ProjectID != uuid.Nil {
		query += " AND i.project_id = ?"
		countQuery += " AND i.project_id = ?"
		args = append(args, filter.ProjectID)
	}

	if filter.Status != "" {
		query += " AND i.status = ?"
		countQuery += " AND i.status = ?"
		args = append(args, filter.Status)
	}

	if filter.AssigneeID != nil {
		query += " AND i.assignee_id = ?"
		countQuery += " AND i.assignee_id = ?"
		args = append(args, *filter.AssigneeID)
	}

	if filter.IsIncident != nil {
		query += " AND i.is_incident = ?"
		countQuery += " AND i.is_incident = ?"
		if *filter.IsIncident {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}

	if filter.Search != "" {
		query += " AND (i.title LIKE ? OR i.description LIKE ?)"
		countQuery += " AND (i.title LIKE ? OR i.description LIKE ?)"
		search := "%" + filter.Search + "%"
		args = append(args, search, search)
	}

	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	query += " ORDER BY i.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)

	err = r.db.Select(&issueRows, query, args...)
	if err != nil {
		return nil, 0, err
	}

	issues := make([]domain.IssueListItem, 0, len(issueRows))
	for _, row := range issueRows {
		createdAt, err := parseSQLiteTime(row.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		updatedAt, err := parseSQLiteTime(row.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}

		issues = append(issues, domain.IssueListItem{
			ID:            row.ID,
			ProjectID:     row.ProjectID,
			Title:         row.Title,
			Status:        row.Status,
			Priority:      row.Priority,
			IsIncident:    row.IsIncident,
			Severity:      row.Severity,
			AssigneeID:    row.AssigneeID,
			Assignee:      row.Assignee.String,
			CreatedBy:     row.CreatedBy,
			CreatedByName: row.CreatedByName.String,
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
		})
	}

	return issues, total, nil
}

func (r *IssueRepository) Update(issue *domain.Issue) error {
	issue.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	_, err := r.db.Exec(
		`UPDATE issues SET title = ?, description = ?, status = ?, assignee_id = ?, 
		 priority = ?, is_incident = ?, severity = ?, updated_at = ? WHERE id = ?`,
		issue.Title, issue.Description, issue.Status, issue.AssigneeID,
		issue.Priority, issue.IsIncident, issue.Severity, issue.UpdatedAt, issue.ID,
	)
	return err
}

func (r *IssueRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM issues WHERE id = ?", id)
	return err
}

func (r *IssueRepository) GetStats(projectID uuid.UUID) (*domain.ProjectStats, error) {
	var stats domain.ProjectStats

	err := r.db.Get(&stats.Total, "SELECT COUNT(*) FROM issues WHERE project_id = ?", projectID)
	if err != nil {
		return nil, err
	}

	err = r.db.Get(&stats.OpenIssues, `SELECT COUNT(*) FROM issues i 
		JOIN workflow_states ws ON i.status = ws.name AND ws.workflow_id = (
			SELECT id FROM workflows WHERE project_id = i.project_id
		)
		WHERE i.project_id = ? AND ws.is_final = 0`, projectID)
	if err != nil {
		stats.OpenIssues = stats.Total
	}

	err = r.db.Get(&stats.IncidentCount, "SELECT COUNT(*) FROM issues WHERE project_id = ? AND is_incident = 1", projectID)
	if err != nil {
		return nil, err
	}

	err = r.db.Get(&stats.CriticalCount, "SELECT COUNT(*) FROM issues WHERE project_id = ? AND severity = 'critical'", projectID)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
