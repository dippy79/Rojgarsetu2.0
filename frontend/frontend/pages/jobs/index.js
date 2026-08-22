// frontend/pages/jobs/index.js - Jobs listing page with filters
import { useState, useEffect } from 'react';
import Link from 'next/link';
import useSWR from 'swr';
import { jobsAPI } from '../../lib/api';

// Loading skeleton component
function JobCardSkeleton() {
    return (
        <div className="job-card-skeleton">
            <div className="skeleton-title"></div>
            <div className="skeleton-company"></div>
            <div className="skeleton-details">
                <div className="skeleton-tag"></div>
                <div className="skeleton-tag"></div>
                <div className="skeleton-tag"></div>
            </div>
            <div className="skeleton-footer">
                <div className="skeleton-date"></div>
                <div className="skeleton-btn"></div>
            </div>
        </div>
    );
}

// Job card component
function JobCard({ job }) {
    const formatDate = (dateStr) => {
        if (!dateStr) return '';
        const date = new Date(dateStr);
        return date.toLocaleDateString('en-IN', { day: 'numeric', month: 'short', year: 'numeric' });
    };

    const getDaysLeft = (dateStr) => {
        if (!dateStr) return null;
        const lastDate = new Date(dateStr);
        const today = new Date();
        const diffTime = lastDate - today;
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
        return diffDays;
    };

    const daysLeft = getDaysLeft(job.last_date);

    return (
        <div className="job-card">
            <div className="job-card-header">
                <h3 className="job-title">
                    <Link href={`/jobs/${job.id}`}>{job.title}</Link>
                </h3>
                {job.is_featured && <span className="featured-badge">Featured</span>}
            </div>
            
            <div className="job-company">
                {job.company_name || job.source || 'RojgarSetu'}
            </div>
            
            <div className="job-meta">
                <span className="meta-item">
                    <i className="icon-location"></i>
                    {job.location || 'All India'}
                </span>
                <span className="meta-item">
                    <i className="icon-briefcase"></i>
                    {job.type === 'full-time' ? 'Full Time' : job.type === 'part-time' ? 'Part Time' : job.type || 'Full Time'}
                </span>
                <span className="meta-item">
                    <i className="icon-graduation"></i>
                    {job.experience_required || 'Freshers'}
                </span>
            </div>
            
            <div className="job-tags">
                {job.category && <span className="tag">{job.category}</span>}
                {job.education_required && <span className="tag">{job.education_required}</span>}
            </div>
            
            <div className="job-card-footer">
                <div className="job-deadline">
                    {daysLeft !== null && (
                        <span className={`days-left ${daysLeft < 7 ? 'urgent' : daysLeft < 30 ? 'warning' : ''}`}>
                            {daysLeft > 0 ? `${daysLeft} days left` : 'Expired'}
                        </span>
                    )}
                    <span className="last-date">Apply by: {formatDate(job.last_date)}</span>
                </div>
                <Link href={`/jobs/${job.id}`} className="apply-btn">
                    View Details
                </Link>
            </div>
        </div>
    );
}

// Filter sidebar component
function FilterSidebar({ filters, setFilters, categories, jobTypes }) {
    const handleChange = (key, value) => {
        setFilters(prev => ({ ...prev, [key]: value }));
    };

    return (
        <div className="filter-sidebar">
            <div className="filter-section">
                <h4>Search</h4>
                <input
                    type="text"
                    placeholder="Search jobs..."
                    value={filters.search}
                    onChange={(e) => handleChange('search', e.target.value)}
                    className="filter-input"
                />
            </div>

            <div className="filter-section">
                <h4>Category</h4>
                <select
                    value={filters.category}
                    onChange={(e) => handleChange('category', e.target.value)}
                    className="filter-select"
                >
                    <option value="">All Categories</option>
                    {categories?.map(cat => (
                        <option key={cat} value={cat}>{cat}</option>
                    ))}
                </select>
            </div>

            <div className="filter-section">
                <h4>Job Type</h4>
                <div className="filter-options">
                    {jobTypes?.map(type => (
                        <label key={type.value} className="filter-checkbox">
                            <input
                                type="checkbox"
                                checked={filters.type === type.value}
                                onChange={() => handleChange('type', filters.type === type.value ? '' : type.value)}
                            />
                            <span>{type.label}</span>
                        </label>
                    ))}
                </div>
            </div>

            <div className="filter-section">
                <h4>Location</h4>
                <input
                    type="text"
                    placeholder="Enter location..."
                    value={filters.location}
                    onChange={(e) => handleChange('location', e.target.value)}
                    className="filter-input"
                />
            </div>

            <div className="filter-section">
                <h4>Experience</h4>
                <select
                    value={filters.experience}
                    onChange={(e) => handleChange('experience', e.target.value)}
                    className="filter-select"
                >
                    <option value="">Any Experience</option>
                    <option value="0-1">0-1 Years</option>
                    <option value="1-3">1-3 Years</option>
                    <option value="3-5">3-5 Years</option>
                    <option value="5-10">5-10 Years</option>
                    <option value="10+">10+ Years</option>
                    <option value="freshers">Freshers Only</option>
                </select>
            </div>

            <button 
                className="clear-filters-btn"
                onClick={() => setFilters({
                    search: '',
                    category: '',
                    type: '',
                    location: '',
                    experience: '',
                    page: 1
                })}
            >
                Clear All Filters
            </button>
        </div>
    );
}

// Main Jobs page component
export default function JobsPage() {
    const [filters, setFilters] = useState({
        search: '',
        category: '',
        type: '',
        location: '',
        experience: '',
        page: 1,
        limit: 10
    });

    // Build query string from filters
    const buildQueryString = () => {
        const params = new URLSearchParams();
        if (filters.search) params.append('q', filters.search);
        if (filters.category) params.append('category', filters.category);
        if (filters.type) params.append('type', filters.type);
        if (filters.location) params.append('location', filters.location);
        if (filters.experience) params.append('experience', filters.experience);
        params.append('page', filters.page);
        params.append('limit', filters.limit);
        return params.toString();
    };

    // Fetch jobs with filters
    const { data: jobsData, error: jobsError, isLoading } = useSWR(
        `/api/jobs?${buildQueryString()}`,
        (url) => jobsAPI.getJobs({ 
            q: filters.search || undefined,
            category: filters.category || undefined,
            type: filters.type || undefined,
            location: filters.location || undefined,
            experience: filters.experience || undefined,
            page: filters.page,
            limit: filters.limit
        }).then(res => res.data),
        { 
            revalidateOnFocus: false,
            dedupingInterval: 5000
        }
    );

    // Fetch categories
    const { data: categoriesData } = useSWR(
        '/api/jobs/categories',
        () => jobsAPI.getCategories().then(res => res.data),
        { revalidateOnFocus: false }
    );

    const categories = categoriesData?.categories || [];
    const jobs = jobsData?.jobs || [];
    const pagination = jobsData?.pagination || { page: 1, limit: 10, total: 0, totalPages: 0 };

    const jobTypes = [
        { value: 'full-time', label: 'Full Time' },
        { value: 'part-time', label: 'Part Time' },
        { value: 'contract', label: 'Contract' },
        { value: 'internship', label: 'Internship' }
    ];

    const handlePageChange = (newPage) => {
        setFilters(prev => ({ ...prev, page: newPage }));
        window.scrollTo({ top: 0, behavior: 'smooth' });
    };

    return (
        <div className="jobs-page">
            <div className="jobs-container">
                <aside className="jobs-sidebar">
                    <FilterSidebar 
                        filters={filters} 
                        setFilters={setFilters}
                        categories={categories}
                        jobTypes={jobTypes}
                    />
                </aside>
                
                <main className="jobs-main">
                    <div className="jobs-header">
                        <h1>Find Your Dream Job</h1>
                        <p>{pagination.total} jobs available</p>
                    </div>

                    {jobsError && (
                        <div className="error-message">
                            <p>Failed to load jobs. Please try again.</p>
                            <button onClick={() => window.location.reload()}>Retry</button>
                        </div>
                    )}

                    {isLoading ? (
                        <div className="jobs-list">
                            {[...Array(5)].map((_, i) => (
                                <JobCardSkeleton key={i} />
                            ))}
                        </div>
                    ) : jobs.length === 0 ? (
                        <div className="no-jobs">
                            <h3>No jobs found</h3>
                            <p>Try adjusting your filters or search criteria</p>
                            <button 
                                onClick={() => setFilters({
                                    search: '',
                                    category: '',
                                    type: '',
                                    location: '',
                                    experience: '',
                                    page: 1
                                })}
                            >
                                Clear Filters
                            </button>
                        </div>
                    ) : (
                        <>
                            <div className="jobs-list">
                                {jobs.map(job => (
                                    <JobCard key={job.id} job={job} />
                                ))}
                            </div>

                            {pagination.totalPages > 1 && (
                                <div className="pagination">
                                    <button
                                        disabled={pagination.page <= 1}
                                        onClick={() => handlePageChange(pagination.page - 1)}
                                        className="page-btn"
                                    >
                                        Previous
                                    </button>
                                    
                                    <span className="page-info">
                                        Page {pagination.page} of {pagination.totalPages}
                                    </span>
                                    
                                    <button
                                        disabled={pagination.page >= pagination.totalPages}
                                        onClick={() => handlePageChange(pagination.page + 1)}
                                        className="page-btn"
                                    >
                                        Next
                                    </button>
                                </div>
                            )}
                        </>
                    )}
                </main>
            </div>

            <style jsx>{`
                .jobs-page {
                    min-height: 100vh;
                    background: #f3f4f6;
                }

                .jobs-container {
                    max-width: 1400px;
                    margin: 0 auto;
                    padding: 20px;
                    display: grid;
                    grid-template-columns: 280px 1fr;
                    gap: 24px;
                }

                @media (max-width: 1024px) {
                    .jobs-container {
                        grid-template-columns: 1fr;
                    }
                }

                .jobs-sidebar {
                    position: sticky;
                    top: 20px;
                    height: fit-content;
                }

                .jobs-main {
                    min-width: 0;
                }

                .jobs-header {
                    margin-bottom: 24px;
                }

                .jobs-header h1 {
                    font-size: 28px;
                    font-weight: 700;
                    color: #1f2937;
                    margin: 0 0 8px 0;
                }

                .jobs-header p {
                    color: #6b7280;
                    margin: 0;
                }

                .jobs-list {
                    display: flex;
                    flex-direction: column;
                    gap: 16px;
                }

                .error-message, .no-jobs {
                    text-align: center;
                    padding: 48px 24px;
                    background: white;
                    border-radius: 12px;
                    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
                }

                .error-message button, .no-jobs button {
                    margin-top: 16px;
                    padding: 10px 24px;
                    background: #4f46e5;
                    color: white;
                    border: none;
                    border-radius: 6px;
                    cursor: pointer;
                    font-weight: 500;
                }

                .pagination {
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    gap: 16px;
                    margin-top: 32px;
                    padding: 16px;
                    background: white;
                    border-radius: 8px;
                }

                .page-btn {
                    padding: 8px 16px;
                    background: #4f46e5;
                    color: white;
                    border: none;
                    border-radius: 6px;
                    cursor: pointer;
                    transition: background 0.2s;
                }

                .page-btn:disabled {
                    background: #d1d5db;
                    cursor: not-allowed;
                }

                .page-btn:not(:disabled):hover {
                    background: #4338ca;
                }

                .page-info {
                    color: #6b7280;
                    font-weight: 500;
                }

                /* Job Card Styles */
                .job-card {
                    background: white;
                    border-radius: 12px;
                    padding: 20px;
                    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
                    transition: box-shadow 0.2s, transform 0.2s;
                }

                .job-card:hover {
                    box-shadow: 0 4px 12px rgba(0,0,0,0.15);
                    transform: translateY(-2px);
                }

                .job-card-header {
                    display: flex;
                    justify-content: space-between;
                    align-items: flex-start;
                    margin-bottom: 8px;
                }

                .job-title {
                    margin: 0;
                    font-size: 18px;
                    font-weight: 600;
                }

                .job-title a {
                    color: #1f2937;
                    text-decoration: none;
                }

                .job-title a:hover {
                    color: #4f46e5;
                }

                .featured-badge {
                    background: #fef3c7;
                    color: #d97706;
                    padding: 4px 8px;
                    border-radius: 4px;
                    font-size: 12px;
                    font-weight: 600;
                }

                .job-company {
                    color: #6b7280;
                    font-size: 14px;
                    margin-bottom: 12px;
                }

                .job-meta {
                    display: flex;
                    flex-wrap: wrap;
                    gap: 16px;
                    margin-bottom: 12px;
                    color: #6b7280;
                    font-size: 14px;
                }

                .job-tags {
                    display: flex;
                    flex-wrap: wrap;
                    gap: 8px;
                    margin-bottom: 16px;
                }

                .tag {
                    background: #e5e7eb;
                    color: #4b5563;
                    padding: 4px 10px;
                    border-radius: 4px;
                    font-size: 12px;
                }

                .job-card-footer {
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                    padding-top: 16px;
                    border-top: 1px solid #e5e7eb;
                }

                .job-deadline {
                    display: flex;
                    flex-direction: column;
                    gap: 4px;
                }

                .days-left {
                    font-weight: 600;
                    color: #059669;
                }

                .days-left.warning {
                    color: #d97706;
                }

                .days-left.urgent {
                    color: #dc2626;
                }

                .last-date {
                    font-size: 12px;
                    color: #9ca3af;
                }

                .apply-btn {
                    background: #4f46e5;
                    color: white;
                    padding: 10px 20px;
                    border-radius: 6px;
                    text-decoration: none;
                    font-weight: 500;
                    transition: background 0.2s;
                }

                .apply-btn:hover {
                    background: #4338ca;
                }

                /* Skeleton styles */
                .job-card-skeleton {
                    background: white;
                    border-radius: 12px;
                    padding: 20px;
                    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
                }

                .skeleton-title, .skeleton-company, .skeleton-tag, .skeleton-date, .skeleton-btn {
                    background: linear-gradient(90deg, #e5e7eb 25%, #f3f4f6 50%, #e5e7eb 75%);
                    background-size: 200% 100%;
                    animation: skeleton-loading 1.5s infinite;
                    border-radius: 4px;
                }

                .skeleton-title {
                    height: 24px;
                    width: 70%;
                    margin-bottom: 12px;
                }

                .skeleton-company {
                    height: 16px;
                    width: 40%;
                    margin-bottom: 16px;
                }

                .skeleton-details {
                    display: flex;
                    gap: 8px;
                    margin-bottom: 16px;
                }

                .skeleton-tag {
                    height: 24px;
                    width: 80px;
                }

                .skeleton-footer {
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                    padding-top: 16px;
                    border-top: 1px solid #e5e7eb;
                }

                .skeleton-date {
                    height: 16px;
                    width: 100px;
                }

                .skeleton-btn {
                    height: 36px;
                    width: 100px;
                }

                @keyframes skeleton-loading {
                    0% { background-position: 200% 0; }
                    100% { background-position: -200% 0; }
                }

                /* Filter Sidebar Styles */
                .filter-sidebar {
                    background: white;
                    border-radius: 12px;
                    padding: 20px;
                    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
                }

                .filter-section {
                    margin-bottom: 24px;
                }

                .filter-section h4 {
                    margin: 0 0 12px 0;
                    font-size: 14px;
                    font-weight: 600;
                    color: #374151;
                }

                .filter-input, .filter-select {
                    width: 100%;
                    padding: 10px 12px;
                    border: 1px solid #d1d5db;
                    border-radius: 6px;
                    font-size: 14px;
                    transition: border-color 0.2s;
                }

                .filter-input:focus, .filter-select:focus {
                    outline: none;
                    border-color: #4f46e5;
                }

                .filter-options {
                    display: flex;
                    flex-direction: column;
                    gap: 8px;
                }

                .filter-checkbox {
                    display: flex;
                    align-items: center;
                    gap: 8px;
                    cursor: pointer;
                    font-size: 14px;
                    color: #4b5563;
                }

                .filter-checkbox input {
                    width: 16px;
                    height: 16px;
                    accent-color: #4f46e5;
                }

                .clear-filters-btn {
                    width: 100%;
                    padding: 10px;
                    background: transparent;
                    border: 1px solid #d1d5db;
                    border-radius: 6px;
                    color: #6b7280;
                    cursor: pointer;
                    font-size: 14px;
                    transition: all 0.2s;
                }

                .clear-filters-btn:hover {
                    background: #f3f4f6;
                    color: #374151;
                }
            `}</style>
        </div>
    );
}
