package app

import (
	"errors"
	"path/filepath"
	"testing"

	"indus-task-manager/internal/db"
	"indus-task-manager/internal/domain"

	"github.com/google/uuid"
)

func newTestApp(t *testing.T) *App {
	t.Helper()

	databasePath := filepath.Join(t.TempDir(), "indus-task-test.db")
	database, err := db.NewDatabase(db.Config{DSN: databasePath})
	if err != nil {
		t.Fatalf("failed to initialize test database: %v", err)
	}
	t.Cleanup(func() {
		_ = database.Close()
	})

	return NewApp(database)
}

func findUserID(users []domain.User, username string) (uuid.UUID, bool) {
	for _, user := range users {
		if user.Username == username {
			return user.ID, true
		}
	}
	return uuid.Nil, false
}

func TestAdminProjectIssueWorkflowLifecycle(t *testing.T) {
	application := newTestApp(t)

	admin, err := application.Login("admin")
	if err != nil {
		t.Fatalf("admin login failed: %v", err)
	}
	if admin.Role != domain.RoleAdmin {
		t.Fatalf("expected admin role, got %s", admin.Role)
	}

	users, err := application.ListUsers()
	if err != nil {
		t.Fatalf("list users failed: %v", err)
	}

	developerID, ok := findUserID(users, "developer")
	if !ok {
		t.Fatal("developer user not found in seed data")
	}

	project, err := application.ProjectService.Create(domain.CreateProjectRequest{
		Name:        "Lifecycle Project",
		Description: "Integration test project",
	}, application.GetCurrentUserID())
	if err != nil {
		t.Fatalf("create project failed: %v", err)
	}

	updatedProject, err := application.ProjectService.Update(domain.UpdateProjectRequest{
		ID:          project.ID,
		Name:        "Lifecycle Project Updated",
		Description: "Updated integration test project",
	}, application.GetCurrentUserID())
	if err != nil {
		t.Fatalf("update project failed: %v", err)
	}
	if updatedProject.Name != "Lifecycle Project Updated" {
		t.Fatalf("project name not updated, got %q", updatedProject.Name)
	}

	workflow, err := application.WorkflowService.Create(domain.CreateWorkflowRequest{
		ProjectID: project.ID,
		Name:      "Default",
	})
	if err != nil {
		t.Fatalf("create workflow failed: %v", err)
	}

	todoState, err := application.WorkflowService.CreateState(domain.CreateWorkflowStateRequest{
		WorkflowID: workflow.ID,
		Name:       "Todo",
		IsInitial:  true,
	})
	if err != nil {
		t.Fatalf("create initial state failed: %v", err)
	}

	inProgressState, err := application.WorkflowService.CreateState(domain.CreateWorkflowStateRequest{
		WorkflowID: workflow.ID,
		Name:       "In Progress",
	})
	if err != nil {
		t.Fatalf("create target state failed: %v", err)
	}

	transitionName := "Start Progress"
	_, err = application.WorkflowService.CreateTransition(domain.CreateWorkflowTransitionRequest{
		WorkflowID:  workflow.ID,
		FromStateID: todoState.ID,
		ToStateID:   inProgressState.ID,
		Name:        transitionName,
	})
	if err != nil {
		t.Fatalf("create transition failed: %v", err)
	}

	issue, err := application.IssueService.Create(domain.CreateIssueRequest{
		ProjectID:   project.ID,
		Title:       "Critical test issue",
		Description: "Verify full lifecycle",
		Priority:    domain.PriorityHigh,
		IsIncident:  true,
		Severity:    nil,
	}, application.GetCurrentUserID())
	if err != nil {
		t.Fatalf("create issue failed: %v", err)
	}
	if issue.Status != "Todo" {
		t.Fatalf("expected initial issue status Todo, got %q", issue.Status)
	}

	assigned, err := application.IssueService.Assign(domain.AssignIssueRequest{
		IssueID:    issue.ID,
		AssigneeID: &developerID,
	}, application.GetCurrentUserID())
	if err != nil {
		t.Fatalf("assign issue failed: %v", err)
	}
	if assigned.AssigneeID == nil || *assigned.AssigneeID != developerID {
		t.Fatal("issue assignee not updated")
	}

	transitioned, err := application.IssueService.Transition(domain.TransitionIssueRequest{
		IssueID:    issue.ID,
		Transition: transitionName,
	}, application.GetCurrentUserID())
	if err != nil {
		t.Fatalf("transition issue failed: %v", err)
	}
	if transitioned.Status != "In Progress" {
		t.Fatalf("expected status In Progress, got %q", transitioned.Status)
	}

	_, err = application.IssueService.CreateComment(domain.CreateCommentRequest{
		IssueID: issue.ID,
		Content: "Integration test comment",
	}, application.GetCurrentUserID())
	if err != nil {
		t.Fatalf("create comment failed: %v", err)
	}

	comments, err := application.IssueService.GetComments(issue.ID)
	if err != nil {
		t.Fatalf("get comments failed: %v", err)
	}
	if len(comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(comments))
	}

	err = application.ProjectService.AddMember(domain.AddProjectMemberRequest{
		ProjectID: project.ID,
		UserID:    developerID,
		Role:      domain.RoleDeveloper,
	}, application.GetCurrentUserID())
	if err != nil {
		t.Fatalf("add project member failed: %v", err)
	}

	members, err := application.ProjectService.GetMembers(project.ID)
	if err != nil {
		t.Fatalf("get project members failed: %v", err)
	}
	if len(members) != 1 {
		t.Fatalf("expected 1 project member, got %d", len(members))
	}

	err = application.ProjectService.RemoveMember(project.ID, members[0].ID, application.GetCurrentUserID())
	if err != nil {
		t.Fatalf("remove project member failed: %v", err)
	}

	err = application.IssueService.Delete(issue.ID, application.GetCurrentUserID())
	if err != nil {
		t.Fatalf("delete issue failed: %v", err)
	}

	_, err = application.IssueService.GetByID(issue.ID)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected issue not found after delete, got: %v", err)
	}

	err = application.ProjectService.Delete(project.ID, application.GetCurrentUserID())
	if err != nil {
		t.Fatalf("delete project failed: %v", err)
	}
}

func TestRolePermissionsAndNormalizedLogin(t *testing.T) {
	application := newTestApp(t)

	user, err := application.Login("  ADMIN  ")
	if err != nil {
		t.Fatalf("normalized admin login failed: %v", err)
	}
	if user.Role != domain.RoleAdmin {
		t.Fatalf("expected admin role after normalized login, got %s", user.Role)
	}

	_, err = application.Login("developer")
	if err != nil {
		t.Fatalf("developer login failed: %v", err)
	}

	_, err = application.ProjectService.Create(domain.CreateProjectRequest{
		Name:        "Forbidden Project",
		Description: "Developers should not create projects",
	}, application.GetCurrentUserID())
	if err == nil {
		t.Fatal("expected forbidden error when developer creates project")
	}
	if appErr, ok := err.(*domain.AppError); !ok || appErr.Code != "FORBIDDEN" {
		t.Fatalf("expected FORBIDDEN app error, got: %v", err)
	}
}
