package permissions

import (
	"testing"

	"indus-task-manager/internal/domain"

	"github.com/google/uuid"
)

func TestReporterPermissions(t *testing.T) {
	pc := NewPermissionChecker(domain.RoleReporter)

	if pc.CanCreateProject() {
		t.Error("reporter should not be able to create projects")
	}

	if pc.CanUpdateProject() {
		t.Error("reporter should not be able to update projects")
	}

	if pc.CanDeleteProject() {
		t.Error("reporter should not be able to delete projects")
	}

	if !pc.CanCreateIssue("") {
		t.Error("reporter should be able to create issues")
	}

	if !pc.CanViewIssue() {
		t.Error("reporter should be able to view issues")
	}

	if !pc.CanCreateComment() {
		t.Error("reporter should be able to comment")
	}

	if pc.CanAssignIssue() {
		t.Error("reporter should not be able to assign issues")
	}

	if pc.CanTransitionIssue() {
		t.Error("reporter should not be able to transition issues")
	}

	if pc.CanManageWorkflow() {
		t.Error("reporter should not be able to manage workflows")
	}

	if pc.CanManageProjectMembers() {
		t.Error("reporter should not be able to manage project members")
	}
}

func TestDeveloperPermissions(t *testing.T) {
	pc := NewPermissionChecker(domain.RoleDeveloper)

	if pc.CanCreateProject() {
		t.Error("developer should not be able to create projects")
	}

	if !pc.CanCreateIssue("") {
		t.Error("developer should be able to create issues")
	}

	if !pc.CanAssignIssue() {
		t.Error("developer should be able to assign issues")
	}

	if !pc.CanTransitionIssue() {
		t.Error("developer should be able to transition issues")
	}

	if !pc.CanUpdateIssue(&domain.Issue{
		ID:        uuid.New(),
		CreatedBy: uuid.New(),
	}, uuid.New().String()) {
		t.Error("developer should be able to update any issue")
	}

	if pc.CanManageWorkflow() {
		t.Error("developer should not be able to manage workflows")
	}
}

func TestAdminPermissions(t *testing.T) {
	pc := NewPermissionChecker(domain.RoleAdmin)

	if !pc.CanCreateProject() {
		t.Error("admin should be able to create projects")
	}

	if !pc.CanUpdateProject() {
		t.Error("admin should be able to update projects")
	}

	if !pc.CanDeleteProject() {
		t.Error("admin should be able to delete projects")
	}

	if !pc.CanManageWorkflow() {
		t.Error("admin should be able to manage workflows")
	}

	if !pc.CanManageProjectMembers() {
		t.Error("admin should be able to manage project members")
	}

	issueID := uuid.New()
	userID := uuid.New()
	if !pc.CanUpdateIssue(&domain.Issue{
		ID:        issueID,
		CreatedBy: userID,
	}, userID.String()) {
		t.Error("admin should be able to update any issue")
	}
}

func TestCheckPermission(t *testing.T) {
	tests := []struct {
		role       domain.Role
		permission string
		expected   bool
	}{
		{domain.RoleReporter, "create_project", false},
		{domain.RoleDeveloper, "create_project", false},
		{domain.RoleAdmin, "create_project", true},
		{domain.RoleAdmin, "manage_workflow", true},
		{domain.RoleReporter, "manage_workflow", false},
	}

	for _, tt := range tests {
		result := CheckPermission(tt.role, tt.permission)
		if result != tt.expected {
			t.Errorf("CheckPermission(%s, %s) = %v; want %v", tt.role, tt.permission, result, tt.expected)
		}
	}
}

func TestIssueUpdatePermissions(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	issue := &domain.Issue{
		ID:        uuid.New(),
		CreatedBy: userID,
	}

	reporter := NewPermissionChecker(domain.RoleReporter)
	developer := NewPermissionChecker(domain.RoleDeveloper)
	admin := NewPermissionChecker(domain.RoleAdmin)

	if !reporter.CanUpdateIssue(issue, userID.String()) {
		t.Error("reporter should be able to update their own issue")
	}

	if reporter.CanUpdateIssue(issue, otherUserID.String()) {
		t.Error("reporter should not be able to update others' issues")
	}

	if !developer.CanUpdateIssue(issue, otherUserID.String()) {
		t.Error("developer should be able to update any issue")
	}

	if !admin.CanUpdateIssue(issue, otherUserID.String()) {
		t.Error("admin should be able to update any issue")
	}
}
