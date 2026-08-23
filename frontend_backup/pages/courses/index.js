// frontend/pages/courses/index.js - Courses listing page with filters
import { useState } from 'react';
import Link from 'next/link';
import useSWR from 'swr';
import { coursesAPI } from '../../lib/api';

// Loading skeleton component
function CourseCardSkeleton() {
    return (
        <div className="course-card-skeleton">
            <div className="skeleton-title"></div>
            <div className="skeleton-provider"></div>
            <div className="skeleton-details">
                <div className="skeleton-tag"></div>
                <div className="skeleton-tag"></div>
            </div>
            <div className="skeleton-footer">
                <div className="skeleton-price"></div>
                <div className="skeleton-btn"></div>
            </div>
        </div>
    );
}

// Course card component
function CourseCard({ course }) {
    const formatFees = (fees) => {
        if (!fees || fees === '0' || fees.toLowerCase() === 'free') return 'Free';
        return '₹' + fees;
    };

    return (
        <div className="course-card">
            <div className="course-card-header">
                <h3 className="course-title">
                    <Link href={'/courses/' + course.id}>{course.name}</Link>
                </h3>
                {course.is_featured && <span className="featured-badge">Featured</span>}
            </div>
            
            <div className="course-provider">
                {course.provider || course.company_name || 'RojgarSetu'}
            </div>
            
            <div className="course-meta">
                <span className="meta-item">
                    <i className="icon-clock"></i>
                    {course.duration || 'Self-paced'}
                </span>
                <span className="meta-item">
                    <i className="icon-level"></i>
                    {course.level || 'Beginner'}
                </span>
                <span className="meta-item">
                    <i className="icon-certificate"></i>
                    {course.certificate ? 'Certificate Included' : 'No Certificate'}
                </span>
            </div>
            
            <div className="course-tags">
                {course.category && <span className="tag">{course.category}</span>}
                {course.mode && <span className="tag">{course.mode}</span>}
            </div>
            
            <div className="course-description">
                {(course.description ? course.description.substring(0, 150) : '') + '...'}
            </div>
            
            <div className="course-card-footer">
                <div className="course-fees">
                    <span className="fees-label">Fees:</span>
                    <span className="fees-amount">{formatFees(course.fees_structure)}</span>
                </div>
                <div className="course-actions">
                    <Link href={'/courses/' + course.id} className="details-btn">
                        View Details
                    </Link>
                </div>
            </div>
        </div>
    );
}

// Filter sidebar component
function FilterSidebar({ filters, setFilters, categories, modes, levels }) {
    const handleChange = (key, value) => {
        setFilters(prev => ({ ...prev, [key]: value }));
    };

    return (
        <div className="filter-sidebar">
            <div className="filter-section">
                <h4>Search</h4>
                <input
                    type="text"
                    placeholder="Search courses..."
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
                <h4>Mode</h4>
                <div className="filter-options">
                    {modes?.map(mode => (
                        <label key={mode.value} className="filter-checkbox">
                            <input
                                type="checkbox"
                                checked={filters.mode === mode.value}
                                onChange={() => handleChange('mode', filters.mode === mode.value ? '' : mode.value)}
                            />
                            <span>{mode.label}</span>
                        </label>
                    ))}
                </div>
            </div>

            <div className="filter-section">
                <h4>Level</h4>
                <select
                    value={filters.level}
                    onChange={(e) => handleChange('level', e.target.value)}
                    className="filter-select"
                >
                    <option value="">All Levels</option>
                    {levels?.map(level => (
                        <option key={level} value={level}>{level}</option>
                    ))}
                </select>
            </div>

            <div className="filter-section">
                <h4>Duration</h4>
                <select
                    value={filters.duration}
                    onChange={(e) => handleChange('duration', e.target.value)}
                    className="filter-select"
                >
                    <option value="">Any Duration</option>
                    <option value="1-month">Within 1 Month</option>
                    <option value="1-3-months">1-3 Months</option>
                    <option value="3-6-months">3-6 Months</option>
                    <option value="6-months+">6+ Months</option>
                </select>
            </div>

            <div className="filter-section">
                <h4>Fees</h4>
                <select
                    value={filters.fees}
                    onChange={(e) => handleChange('fees', e.target.value)}
                    className="filter-select"
                >
                    <option value="">Any Fees</option>
                    <option value="free">Free</option>
                    <option value="paid">Paid</option>
                </select>
            </div>

            <button 
                className="clear-filters-btn"
                onClick={() => setFilters({
                    search: '',
                    category: '',
                    mode: '',
                    level: '',
                    duration: '',
                    fees: '',
                    page: 1
                })}
            >
                Clear All Filters
            </button>
        </div>
    );
}

// Main Courses page component
export default function CoursesPage() {
    const [filters, setFilters] = useState({
        search: '',
        category: '',
        mode: '',
        level: '',
        duration: '',
        fees: '',
        page: 1,
        limit: 10
    });

    // Fetch courses with filters
    const { data: coursesData, error: coursesError, isLoading } = useSWR(
        '/api/courses',
        () => coursesAPI.getCourses({ 
            q: filters.search || undefined,
            category: filters.category || undefined,
            mode: filters.mode || undefined,
            level: filters.level || undefined,
            duration: filters.duration || undefined,
            fees: filters.fees || undefined,
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
        '/api/courses/categories',
        () => coursesAPI.getCategories().then(res => res.data),
        { revalidateOnFocus: false }
    );

    const categories = categoriesData?.categories || [];
    const courses = coursesData?.courses || [];
    const pagination = coursesData?.pagination || { page: 1, limit: 10, total: 0, totalPages: 0 };

    const modes = [
        { value: 'online', label: 'Online' },
        { value: 'offline', label: 'Offline' },
        { value: 'hybrid', label: 'Hybrid' }
    ];

    const levels = ['Beginner', 'Intermediate', 'Advanced', 'Expert'];

    const handlePageChange = (newPage) => {
        setFilters(prev => ({ ...prev, page: newPage }));
        window.scrollTo({ top: 0, behavior: 'smooth' });
    };

    return (
        <div className="courses-page">
            <div className="courses-container">
                <aside className="courses-sidebar">
                    <FilterSidebar 
                        filters={filters} 
                        setFilters={setFilters}
                        categories={categories}
                        modes={modes}
                        levels={levels}
                    />
                </aside>
                
                <main className="courses-main">
                    <div className="courses-header">
                        <h1>Upskill with Our Courses</h1>
                        <p>{pagination.total} courses available</p>
                    </div>

                    {coursesError && (
                        <div className="error-message">
                            <p>Failed to load courses. Please try again.</p>
                            <button onClick={() => window.location.reload()}>Retry</button>
                        </div>
                    )}

                    {isLoading ? (
                        <div className="courses-list">
                            {[...Array(5)].map((_, i) => (
                                <CourseCardSkeleton key={i} />
                            ))}
                        </div>
                    ) : courses.length === 0 ? (
                        <div className="no-courses">
                            <h3>No courses found</h3>
                            <p>Try adjusting your filters or search criteria</p>
                            <button 
                                onClick={() => setFilters({
                                    search: '',
                                    category: '',
                                    mode: '',
                                    level: '',
                                    duration: '',
                                    fees: '',
                                    page: 1
                                })}
                            >
                                Clear Filters
                            </button>
                        </div>
                    ) : (
                        <>
                            <div className="courses-list">
                                {courses.map(course => (
                                    <CourseCard key={course.id} course={course} />
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
                .courses-page {
                    min-height: 100vh;
                    background: #f3f4f6;
                }

                .courses-container {
                    max-width: 1400px;
                    margin: 0 auto;
                    padding: 20px;
                    display: grid;
                    grid-template-columns: 280px 1fr;
                    gap: 24px;
                }

                @media (max-width: 1024px) {
                    .courses-container {
                        grid-template-columns: 1fr;
                    }
                }

                .courses-sidebar {
                    position: sticky;
                    top: 20px;
                    height: fit-content;
                }

                .courses-main {
                    min-width: 0;
                }

                .courses-header {
                    margin-bottom: 24px;
                }

                .courses-header h1 {
                    font-size: 28px;
                    font-weight: 700;
                    color: #1f2937;
                    margin: 0 0 8px 0;
                }

                .courses-header p {
                    color: #6b7280;
                    margin: 0;
                }

                .courses-list {
                    display: flex;
                    flex-direction: column;
                    gap: 16px;
                }

                .error-message, .no-courses {
                    text-align: center;
                    padding: 48px 24px;
                    background: white;
                    border-radius: 12px;
                    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
                }

                .error-message button, .no-courses button {
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

                .course-card {
                    background: white;
                    border-radius: 12px;
                    padding: 20px;
                    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
                    transition: box-shadow 0.2s, transform 0.2s;
                }

                .course-card:hover {
                    box-shadow: 0 4px 12px rgba(0,0,0,0.15);
                    transform: translateY(-2px);
                }

                .course-card-header {
                    display: flex;
                    justify-content: space-between;
                    align-items: flex-start;
                    margin-bottom: 8px;
                }

                .course-title {
                    margin: 0;
                    font-size: 18px;
                    font-weight: 600;
                }

                .course-title a {
                    color: #1f2937;
                    text-decoration: none;
                }

                .course-title a:hover {
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

                .course-provider {
                    color: #6b7280;
                    font-size: 14px;
                    margin-bottom: 12px;
                }

                .course-meta {
                    display: flex;
                    flex-wrap: wrap;
                    gap: 16px;
                    margin-bottom: 12px;
                    color: #6b7280;
                    font-size: 14px;
                }

                .course-tags {
                    display: flex;
                    flex-wrap: wrap;
                    gap: 8px;
                    margin-bottom: 12px;
                }

                .tag {
                    background: #e5e7eb;
                    color: #4b5563;
                    padding: 4px 10px;
                    border-radius: 4px;
                    font-size: 12px;
                }

                .course-description {
                    color: #6b7280;
                    font-size: 14px;
                    line-height: 1.5;
                    margin-bottom: 16px;
                }

                .course-card-footer {
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                    padding-top: 16px;
                    border-top: 1px solid #e5e7eb;
                }

                .course-fees {
                    display: flex;
                    flex-direction: column;
                }

                .fees-label {
                    font-size: 12px;
                    color: #9ca3af;
                }

                .fees-amount {
                    font-size: 20px;
                    font-weight: 700;
                    color: #059669;
                }

                .details-btn {
                    background: #4f46e5;
                    color: white;
                    padding: 10px 20px;
                    border-radius: 6px;
                    text-decoration: none;
                    font-weight: 500;
                    transition: background 0.2s;
                }

                .details-btn:hover {
                    background: #4338ca;
                }

                .course-card-skeleton {
                    background: white;
                    border-radius: 12px;
                    padding: 20px;
                    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
                }

                .skeleton-title, .skeleton-provider, .skeleton-tag, .skeleton-price, .skeleton-btn {
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

                .skeleton-provider {
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

                .skeleton-price {
                    height: 20px;
                    width: 60px;
                }

                .skeleton-btn {
                    height: 36px;
                    width: 100px;
                }

                @keyframes skeleton-loading {
                    0% { background-position: 200% 0; }
                    100% { background-position: -200% 0; }
                }

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
