import React, { useState } from 'react';
import { useApp } from '../state/AppContext';

export default function LoginPage() {
  const { login, loading, error, clearError } = useApp();
  const [username, setUsername] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (username.trim()) {
      try {
        await login(username.trim());
      } catch (err) {
        console.error('Login failed:', err);
      }
    }
  };

  return (
    <div className="login-container">
      <div className="card login-card">
        <div className="card-header text-center">
          <h2>Indus Task Manager</h2>
        </div>
        <div className="card-body">
          {error && (
            <div className="error-message">
              {error}
              <button onClick={clearError} style={{ marginLeft: '8px' }}>x</button>
            </div>
          )}
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label className="form-label">Username</label>
              <input
                type="text"
                className="form-input"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="Enter username (admin, developer, reporter)"
                disabled={loading}
              />
            </div>
            <button type="submit" className="btn btn-primary" style={{ width: '100%' }} disabled={loading}>
              {loading ? 'Logging in...' : 'Login'}
            </button>
          </form>
          <div className="text-center mt-4 text-muted">
            <small>Default users: admin, developer, reporter</small>
          </div>
        </div>
      </div>
    </div>
  );
}
