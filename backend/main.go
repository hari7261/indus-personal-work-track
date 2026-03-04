package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"indus-task-manager/internal/app"
	"indus-task-manager/internal/db"
	"indus-task-manager/internal/domain"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

type App struct {
	ctx         context.Context
	application *app.App
}

func NewApp() *App {
	appDataDir := os.Getenv("APPDATA")
	if appDataDir == "" {
		appDataDir = "."
	}
	userDataDir := filepath.Join(appDataDir, "indus-task")
	dbPath := filepath.Join(userDataDir, "indus-task.db")

	if err := os.MkdirAll(userDataDir, 0o755); err != nil {
		panic(fmt.Sprintf("failed to create app data directory: %v", err))
	}

	database, err := db.NewDatabase(db.Config{DSN: dbPath})
	if err != nil {
		panic(fmt.Sprintf("failed to initialize database: %v", err))
	}

	application := app.NewApp(database)

	return &App{
		application: application,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Login(username string) (*domain.User, error) {
	return a.application.Login(username)
}

func (a *App) GetCurrentUser() *domain.User {
	return a.application.GetCurrentUser()
}

func (a *App) ListUsers() ([]domain.User, error) {
	return a.application.ListUsers()
}

func (a *App) CreateProject(req domain.CreateProjectRequest) (*domain.Project, error) {
	return a.application.ProjectService.Create(req, a.application.GetCurrentUserID())
}

func (a *App) ListProjects() ([]domain.Project, error) {
	return a.application.ProjectService.List()
}

func (a *App) GetProject(id uuid.UUID) (*domain.Project, error) {
	return a.application.ProjectService.GetByID(id)
}

func (a *App) UpdateProject(req domain.UpdateProjectRequest) (*domain.Project, error) {
	return a.application.ProjectService.Update(req, a.application.GetCurrentUserID())
}

func (a *App) DeleteProject(id uuid.UUID) error {
	return a.application.ProjectService.Delete(id, a.application.GetCurrentUserID())
}

func (a *App) CreateIssue(req domain.CreateIssueRequest) (*domain.Issue, error) {
	return a.application.IssueService.Create(req, a.application.GetCurrentUserID())
}

func (a *App) GetIssue(id uuid.UUID) (*domain.Issue, error) {
	return a.application.IssueService.GetByID(id)
}

func (a *App) ListIssues(filter domain.IssueFilter) (map[string]interface{}, error) {
	issues, total, err := a.application.IssueService.List(filter)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"data":  issues,
		"total": total,
		"page":  filter.Page,
	}, nil
}

func (a *App) UpdateIssue(req domain.UpdateIssueRequest) (*domain.Issue, error) {
	return a.application.IssueService.Update(req, a.application.GetCurrentUserID())
}

func (a *App) DeleteIssue(id uuid.UUID) error {
	return a.application.IssueService.Delete(id, a.application.GetCurrentUserID())
}

func (a *App) TransitionIssue(req domain.TransitionIssueRequest) (*domain.Issue, error) {
	return a.application.IssueService.Transition(req, a.application.GetCurrentUserID())
}

func (a *App) AssignIssue(req domain.AssignIssueRequest) (*domain.Issue, error) {
	return a.application.IssueService.Assign(req, a.application.GetCurrentUserID())
}

func (a *App) GetIssueStats(projectID uuid.UUID) (*domain.ProjectStats, error) {
	return a.application.IssueService.GetStats(projectID)
}

func (a *App) GetIssueComments(issueID uuid.UUID) ([]domain.Comment, error) {
	return a.application.IssueService.GetComments(issueID)
}

func (a *App) CreateComment(req domain.CreateCommentRequest) (*domain.Comment, error) {
	return a.application.IssueService.CreateComment(req, a.application.GetCurrentUserID())
}

func (a *App) GetAvailableTransitions(issueID uuid.UUID) ([]domain.WorkflowTransition, error) {
	return a.application.IssueService.GetAvailableTransitions(issueID)
}

func (a *App) CreateWorkflow(req domain.CreateWorkflowRequest) (*domain.Workflow, error) {
	return a.application.WorkflowService.Create(req)
}

func (a *App) GetWorkflow(projectID uuid.UUID) (*domain.Workflow, error) {
	return a.application.WorkflowService.GetByProjectID(projectID)
}

func (a *App) CreateWorkflowState(req domain.CreateWorkflowStateRequest) (*domain.WorkflowState, error) {
	return a.application.WorkflowService.CreateState(req)
}

func (a *App) CreateWorkflowTransition(req domain.CreateWorkflowTransitionRequest) (*domain.WorkflowTransition, error) {
	return a.application.WorkflowService.CreateTransition(req)
}

func (a *App) DeleteWorkflowState(stateID uuid.UUID) error {
	return a.application.WorkflowService.DeleteState(stateID)
}

func (a *App) DeleteWorkflowTransition(transitionID uuid.UUID) error {
	return a.application.WorkflowService.DeleteTransition(transitionID)
}

func (a *App) GetProjectMembers(projectID uuid.UUID) ([]domain.ProjectMember, error) {
	return a.application.ProjectService.GetMembers(projectID)
}

func (a *App) AddProjectMember(req domain.AddProjectMemberRequest) error {
	return a.application.ProjectService.AddMember(req, a.application.GetCurrentUserID())
}

func (a *App) RemoveProjectMember(projectID, memberID uuid.UUID) error {
	return a.application.ProjectService.RemoveMember(projectID, memberID, a.application.GetCurrentUserID())
}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "Indus Task Manager",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
