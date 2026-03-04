package audit

import (
	"encoding/json"

	"indus-task-manager/internal/domain"
	"indus-task-manager/internal/repository"

	"github.com/google/uuid"
)

type Engine struct {
	auditRepo repository.AuditRepository
}

func NewEngine(auditRepo repository.AuditRepository) *Engine {
	return &Engine{auditRepo: auditRepo}
}

func (e *Engine) Log(entityType, action string, entityID, userID uuid.UUID, oldValue, newValue interface{}) error {
	oldJSON, _ := json.Marshal(oldValue)
	newJSON, _ := json.Marshal(newValue)

	log := &domain.AuditLog{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		UserID:     userID,
		OldValue:   nil,
		NewValue:   nil,
	}

	if len(oldJSON) > 0 {
		var old map[string]interface{}
		json.Unmarshal(oldJSON, &old)
		log.OldValue = old
	}

	if len(newJSON) > 0 {
		var new map[string]interface{}
		json.Unmarshal(newJSON, &new)
		log.NewValue = new
	}

	return e.auditRepo.Create(log)
}

func (e *Engine) LogProjectCreated(project *domain.Project, userID uuid.UUID) error {
	return e.Log("project", "created", project.ID, userID, nil, project)
}

func (e *Engine) LogProjectUpdated(oldProject, newProject *domain.Project, userID uuid.UUID) error {
	return e.Log("project", "updated", oldProject.ID, userID, oldProject, newProject)
}

func (e *Engine) LogProjectDeleted(project *domain.Project, userID uuid.UUID) error {
	return e.Log("project", "deleted", project.ID, userID, project, nil)
}

func (e *Engine) LogIssueCreated(issue *domain.Issue, userID uuid.UUID) error {
	return e.Log("issue", "created", issue.ID, userID, nil, issue)
}

func (e *Engine) LogIssueUpdated(oldIssue, newIssue *domain.Issue, userID uuid.UUID) error {
	return e.Log("issue", "updated", oldIssue.ID, userID, oldIssue, newIssue)
}

func (e *Engine) LogIssueDeleted(issue *domain.Issue, userID uuid.UUID) error {
	return e.Log("issue", "deleted", issue.ID, userID, issue, nil)
}

func (e *Engine) LogIssueTransitioned(issue *domain.Issue, fromStatus, toStatus string, userID uuid.UUID) error {
	return e.Log("issue", "transitioned", issue.ID, userID, map[string]string{"status": fromStatus}, map[string]string{"status": toStatus})
}

func (e *Engine) LogIssueAssigned(issue *domain.Issue, oldAssigneeID, newAssigneeID *uuid.UUID, userID uuid.UUID) error {
	return e.Log("issue", "assigned", issue.ID, userID, oldAssigneeID, newAssigneeID)
}

func (e *Engine) LogCommentCreated(comment *domain.Comment, userID uuid.UUID) error {
	return e.Log("comment", "created", comment.ID, userID, nil, comment)
}

func (e *Engine) LogCommentDeleted(comment *domain.Comment, userID uuid.UUID) error {
	return e.Log("comment", "deleted", comment.ID, userID, comment, nil)
}

func (e *Engine) GetEntityHistory(entityType string, entityID uuid.UUID) ([]domain.AuditLog, error) {
	return e.auditRepo.GetByEntityID(entityType, entityID)
}
