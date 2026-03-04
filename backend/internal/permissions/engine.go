package permissions

import (
	"indus-task-manager/internal/domain"
)

type PermissionChecker struct {
	role domain.Role
}

func NewPermissionChecker(role domain.Role) *PermissionChecker {
	return &PermissionChecker{role: role}
}

func (p *PermissionChecker) CanCreateProject() bool {
	return p.role == domain.RoleAdmin
}

func (p *PermissionChecker) CanUpdateProject() bool {
	return p.role == domain.RoleAdmin
}

func (p *PermissionChecker) CanDeleteProject() bool {
	return p.role == domain.RoleAdmin
}

func (p *PermissionChecker) CanManageProjectMembers() bool {
	return p.role == domain.RoleAdmin
}

func (p *PermissionChecker) CanManageWorkflow() bool {
	return p.role == domain.RoleAdmin
}

func (p *PermissionChecker) CanCreateIssue(projectID string) bool {
	return p.role == domain.RoleReporter || p.role == domain.RoleDeveloper || p.role == domain.RoleAdmin
}

func (p *PermissionChecker) CanUpdateIssue(issue *domain.Issue, userID string) bool {
	if p.role == domain.RoleAdmin {
		return true
	}
	if p.role == domain.RoleDeveloper {
		return true
	}
	return issue.CreatedBy.String() == userID
}

func (p *PermissionChecker) CanDeleteIssue(issue *domain.Issue) bool {
	return p.role == domain.RoleAdmin || p.role == domain.RoleDeveloper
}

func (p *PermissionChecker) CanAssignIssue() bool {
	return p.role == domain.RoleDeveloper || p.role == domain.RoleAdmin
}

func (p *PermissionChecker) CanTransitionIssue() bool {
	return p.role == domain.RoleDeveloper || p.role == domain.RoleAdmin
}

func (p *PermissionChecker) CanViewIssue() bool {
	return p.role == domain.RoleReporter || p.role == domain.RoleDeveloper || p.role == domain.RoleAdmin
}

func (p *PermissionChecker) CanCreateComment() bool {
	return p.role == domain.RoleReporter || p.role == domain.RoleDeveloper || p.role == domain.RoleAdmin
}

func (p *PermissionChecker) CanDeleteComment(comment *domain.Comment, userID string) bool {
	if p.role == domain.RoleAdmin {
		return true
	}
	return comment.UserID.String() == userID
}

func (p *PermissionChecker) CanManageUsers() bool {
	return p.role == domain.RoleAdmin
}

func CheckPermission(role domain.Role, permission string) bool {
	pc := NewPermissionChecker(role)
	switch permission {
	case "create_project":
		return pc.CanCreateProject()
	case "update_project":
		return pc.CanUpdateProject()
	case "delete_project":
		return pc.CanDeleteProject()
	case "manage_project_members":
		return pc.CanManageProjectMembers()
	case "manage_workflow":
		return pc.CanManageWorkflow()
	default:
		return false
	}
}
