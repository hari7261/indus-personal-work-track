import React, { useState } from 'react';
import { useApp } from '../state/AppContext';

interface Props {
  onBack: () => void;
  onCreated: () => void;
}

export default function CreateIssuePage({ onBack, onCreated }: Props) {
  const { currentProject, createIssue } = useApp();
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [priority, setPriority] = useState('medium');
  const [isIncident, setIsIncident] = useState(false);
  const [severity, setSeverity] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim() || !currentProject) return;

    setLoading(true);
    try {
      await createIssue({
        project_id: currentProject.id,
        title: title.trim(),
        description: description.trim(),
        priority,
        is_incident: isIncident,
        severity: isIncident ? severity : null,
      });
      onCreated();
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
      </div>

      <div className="card" style={{ maxWidth: '600px' }}>
        <div className="card-header">
          <h2>Create New Issue</h2>
        </div>
        <div className="card-body">
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label className="form-label">Title *</label>
              <input
                type="text"
                className="form-input"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                required
                maxLength={200}
              />
            </div>

            <div className="form-group">
              <label className="form-label">Description</label>
              <textarea
                className="form-textarea"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />
            </div>

            <div className="form-group">
              <label className="form-label">Priority</label>
              <select
                className="form-select"
                value={priority}
                onChange={(e) => setPriority(e.target.value)}
              >
                <option value="low">Low</option>
                <option value="medium">Medium</option>
                <option value="high">High</option>
                <option value="critical">Critical</option>
              </select>
            </div>

            <div className="form-group">
              <label className="flex gap-2" style={{ alignItems: 'center' }}>
                <input
                  type="checkbox"
                  checked={isIncident}
                  onChange={(e) => {
                    setIsIncident(e.target.checked);
                    if (!e.target.checked) setSeverity(null);
                  }}
                />
                This is an incident
              </label>
            </div>

            {isIncident && (
              <div className="form-group">
                <label className="form-label">Severity</label>
                <select
                  className="form-select"
                  value={severity || ''}
                  onChange={(e) => setSeverity(e.target.value || null)}
                >
                  <option value="">Select severity</option>
                  <option value="minor">Minor</option>
                  <option value="major">Major</option>
                  <option value="critical">Critical</option>
                </select>
              </div>
            )}

            <div className="flex gap-2">
              <button type="submit" className="btn btn-primary" disabled={loading || !title.trim()}>
                {loading ? 'Creating...' : 'Create Issue'}
              </button>
              <button type="button" className="btn btn-secondary" onClick={onBack}>
                Cancel
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
