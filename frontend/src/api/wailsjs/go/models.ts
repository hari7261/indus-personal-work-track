export namespace domain {
	
	export class AddProjectMemberRequest {
	    project_id: number[];
	    user_id: number[];
	    role: string;
	
	    static createFrom(source: any = {}) {
	        return new AddProjectMemberRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.project_id = source["project_id"];
	        this.user_id = source["user_id"];
	        this.role = source["role"];
	    }
	}
	export class AssignIssueRequest {
	    issue_id: number[];
	    assignee_id?: number[];
	
	    static createFrom(source: any = {}) {
	        return new AssignIssueRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.issue_id = source["issue_id"];
	        this.assignee_id = source["assignee_id"];
	    }
	}
	export class Comment {
	    id: number[];
	    issue_id: number[];
	    user_id: number[];
	    content: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Comment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.issue_id = source["issue_id"];
	        this.user_id = source["user_id"];
	        this.content = source["content"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CreateCommentRequest {
	    issue_id: number[];
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateCommentRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.issue_id = source["issue_id"];
	        this.content = source["content"];
	    }
	}
	export class CreateIssueRequest {
	    project_id: number[];
	    title: string;
	    description: string;
	    priority: string;
	    is_incident: boolean;
	    severity?: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateIssueRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.project_id = source["project_id"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.priority = source["priority"];
	        this.is_incident = source["is_incident"];
	        this.severity = source["severity"];
	    }
	}
	export class CreateProjectRequest {
	    name: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateProjectRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	    }
	}
	export class CreateWorkflowRequest {
	    project_id: number[];
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateWorkflowRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.project_id = source["project_id"];
	        this.name = source["name"];
	    }
	}
	export class CreateWorkflowStateRequest {
	    workflow_id: number[];
	    name: string;
	    is_initial: boolean;
	    is_final: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CreateWorkflowStateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.workflow_id = source["workflow_id"];
	        this.name = source["name"];
	        this.is_initial = source["is_initial"];
	        this.is_final = source["is_final"];
	    }
	}
	export class CreateWorkflowTransitionRequest {
	    workflow_id: number[];
	    from_state_id: number[];
	    to_state_id: number[];
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateWorkflowTransitionRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.workflow_id = source["workflow_id"];
	        this.from_state_id = source["from_state_id"];
	        this.to_state_id = source["to_state_id"];
	        this.name = source["name"];
	    }
	}
	export class Issue {
	    id: number[];
	    project_id: number[];
	    title: string;
	    description: string;
	    status: string;
	    assignee_id?: number[];
	    priority: string;
	    is_incident: boolean;
	    severity?: string;
	    created_by: number[];
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Issue(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.project_id = source["project_id"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.status = source["status"];
	        this.assignee_id = source["assignee_id"];
	        this.priority = source["priority"];
	        this.is_incident = source["is_incident"];
	        this.severity = source["severity"];
	        this.created_by = source["created_by"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class IssueFilter {
	    project_id: number[];
	    status: string;
	    assignee_id?: number[];
	    is_incident?: boolean;
	    search: string;
	    page: number;
	    page_size: number;
	
	    static createFrom(source: any = {}) {
	        return new IssueFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.project_id = source["project_id"];
	        this.status = source["status"];
	        this.assignee_id = source["assignee_id"];
	        this.is_incident = source["is_incident"];
	        this.search = source["search"];
	        this.page = source["page"];
	        this.page_size = source["page_size"];
	    }
	}
	export class Project {
	    id: number[];
	    name: string;
	    description: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ProjectMember {
	    id: number[];
	    project_id: number[];
	    user_id: number[];
	    role: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new ProjectMember(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.project_id = source["project_id"];
	        this.user_id = source["user_id"];
	        this.role = source["role"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ProjectStats {
	    total: number;
	    open_issues: number;
	    incident_count: number;
	    critical_count: number;
	
	    static createFrom(source: any = {}) {
	        return new ProjectStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.open_issues = source["open_issues"];
	        this.incident_count = source["incident_count"];
	        this.critical_count = source["critical_count"];
	    }
	}
	export class TransitionIssueRequest {
	    issue_id: number[];
	    transition: string;
	
	    static createFrom(source: any = {}) {
	        return new TransitionIssueRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.issue_id = source["issue_id"];
	        this.transition = source["transition"];
	    }
	}
	export class UpdateIssueRequest {
	    id: number[];
	    title: string;
	    description: string;
	    priority: string;
	    assignee_id?: number[];
	    is_incident: boolean;
	    severity?: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateIssueRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.priority = source["priority"];
	        this.assignee_id = source["assignee_id"];
	        this.is_incident = source["is_incident"];
	        this.severity = source["severity"];
	    }
	}
	export class UpdateProjectRequest {
	    id: number[];
	    name: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateProjectRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	    }
	}
	export class User {
	    id: number[];
	    username: string;
	    role: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.username = source["username"];
	        this.role = source["role"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class WorkflowTransition {
	    id: number[];
	    workflow_id: number[];
	    from_state_id: number[];
	    to_state_id: number[];
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new WorkflowTransition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.workflow_id = source["workflow_id"];
	        this.from_state_id = source["from_state_id"];
	        this.to_state_id = source["to_state_id"];
	        this.name = source["name"];
	    }
	}
	export class WorkflowState {
	    id: number[];
	    workflow_id: number[];
	    name: string;
	    is_initial: boolean;
	    is_final: boolean;
	
	    static createFrom(source: any = {}) {
	        return new WorkflowState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.workflow_id = source["workflow_id"];
	        this.name = source["name"];
	        this.is_initial = source["is_initial"];
	        this.is_final = source["is_final"];
	    }
	}
	export class Workflow {
	    id: number[];
	    project_id: number[];
	    name: string;
	    states?: WorkflowState[];
	    transitions?: WorkflowTransition[];
	
	    static createFrom(source: any = {}) {
	        return new Workflow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.project_id = source["project_id"];
	        this.name = source["name"];
	        this.states = this.convertValues(source["states"], WorkflowState);
	        this.transitions = this.convertValues(source["transitions"], WorkflowTransition);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

