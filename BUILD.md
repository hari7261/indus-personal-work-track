# Build Instructions

## Prerequisites

- Go 1.21 or later
- Node.js 18 or later
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

## Backend Setup

```bash
cd backend
go mod tidy
```

## Frontend Setup

```bash
cd frontend
npm install
```

## Development

### Option 1: Run with Wails (recommended)

```bash
cd backend
wails dev
```

This will:
1. Build the frontend
2. Start the Go backend
3. Launch the desktop application

### Option 2: Run frontend separately for development

```bash
# Terminal 1 - Backend
cd backend
wails dev

# Terminal 2 - Frontend (for hot reload)
cd frontend
npm run dev
```

## Production Build

```bash
cd backend
wails build -clean -nsis
```

The executable and installer will be created in:
- Windows app: `backend/build/bin/IndusTaskManager.exe`
- Windows installer: `backend/build/bin/indus-task-manager-amd64-installer.exe`
- macOS: `backend/build/bin/Indus Task Manager.app`
- Linux: `backend/build/bin/indus-task-manager`

## Database

The application uses SQLite and creates a database file at:
- Windows: `%APPDATA%/indus-task/indus-task.db`
- macOS: `~/Library/Application Support/indus-task/indus-task.db`
- Linux: `~/.config/indus-task/indus-task.db`

## Default Users

The application seeds three default users on first run:
- `admin` - Admin role (full access)
- `developer` - Developer role
- `reporter` - Reporter role

## Architecture Notes

- The backend is designed to be reusable as a CLI tool or HTTP service
- Wails is only an adapter layer - all business logic is in the `internal` packages
- The domain layer has no dependencies on Wails or any UI framework
- Permissions are centralized in `internal/permissions`
- Workflow transitions are validated by the `internal/workflow` engine
- All important changes produce audit logs via `internal/audit`
