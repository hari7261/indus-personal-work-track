import { createContext, useContext, useState, useCallback, ReactNode } from 'react';
import type { User, Project, Issue, IssueListItem, Comment, Workflow, ProjectMember, ProjectStats, IssueFilter } from '../types';

interface AppState {
  currentUser: User | null;
  isLoggedIn: boolean;
  users: User[];
  projects: Project[];
  currentProject: Project | null;
  issues: IssueListItem[];
  totalIssues: number;
  currentIssue: Issue | null;
  comments: Comment[];
  workflow: Workflow | null;
  members: ProjectMember[];
  stats: ProjectStats | null;
  loading: boolean;
  error: string | null;
}

interface AppContextType extends AppState {
  login: (username: string) => Promise<void>;
  logout: () => void;
  loadUsers: () => Promise<void>;
  loadProjects: () => Promise<void>;
  createProject: (name: string, description: string) => Promise<Project>;
  selectProject: (project: Project) => Promise<void>;
  loadIssues: (filter: IssueFilter) => Promise<void>;
  createIssue: (data: CreateIssueData) => Promise<Issue>;
  selectIssue: (issueId: string) => Promise<void>;
  updateIssue: (data: UpdateIssueData) => Promise<Issue>;
  deleteIssue: (issueId: string) => Promise<void>;
  transitionIssue: (issueId: string, transition: string) => Promise<Issue>;
  assignIssue: (issueId: string, assigneeId: string | null) => Promise<Issue>;
  loadComments: (issueId: string) => Promise<void>;
  createComment: (issueId: string, content: string) => Promise<Comment>;
  loadWorkflow: (projectId: string) => Promise<void>;
  loadMembers: (projectId: string) => Promise<void>;
  addMember: (projectId: string, userId: string, role: string) => Promise<void>;
  removeMember: (projectId: string, memberId: string) => Promise<void>;
  clearError: () => void;
}

interface CreateIssueData {
  project_id: string;
  title: string;
  description: string;
  priority: string;
  is_incident: boolean;
  severity: string | null;
}

interface UpdateIssueData {
  id: string;
  title: string;
  description: string;
  priority: string;
  assignee_id: string | null;
  is_incident: boolean;
  severity: string | null;
}

const AppContext = createContext<AppContextType | undefined>(undefined);

export function AppProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<AppState>({
    currentUser: null,
    isLoggedIn: false,
    users: [],
    projects: [],
    currentProject: null,
    issues: [],
    totalIssues: 0,
    currentIssue: null,
    comments: [],
    workflow: null,
    members: [],
    stats: null,
    loading: false,
    error: null,
  });

  const setLoading = (loading: boolean) => setState(s => ({ ...s, loading }));
  const clearError = useCallback(() => setState(s => ({ ...s, error: null })), []);
  const getErrorMessage = (err: unknown) => {
    if (err instanceof Error && err.message) return err.message;
    if (typeof err === 'string') return err;
    if (err && typeof err === 'object' && 'message' in err && typeof (err as { message: unknown }).message === 'string') {
      return (err as { message: string }).message;
    }
    return 'Login failed. Please check username and try again.';
  };

  const login = async (username: string) => {
    clearError();
    setLoading(true);
    try {
      const user = await window.go.main.App.Login(username);
      setState(s => ({ ...s, currentUser: user, isLoggedIn: true, error: null }));
    } catch (err) {
      setState(s => ({ ...s, currentUser: null, isLoggedIn: false, error: getErrorMessage(err) }));
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const logout = () => {
    setState(s => ({
      ...s,
      currentUser: null,
      isLoggedIn: false,
      currentProject: null,
      issues: [],
      currentIssue: null,
    }));
  };

  const loadUsers = async () => {
    const users = await window.go.main.App.ListUsers();
    setState(s => ({ ...s, users }));
  };

  const loadProjects = async () => {
    const projects = await window.go.main.App.ListProjects();
    setState(s => ({ ...s, projects }));
  };

  const createProject = async (name: string, description: string) => {
    const project = await window.go.main.App.CreateProject({ name, description });
    setState(s => ({ ...s, projects: [...s.projects, project] }));
    return project;
  };

  const selectProject = async (project: Project) => {
    setState(s => ({ ...s, currentProject: project, issues: [], currentIssue: null }));
    await loadIssues({ project_id: project.id, status: '', assignee_id: null, is_incident: null, search: '', page: 1, page_size: 20 });
    const stats = await window.go.main.App.GetIssueStats(project.id);
    setState(s => ({ ...s, stats }));
    await loadWorkflow(project.id);
    await loadMembers(project.id);
  };

  const loadIssues = async (filter: IssueFilter) => {
    const result = await window.go.main.App.ListIssues(filter);
    setState(s => ({ ...s, issues: result.data as IssueListItem[], totalIssues: result.total }));
  };

  const createIssue = async (data: CreateIssueData) => {
    const issue = await window.go.main.App.CreateIssue(data);
    if (state.currentProject) {
      await loadIssues({ project_id: state.currentProject.id, status: '', assignee_id: null, is_incident: null, search: '', page: 1, page_size: 20 });
      const stats = await window.go.main.App.GetIssueStats(state.currentProject.id);
      setState(s => ({ ...s, stats }));
    }
    return issue;
  };

  const selectIssue = async (issueId: string) => {
    const issue = await window.go.main.App.GetIssue(issueId);
    const comments = await window.go.main.App.GetIssueComments(issueId);
    setState(s => ({ ...s, currentIssue: issue, comments }));
  };

  const updateIssue = async (data: UpdateIssueData) => {
    const issue = await window.go.main.App.UpdateIssue(data);
    if (state.currentProject) {
      await loadIssues({ project_id: state.currentProject.id, status: '', assignee_id: null, is_incident: null, search: '', page: 1, page_size: 20 });
    }
    setState(s => ({ ...s, currentIssue: issue }));
    return issue;
  };

  const deleteIssue = async (issueId: string) => {
    await window.go.main.App.DeleteIssue(issueId);
    setState(s => ({ ...s, currentIssue: null }));
    if (state.currentProject) {
      await loadIssues({ project_id: state.currentProject.id, status: '', assignee_id: null, is_incident: null, search: '', page: 1, page_size: 20 });
    }
  };

  const transitionIssue = async (issueId: string, transition: string) => {
    const issue = await window.go.main.App.TransitionIssue({ issue_id: issueId, transition });
    setState(s => ({ ...s, currentIssue: issue }));
    return issue;
  };

  const assignIssue = async (issueId: string, assigneeId: string | null) => {
    const issue = await window.go.main.App.AssignIssue({ issue_id: issueId, assignee_id: assigneeId });
    setState(s => ({ ...s, currentIssue: issue }));
    return issue;
  };

  const loadComments = async (issueId: string) => {
    const comments = await window.go.main.App.GetIssueComments(issueId);
    setState(s => ({ ...s, comments }));
  };

  const createComment = async (issueId: string, content: string) => {
    const comment = await window.go.main.App.CreateComment({ issue_id: issueId, content });
    setState(s => ({ ...s, comments: [...s.comments, comment] }));
    return comment;
  };

  const loadWorkflow = async (projectId: string) => {
    try {
      const workflow = await window.go.main.App.GetWorkflow(projectId);
      setState(s => ({ ...s, workflow }));
    } catch {
      setState(s => ({ ...s, workflow: null }));
    }
  };

  const loadMembers = async (projectId: string) => {
    const members = await window.go.main.App.GetProjectMembers(projectId);
    setState(s => ({ ...s, members }));
  };

  const addMember = async (projectId: string, userId: string, role: string) => {
    await window.go.main.App.AddProjectMember({ project_id: projectId, user_id: userId, role });
    await loadMembers(projectId);
  };

  const removeMember = async (projectId: string, memberId: string) => {
    await window.go.main.App.RemoveProjectMember(projectId, memberId);
    await loadMembers(projectId);
  };

  return (
    <AppContext.Provider value={{
      ...state,
      login,
      logout,
      loadUsers,
      loadProjects,
      createProject,
      selectProject,
      loadIssues,
      createIssue,
      selectIssue,
      updateIssue,
      deleteIssue,
      transitionIssue,
      assignIssue,
      loadComments,
      createComment,
      loadWorkflow,
      loadMembers,
      addMember,
      removeMember,
      clearError,
    }}>
      {children}
    </AppContext.Provider>
  );
}

export function useApp() {
  const context = useContext(AppContext);
  if (!context) {
    throw new Error('useApp must be used within AppProvider');
  }
  return context;
}
