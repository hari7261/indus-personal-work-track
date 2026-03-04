package workflow

import (
	"indus-task-manager/internal/domain"
	"indus-task-manager/internal/repository"

	"github.com/google/uuid"
)

type Engine struct {
	workflowRepo repository.WorkflowRepository
}

func NewEngine(workflowRepo repository.WorkflowRepository) *Engine {
	return &Engine{workflowRepo: workflowRepo}
}

type TransitionResult struct {
	Success   bool   `json:"success"`
	NewStatus string `json:"new_status"`
	Error     string `json:"error,omitempty"`
}

func (e *Engine) CanTransition(workflowID uuid.UUID, currentState string, transitionName string) (bool, *domain.WorkflowTransition, error) {
	workflow, err := e.workflowRepo.GetWorkflowByID(workflowID)
	if err != nil {
		return false, nil, err
	}

	for _, t := range workflow.Transitions {
		fromState, err := e.workflowRepo.GetStateByID(t.FromStateID)
		if err != nil {
			continue
		}
		if fromState.Name == currentState && t.Name == transitionName {
			return true, &t, nil
		}
	}

	return false, nil, nil
}

func (e *Engine) GetAvailableTransitions(workflowID uuid.UUID, currentState string) ([]domain.WorkflowTransition, error) {
	workflow, err := e.workflowRepo.GetWorkflowByID(workflowID)
	if err != nil {
		return nil, err
	}

	var available []domain.WorkflowTransition
	for _, t := range workflow.Transitions {
		fromState, err := e.workflowRepo.GetStateByID(t.FromStateID)
		if err != nil {
			continue
		}
		if fromState.Name == currentState {
			toState, err := e.workflowRepo.GetStateByID(t.ToStateID)
			if err != nil {
				continue
			}
			t.Name = toState.Name
			available = append(available, t)
		}
	}

	return available, nil
}

func (e *Engine) GetInitialState(workflowID uuid.UUID) (string, error) {
	states, err := e.workflowRepo.GetStatesByWorkflowID(workflowID)
	if err != nil {
		return "", err
	}

	for _, s := range states {
		if s.IsInitial {
			return s.Name, nil
		}
	}

	return "", domain.ErrNotFound
}

func (e *Engine) ValidateWorkflow(workflow *domain.Workflow) error {
	stateMap := make(map[string]uuid.UUID)
	initialCount := 0
	finalCount := 0

	for _, s := range workflow.States {
		stateMap[s.Name] = s.ID
		if s.IsInitial {
			initialCount++
		}
		if s.IsFinal {
			finalCount++
		}
	}

	if initialCount != 1 {
		return domain.NewInvalidInputError("workflow must have exactly one initial state")
	}

	if finalCount < 1 {
		return domain.NewInvalidInputError("workflow must have at least one final state")
	}

	for _, t := range workflow.Transitions {
		if _, ok := stateMap[t.Name]; !ok {
			fromState, _ := e.workflowRepo.GetStateByID(t.FromStateID)
			toState, _ := e.workflowRepo.GetStateByID(t.ToStateID)
			if fromState == nil || toState == nil {
				return domain.NewInvalidInputError("invalid transition: referenced state does not exist")
			}
		}
	}

	return e.checkCircularDependencies(workflow)
}

func (e *Engine) checkCircularDependencies(workflow *domain.Workflow) error {
	stateGraph := make(map[string][]string)

	for _, s := range workflow.States {
		stateGraph[s.Name] = []string{}
	}

	for _, t := range workflow.Transitions {
		fromState, _ := e.workflowRepo.GetStateByID(t.FromStateID)
		toState, _ := e.workflowRepo.GetStateByID(t.ToStateID)
		if fromState != nil && toState != nil {
			stateGraph[fromState.Name] = append(stateGraph[fromState.Name], toState.Name)
		}
	}

	visited := make(map[string]bool)
	path := make(map[string]bool)

	var dfs func(state string) bool
	dfs = func(state string) bool {
		visited[state] = true
		path[state] = true

		for _, next := range stateGraph[state] {
			if !visited[next] {
				if dfs(next) {
					return true
				}
			} else if path[next] {
				return true
			}
		}

		path[state] = false
		return false
	}

	for state := range stateGraph {
		if !visited[state] {
			if dfs(state) {
				return domain.ErrCircularDependency
			}
		}
	}

	return nil
}
