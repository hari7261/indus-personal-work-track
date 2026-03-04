# Indus Task Manager - Specification Document

## 1. Project Overview

**Project Name:** Indus Task Manager  
**Type:** Desktop Application (Internal Developer Tool)  
**Core Summary:** A lightweight, offline-first task and incident management system for development teams  
**Target Users:** Internal development teams (Reporters, Developers, Admins)

## 2. Technology Stack

### Backend
- **Language:** Go
- **Desktop Shell:** Wails
- **Database:** SQLite (default), PostgreSQL-ready architecture
- **Architecture:** Layered (Domain → Services → Repository → Database)

### Frontend
- **Framework:** React with Vite
- **Language:** TypeScript
- **Styling:** CSS Modules / Tailwind CSS
- **State Management:** React Context + Hooks

### Constraints
- No Electron
- No Node backend
- No Next.js SSR
- No cloud assumptions

## 3. Roles & Permissions

### Roles (Hard-coded)
1. **REPORTER**
   - Create issue
   - View issue
   - Comment on issue

2. **DEVELOPER**
   - All Reporter permissions
   - Assign issues
   - Transition issue status
   - Update issue details

3. **ADMIN**
   - Full access
   - Manage workflows
   - Manage project members

### Permission Model
- Centralized permission checks (never scattered)
- Role-based access control (RBAC)
- Permission interface in `/internal/permissions`

## 4. Domain Features

### Projects
- `id` (UUID)
- `name` (string, required)
- `description` (text)
- `created_at` (timestamp)
- `updated_at` (timestamp)

### Issues
- `id` (UUID)
- `project_id` (UUID, FK)
- `title` (string, required, max 200)
- `description` (text)
- `status` (string, FK to workflow state)
- `assignee_id` (UUID, FK to user, nullable)
- `priority` (enum: low, medium, high, critical)
- `is_incident` (boolean)
- `severity` (enum: null, minor, major, critical - nullable, only for incidents)
- `created_by` (UUID, FK to user)
- `created_at` (timestamp)
- `updated_at` (timestamp)

### Comments
- `id` (UUID)
- `issue_id` (UUID, FK)
- `user_id` (UUID, FK)
- `content` (text)
- `created_at` (timestamp)

### Users (Internal)
- `id` (UUID)
- `username` (string, unique)
- `role` (enum: reporter, developer, admin)
- `created_at` (timestamp)

### Workflow (Per Project)
- `id` (UUID)
- `project_id` (UUID, FK)
- `name` (string)
- States: `id`, `workflow_id`, `name`, `is_initial`, `is_final`
- Transitions: `id`, `workflow_id`, `from`, `to_state_state_id_id`, `name`

### Audit Log
- `id` (UUID)
- `entity_type` (string: issue, project, comment, workflow)
- `entity_id` (UUID)
- `action` (string)
- `user_id` (UUID)
- `old_value` (JSON)
- `new_value` (JSON)
- `created_at` (timestamp)

## 5. Database Schema

### Indexes
- `idx_issues_project_id` on issues(project_id)
- `idx_issues_status` on issues(status)
- `idx_issues_assignee` on issues(assignee_id)
- `idx_issues_created_at` on issues(created_at)
- `idx_issues_is_incident` on issues(is_incident)
- `idx_comments_issue_id` on comments(issue_id)

### Foreign Keys
- All foreign keys enabled
- CASCADE deletes where appropriate

## 6. Frontend Screens

1. **Login/Init** - First-time setup or role selection
2. **Project List** - View all accessible projects
3. **Project Dashboard** - Project overview with stats
4. **Issue List** - Filterable list (status, assignee, incident flag)
5. **Issue Details** - Full issue view with comments
6. **Create/Edit Issue** - Form for issue management
7. **Workflow Editor** - ADMIN only - manage states/transitions
8. **Project Members** - ADMIN only - manage users and roles

## 7. Acceptance Criteria

### Backend
- [ ] All CRUD operations work correctly
- [ ] Permission checks block unauthorized actions
- [ ] Workflow transitions validate state machine rules
- [ ] Audit log captures all important changes
- [ ] Database transactions ensure consistency
- [ ] Unit tests pass for permissions and workflow

### Frontend
- [ ] Responsive layout works on desktop
- [ ] Issue list is paginated/virtualized
- [ ] Optimistic UI updates for comments
- [ ] Clear error messages from backend
- [ ] Keyboard navigation supported
- [ ] Incident quick filter works

### Desktop
- [ ] Application starts quickly
- [ ] Memory usage stays low
- [ ] No HTTP server exposed
- [ ] Wails bindings only call services

## 8. Build Instructions

### Prerequisites
- Go 1.21+
- Node.js 18+
- Wails CLI

### Build Commands
```bash
# Backend
cd backend
go mod tidy
wails dev  # Development
wails build  # Production

# Frontend
cd frontend
npm install
npm run dev
```

## 9. Project Structure

```
/backend
  /cmd
    /app
  /internal
    /domain
    /services
    /repository
    /permissions
    /workflow
    /audit
    /app
    /db
  /migrations
  /wails

/frontend
  /src
    /pages
    /components
    /api
    /state
    /validation
  /public
```
