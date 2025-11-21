import axios, { type AxiosResponse } from 'axios';
import type {
    Source,
    CreateSourceRequest,
    UpdateSourceRequest,
    Target,
    CreateTargetRequest,
    Snapshot,
    Job,
    BackupResponse,
} from '@/types/api';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

const axiosInstance = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Helper to unwrap API response
function unwrap<T>(promise: Promise<AxiosResponse<{ data: T }>>): Promise<T> {
    return promise.then((response) => response.data.data);
}

export const sourcesApi = {
    list: () => unwrap<Source[]>(axiosInstance.get('/sources')),
    get: (id: number) => unwrap<Source>(axiosInstance.get(`/sources/${id}`)),
    create: (data: CreateSourceRequest) => unwrap<Source>(axiosInstance.post('/sources', data)),
    update: (id: number, data: UpdateSourceRequest) =>
        unwrap<Source>(axiosInstance.put(`/sources/${id}`, data)),
    delete: (id: number) => axiosInstance.delete(`/sources/${id}`),
    run: (id: number) => unwrap<BackupResponse>(axiosInstance.post(`/sources/${id}/run`)),
};

export const targetsApi = {
    list: () => unwrap<Target[]>(axiosInstance.get('/targets')),
    get: (id: number) => unwrap<Target>(axiosInstance.get(`/targets/${id}`)),
    create: (data: CreateTargetRequest) => unwrap<Target>(axiosInstance.post('/targets', data)),
    update: (id: number, data: Partial<CreateTargetRequest>) =>
        unwrap<Target>(axiosInstance.put(`/targets/${id}`, data)),
    delete: (id: number) => axiosInstance.delete(`/targets/${id}`),
};

export const snapshotsApi = {
    list: () => unwrap<Snapshot[]>(axiosInstance.get('/snapshots')),
    get: (id: number) => unwrap<Snapshot>(axiosInstance.get(`/snapshots/${id}`)),
    restore: (id: number) => axiosInstance.post(`/snapshots/${id}/restore`),
};

export const jobsApi = {
    list: () => unwrap<Job[]>(axiosInstance.get('/jobs')),
    get: (id: number) => unwrap<Job>(axiosInstance.get(`/jobs/${id}`)),
};

export const healthApi = {
    check: () => axiosInstance.get('/health'),
};
