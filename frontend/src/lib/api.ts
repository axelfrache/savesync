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

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const axiosInstance = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Add token to requests
axiosInstance.interceptors.request.use((config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// Helper to unwrap API response
function unwrap<T>(promise: Promise<AxiosResponse<{ data: T }>>): Promise<T> {
    return promise.then((response) => response.data.data);
}

// Auth API
export const authApi = {
    login: (email: string, password: string) =>
        axiosInstance.post('/auth/login', { email, password }).then((res) => res.data.data),
    register: (email: string, password: string) =>
        axiosInstance.post('/auth/register', { email, password }).then((res) => res.data.data),
    me: () => axiosInstance.get('/auth/me').then((res) => res.data.data),
};

export const sourcesApi = {
    list: () => unwrap<Source[]>(axiosInstance.get('/api/sources')),
    get: (id: number) => unwrap<Source>(axiosInstance.get(`/api/sources/${id}`)),
    create: (data: CreateSourceRequest) => unwrap<Source>(axiosInstance.post('/api/sources', data)),
    update: (id: number, data: UpdateSourceRequest) =>
        unwrap<Source>(axiosInstance.put(`/api/sources/${id}`, data)),
    delete: (id: number) => axiosInstance.delete(`/api/sources/${id}`),
    run: (id: number) => unwrap<BackupResponse>(axiosInstance.post(`/api/sources/${id}/run`)),
};

export const targetsApi = {
    list: () => unwrap<Target[]>(axiosInstance.get('/api/targets')),
    get: (id: number) => unwrap<Target>(axiosInstance.get(`/api/targets/${id}`)),
    create: (data: CreateTargetRequest) => unwrap<Target>(axiosInstance.post('/api/targets', data)),
    update: (id: number, data: Partial<CreateTargetRequest>) =>
        unwrap<Target>(axiosInstance.put(`/api/targets/${id}`, data)),
    delete: (id: number) => axiosInstance.delete(`/api/targets/${id}`),
};

export const snapshotsApi = {
    list: () => unwrap<Snapshot[]>(axiosInstance.get('/api/snapshots')),
    get: (id: number) => unwrap<Snapshot>(axiosInstance.get(`/api/snapshots/${id}`)),
    restore: (id: number) => axiosInstance.post(`/api/snapshots/${id}/restore`),
};

export const jobsApi = {
    list: () => unwrap<Job[]>(axiosInstance.get('/api/jobs')),
    get: (id: number) => unwrap<Job>(axiosInstance.get(`/api/jobs/${id}`)),
};

export const healthApi = {
    check: () => axiosInstance.get('/health'),
};
