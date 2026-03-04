import React, { useState } from 'react';
import { useApp } from '../state/AppContext';

interface Props {
  onSelectProject: () => void;
}

export default function ProjectsPage({ onSelectProject }: Props) {
  const { projects, createProject, selectProject, currentUser, loading } = useApp();
  const [showCreate, setShowCreate] = useState(false);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (name.trim()) {
      await createProject(name.trim(), description.trim());
      setName('');
      setDescription('');
      setShowCreate(false);
    }
  };

  const handleSelect = async (project: any) => {
    await selectProject(project);
    onSelectProject();
  };

  const isAdmin = currentUser?.role === 'admin';

  return (
    <div className="content">
      <div className="flex mb-4" style={{ justifyContent: 'space-between', alignItems: 'center' }}>
        <h2>Projects</h2>
        {isAdmin && (
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            New Project
          </button>
        )}
      </div>

      {showCreate && (
        <div className="card mb-4">
          <div className="card-header">Create Project</div>
          <div className="card-body">
            <form onSubmit={handleCreate}>
              <div className="form-group">
                <label className="form-label">Name</label>
                <input
                  type="text"
                  className="form-input"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  required
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
              <div className="flex gap-2">
                <button type="submit" className="btn btn-primary" disabled={loading}>
                  Create
                </button>
                <button type="button" className="btn btn-secondary" onClick={() => setShowCreate(false)}>
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {projects.length === 0 ? (
        <div className="empty-state">
          <p>No projects yet</p>
        </div>
      ) : (
        <div className="issue-list">
          {projects.map((project) => (
            <div
              key={project.id}
              className="issue-item"
              onClick={() => handleSelect(project)}
            >
              <div className="issue-item-info">
                <div className="issue-item-title">{project.name}</div>
                <div className="issue-item-meta">
                  <span>{project.description || 'No description'}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
