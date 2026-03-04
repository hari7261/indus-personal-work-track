package app

import (
	"indus-task-manager/internal/audit"
	"indus-task-manager/internal/db"
	"indus-task-manager/internal/domain"
	"indus-task-manager/internal/permissions"
	"indus-task-manager/internal/repository/sqlite"
	"indus-task-manager/internal/services"
	"indus-task-manager/internal/workflow"

	"github.com/google/uuid"
)

type App struct {
	DB              *db.DB
	AuthService     *services.AuthService
	ProjectService  *services.ProjectService
	IssueService    *services.IssueService
	WorkflowService *services.WorkflowService
	currentUserID   uuid.UUID
	currentRole     domain.Role
}

func NewApp(db *db.DB) *App {
	userRepo := sqlite.NewUserRepository(db.DB)
	projectRepo := sqlite.NewProjectRepository(db.DB)
	issueRepo := sqlite.NewIssueRepository(db.DB)
	commentRepo := sqlite.NewCommentRepository(db.DB)
	workflowRepo := sqlite.NewWorkflowRepository(db.DB)
	auditRepo := sqlite.NewAuditRepository(db.DB)
	memberRepo := sqlite.NewProjectMemberRepository(db.DB)

	auditEngine := audit.NewEngine(auditRepo)

	defaultRole := domain.RoleReporter
	permChecker := permissions.NewPermissionChecker(defaultRole)
	workflowEngine := workflow.NewEngine(workflowRepo)

	authService := services.NewAuthService(userRepo)
	projectService := services.NewProjectService(projectRepo, memberRepo, auditEngine, permChecker)
	issueService := services.NewIssueService(issueRepo, commentRepo, workflowRepo, auditEngine, workflowEngine, permChecker)
	workflowService := services.NewWorkflowService(workflowRepo, permChecker)

	return &App{
		DB:              db,
		AuthService:     authService,
		ProjectService:  projectService,
		IssueService:    issueService,
		WorkflowService: workflowService,
		currentUserID:   uuid.Nil,
		currentRole:     defaultRole,
	}
}

func (a *App) Login(username string) (*domain.User, error) {
	user, err := a.AuthService.Login(username)
	if err != nil {
		return nil, err
	}

	a.currentUserID = user.ID
	a.currentRole = user.Role

	projectRepo := sqlite.NewProjectRepository(a.DB.DB)
	issueRepo := sqlite.NewIssueRepository(a.DB.DB)
	commentRepo := sqlite.NewCommentRepository(a.DB.DB)
	workflowRepo := sqlite.NewWorkflowRepository(a.DB.DB)
	auditRepo := sqlite.NewAuditRepository(a.DB.DB)
	memberRepo := sqlite.NewProjectMemberRepository(a.DB.DB)

	auditEngine := audit.NewEngine(auditRepo)
	permChecker := permissions.NewPermissionChecker(user.Role)
	workflowEngine := workflow.NewEngine(workflowRepo)

	a.ProjectService = services.NewProjectService(projectRepo, memberRepo, auditEngine, permChecker)
	a.IssueService = services.NewIssueService(issueRepo, commentRepo, workflowRepo, auditEngine, workflowEngine, permChecker)
	a.WorkflowService = services.NewWorkflowService(workflowRepo, permChecker)

	return user, nil
}

func (a *App) GetCurrentUser() *domain.User {
	return &domain.User{
		ID:       a.currentUserID,
		Username: "",
		Role:     a.currentRole,
	}
}

func (a *App) GetCurrentUserID() uuid.UUID {
	return a.currentUserID
}

func (a *App) ListUsers() ([]domain.User, error) {
	return a.AuthService.ListUsers()
}
