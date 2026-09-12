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
  'https://api.rojgarsetu.in',
];

function normalizeBase(base?: string) {
  if (!base) return null;
  let b = String(base).replace(/\/+$/, '');
  return b;
}

export function getApiBaseUrl(): string {
  // Use a stable default for both client and server during SSR/Build
  const defaultUrl = 'http://localhost:3001';

  if (typeof window !== 'undefined') {
    // Client-side: try to find injected global configs
    // @ts-ignore
    if (window.__ROJGAR_API__) return normalizeBase(window.__ROJGAR_API__) || defaultUrl;
    // @ts-ignore
    if (window.__ROJGAR_API_ENV__) return normalizeBase(window.__ROJGAR_API_ENV__) || defaultUrl;
  }

  // Check process.env (Next.js will replace these at build time for client,
  // or use them at runtime for server-side)
  const envUrl = process.env.NEXT_PUBLIC_API_URL || process.env.NEXT_PUBLIC_API_BASE;

  if (!envUrl) {
    if (typeof window === 'undefined') {
        // This is normal during build if .env is missing
        return defaultUrl;
    }
    return defaultUrl;
  }

  return normalizeBase(envUrl) || defaultUrl;
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
        api.post<ApiResponse<{ user: User }>>('/api/v1/auth/register', data),

    login: (data: { email: string; password: string }) =>
        api.post<ApiResponse<{ user: User; token?: string }>>('/api/v1/auth/login', data),

    logout: () => api.post<ApiResponse<any>>('/api/v1/auth/logout'),

    getProfile: () => api.get<ApiResponse<UserProfile>>('/api/v1/auth/me'),

    updateProfile: (data: Partial<UserProfile>) =>
        api.put<ApiResponse<any>>('/api/v1/candidates/me', data),

    refreshToken: () =>
        api.post<ApiResponse<any>>('/api/v1/auth/refresh')
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
