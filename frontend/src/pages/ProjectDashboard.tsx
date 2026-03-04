import { useState } from 'react';
import { useApp } from '../state/AppContext';

interface Props {
  onBack: () => void;
  onViewIssue: (id: string) => void;
  onCreateIssue: () => void;
  onManageWorkflow: () => void;
  onManageMembers: () => void;
}

export default function ProjectDashboard({ onBack, onViewIssue, onCreateIssue, onManageWorkflow, onManageMembers }: Props) {
  const { currentProject, issues, stats, currentUser, selectIssue, loadIssues } = useApp();
  const [filterStatus, setFilterStatus] = useState('');
  const [filterIncident, setFilterIncident] = useState<boolean | null>(null);
  const [search, setSearch] = useState('');

  const isAdmin = currentUser?.role === 'admin';

  const handleFilter = async () => {
    if (currentProject) {
      await loadIssues({
        project_id: currentProject.id,
        status: filterStatus,
        assignee_id: null,
        is_incident: filterIncident,
        search: search,
        page: 1,
        page_size: 20,
      });
    }
  };

  const handleViewIssue = async (issueId: string) => {
    await selectIssue(issueId);
    onViewIssue(issueId);
  };

  const getPriorityClass = (priority: string) => `badge badge-${priority}`;

  return (
    <div className="content">
      <div className="flex mb-4" style={{ justifyContent: 'space-between', alignItems: 'center' }}>
        <div className="flex gap-2" style={{ alignItems: 'center' }}>
          <button className="btn btn-secondary btn-sm" onClick={onBack}>
            ← Back
          </button>
          <h2>{currentProject?.name}</h2>
        </div>
        <div className="flex gap-2">
          {isAdmin && (
            <>
              <button className="btn btn-secondary btn-sm" onClick={onManageWorkflow}>
                Workflow
              </button>
              <button className="btn btn-secondary btn-sm" onClick={onManageMembers}>
                Members
              </button>
            </>
          )}
          <button className="btn btn-primary btn-sm" onClick={onCreateIssue}>
            New Issue
          </button>
        </div>
      </div>

      {stats && (
        <div className="stats-grid">
          <div className="stat-card">
            <div className="stat-value">{stats.total}</div>
            <div className="stat-label">Total Issues</div>
          </div>
          <div className="stat-card">
            <div className="stat-value">{stats.open_issues}</div>
            <div className="stat-label">Open Issues</div>
          </div>
          <div className="stat-card">
            <div className="stat-value" style={{ color: '#dc3545' }}>{stats.incident_count}</div>
            <div className="stat-label">Incidents</div>
          </div>
          <div className="stat-card">
            <div className="stat-value" style={{ color: '#dc3545' }}>{stats.critical_count}</div>
            <div className="stat-label">Critical</div>
          </div>
        </div>
      )}

      <div className="card mb-4">
        <div className="card-body">
          <div className="flex gap-2">
            <input
              type="text"
              className="form-input"
              placeholder="Search issues..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              style={{ maxWidth: '300px' }}
            />
            <select
              className="form-select"
              value={filterStatus}
              onChange={(e) => setFilterStatus(e.target.value)}
              style={{ maxWidth: '150px' }}
            >
              <option value="">All Status</option>
              <option value="Open">Open</option>
              <option value="In Progress">In Progress</option>
              <option value="Resolved">Resolved</option>
              <option value="Closed">Closed</option>
            </select>
            <select
              className="form-select"
              value={filterIncident === null ? '' : filterIncident.toString()}
              onChange={(e) => setFilterIncident(e.target.value === '' ? null : e.target.value === 'true')}
              style={{ maxWidth: '150px' }}
            >
              <option value="">All Types</option>
              <option value="false">Regular Issue</option>
              <option value="true">Incident</option>
            </select>
            <button className="btn btn-secondary" onClick={handleFilter}>
              Filter
            </button>
          </div>
        </div>
      </div>

      {issues.length === 0 ? (
        <div className="empty-state">
          <p>No issues found</p>
        </div>
      ) : (
        <div className="issue-list">
          {issues.map((issue) => (
            <div
              key={issue.id}
              className="issue-item"
              onClick={() => handleViewIssue(issue.id)}
            >
              <div className="issue-item-info">
                <div className="issue-item-title">
                  {issue.is_incident && <span className="badge badge-incident" style={{ marginRight: '8px' }}>INC</span>}
                  {issue.title}
                </div>
                <div className="issue-item-meta">
                  <span className={getPriorityClass(issue.priority)}>{issue.priority}</span>
                  <span>{issue.status}</span>
                  <span>{issue.assignee || 'Unassigned'}</span>
                  <span>{new Date(issue.created_at).toLocaleDateString()}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
