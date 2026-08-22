// frontend/src/lib/api.ts - API utilities for frontend with TypeScript
import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import type { 
  ApiResponse, 
  Job, 
  Course, 
  User, 
  UserProfile,
  Pagination 
} from '../types';

const API_BASE = process.env.NEXT_PUBLIC_API_BASE || 'http://localhost:3000';

// Create axios instance with default config
const api: AxiosInstance = axios.create({
    baseURL: API_BASE,
    headers: {
        'Content-Type': 'application/json'
    }
});

// Add token to requests if available
api.interceptors.request.use((config) => {
    if (typeof window !== 'undefined') {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
    }
    return config;
});

// Handle response errors
api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            if (typeof window !== 'undefined') {
                localStorage.removeItem('token');
                localStorage.removeItem('user');
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
        api.post<ApiResponse<{ user: User; accessToken: string; refreshToken: string }>>('/api/auth/register', data),
    
    login: (data: { email: string; password: string }) => 
        api.post<ApiResponse<{ user: User; accessToken: string; refreshToken: string }>>('/api/auth/login', data),
    
    logout: () => api.post<ApiResponse>('/api/auth/logout'),
    
    getProfile: () => api.get<ApiResponse<UserProfile>>('/api/auth/profile'),
    
    updateProfile: (data: Partial<UserProfile>) => 
        api.put<ApiResponse>('/api/auth/profile', data),
    
    refreshToken: (data: { refreshToken: string }) => 
        api.post<ApiResponse<{ accessToken: string; refreshToken: string }>>('/api/auth/refresh', data)
};

// Jobs API
export const jobsAPI = {
    // Get all jobs with filters
    getJobs: (params?: { 
        category?: string; 
        type?: string; 
        location?: string; 
        salaryMin?: number; 
        salaryMax?: number; 
        search?: string; 
        page?: number; 
        limit?: number; 
        sortBy?: string; 
        sortOrder?: string;
    }) => api.get<ApiResponse<Job[], Pagination>>('/api/jobs', { params }),
    
    // Get single job
    getJob: (id: string) => api.get<ApiResponse<Job>>(`/api/jobs/${id}`),
    
    // Search jobs
    searchJobs: (q: string, page?: number, limit?: number) => 
        api.get<ApiResponse<Job[], Pagination>>('/api/jobs/search', { params: { q, page, limit } }),
    
    // Get featured jobs
    getFeaturedJobs: (limit = 5) => 
        api.get<ApiResponse<Job[]>>('/api/jobs/featured', { params: { limit } }),
    
    // Get job categories
    getCategories: () => api.get<ApiResponse<Array<{ category: string; count: number }>>>('/api/jobs/categories'),
    
    // Get similar jobs
    getSimilarJobs: (id: string, limit?: number) => 
        api.get<ApiResponse<Job[]>>(`/api/jobs/${id}/similar`, { params: { limit } }),
    
    // Create job (company/admin)
    createJob: (data: Partial<Job>) => api.post<ApiResponse<Job>>('/api/jobs', data),
    
    // Update job (company/admin)
    updateJob: (id: string, data: Partial<Job>) => api.put<ApiResponse<Job>>(`/api/jobs/${id}`, data),
    
    // Delete job (company/admin)
    deleteJob: (id: string) => api.delete<ApiResponse>(`/api/jobs/${id}`),
    
    // Save job (candidate)
    saveJob: (jobId: string) => api.post<ApiResponse>(`/api/jobs/${jobId}/save`),
    
    // Unsave job (candidate)
    unsaveJob: (jobId: string) => api.delete<ApiResponse>(`/api/jobs/${jobId}/save`),
    
    // Get saved jobs (candidate)
    getSavedJobs: () => api.get<ApiResponse<Job[]>>('/api/jobs/saved'),
    
    // Apply to job (candidate)
    applyToJob: (jobId: string, data: { coverLetter?: string; resumeUrl?: string }) => 
        api.post<ApiResponse>(`/api/jobs/${jobId}/apply`, data),
    
    // Withdraw application (candidate)
    withdrawApplication: (jobId: string) => api.delete<ApiResponse>(`/api/jobs/${jobId}/apply`),
    
    // Get my applications (candidate)
    getMyApplications: () => api.get<ApiResponse<any[]>>('/api/jobs/applications')
};

// Courses API
export const coursesAPI = {
    // Get all courses with filters
    getCourses: (params?: { 
        category?: string; 
        duration?: string; 
        mode?: string; 
        feesMin?: number; 
        feesMax?: number; 
        search?: string; 
        page?: number; 
        limit?: number; 
        sortBy?: string; 
        sortOrder?: string;
    }) => api.get<ApiResponse<Course[], Pagination>>('/api/courses', { params }),
    
    // Get single course
    getCourse: (id: string) => api.get<ApiResponse<Course>>(`/api/courses/${id}`),
    
    // Search courses
    searchCourses: (q: string, page?: number, limit?: number) => 
        api.get<ApiResponse<Course[], Pagination>>('/api/courses/search', { params: { q, page, limit } }),
    
    // Get featured courses
    getFeaturedCourses: (limit = 5) => 
        api.get<ApiResponse<Course[]>>('/api/courses/featured', { params: { limit } }),
    
    // Get course categories
    getCategories: () => api.get<ApiResponse<Array<{ category: string; count: number }>>>('/api/courses/categories'),
    
    // Get similar courses
    getSimilarCourses: (id: string, limit?: number) => 
        api.get<ApiResponse<Course[]>>(`/api/courses/${id}/similar`, { params: { limit } }),
    
    // Create course (company/admin)
    createCourse: (data: Partial<Course>) => api.post<ApiResponse<Course>>('/api/courses', data),
    
    // Update course (company/admin)
    updateCourse: (id: string, data: Partial<Course>) => api.put<ApiResponse<Course>>(`/api/courses/${id}`, data),
    
    // Delete course (company/admin)
    deleteCourse: (id: string) => api.delete<ApiResponse>(`/api/courses/${id}`)
};

// Companies API
export const companiesAPI = {
    // Get all companies
    getCompanies: (params?: any) => api.get<ApiResponse<any[]>>('/api/companies', { params }),
    
    // Get single company
    getCompany: (id: string) => api.get<ApiResponse<any>>(`/api/companies/${id}`),
    
    // Get company jobs
    getCompanyJobs: (id: string) => api.get<ApiResponse<Job[]>>(`/api/companies/${id}/jobs`),
    
    // Get company courses
    getCompanyCourses: (id: string) => api.get<ApiResponse<Course[]>>(`/api/companies/${id}/courses`),
    
    // Update company (company admin)
    updateCompany: (id: string, data: any) => api.put<ApiResponse>(`/api/companies/${id}`, data),
    
    // Get company applicants (company admin)
    getCompanyApplicants: (companyId: string, params?: any) => 
        api.get<ApiResponse<any[]>>(`/api/companies/${companyId}/applicants`, { params }),
    
    // Update applicant status (company admin)
    updateApplicantStatus: (companyId: string, applicationId: string, status: string) => 
        api.put<ApiResponse>(`/api/companies/${companyId}/applicants/${applicationId}`, { status })
};

// Applications API
export const applicationsAPI = {
    // Get all applications (admin/company)
    getApplications: (params?: any) => api.get<ApiResponse<any[]>>('/api/applications', { params }),
    
    // Get single application
    getApplication: (id: string) => api.get<ApiResponse<any>>(`/api/applications/${id}`),
    
    // Update application status
    updateApplication: (id: string, status: string, notes?: string) => 
        api.put<ApiResponse>(`/api/applications/${id}`, { status, notes })
};

// Notifications API
export const notificationsAPI = {
    // Get my notifications
    getNotifications: () => api.get<ApiResponse<any[]>>('/api/notifications'),
    
    // Mark notification as read
    markAsRead: (id: string) => api.put<ApiResponse>(`/api/notifications/${id}/read`),
    
    // Mark all as read
    markAllAsRead: () => api.put<ApiResponse>('/api/notifications/read-all'),
    
    // Delete notification
    deleteNotification: (id: string) => api.delete<ApiResponse>(`/api/notifications/${id}`),
    
    // Get unread count
    getUnreadCount: () => api.get<ApiResponse<{ count: number }>>('/api/notifications/unread-count')
};

export default api;