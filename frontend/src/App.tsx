import { useState, useEffect } from 'react';
import { AppProvider, useApp } from './state/AppContext';
import LoginPage from './pages/LoginPage';
import ProjectsPage from './pages/ProjectsPage';
import ProjectDashboard from './pages/ProjectDashboard';
import IssueDetailPage from './pages/IssueDetailPage';
import CreateIssuePage from './pages/CreateIssuePage';
import WorkflowPage from './pages/WorkflowPage';
import MembersPage from './pages/MembersPage';

type Page = 'projects' | 'dashboard' | 'issue' | 'create-issue' | 'workflow' | 'members';

function AppContent() {
  const { isLoggedIn, currentUser, logout, loadUsers, loadProjects, currentProject, currentIssue } = useApp();
  const [currentPage, setCurrentPage] = useState<Page>('projects');

  useEffect(() => {
    if (isLoggedIn) {
      loadUsers();
      loadProjects();
    }
  }, [isLoggedIn]);

  if (!isLoggedIn) {
    return <LoginPage />;
  }

  const renderPage = () => {
    switch (currentPage) {
      case 'projects':
        return <ProjectsPage onSelectProject={() => setCurrentPage('dashboard')} />;
      case 'dashboard':
        return (
          <ProjectDashboard
            onBack={() => setCurrentPage('projects')}
            onViewIssue={(_id: string) => { setCurrentPage('issue'); }}
            onCreateIssue={() => setCurrentPage('create-issue')}
            onManageWorkflow={() => setCurrentPage('workflow')}
            onManageMembers={() => setCurrentPage('members')}
          />
        );
      case 'issue':
        return currentIssue ? (
          <IssueDetailPage onBack={() => setCurrentPage('dashboard')} />
        ) : (
          <ProjectDashboard
            onBack={() => setCurrentPage('projects')}
            onViewIssue={(_id: string) => {}}
            onCreateIssue={() => setCurrentPage('create-issue')}
            onManageWorkflow={() => setCurrentPage('workflow')}
            onManageMembers={() => setCurrentPage('members')}
          />
        );
      case 'create-issue':
        return <CreateIssuePage onBack={() => setCurrentPage('dashboard')} onCreated={() => setCurrentPage('dashboard')} />;
      case 'workflow':
        return <WorkflowPage onBack={() => setCurrentPage('dashboard')} />;
      case 'members':
        return <MembersPage onBack={() => setCurrentPage('dashboard')} />;
      default:
        return <ProjectsPage onSelectProject={() => setCurrentPage('dashboard')} />;
    }
  };

  return (
    <div className="app">
      <header className="header">
        <div className="header-title">Indus Task Manager</div>
        <div className="header-user">
          <span>{currentUser?.username} ({currentUser?.role})</span>
          {currentProject && (
            <span className="text-muted">{currentProject.name}</span>
          )}
          <button className="btn btn-secondary btn-sm" onClick={logout}>
            Logout
          </button>
        </div>
      </header>
      <main className="main">
        {renderPage()}
      </main>
    </div>
  );
}

export default function App() {
  return (
    <AppProvider>
      <AppContent />
    </AppProvider>
  );
}
