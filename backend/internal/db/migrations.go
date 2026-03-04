package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func migrate(db *sqlx.DB) error {
	schema := `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL CHECK(role IN ('reporter', 'developer', 'admin')),
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS project_members (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role TEXT NOT NULL CHECK(role IN ('reporter', 'developer', 'admin')),
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(project_id, user_id)
);

CREATE TABLE IF NOT EXISTS workflows (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    UNIQUE(project_id)
);

CREATE TABLE IF NOT EXISTS workflow_states (
    id TEXT PRIMARY KEY,
    workflow_id TEXT NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    is_initial INTEGER NOT NULL DEFAULT 0,
    is_final INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS workflow_transitions (
    id TEXT PRIMARY KEY,
    workflow_id TEXT NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    from_state_id TEXT NOT NULL REFERENCES workflow_states(id) ON DELETE CASCADE,
    to_state_id TEXT NOT NULL REFERENCES workflow_states(id) ON DELETE CASCADE,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS issues (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL,
    assignee_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    priority TEXT NOT NULL CHECK(priority IN ('low', 'medium', 'high', 'critical')),
    is_incident INTEGER NOT NULL DEFAULT 0,
    severity TEXT CHECK(severity IN ('', 'minor', 'major', 'critical')),
    created_by TEXT NOT NULL REFERENCES users(id),
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_issues_project_id ON issues(project_id);
CREATE INDEX IF NOT EXISTS idx_issues_status ON issues(status);
CREATE INDEX IF NOT EXISTS idx_issues_assignee ON issues(assignee_id);
CREATE INDEX IF NOT EXISTS idx_issues_created_at ON issues(created_at);
CREATE INDEX IF NOT EXISTS idx_issues_is_incident ON issues(is_incident);

CREATE TABLE IF NOT EXISTS comments (
    id TEXT PRIMARY KEY,
    issue_id TEXT NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_comments_issue_id ON comments(issue_id);

CREATE TABLE IF NOT EXISTS audit_logs (
    id TEXT PRIMARY KEY,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    action TEXT NOT NULL,
    user_id TEXT NOT NULL REFERENCES users(id),
    old_value TEXT,
    new_value TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_audit_entity ON audit_logs(entity_type, entity_id);
`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	if err := seedDefaultData(db); err != nil {
		return fmt.Errorf("seeding failed: %w", err)
	}

	return nil
}

func seedDefaultData(db *sqlx.DB) error {
	defaultUsers := []struct {
		ID       string
		Username string
		Role     string
	}{
		{ID: "00000000-0000-0000-0000-000000000001", Username: "admin", Role: "admin"},
		{ID: "00000000-0000-0000-0000-000000000002", Username: "developer", Role: "developer"},
		{ID: "00000000-0000-0000-0000-000000000003", Username: "reporter", Role: "reporter"},
	}

	for _, user := range defaultUsers {
		if _, err := db.Exec(
			"INSERT OR IGNORE INTO users (id, username, role, created_at) VALUES (?, ?, ?, datetime('now'))",
			user.ID,
			user.Username,
			user.Role,
		); err != nil {
			return err
		}

		if _, err := db.Exec(
			"UPDATE users SET role = ? WHERE username = ?",
			user.Role,
			user.Username,
		); err != nil {
			return err
		}
	}

	return nil
}
