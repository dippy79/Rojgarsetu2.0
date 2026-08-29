// frontend/src/lib/api.ts - API utilities for frontend with TypeScript
import axios, { AxiosInstance } from 'axios';
import type {
  ApiResponse,
  Job,
  Course,
  User,
  UserProfile,
} from '../types';

// Central API base-URL resolver (Consolidated from apiConfig.js)
const NORMALIZED_DEFAULTS = [
  'http://localhost:3001',
  'http://localhost:8083',
  'https://api.rojgarsetu.in',
];

function normalizeBase(base?: string) {
  if (!base) return null;
  let b = String(base).replace(/\/+$/, '');
  return b;
}

export function getApiBaseUrl(): string {
  if (typeof window !== 'undefined') {
    // @ts-ignore
    if (window.__ROJGAR_API__) return normalizeBase(window.__ROJGAR_API__) || 'http://localhost:3001';
    // @ts-ignore
    if (window.__ROJGAR_API_ENV__) return normalizeBase(window.__ROJGAR_API_ENV__) || 'http://localhost:3001';
  }

  const envUrl = process.env.NEXT_PUBLIC_API_URL || process.env.NEXT_PUBLIC_API_BASE;
  if (envUrl) return normalizeBase(envUrl) || 'http://localhost:3001';

  return 'http://localhost:3001'; // Default to API Gateway
}

const API_BASE = getApiBaseUrl();

// Create axios instance with default config
const api: AxiosInstance = axios.create({
    baseURL: API_BASE,
    withCredentials: true, // IMPORTANT: Enable cookie-based auth
    headers: {
        'Content-Type': 'application/json'
    }
});

// Handle response errors
api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            if (typeof window !== 'undefined') {
                window.location.href = '/login';
            }
        }
        return Promise.reject(error);
    }
);

// Basic fetcher for React Query
export const fetcher = async (url: string) => {
    const response = await api.get(url);
    return response.data;
};

// Auth API
export const authAPI = {
    register: (data: { email: string; password: string; role: string; firstName?: string; lastName?: string; companyName?: string }) =>
        api.post<ApiResponse<{ user: User }>>('/api/auth/register', data),

    login: (data: { email: string; password: string }) =>
        api.post<ApiResponse<{ user: User }>>('/api/auth/login', data),

    logout: () => api.post<ApiResponse<any>>('/api/auth/logout'),

    getProfile: () => api.get<ApiResponse<UserProfile>>('/api/v1/candidates/me'),

    updateProfile: (data: Partial<UserProfile>) =>
        api.put<ApiResponse<any>>('/api/v1/candidates/me', data),

    refreshToken: () =>
        api.post<ApiResponse<any>>('/api/auth/refresh')
};

// Jobs API
export const jobsAPI = {
    getJobs: (params?: any) => api.get<ApiResponse<Job[]>>('/api/v1/jobs', { params }),
    getJob: (id: string) => api.get<ApiResponse<Job>>(`/api/v1/jobs/${id}`),
    searchJobs: (q: string, page?: number, limit?: number) =>
        api.get<ApiResponse<Job[]>>('/api/v1/jobs/search', { params: { q, page, limit } }),
};

// Courses API
export const coursesAPI = {
    getCourses: (params?: any) => api.get<ApiResponse<Course[]>>('/api/v1/courses', { params }),
    getCourse: (id: string) => api.get<ApiResponse<Course>>(`/api/v1/courses/${id}`),
};

export default api;
