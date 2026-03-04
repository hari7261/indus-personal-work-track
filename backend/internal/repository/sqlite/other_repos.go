package sqlite

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"indus-task-manager/internal/domain"
)

type CommentRepository struct {
	db *sqlx.DB
}

type commentRow struct {
	ID        uuid.UUID `db:"id"`
	IssueID   uuid.UUID `db:"issue_id"`
	UserID    uuid.UUID `db:"user_id"`
	Content   string    `db:"content"`
	CreatedAt string    `db:"created_at"`
}

func NewCommentRepository(db *sqlx.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *domain.Comment) error {
	comment.ID = uuid.New()
	comment.CreatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	_, err := r.db.Exec(
		"INSERT INTO comments (id, issue_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)",
		comment.ID, comment.IssueID, comment.UserID, comment.Content, comment.CreatedAt,
	)
	return err
}

func (r *CommentRepository) GetByIssueID(issueID uuid.UUID) ([]domain.Comment, error) {
	var rows []commentRow
	query := `SELECT c.id, c.issue_id, c.user_id, c.content, c.created_at 
			  FROM comments c WHERE c.issue_id = ? ORDER BY c.created_at ASC`
	if err := r.db.Select(&rows, query, issueID); err != nil {
		return nil, err
	}

	comments := make([]domain.Comment, 0, len(rows))
	for _, row := range rows {
		createdAt, err := parseSQLiteTime(row.CreatedAt)
		if err != nil {
			return nil, err
		}

		comments = append(comments, domain.Comment{
			ID:        row.ID,
			IssueID:   row.IssueID,
			UserID:    row.UserID,
			Content:   row.Content,
			CreatedAt: createdAt,
		})
	}

	return comments, nil
}

func (r *CommentRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM comments WHERE id = ?", id)
	return err
}

type WorkflowRepository struct {
	db *sqlx.DB
}

func NewWorkflowRepository(db *sqlx.DB) *WorkflowRepository {
	return &WorkflowRepository{db: db}
}

func (r *WorkflowRepository) CreateWorkflow(workflow *domain.Workflow) error {
	workflow.ID = uuid.New()
	_, err := r.db.Exec(
		"INSERT INTO workflows (id, project_id, name) VALUES (?, ?, ?)",
		workflow.ID, workflow.ProjectID, workflow.Name,
	)
	return err
}

func (r *WorkflowRepository) GetWorkflowByID(id uuid.UUID) (*domain.Workflow, error) {
	var workflow domain.Workflow
	err := r.db.Get(&workflow, "SELECT id, project_id, name FROM workflows WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	states, err := r.GetStatesByWorkflowID(workflow.ID)
	if err != nil {
		return nil, err
	}
	workflow.States = states
	transitions, err := r.GetTransitionsByWorkflowID(workflow.ID)
	if err != nil {
		return nil, err
	}
	workflow.Transitions = transitions
	return &workflow, nil
}

func (r *WorkflowRepository) GetWorkflowByProjectID(projectID uuid.UUID) (*domain.Workflow, error) {
	var workflow domain.Workflow
	err := r.db.Get(&workflow, "SELECT id, project_id, name FROM workflows WHERE project_id = ?", projectID)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	states, err := r.GetStatesByWorkflowID(workflow.ID)
	if err != nil {
		return nil, err
	}
	workflow.States = states
	transitions, err := r.GetTransitionsByWorkflowID(workflow.ID)
	if err != nil {
		return nil, err
	}
	workflow.Transitions = transitions
	return &workflow, nil
}

func (r *WorkflowRepository) UpdateWorkflow(workflow *domain.Workflow) error {
	_, err := r.db.Exec("UPDATE workflows SET name = ? WHERE id = ?", workflow.Name, workflow.ID)
	return err
}

func (r *WorkflowRepository) DeleteWorkflow(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM workflows WHERE id = ?", id)
	return err
}

func (r *WorkflowRepository) CreateState(state *domain.WorkflowState) error {
	state.ID = uuid.New()
	_, err := r.db.Exec(
		"INSERT INTO workflow_states (id, workflow_id, name, is_initial, is_final) VALUES (?, ?, ?, ?, ?)",
		state.ID, state.WorkflowID, state.Name, state.IsInitial, state.IsFinal,
	)
	return err
}

func (r *WorkflowRepository) GetStateByID(id uuid.UUID) (*domain.WorkflowState, error) {
	var state domain.WorkflowState
	err := r.db.Get(&state, "SELECT id, workflow_id, name, is_initial, is_final FROM workflow_states WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &state, err
}

func (r *WorkflowRepository) GetStatesByWorkflowID(workflowID uuid.UUID) ([]domain.WorkflowState, error) {
	var states []domain.WorkflowState
	err := r.db.Select(&states, "SELECT id, workflow_id, name, is_initial, is_final FROM workflow_states WHERE workflow_id = ?", workflowID)
	return states, err
}

func (r *WorkflowRepository) DeleteState(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM workflow_states WHERE id = ?", id)
	return err
}

func (r *WorkflowRepository) CreateTransition(transition *domain.WorkflowTransition) error {
	transition.ID = uuid.New()
	_, err := r.db.Exec(
		"INSERT INTO workflow_transitions (id, workflow_id, from_state_id, to_state_id, name) VALUES (?, ?, ?, ?, ?)",
		transition.ID, transition.WorkflowID, transition.FromStateID, transition.ToStateID, transition.Name,
	)
	return err
}

func (r *WorkflowRepository) GetTransitionsByWorkflowID(workflowID uuid.UUID) ([]domain.WorkflowTransition, error) {
	var transitions []domain.WorkflowTransition
	err := r.db.Select(&transitions, "SELECT id, workflow_id, from_state_id, to_state_id, name FROM workflow_transitions WHERE workflow_id = ?", workflowID)
	return transitions, err
}

func (r *WorkflowRepository) DeleteTransition(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM workflow_transitions WHERE id = ?", id)
	return err
}

type AuditRepository struct {
	db *sqlx.DB
}

type auditRow struct {
	ID         uuid.UUID `db:"id"`
	EntityType string    `db:"entity_type"`
	EntityID   uuid.UUID `db:"entity_id"`
	Action     string    `db:"action"`
	UserID     uuid.UUID `db:"user_id"`
	CreatedAt  string    `db:"created_at"`
}

func NewAuditRepository(db *sqlx.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(log *domain.AuditLog) error {
	log.ID = uuid.New()
	log.CreatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	_, err := r.db.Exec(
		"INSERT INTO audit_logs (id, entity_type, entity_id, action, user_id, old_value, new_value, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		log.ID, log.EntityType, log.EntityID, log.Action, log.UserID, log.OldValue, log.NewValue, log.CreatedAt,
	)
	return err
}

func (r *AuditRepository) GetByEntityID(entityType string, entityID uuid.UUID) ([]domain.AuditLog, error) {
	var rows []auditRow
	err := r.db.Select(&rows,
		"SELECT id, entity_type, entity_id, action, user_id, created_at FROM audit_logs WHERE entity_type = ? AND entity_id = ? ORDER BY created_at DESC",
		entityType,
		entityID,
	)
	if err != nil {
		return nil, err
	}

	logs := make([]domain.AuditLog, 0, len(rows))
	for _, row := range rows {
		createdAt, err := parseSQLiteTime(row.CreatedAt)
		if err != nil {
			return nil, err
		}

		logs = append(logs, domain.AuditLog{
			ID:         row.ID,
			EntityType: row.EntityType,
			EntityID:   row.EntityID,
			Action:     row.Action,
			UserID:     row.UserID,
			CreatedAt:  createdAt,
		})
	}

	return logs, nil
}

type ProjectMemberRepository struct {
	db *sqlx.DB
}

type projectMemberRow struct {
	ID        uuid.UUID   `db:"id"`
	ProjectID uuid.UUID   `db:"project_id"`
	UserID    uuid.UUID   `db:"user_id"`
	Role      domain.Role `db:"role"`
	CreatedAt string      `db:"created_at"`
}

func NewProjectMemberRepository(db *sqlx.DB) *ProjectMemberRepository {
	return &ProjectMemberRepository{db: db}
}

func (r *ProjectMemberRepository) Create(member *domain.ProjectMember) error {
	member.ID = uuid.New()
	member.CreatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	_, err := r.db.Exec(
		"INSERT INTO project_members (id, project_id, user_id, role, created_at) VALUES (?, ?, ?, ?, ?)",
		member.ID, member.ProjectID, member.UserID, member.Role, member.CreatedAt,
	)
	return err
}

func (r *ProjectMemberRepository) GetByProjectID(projectID uuid.UUID) ([]domain.ProjectMember, error) {
	var rows []projectMemberRow
	if err := r.db.Select(&rows, "SELECT id, project_id, user_id, role, created_at FROM project_members WHERE project_id = ?", projectID); err != nil {
		return nil, err
	}

	members := make([]domain.ProjectMember, 0, len(rows))
	for _, row := range rows {
		createdAt, err := parseSQLiteTime(row.CreatedAt)
		if err != nil {
			return nil, err
		}

		members = append(members, domain.ProjectMember{
			ID:        row.ID,
			ProjectID: row.ProjectID,
			UserID:    row.UserID,
			Role:      row.Role,
			CreatedAt: createdAt,
		})
	}

	return members, nil
}

func (r *ProjectMemberRepository) GetByUserAndProject(userID, projectID uuid.UUID) (*domain.ProjectMember, error) {
	var row projectMemberRow
	err := r.db.Get(&row, "SELECT id, project_id, user_id, role, created_at FROM project_members WHERE user_id = ? AND project_id = ?", userID, projectID)
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

	return &domain.ProjectMember{
		ID:        row.ID,
		ProjectID: row.ProjectID,
		UserID:    row.UserID,
		Role:      row.Role,
		CreatedAt: createdAt,
	}, nil
}

func (r *ProjectMemberRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM project_members WHERE id = ?", id)
	return err
}

func (r *ProjectMemberRepository) UpdateRole(id uuid.UUID, role domain.Role) error {
	_, err := r.db.Exec("UPDATE project_members SET role = ? WHERE id = ?", role, id)
	return err
}
