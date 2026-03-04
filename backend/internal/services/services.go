package services

import (
	"errors"
	"indus-task-manager/internal/audit"
	"indus-task-manager/internal/domain"
	"indus-task-manager/internal/permissions"
	"indus-task-manager/internal/repository"
	"indus-task-manager/internal/workflow"
	"strings"

	"github.com/google/uuid"
)

type AppService struct {
	userRepo          repository.UserRepository
	projectRepo       repository.ProjectRepository
	issueRepo         repository.IssueRepository
	commentRepo       repository.CommentRepository
	workflowRepo      repository.WorkflowRepository
	auditRepo         repository.AuditRepository
	memberRepo        repository.ProjectMemberRepository
	permissionChecker *permissions.PermissionChecker
	workflowEngine    *workflow.Engine
	auditEngine       *audit.Engine
}

func NewAppService(
	userRepo repository.UserRepository,
	projectRepo repository.ProjectRepository,
	issueRepo repository.IssueRepository,
	commentRepo repository.CommentRepository,
	workflowRepo repository.WorkflowRepository,
	auditRepo repository.AuditRepository,
	memberRepo repository.ProjectMemberRepository,
	role domain.Role,
) *AppService {
	return &AppService{
		userRepo:          userRepo,
		projectRepo:       projectRepo,
		issueRepo:         issueRepo,
		commentRepo:       commentRepo,
		workflowRepo:      workflowRepo,
		auditRepo:         auditRepo,
		memberRepo:        memberRepo,
		permissionChecker: permissions.NewPermissionChecker(role),
		workflowEngine:    workflow.NewEngine(workflowRepo),
		auditEngine:       audit.NewEngine(auditRepo),
	}
}

func (s *AppService) GetCurrentUser() (*domain.User, error) {
	return nil, nil
}

func (s *AppService) SetCurrentUser(userID uuid.UUID, role domain.Role) {
	s.permissionChecker = permissions.NewPermissionChecker(role)
	s.workflowEngine = workflow.NewEngine(s.workflowRepo)
	s.auditEngine = audit.NewEngine(s.auditRepo)
}

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Login(username string) (*domain.User, error) {
	normalized := strings.ToLower(strings.TrimSpace(username))
	if normalized == "" {
		return nil, domain.NewUnauthorizedError("invalid credentials")
	}

	user, err := s.userRepo.GetByUsername(normalized)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.NewUnauthorizedError("invalid credentials")
		}
		return nil, domain.NewInvalidInputError("login failed: " + err.Error())
	}
	if user == nil {
		return nil, domain.NewUnauthorizedError("invalid credentials")
	}
	return user, nil
}

func (s *AuthService) ListUsers() ([]domain.User, error) {
	return s.userRepo.List()
}

type ProjectService struct {
	repo        repository.ProjectRepository
	memberRepo  repository.ProjectMemberRepository
	auditEngine *audit.Engine
	permChecker *permissions.PermissionChecker
}

func NewProjectService(
	repo repository.ProjectRepository,
	memberRepo repository.ProjectMemberRepository,
	auditEngine *audit.Engine,
	permChecker *permissions.PermissionChecker,
) *ProjectService {
	return &ProjectService{repo: repo, memberRepo: memberRepo, auditEngine: auditEngine, permChecker: permChecker}
}

func (s *ProjectService) Create(req domain.CreateProjectRequest, userID uuid.UUID) (*domain.Project, error) {
	if !s.permChecker.CanCreateProject() {
		return nil, domain.NewForbiddenError("only admins can create projects")
	}

	project := &domain.Project{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.Create(project); err != nil {
		return nil, err
	}

	s.auditEngine.LogProjectCreated(project, userID)
	return project, nil
}

func (s *ProjectService) GetByID(id uuid.UUID) (*domain.Project, error) {
	return s.repo.GetByID(id)
}

func (s *ProjectService) List() ([]domain.Project, error) {
	return s.repo.List()
}

func (s *ProjectService) Update(req domain.UpdateProjectRequest, userID uuid.UUID) (*domain.Project, error) {
	if !s.permChecker.CanUpdateProject() {
		return nil, domain.NewForbiddenError("only admins can update projects")
	}

	project, err := s.repo.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	oldProject := *project
	project.Name = req.Name
	project.Description = req.Description

	if err := s.repo.Update(project); err != nil {
		return nil, err
	}

	s.auditEngine.LogProjectUpdated(&oldProject, project, userID)
	return project, nil
}

func (s *ProjectService) Delete(id uuid.UUID, userID uuid.UUID) error {
	if !s.permChecker.CanDeleteProject() {
		return domain.NewForbiddenError("only admins can delete projects")
	}

	project, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	s.auditEngine.LogProjectDeleted(project, userID)
	return nil
}

func (s *ProjectService) GetMembers(projectID uuid.UUID) ([]domain.ProjectMember, error) {
	return s.memberRepo.GetByProjectID(projectID)
}

func (s *ProjectService) AddMember(req domain.AddProjectMemberRequest, userID uuid.UUID) error {
	if !s.permChecker.CanManageProjectMembers() {
		return domain.NewForbiddenError("only admins can manage project members")
	}

	member := &domain.ProjectMember{
		ProjectID: req.ProjectID,
		UserID:    req.UserID,
		Role:      req.Role,
	}

	return s.memberRepo.Create(member)
}

func (s *ProjectService) RemoveMember(projectID, memberID uuid.UUID, userID uuid.UUID) error {
	if !s.permChecker.CanManageProjectMembers() {
		return domain.NewForbiddenError("only admins can manage project members")
	}

	return s.memberRepo.Delete(memberID)
}

type IssueService struct {
	repo           repository.IssueRepository
	commentRepo    repository.CommentRepository
	workflowRepo   repository.WorkflowRepository
	auditEngine    *audit.Engine
	workflowEngine *workflow.Engine
	permChecker    *permissions.PermissionChecker
}

func NewIssueService(
	repo repository.IssueRepository,
	commentRepo repository.CommentRepository,
	workflowRepo repository.WorkflowRepository,
	auditEngine *audit.Engine,
	workflowEngine *workflow.Engine,
	permChecker *permissions.PermissionChecker,
) *IssueService {
	return &IssueService{
		repo: repo, commentRepo: commentRepo, workflowRepo: workflowRepo,
		auditEngine: auditEngine, workflowEngine: workflowEngine, permChecker: permChecker,
	}
}

func (s *IssueService) Create(req domain.CreateIssueRequest, userID uuid.UUID) (*domain.Issue, error) {
	if !s.permChecker.CanCreateIssue(req.ProjectID.String()) {
		return nil, domain.NewForbiddenError("you cannot create issues")
	}

	workflow, err := s.workflowRepo.GetWorkflowByProjectID(req.ProjectID)
	if err != nil {
		return nil, domain.NewInvalidInputError("project has no workflow configured")
	}

	initialStatus, err := s.workflowEngine.GetInitialState(workflow.ID)
	if err != nil {
		return nil, domain.NewInvalidInputError("workflow has no initial state")
	}

	issue := &domain.Issue{
		ProjectID:   req.ProjectID,
		Title:       req.Title,
		Description: req.Description,
		Status:      initialStatus,
		Priority:    req.Priority,
		IsIncident:  req.IsIncident,
		Severity:    req.Severity,
		CreatedBy:   userID,
	}

	if issue.Severity == nil && issue.IsIncident {
		sev := domain.SeverityMinor
		issue.Severity = &sev
	}

	if err := s.repo.Create(issue); err != nil {
		return nil, err
	}

	s.auditEngine.LogIssueCreated(issue, userID)
	return issue, nil
}

func (s *IssueService) GetByID(id uuid.UUID) (*domain.Issue, error) {
	issue, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if !s.permChecker.CanViewIssue() {
		return nil, domain.NewForbiddenError("you cannot view issues")
	}

	return issue, nil
}

func (s *IssueService) List(filter domain.IssueFilter) ([]domain.IssueListItem, int, error) {
	if !s.permChecker.CanViewIssue() {
		return nil, 0, domain.NewForbiddenError("you cannot view issues")
	}
	return s.repo.List(filter)
}

func (s *IssueService) Update(req domain.UpdateIssueRequest, userID uuid.UUID) (*domain.Issue, error) {
	issue, err := s.repo.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	if !s.permChecker.CanUpdateIssue(issue, userID.String()) {
		return nil, domain.NewForbiddenError("you cannot update this issue")
	}

	oldIssue := *issue
	issue.Title = req.Title
	issue.Description = req.Description
	issue.Priority = req.Priority
	issue.AssigneeID = req.AssigneeID
	issue.IsIncident = req.IsIncident
	issue.Severity = req.Severity

	if err := s.repo.Update(issue); err != nil {
		return nil, err
	}

	s.auditEngine.LogIssueUpdated(&oldIssue, issue, userID)
	return issue, nil
}

func (s *IssueService) Delete(id uuid.UUID, userID uuid.UUID) error {
	issue, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if !s.permChecker.CanDeleteIssue(issue) {
		return domain.NewForbiddenError("you cannot delete this issue")
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	s.auditEngine.LogIssueDeleted(issue, userID)
	return nil
}

func (s *IssueService) Transition(req domain.TransitionIssueRequest, userID uuid.UUID) (*domain.Issue, error) {
	if !s.permChecker.CanTransitionIssue() {
		return nil, domain.NewForbiddenError("you cannot transition issues")
	}

	issue, err := s.repo.GetByID(req.IssueID)
	if err != nil {
		return nil, err
	}

	workflow, err := s.workflowRepo.GetWorkflowByProjectID(issue.ProjectID)
	if err != nil {
		return nil, domain.NewInvalidInputError("project has no workflow")
	}

	canTransition, transition, err := s.workflowEngine.CanTransition(workflow.ID, issue.Status, req.Transition)
	if err != nil {
		return nil, err
	}

	if !canTransition {
		return nil, domain.NewInvalidTransitionError("invalid transition from " + issue.Status)
	}

	toState, err := s.workflowRepo.GetStateByID(transition.ToStateID)
	if err != nil {
		return nil, err
	}

	oldStatus := issue.Status
	issue.Status = toState.Name

	if err := s.repo.Update(issue); err != nil {
		return nil, err
	}

	s.auditEngine.LogIssueTransitioned(issue, oldStatus, issue.Status, userID)
	return issue, nil
}

func (s *IssueService) Assign(req domain.AssignIssueRequest, userID uuid.UUID) (*domain.Issue, error) {
	if !s.permChecker.CanAssignIssue() {
		return nil, domain.NewForbiddenError("you cannot assign issues")
	}

	issue, err := s.repo.GetByID(req.IssueID)
	if err != nil {
		return nil, err
	}

	oldAssigneeID := issue.AssigneeID
	issue.AssigneeID = req.AssigneeID

	if err := s.repo.Update(issue); err != nil {
		return nil, err
	}

	s.auditEngine.LogIssueAssigned(issue, oldAssigneeID, req.AssigneeID, userID)
	return issue, nil
}

func (s *IssueService) GetStats(projectID uuid.UUID) (*domain.ProjectStats, error) {
	return s.repo.GetStats(projectID)
}

func (s *IssueService) GetComments(issueID uuid.UUID) ([]domain.Comment, error) {
	return s.commentRepo.GetByIssueID(issueID)
}

func (s *IssueService) CreateComment(req domain.CreateCommentRequest, userID uuid.UUID) (*domain.Comment, error) {
	if !s.permChecker.CanCreateComment() {
		return nil, domain.NewForbiddenError("you cannot comment")
	}

	comment := &domain.Comment{
		IssueID: req.IssueID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	s.auditEngine.LogCommentCreated(comment, userID)
	return comment, nil
}

func (s *IssueService) GetAvailableTransitions(issueID uuid.UUID) ([]domain.WorkflowTransition, error) {
	issue, err := s.repo.GetByID(issueID)
	if err != nil {
		return nil, err
	}

	workflow, err := s.workflowRepo.GetWorkflowByProjectID(issue.ProjectID)
	if err != nil {
		return nil, err
	}

	return s.workflowEngine.GetAvailableTransitions(workflow.ID, issue.Status)
}

type WorkflowService struct {
	repo        repository.WorkflowRepository
	permChecker *permissions.PermissionChecker
}

func NewWorkflowService(repo repository.WorkflowRepository, permChecker *permissions.PermissionChecker) *WorkflowService {
	return &WorkflowService{repo: repo, permChecker: permChecker}
}

func (s *WorkflowService) Create(req domain.CreateWorkflowRequest) (*domain.Workflow, error) {
	if !s.permChecker.CanManageWorkflow() {
		return nil, domain.NewForbiddenError("only admins can manage workflows")
	}

	workflow := &domain.Workflow{
		ProjectID: req.ProjectID,
		Name:      req.Name,
	}

	if err := s.repo.CreateWorkflow(workflow); err != nil {
		return nil, err
	}

	return workflow, nil
}

func (s *WorkflowService) GetByProjectID(projectID uuid.UUID) (*domain.Workflow, error) {
	return s.repo.GetWorkflowByProjectID(projectID)
}

func (s *WorkflowService) CreateState(req domain.CreateWorkflowStateRequest) (*domain.WorkflowState, error) {
	if !s.permChecker.CanManageWorkflow() {
		return nil, domain.NewForbiddenError("only admins can manage workflows")
	}

	state := &domain.WorkflowState{
		WorkflowID: req.WorkflowID,
		Name:       req.Name,
		IsInitial:  req.IsInitial,
		IsFinal:    req.IsFinal,
	}

	if err := s.repo.CreateState(state); err != nil {
		return nil, err
	}

	return state, nil
}

func (s *WorkflowService) CreateTransition(req domain.CreateWorkflowTransitionRequest) (*domain.WorkflowTransition, error) {
	if !s.permChecker.CanManageWorkflow() {
		return nil, domain.NewForbiddenError("only admins can manage workflows")
	}

	transition := &domain.WorkflowTransition{
		WorkflowID:  req.WorkflowID,
		FromStateID: req.FromStateID,
		ToStateID:   req.ToStateID,
		Name:        req.Name,
	}

	if err := s.repo.CreateTransition(transition); err != nil {
		return nil, err
	}

	return transition, nil
}

func (s *WorkflowService) DeleteState(stateID uuid.UUID) error {
	if !s.permChecker.CanManageWorkflow() {
		return domain.NewForbiddenError("only admins can manage workflows")
	}
	return s.repo.DeleteState(stateID)
}

func (s *WorkflowService) DeleteTransition(transitionID uuid.UUID) error {
	if !s.permChecker.CanManageWorkflow() {
		return domain.NewForbiddenError("only admins can manage workflows")
	}
	return s.repo.DeleteTransition(transitionID)
}
