package domain

import "github.com/google/uuid"

type Role string

const (
	RoleReporter  Role = "reporter"
	RoleDeveloper Role = "developer"
	RoleAdmin     Role = "admin"
)

type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

type Severity string

const (
	SeverityNull     Severity = ""
	SeverityMinor    Severity = "minor"
	SeverityMajor    Severity = "major"
	SeverityCritical Severity = "critical"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Role      Role      `db:"role" json:"role"`
	CreatedAt string    `db:"created_at" json:"created_at"`
}

type Project struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   string    `db:"created_at" json:"created_at"`
	UpdatedAt   string    `db:"updated_at" json:"updated_at"`
}

type Issue struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	ProjectID   uuid.UUID  `db:"project_id" json:"project_id"`
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	Status      string     `db:"status" json:"status"`
	AssigneeID  *uuid.UUID `db:"assignee_id" json:"assignee_id,omitempty"`
	Priority    Priority   `db:"priority" json:"priority"`
	IsIncident  bool       `db:"is_incident" json:"is_incident"`
	Severity    *Severity  `db:"severity" json:"severity,omitempty"`
	CreatedBy   uuid.UUID  `db:"created_by" json:"created_by"`
	CreatedAt   string     `db:"created_at" json:"created_at"`
	UpdatedAt   string     `db:"updated_at" json:"updated_at"`
}

type Comment struct {
	ID        uuid.UUID `db:"id" json:"id"`
	IssueID   uuid.UUID `db:"issue_id" json:"issue_id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Content   string    `db:"content" json:"content"`
	CreatedAt string    `db:"created_at" json:"created_at"`
}

type WorkflowState struct {
	ID         uuid.UUID `db:"id" json:"id"`
	WorkflowID uuid.UUID `db:"workflow_id" json:"workflow_id"`
	Name       string    `db:"name" json:"name"`
	IsInitial  bool      `db:"is_initial" json:"is_initial"`
	IsFinal    bool      `db:"is_final" json:"is_final"`
}

type WorkflowTransition struct {
	ID          uuid.UUID `db:"id" json:"id"`
	WorkflowID  uuid.UUID `db:"workflow_id" json:"workflow_id"`
	FromStateID uuid.UUID `db:"from_state_id" json:"from_state_id"`
	ToStateID   uuid.UUID `db:"to_state_id" json:"to_state_id"`
	Name        string    `db:"name" json:"name"`
}

type Workflow struct {
	ID          uuid.UUID            `db:"id" json:"id"`
	ProjectID   uuid.UUID            `db:"project_id" json:"project_id"`
	Name        string               `db:"name" json:"name"`
	States      []WorkflowState      `json:"states,omitempty"`
	Transitions []WorkflowTransition `json:"transitions,omitempty"`
}

type AuditLog struct {
	ID         uuid.UUID              `db:"id" json:"id"`
	EntityType string                 `db:"entity_type" json:"entity_type"`
	EntityID   uuid.UUID              `db:"entity_id" json:"entity_id"`
	Action     string                 `db:"action" json:"action"`
	UserID     uuid.UUID              `db:"user_id" json:"user_id"`
	OldValue   map[string]interface{} `db:"-" json:"old_value,omitempty"`
	NewValue   map[string]interface{} `db:"-" json:"new_value,omitempty"`
	CreatedAt  string                 `db:"created_at" json:"created_at"`
}

type ProjectMember struct {
	ID        uuid.UUID `db:"id" json:"id"`
	ProjectID uuid.UUID `db:"project_id" json:"project_id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Role      Role      `db:"role" json:"role"`
	CreatedAt string    `db:"created_at" json:"created_at"`
}
