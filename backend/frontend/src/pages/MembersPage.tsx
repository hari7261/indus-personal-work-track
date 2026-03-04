import React, { useState } from 'react';
import { useApp } from '../state/AppContext';

interface Props {
  onBack: () => void;
}

export default function MembersPage({ onBack }: Props) {
  const { currentProject, members, users, addMember, removeMember } = useApp();
  const [showAddForm, setShowAddForm] = useState(false);
  const [selectedUser, setSelectedUser] = useState('');
  const [selectedRole, setSelectedRole] = useState('developer');
  const [loading, setLoading] = useState(false);

  const memberUserIds = members.map((m) => m.user_id);
  const availableUsers = users.filter((u) => !memberUserIds.includes(u.id));

  const handleAddMember = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!currentProject || !selectedUser) return;
    setLoading(true);
    try {
      await addMember(currentProject.id, selectedUser, selectedRole);
      setShowAddForm(false);
      setSelectedUser('');
      setSelectedRole('developer');
    } finally {
      setLoading(false);
    }
  };

  const handleRemoveMember = async (memberId: string) => {
    if (!currentProject || !confirm('Remove this member?')) return;
    setLoading(true);
    try {
      await removeMember(currentProject.id, memberId);
    } finally {
      setLoading(false);
    }
  };

  const getRoleBadgeClass = (role: string) => {
    switch (role) {
      case 'admin': return 'badge-critical';
      case 'developer': return 'badge-medium';
      default: return 'badge-low';
    }
  };

  return (
    <div className="content">
      <div className="flex mb-4">
        <button className="btn btn-secondary btn-sm" onClick={onBack}>
          ← Back
        </button>
        <h2 style={{ marginLeft: '16px' }}>Project Members</h2>
      </div>

      <div className="card">
        <div className="card-header flex" style={{ justifyContent: 'space-between' }}>
          <span>Members ({members.length})</span>
          <button className="btn btn-primary btn-sm" onClick={() => setShowAddForm(true)}>
            Add Member
          </button>
        </div>
        <div className="card-body">
          {showAddForm && (
            <form onSubmit={handleAddMember} className="mb-4" style={{ padding: '12px', background: '#f8f9fa', borderRadius: '6px' }}>
              <div className="form-group">
                <select
                  className="form-select"
                  value={selectedUser}
                  onChange={(e) => setSelectedUser(e.target.value)}
                  required
                >
                  <option value="">Select user</option>
                  {availableUsers.map((u) => (
                    <option key={u.id} value={u.id}>{u.username} ({u.role})</option>
                  ))}
                </select>
              </div>
              <div className="form-group">
                <select
                  className="form-select"
                  value={selectedRole}
                  onChange={(e) => setSelectedRole(e.target.value)}
                >
                  <option value="reporter">Reporter</option>
                  <option value="developer">Developer</option>
                  <option value="admin">Admin</option>
                </select>
              </div>
              <div className="flex gap-2">
                <button type="submit" className="btn btn-primary btn-sm" disabled={loading}>
                  Add
                </button>
                <button type="button" className="btn btn-secondary btn-sm" onClick={() => setShowAddForm(false)}>
                  Cancel
                </button>
              </div>
            </form>
          )}

          {members.length === 0 ? (
            <div className="empty-state">
              <p>No members yet</p>
            </div>
          ) : (
            <table className="table">
              <thead>
                <tr>
                  <th>Username</th>
                  <th>Role</th>
                  <th>Added</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {members.map((member) => {
                  const user = users.find((u) => u.id === member.user_id);
                  return (
                    <tr key={member.id}>
                      <td>{user?.username || 'Unknown'}</td>
                      <td>
                        <span className={`badge ${getRoleBadgeClass(member.role)}`}>
                          {member.role}
                        </span>
                      </td>
                      <td>{new Date(member.created_at).toLocaleDateString()}</td>
                      <td>
                        <button
                          className="btn btn-danger btn-sm"
                          onClick={() => handleRemoveMember(member.id)}
                          disabled={loading}
                        >
                          Remove
                        </button>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  );
}
