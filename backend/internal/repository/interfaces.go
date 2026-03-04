package repository

import (
	"github.com/google/uuid"
	"indus-task-manager/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id uuid.UUID) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	List() ([]domain.User, error)
	Delete(id uuid.UUID) error
}

type ProjectRepository interface {
	Create(project *domain.Project) error
	GetByID(id uuid.UUID) (*domain.Project, error)
	List() ([]domain.Project, error)
	Update(project *domain.Project) error
	Delete(id uuid.UUID) error
}

type IssueRepository interface {
	Create(issue *domain.Issue) error
	GetByID(id uuid.UUID) (*domain.Issue, error)
	List(filter domain.IssueFilter) ([]domain.IssueListItem, int, error)
	Update(issue *domain.Issue) error
	Delete(id uuid.UUID) error
	GetStats(projectID uuid.UUID) (*domain.ProjectStats, error)
}

type CommentRepository interface {
	Create(comment *domain.Comment) error
	GetByIssueID(issueID uuid.UUID) ([]domain.Comment, error)
	Delete(id uuid.UUID) error
}

type WorkflowRepository interface {
	CreateWorkflow(workflow *domain.Workflow) error
	GetWorkflowByID(id uuid.UUID) (*domain.Workflow, error)
	GetWorkflowByProjectID(projectID uuid.UUID) (*domain.Workflow, error)
	UpdateWorkflow(workflow *domain.Workflow) error
	DeleteWorkflow(id uuid.UUID) error

	CreateState(state *domain.WorkflowState) error
	GetStateByID(id uuid.UUID) (*domain.WorkflowState, error)
	GetStatesByWorkflowID(workflowID uuid.UUID) ([]domain.WorkflowState, error)
	DeleteState(id uuid.UUID) error

	CreateTransition(transition *domain.WorkflowTransition) error
	GetTransitionsByWorkflowID(workflowID uuid.UUID) ([]domain.WorkflowTransition, error)
	DeleteTransition(id uuid.UUID) error
}

type AuditRepository interface {
	Create(log *domain.AuditLog) error
	GetByEntityID(entityType string, entityID uuid.UUID) ([]domain.AuditLog, error)
}

type ProjectMemberRepository interface {
	Create(member *domain.ProjectMember) error
	GetByProjectID(projectID uuid.UUID) ([]domain.ProjectMember, error)
	GetByUserAndProject(userID, projectID uuid.UUID) (*domain.ProjectMember, error)
	Delete(id uuid.UUID) error
	UpdateRole(id uuid.UUID, role domain.Role) error
}
