// frontend/lib/api.js - API utilities for frontend
import axios from 'axios';

const API_BASE = process.env.NEXT_PUBLIC_API_BASE || 'http://localhost:3000';

// Create axios instance with default config
const api = axios.create({
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

// Basic fetcher for SWR
export const fetcher = (url) => api.get(url).then(r => r.data);

// Auth API
export const authAPI = {
    register: (data) => api.post('/api/auth/register', data),
    login: (data) => api.post('/api/auth/login', data),
    logout: () => api.post('/api/auth/logout'),
    getProfile: () => api.get('/api/auth/profile'),
    updateProfile: (data) => api.put('/api/auth/profile', data),
    refreshToken: () => api.post('/api/auth/refresh')
};

// Jobs API
export const jobsAPI = {
    // Get all jobs with filters
    getJobs: (params) => api.get('/api/jobs', { params }),
    
    // Get single job
    getJob: (id) => api.get(`/api/jobs/${id}`),
    
    // Search jobs
    searchJobs: (q) => api.get('/api/jobs/search', { params: { q } }),
    
    // Get featured jobs
    getFeaturedJobs: (limit = 5) => api.get('/api/jobs/featured', { params: { limit } }),
    
    // Get job categories
    getCategories: () => api.get('/api/jobs/categories'),
    
    // Get similar jobs
    getSimilarJobs: (id) => api.get(`/api/jobs/${id}/similar`),
    
    // Create job (company/admin)
    createJob: (data) => api.post('/api/jobs', data),
    
    // Update job (company/admin)
    updateJob: (id, data) => api.put(`/api/jobs/${id}`, data),
    
    // Delete job (company/admin)
    deleteJob: (id) => api.delete(`/api/jobs/${id}`),
    
    // Save job (candidate)
    saveJob: (jobId) => api.post('/api/jobs/${jobId}/save'),
    
    // Unsave job (candidate)
    unsaveJob: (jobId) => api.delete(`/api/jobs/${jobId}/save`),
    
    // Get saved jobs (candidate)
    getSavedJobs: () => api.get('/api/jobs/saved'),
    
    // Apply to job (candidate)
    applyToJob: (jobId, data) => api.post(`/api/jobs/${jobId}/apply`, data),
    
    // Withdraw application (candidate)
    withdrawApplication: (jobId) => api.delete(`/api/jobs/${jobId}/apply`),
    
    // Get my applications (candidate)
    getMyApplications: () => api.get('/api/jobs/applications')
};

// Courses API
export const coursesAPI = {
    // Get all courses with filters
    getCourses: (params) => api.get('/api/courses', { params }),
    
    // Get single course
    getCourse: (id) => api.get(`/api/courses/${id}`),
    
    // Search courses
    searchCourses: (q) => api.get('/api/courses/search', { params: { q } }),
    
    // Get featured courses
    getFeaturedCourses: (limit = 5) => api.get('/api/courses/featured', { params: { limit } }),
    
    // Get course categories
    getCategories: () => api.get('/api/courses/categories'),
    
    // Get similar courses
    getSimilarCourses: (id) => api.get(`/api/courses/${id}/similar`),
    
    // Create course (company/admin)
    createCourse: (data) => api.post('/api/courses', data),
    
    // Update course (company/admin)
    updateCourse: (id, data) => api.put(`/api/courses/${id}`, data),
    
    // Delete course (company/admin)
    deleteCourse: (id) => api.delete(`/api/courses/${id}`)
};

// Companies API
export const companiesAPI = {
    // Get all companies
    getCompanies: (params) => api.get('/api/companies', { params }),
    
    // Get single company
    getCompany: (id) => api.get(`/api/companies/${id}`),
    
    // Get company jobs
    getCompanyJobs: (id) => api.get(`/api/companies/${id}/jobs`),
    
    // Get company courses
    getCompanyCourses: (id) => api.get(`/api/companies/${id}/courses`),
    
    // Update company (company admin)
    updateCompany: (id, data) => api.put(`/api/companies/${id}`, data),
    
    // Get company applicants (company admin)
    getCompanyApplicants: (companyId, params) => api.get(`/api/companies/${companyId}/applicants`, { params }),
    
    // Update applicant status (company admin)
    updateApplicantStatus: (companyId, applicationId, status) => 
        api.put(`/api/companies/${companyId}/applicants/${applicationId}`, { status })
};

// Applications API
export const applicationsAPI = {
    // Get all applications (admin/company)
    getApplications: (params) => api.get('/api/applications', { params }),
    
    // Get single application
    getApplication: (id) => api.get(`/api/applications/${id}`),
    
    // Update application status
    updateApplication: (id, status, notes) => api.put(`/api/applications/${id}`, { status, notes })
};

// Notifications API
export const notificationsAPI = {
    // Get my notifications
    getNotifications: () => api.get('/api/notifications'),
    
    // Mark notification as read
    markAsRead: (id) => api.put(`/api/notifications/${id}/read`),
    
    // Mark all as read
    markAllAsRead: () => api.put('/api/notifications/read-all'),
    
    // Delete notification
    deleteNotification: (id) => api.delete(`/api/notifications/${id}`),
    
    // Get unread count
    getUnreadCount: () => api.get('/api/notifications/unread-count')
};

export default api;
