import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import './index.css'

declare global {
  interface Window {
    go: {
      main: {
        App: {
          Login: (username: string) => Promise<any>;
          GetCurrentUser: () => Promise<any>;
          ListUsers: () => Promise<any[]>;
          ListProjects: () => Promise<any[]>;
          CreateProject: (data: { name: string; description: string }) => Promise<any>;
          GetProject: (id: string) => Promise<any>;
          UpdateProject: (data: any) => Promise<any>;
          DeleteProject: (id: string) => Promise<void>;
          CreateIssue: (data: any) => Promise<any>;
          GetIssue: (id: string) => Promise<any>;
          ListIssues: (filter: any) => Promise<{ data: any[]; total: number; page: number }>;
          UpdateIssue: (data: any) => Promise<any>;
          DeleteIssue: (id: string) => Promise<void>;
          TransitionIssue: (data: { issue_id: string; transition: string }) => Promise<any>;
          AssignIssue: (data: { issue_id: string; assignee_id: string | null }) => Promise<any>;
          GetIssueStats: (projectId: string) => Promise<any>;
          GetIssueComments: (issueId: string) => Promise<any[]>;
          CreateComment: (data: { issue_id: string; content: string }) => Promise<any>;
          GetAvailableTransitions: (issueId: string) => Promise<any[]>;
          CreateWorkflow: (data: any) => Promise<any>;
          GetWorkflow: (projectId: string) => Promise<any>;
          CreateWorkflowState: (data: any) => Promise<any>;
          CreateWorkflowTransition: (data: any) => Promise<any>;
          DeleteWorkflowState: (stateId: string) => Promise<void>;
          DeleteWorkflowTransition: (transitionId: string) => Promise<void>;
          GetProjectMembers: (projectId: string) => Promise<any[]>;
          AddProjectMember: (data: any) => Promise<void>;
          RemoveProjectMember: (projectId: string, memberId: string) => Promise<void>;
        };
      };
    };
  }
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
