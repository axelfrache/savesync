export interface Source {
    id: number;
    name: string;
    path: string;
    exclusions: string[];
    target_id: number | null;
    schedule_id: number | null;
    created_at: string;
    updated_at: string;
}

export interface CreateSourceRequest {
    name: string;
    path: string;
    exclusions: string[];
    target_id: number | null;
    schedule_id?: number | null;
}

export interface UpdateSourceRequest {
    name: string;
    path: string;
    exclusions: string[];
    target_id: number | null;
    schedule_id?: number | null;
}

export interface Target {
    id: number;
    name: string;
    type: 'local' | 's3' | 'sftp';
    config: Record<string, string>;
    created_at: string;
    updated_at: string;
}

export interface CreateTargetRequest {
    name: string;
    type: 'local' | 's3' | 'sftp';
    config: Record<string, string>;
}

export interface Snapshot {
    id: number;
    source_id: number;
    target_id: number;
    status: 'pending' | 'running' | 'success' | 'failed';
    file_count: number;
    total_bytes: number;
    delta_bytes: number;
    error?: string;
    created_at: string;
    completed_at?: string;
}

export interface Job {
    id: number;
    type: 'backup' | 'restore';
    source_id?: number;
    snapshot_id?: number;
    status: 'pending' | 'running' | 'success' | 'failed';
    error?: string;
    started_at: string;
    ended_at?: string;
}

export interface BackupResponse {
    job_id: number;
    status: string;
}

export interface ApiResponse<T> {
    data: T;
}

export interface ApiError {
    error: {
        message: string;
        code?: string;
    };
}
