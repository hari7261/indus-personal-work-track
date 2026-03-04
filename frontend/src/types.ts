export interface User {
  id: string;
  username: string;
  role: 'reporter' | 'developer' | 'admin';
  created_at: string;
}

export interface Project {
  id: string;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface Issue {
  id: string;
  project_id: string;
  title: string;
  description: string;
  status: string;
  assignee_id: string | null;
  priority: 'low' | 'medium' | 'high' | 'critical';
  is_incident: boolean;
  severity: 'minor' | 'major' | 'critical' | null;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface IssueListItem {
  id: string;
  project_id: string;
  title: string;
  status: string;
  priority: 'low' | 'medium' | 'high' | 'critical';
  is_incident: boolean;
  severity: 'minor' | 'major' | 'critical' | null;
  assignee_id: string | null;
  assignee: string | null;
  created_by: string;
  created_by_name: string;
  created_at: string;
  updated_at: string;
}

export interface Comment {
  id: string;
  issue_id: string;
  user_id: string;
  content: string;
  created_at: string;
}

export interface Workflow {
  id: string;
  project_id: string;
  name: string;
  states: WorkflowState[];
  transitions: WorkflowTransition[];
}

export interface WorkflowState {
  id: string;
  workflow_id: string;
  name: string;
  is_initial: boolean;
  is_final: boolean;
}

export interface WorkflowTransition {
  id: string;
  workflow_id: string;
  from_state_id: string;
  to_state_id: string;
  name: string;
}

export interface ProjectMember {
  id: string;
  project_id: string;
  user_id: string;
  role: 'reporter' | 'developer' | 'admin';
  created_at: string;
}

export interface ProjectStats {
  total: number;
  open_issues: number;
  incident_count: number;
  critical_count: number;
}

export interface IssueFilter {
  project_id: string;
  status: string;
  assignee_id: string | null;
  is_incident: boolean | null;
  search: string;
  page: number;
  page_size: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
}

export interface AppError {
  code: string;
  message: string;
}

export type Priority = 'low' | 'medium' | 'high' | 'critical';
export type Severity = 'minor' | 'major' | 'critical';
export type Role = 'reporter' | 'developer' | 'admin';
