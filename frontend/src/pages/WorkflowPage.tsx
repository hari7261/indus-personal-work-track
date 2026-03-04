import React, { useState, useEffect } from 'react';
import { useApp } from '../state/AppContext';

interface Props {
  onBack: () => void;
}

export default function WorkflowPage({ onBack }: Props) {
  const { currentProject, workflow, loadWorkflow } = useApp();
  const [newStateName, setNewStateName] = useState('');
  const [newStateInitial, setNewStateInitial] = useState(false);
  const [newStateFinal, setNewStateFinal] = useState(false);
  const [showStateForm, setShowStateForm] = useState(false);
  const [transitionFrom, setTransitionFrom] = useState('');
  const [transitionTo, setTransitionTo] = useState('');
  const [transitionName, setTransitionName] = useState('');
  const [showTransitionForm, setShowTransitionForm] = useState(false);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (currentProject) {
      loadWorkflow(currentProject.id);
    }
  }, [currentProject]);

  const handleCreateWorkflow = async () => {
    if (!currentProject) return;
    setLoading(true);
    try {
      await window.go.main.App.CreateWorkflow({
        project_id: currentProject.id,
        name: 'Default Workflow',
      });
      await loadWorkflow(currentProject.id);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateState = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!workflow) return;
    setLoading(true);
    try {
      await window.go.main.App.CreateWorkflowState({
        workflow_id: workflow.id,
        name: newStateName,
        is_initial: newStateInitial,
        is_final: newStateFinal,
      });
      await loadWorkflow(currentProject!.id);
      setShowStateForm(false);
      setNewStateName('');
      setNewStateInitial(false);
      setNewStateFinal(false);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateTransition = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!workflow) return;
    setLoading(true);
    try {
      await window.go.main.App.CreateWorkflowTransition({
        workflow_id: workflow.id,
        from_state_id: transitionFrom,
        to_state_id: transitionTo,
        name: transitionName,
      });
      await loadWorkflow(currentProject!.id);
      setShowTransitionForm(false);
      setTransitionFrom('');
      setTransitionTo('');
      setTransitionName('');
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteState = async (stateId: string) => {
    if (!confirm('Delete this state?')) return;
    setLoading(true);
    try {
      await window.go.main.App.DeleteWorkflowState(stateId);
      await loadWorkflow(currentProject!.id);
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteTransition = async (transitionId: string) => {
    if (!confirm('Delete this transition?')) return;
    setLoading(true);
    try {
      await window.go.main.App.DeleteWorkflowTransition(transitionId);
      await loadWorkflow(currentProject!.id);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="content">
      <div className="flex mb-4">
        <button className="btn btn-secondary btn-sm" onClick={onBack}>
          ← Back
        </button>
        <h2 style={{ marginLeft: '16px' }}>Workflow Editor</h2>
      </div>

      {!workflow ? (
        <div className="card">
          <div className="card-body text-center">
            <p className="mb-4">No workflow configured for this project</p>
            <button className="btn btn-primary" onClick={handleCreateWorkflow} disabled={loading}>
              Create Workflow
            </button>
          </div>
        </div>
      ) : (
        <div className="flex gap-4">
          <div className="card" style={{ flex: 1 }}>
            <div className="card-header flex" style={{ justifyContent: 'space-between' }}>
              <span>States</span>
              <button className="btn btn-primary btn-sm" onClick={() => setShowStateForm(true)}>
                Add State
              </button>
            </div>
            <div className="card-body">
              {showStateForm && (
                <form onSubmit={handleCreateState} className="mb-4" style={{ padding: '12px', background: '#f8f9fa', borderRadius: '6px' }}>
                  <div className="form-group">
                    <input
                      type="text"
                      className="form-input"
                      placeholder="State name"
                      value={newStateName}
                      onChange={(e) => setNewStateName(e.target.value)}
                      required
                    />
                  </div>
                  <div className="flex gap-2 mb-4">
                    <label>
                      <input
                        type="checkbox"
                        checked={newStateInitial}
                        onChange={(e) => setNewStateInitial(e.target.checked)}
                      /> Initial
                    </label>
                    <label>
                      <input
                        type="checkbox"
                        checked={newStateFinal}
                        onChange={(e) => setNewStateFinal(e.target.checked)}
                      /> Final
                    </label>
                  </div>
                  <div className="flex gap-2">
                    <button type="submit" className="btn btn-primary btn-sm" disabled={loading}>Add</button>
                    <button type="button" className="btn btn-secondary btn-sm" onClick={() => setShowStateForm(false)}>Cancel</button>
                  </div>
                </form>
              )}
              {workflow.states.map((state) => (
                <div key={state.id} className="flex" style={{ justifyContent: 'space-between', padding: '8px', borderBottom: '1px solid #dee2e6' }}>
                  <span>
                    {state.name}
                    {state.is_initial && <span className="badge badge-medium" style={{ marginLeft: '8px' }}>Initial</span>}
                    {state.is_final && <span className="badge badge-critical" style={{ marginLeft: '8px' }}>Final</span>}
                  </span>
                  <button className="btn btn-danger btn-sm" onClick={() => handleDeleteState(state.id)}>×</button>
                </div>
              ))}
            </div>
          </div>

          <div className="card" style={{ flex: 1 }}>
            <div className="card-header flex" style={{ justifyContent: 'space-between' }}>
              <span>Transitions</span>
              <button className="btn btn-primary btn-sm" onClick={() => setShowTransitionForm(true)}>
                Add Transition
              </button>
            </div>
            <div className="card-body">
              {showTransitionForm && (
                <form onSubmit={handleCreateTransition} className="mb-4" style={{ padding: '12px', background: '#f8f9fa', borderRadius: '6px' }}>
                  <div className="form-group">
                    <select className="form-select" value={transitionFrom} onChange={(e) => setTransitionFrom(e.target.value)} required>
                      <option value="">From state</option>
                      {workflow.states.map((s) => (
                        <option key={s.id} value={s.id}>{s.name}</option>
                      ))}
                    </select>
                  </div>
                  <div className="form-group">
                    <select className="form-select" value={transitionTo} onChange={(e) => setTransitionTo(e.target.value)} required>
                      <option value="">To state</option>
                      {workflow.states.map((s) => (
                        <option key={s.id} value={s.id}>{s.name}</option>
                      ))}
                    </select>
                  </div>
                  <div className="form-group">
                    <input
                      type="text"
                      className="form-input"
                      placeholder="Transition name"
                      value={transitionName}
                      onChange={(e) => setTransitionName(e.target.value)}
                      required
                    />
                  </div>
                  <div className="flex gap-2">
                    <button type="submit" className="btn btn-primary btn-sm" disabled={loading}>Add</button>
                    <button type="button" className="btn btn-secondary btn-sm" onClick={() => setShowTransitionForm(false)}>Cancel</button>
                  </div>
                </form>
              )}
              {workflow.transitions.map((t) => (
                <div key={t.id} className="flex" style={{ justifyContent: 'space-between', padding: '8px', borderBottom: '1px solid #dee2e6' }}>
                  <span>{t.name}</span>
                  <button className="btn btn-danger btn-sm" onClick={() => handleDeleteTransition(t.id)}>×</button>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
