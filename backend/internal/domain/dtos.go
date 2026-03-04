package domain

import "github.com/google/uuid"

type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"max=1000"`
}

type UpdateProjectRequest struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required,max=100"`
	Description string    `json:"description" validate:"max=1000"`
}

type CreateIssueRequest struct {
	ProjectID   uuid.UUID `json:"project_id" validate:"required"`
	Title       string    `json:"title" validate:"required,max=200"`
	Description string    `json:"description"`
	Priority    Priority  `json:"priority" validate:"required"`
	IsIncident  bool      `json:"is_incident"`
	Severity    *Severity `json:"severity"`
}

type UpdateIssueRequest struct {
	ID          uuid.UUID  `json:"id" validate:"required"`
	Title       string     `json:"title" validate:"required,max=200"`
	Description string     `json:"description"`
	Priority    Priority   `json:"priority"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	IsIncident  bool       `json:"is_incident"`
	Severity    *Severity  `json:"severity"`
}

type TransitionIssueRequest struct {
	IssueID    uuid.UUID `json:"issue_id" validate:"required"`
	Transition string    `json:"transition" validate:"required"`
}

type AssignIssueRequest struct {
	IssueID    uuid.UUID  `json:"issue_id" validate:"required"`
	AssigneeID *uuid.UUID `json:"assignee_id"`
}

type CreateCommentRequest struct {
	IssueID uuid.UUID `json:"issue_id" validate:"required"`
	Content string    `json:"content" validate:"required,max=5000"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,max=50,alphanum"`
	Role     Role   `json:"role" validate:"required"`
}

type AddProjectMemberRequest struct {
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	Role      Role      `json:"role" validate:"required"`
}

type CreateWorkflowRequest struct {
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
	Name      string    `json:"name" validate:"required,max=100"`
}

type CreateWorkflowStateRequest struct {
	WorkflowID uuid.UUID `json:"workflow_id" validate:"required"`
	Name       string    `json:"name" validate:"required,max=50"`
	IsInitial  bool      `json:"is_initial"`
	IsFinal    bool      `json:"is_final"`
}

type CreateWorkflowTransitionRequest struct {
	WorkflowID  uuid.UUID `json:"workflow_id" validate:"required"`
	FromStateID uuid.UUID `json:"from_state_id" validate:"required"`
	ToStateID   uuid.UUID `json:"to_state_id" validate:"required"`
	Name        string    `json:"name" validate:"required,max=50"`
}

type IssueFilter struct {
	ProjectID  uuid.UUID  `json:"project_id"`
	Status     string     `json:"status"`
	AssigneeID *uuid.UUID `json:"assignee_id"`
	IsIncident *bool      `json:"is_incident"`
	Search     string     `json:"search"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
}

type PaginatedResponse struct {
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Data     interface{} `json:"data"`
}

type IssueListItem struct {
	ID            uuid.UUID  `json:"id"`
	ProjectID     uuid.UUID  `json:"project_id"`
	Title         string     `json:"title"`
	Status        string     `json:"status"`
	Priority      Priority   `json:"priority"`
	IsIncident    bool       `json:"is_incident"`
	Severity      *Severity  `json:"severity,omitempty"`
	AssigneeID    *uuid.UUID `json:"assignee_id,omitempty"`
	Assignee      string     `json:"assignee,omitempty"`
	CreatedBy     uuid.UUID  `json:"created_by"`
	CreatedByName string     `json:"created_by_name"`
	CreatedAt     string     `json:"created_at"`
	UpdatedAt     string     `json:"updated_at"`
}

type ProjectStats struct {
	Total         int `json:"total"`
	OpenIssues    int `json:"open_issues"`
	IncidentCount int `json:"incident_count"`
	CriticalCount int `json:"critical_count"`
}
