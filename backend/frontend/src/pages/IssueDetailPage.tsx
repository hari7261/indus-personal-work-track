import { useState, useEffect } from 'react';
import { useApp } from '../state/AppContext';

interface Props {
  onBack: () => void;
}

export default function IssueDetailPage({ onBack }: Props) {
  const { currentIssue, comments, currentUser, users, transitionIssue, assignIssue, updateIssue, deleteIssue, createComment } = useApp();
  const [isEditing, setIsEditing] = useState(false);
  const [editTitle, setEditTitle] = useState('');
  const [editDescription, setEditDescription] = useState('');
  const [editPriority, setEditPriority] = useState('');
  const [editAssignee, setEditAssignee] = useState<string | null>(null);
  const [editIsIncident, setEditIsIncident] = useState(false);
  const [editSeverity, setEditSeverity] = useState<string | null>(null);
  const [newComment, setNewComment] = useState('');
  const [transitions, setTransitions] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);

  const isDeveloper = currentUser?.role === 'developer' || currentUser?.role === 'admin';

  useEffect(() => {
    if (currentIssue) {
      setEditTitle(currentIssue.title);
      setEditDescription(currentIssue.description);
      setEditPriority(currentIssue.priority);
      setEditAssignee(currentIssue.assignee_id);
      setEditIsIncident(currentIssue.is_incident);
      setEditSeverity(currentIssue.severity);
      loadTransitions();
    }
  }, [currentIssue]);

  const loadTransitions = async () => {
    if (currentIssue) {
      try {
        const t = await window.go.main.App.GetAvailableTransitions(currentIssue.id);
        setTransitions(t);
      } catch (e) {
        setTransitions([]);
      }
    }
  };

  const handleTransition = async (transition: string) => {
    if (currentIssue) {
      setLoading(true);
      try {
        await transitionIssue(currentIssue.id, transition);
        await loadTransitions();
      } finally {
        setLoading(false);
      }
    }
  };

  const handleAssign = async (assigneeId: string | null) => {
    if (currentIssue) {
      await assignIssue(currentIssue.id, assigneeId);
    }
  };

  const handleSave = async () => {
    if (currentIssue) {
      await updateIssue({
        id: currentIssue.id,
        title: editTitle,
        description: editDescription,
        priority: editPriority,
        assignee_id: editAssignee,
        is_incident: editIsIncident,
        severity: editIsIncident ? editSeverity : null,
      });
      setIsEditing(false);
    }
  };

  const handleDelete = async () => {
    if (currentIssue && confirm('Are you sure you want to delete this issue?')) {
      await deleteIssue(currentIssue.id);
      onBack();
    }
  };

  const handleAddComment = async () => {
    if (currentIssue && newComment.trim()) {
      await createComment(currentIssue.id, newComment.trim());
      setNewComment('');
    }
  };

  if (!currentIssue) {
    return <div className="content">Loading...</div>;
  }

  return (
    <div className="content">
      <div className="flex mb-4" style={{ justifyContent: 'space-between', alignItems: 'center' }}>
        <button className="btn btn-secondary btn-sm" onClick={onBack}>
          ← Back
        </button>
        {isDeveloper && !isEditing && (
          <div className="flex gap-2">
            <button className="btn btn-secondary btn-sm" onClick={() => setIsEditing(true)}>
              Edit
            </button>
            <button className="btn btn-danger btn-sm" onClick={handleDelete}>
              Delete
            </button>
          </div>
        )}
      </div>

      <div className="card issue-detail">
        <div className="card-body">
          {isEditing ? (
            <div>
              <div className="form-group">
                <label className="form-label">Title</label>
                <input
                  type="text"
                  className="form-input"
                  value={editTitle}
                  onChange={(e) => setEditTitle(e.target.value)}
                />
              </div>
              <div className="form-group">
                <label className="form-label">Description</label>
                <textarea
                  className="form-textarea"
                  value={editDescription}
                  onChange={(e) => setEditDescription(e.target.value)}
                />
              </div>
              <div className="form-group">
                <label className="form-label">Priority</label>
                <select
                  className="form-select"
                  value={editPriority}
                  onChange={(e) => setEditPriority(e.target.value)}
                >
                  <option value="low">Low</option>
                  <option value="medium">Medium</option>
                  <option value="high">High</option>
                  <option value="critical">Critical</option>
                </select>
              </div>
              <div className="form-group">
                <label className="form-label">Assignee</label>
                <select
                  className="form-select"
                  value={editAssignee || ''}
                  onChange={(e) => setEditAssignee(e.target.value || null)}
                >
                  <option value="">Unassigned</option>
                  {users.map((u) => (
                    <option key={u.id} value={u.id}>{u.username}</option>
                  ))}
                </select>
              </div>
              <div className="form-group">
                <label className="flex gap-2">
                  <input
                    type="checkbox"
                    checked={editIsIncident}
                    onChange={(e) => setEditIsIncident(e.target.checked)}
                  />
                  Is Incident
                </label>
              </div>
              {editIsIncident && (
                <div className="form-group">
                  <label className="form-label">Severity</label>
                  <select
                    className="form-select"
                    value={editSeverity || ''}
                    onChange={(e) => setEditSeverity(e.target.value)}
                  >
                    <option value="">Select severity</option>
                    <option value="minor">Minor</option>
                    <option value="major">Major</option>
                    <option value="critical">Critical</option>
                  </select>
                </div>
              )}
              <div className="flex gap-2">
                <button className="btn btn-primary" onClick={handleSave}>Save</button>
                <button className="btn btn-secondary" onClick={() => setIsEditing(false)}>Cancel</button>
              </div>
            </div>
          ) : (
            <>
              <div className="issue-detail-header">
                <h1 className="issue-detail-title">
                  {currentIssue.is_incident && <span className="badge badge-incident" style={{ marginRight: '8px' }}>INCIDENT</span>}
                  {currentIssue.title}
                </h1>
              </div>

              <div className="issue-detail-meta">
                <div className="issue-detail-field">
                  <span className="issue-detail-label">Status</span>
                  <span>{currentIssue.status}</span>
                </div>
                <div className="issue-detail-field">
                  <span className="issue-detail-label">Priority</span>
                  <span className={`badge badge-${currentIssue.priority}`}>{currentIssue.priority}</span>
                </div>
                <div className="issue-detail-field">
                  <span className="issue-detail-label">Assignee</span>
                  <select
                    className="form-select"
                    value={currentIssue.assignee_id || ''}
                    onChange={(e) => handleAssign(e.target.value || null)}
                    disabled={!isDeveloper}
                    style={{ padding: '4px' }}
                  >
                    <option value="">Unassigned</option>
                    {users.map((u) => (
                      <option key={u.id} value={u.id}>{u.username}</option>
                    ))}
                  </select>
                </div>
                <div className="issue-detail-field">
                  <span className="issue-detail-label">Created</span>
                  <span>{new Date(currentIssue.created_at).toLocaleString()}</span>
                </div>
                {currentIssue.is_incident && currentIssue.severity && (
                  <div className="issue-detail-field">
                    <span className="issue-detail-label">Severity</span>
                    <span className={`badge badge-severity-${currentIssue.severity}`}>{currentIssue.severity}</span>
                  </div>
                )}
              </div>

              {transitions.length > 0 && isDeveloper && (
                <div className="mb-4">
                  <label className="form-label">Actions</label>
                  <div className="flex gap-2">
                    {transitions.map((t) => (
                      <button
                        key={t.id}
                        className="btn btn-success btn-sm"
                        onClick={() => handleTransition(t.name)}
                        disabled={loading}
                      >
                        {t.name}
                      </button>
                    ))}
                  </div>
                </div>
              )}

              <div className="issue-description">
                <h3>Description</h3>
                <p>{currentIssue.description || 'No description'}</p>
              </div>
            </>
          )}
        </div>
      </div>

      <div className="comments-section">
        <h3 className="mb-4">Comments</h3>
        {comments.map((comment) => (
          <div key={comment.id} className="comment">
            <div className="comment-header">
              <span className="comment-author">User</span>
              <span className="comment-date">{new Date(comment.created_at).toLocaleString()}</span>
            </div>
            <p>{comment.content}</p>
          </div>
        ))}

        <div className="comment-form">
          <textarea
            className="form-textarea"
            placeholder="Add a comment..."
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
          />
          <button className="btn btn-primary btn-sm" onClick={handleAddComment} disabled={!newComment.trim()}>
            Add Comment
          </button>
        </div>
      </div>
    </div>
  );
}
