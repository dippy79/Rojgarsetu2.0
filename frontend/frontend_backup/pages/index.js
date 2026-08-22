// frontend/pages/index.js - Home page
import Link from 'next/link';
import useSWR from 'swr';
import { jobsAPI, coursesAPI } from '../lib/api';

export default function Home() {
    // Fetch featured jobs
    const { data: jobsData } = useSWR(
        '/api/jobs?limit=5&is_featured=true',
        () => jobsAPI.getJobs({ limit: 5, is_featured: true }).then(res => res.data),
        { revalidateOnFocus: false }
    );

    // Fetch featured courses
    const { data: coursesData } = useSWR(
        '/api/courses?limit=5&is_featured=true',
        () => coursesAPI.getCourses({ limit: 5, is_featured: true }).then(res => res.data),
        { revalidateOnFocus: false }
    );

    const jobs = jobsData?.jobs || [];
    const courses = coursesData?.courses || [];

    return (
        <div className="home-page">
            {/* Hero Section */}
            <section className="hero">
                <div className="hero-content">
                    <h1>Welcome to RojgarSetu</h1>
                    <p className="hero-subtitle">Your Gateway to Government Jobs & Professional Courses</p>
                    <div className="hero-buttons">
                        <Link href="/jobs" className="btn btn-primary">Browse Jobs</Link>
                        <Link href="/courses" className="btn btn-secondary">Explore Courses</Link>
                    </div>
                </div>
            </section>

            {/* Stats Section */}
            <section className="stats">
                <div className="stat-item">
                    <span className="stat-number">500+</span>
                    <span className="stat-label">Government Jobs</span>
                </div>
                <div className="stat-item">
                    <span className="stat-number">200+</span>
                    <span className="stat-label">Courses</span>
                </div>
                <div className="stat-item">
                    <span className="stat-number">50+</span>
                    <span className="stat-label">Categories</span>
                </div>
                <div className="stat-item">
                    <span className="stat-number">10K+</span>
                    <span className="stat-label">Candidates</span>
                </div>
            </section>

            {/* Jobs Section */}
            <section className="section">
                <div className="section-header">
                    <h2>Latest Government Jobs</h2>
                    <Link href="/jobs" className="view-all">View All Jobs →</Link>
                </div>
                <div className="jobs-grid">
                    {jobs.length > 0 ? (
                        jobs.map(job => (
                            <div key={job.id} className="job-card">
                                <h3>
                                    <Link href={'/jobs/' + job.id}>{job.title}</Link>
                                </h3>
                                <p className="job-company">{job.company_name || job.source || 'RojgarSetu'}</p>
                                <div className="job-meta">
                                    <span>{job.location || 'All India'}</span>
                                    <span>{job.type === 'full-time' ? 'Full Time' : job.type}</span>
                                </div>
                                <Link href={'/jobs/' + job.id} className="apply-link">View Details →</Link>
                            </div>
                        ))
                    ) : (
                        <div className="no-data">
                            <p>No featured jobs available at the moment.</p>
                        </div>
                    )}
                </div>
            </section>

            {/* Courses Section */}
            <section className="section">
                <div className="section-header">
                    <h2>Popular Courses</h2>
                    <Link href="/courses" className="view-all">View All Courses →</Link>
                </div>
                <div className="courses-grid">
                    {courses.length > 0 ? (
                        courses.map(course => (
                            <div key={course.id} className="course-card">
                                <h3>
                                    <Link href={'/courses/' + course.id}>{course.name}</Link>
                                </h3>
                                <p className="course-provider">{course.provider || 'RojgarSetu'}</p>
                                <div className="course-meta">
                                    <span>{course.duration || 'Self-paced'}</span>
                                    <span>{course.level || 'Beginner'}</span>
                                </div>
                                <div className="course-fees">
                                    {course.fees_structure === '0' || !course.fees_structure ? 'Free' : '₹' + course.fees_structure}
                                </div>
                                <Link href={'/courses/' + course.id} className="view-link">View Details →</Link>
                            </div>
                        ))
                    ) : (
                        <div className="no-data">
                            <p>No featured courses available at the moment.</p>
                        </div>
                    )}
                </div>
            </section>

            {/* CTA Section */}
            <section className="cta-section">
                <h2>Ready to Start Your Career?</h2>
                <p>Join thousands of job seekers who found their dream job through RojgarSetu</p>
                <div className="cta-buttons">
                    <Link href="/jobs" className="btn btn-primary">Find Jobs Now</Link>
                    <Link href="/register" className="btn btn-outline">Create Account</Link>
                </div>
            </section>

            <style jsx>{`
                .home-page {
                    min-height: 100vh;
                    background: #f8fafc;
                }

                .hero {
                    background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
                    color: white;
                    padding: 80px 20px;
                    text-align: center;
                }

                .hero-content {
                    max-width: 800px;
                    margin: 0 auto;
                }

                .hero h1 {
                    font-size: 48px;
                    font-weight: 800;
                    margin: 0 0 16px 0;
                }

                .hero-subtitle {
                    font-size: 20px;
                    opacity: 0.9;
                    margin: 0 0 32px 0;
                }

                .hero-buttons {
                    display: flex;
                    gap: 16px;
                    justify-content: center;
                }

                .btn {
                    padding: 14px 32px;
                    border-radius: 8px;
                    font-weight: 600;
                    text-decoration: none;
                    transition: all 0.2s;
                }

                .btn-primary {
                    background: white;
                    color: #4f46e5;
                }

                .btn-primary:hover {
                    transform: translateY(-2px);
                    box-shadow: 0 4px 12px rgba(0,0,0,0.2);
                }

                .btn-secondary {
                    background: transparent;
                    color: white;
                    border: 2px solid white;
                }

                .btn-secondary:hover {
                    background: rgba(255,255,255,0.1);
                }

                .btn-outline {
                    background: transparent;
                    color: #4f46e5;
                    border: 2px solid #4f46e5;
                }

                .stats {
                    display: grid;
                    grid-template-columns: repeat(4, 1fr);
                    gap: 20px;
                    max-width: 1000px;
                    margin: -40px auto 60px;
                    padding: 0 20px;
                }

                .stat-item {
                    background: white;
                    padding: 24px;
                    border-radius: 12px;
                    text-align: center;
                    box-shadow: 0 4px 6px rgba(0,0,0,0.1);
                }

                .stat-number {
                    display: block;
                    font-size: 32px;
                    font-weight: 700;
                    color: #4f46e5;
                }

                .stat-label {
                    color: #6b7280;
                    font-size: 14px;
                }

                .section {
                    max-width: 1200px;
                    margin: 0 auto;
                    padding: 40px 20px;
                }

                .section-header {
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                    margin-bottom: 24px;
                }

                .section-header h2 {
                    font-size: 28px;
                    font-weight: 700;
                    color: #1f2937;
                    margin: 0;
                }

                .view-all {
                    color: #4f46e5;
                    text-decoration: none;
                    font-weight: 500;
                }

                .view-all:hover {
                    text-decoration: underline;
                }

                .jobs-grid, .courses-grid {
                    display: grid;
                    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
                    gap: 20px;
                }

                .job-card, .course-card {
                    background: white;
                    border-radius: 12px;
                    padding: 20px;
                    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
                    transition: all 0.2s;
                }

                .job-card:hover, .course-card:hover {
                    box-shadow: 0 4px 12px rgba(0,0,0,0.15);
                    transform: translateY(-2px);
                }

                .job-card h3, .course-card h3 {
                    margin: 0 0 8px 0;
                    font-size: 18px;
                }

                .job-card h3 a, .course-card h3 a {
                    color: #1f2937;
                    text-decoration: none;
                }

                .job-card h3 a:hover, .course-card h3 a:hover {
                    color: #4f46e5;
                }

                .job-company, .course-provider {
                    color: #6b7280;
                    font-size: 14px;
                    margin-bottom: 12px;
                }

                .job-meta, .course-meta {
                    display: flex;
                    gap: 16px;
                    color: #6b7280;
                    font-size: 14px;
                    margin-bottom: 16px;
                }

                .course-fees {
                    font-size: 20px;
                    font-weight: 700;
                    color: #059669;
                    margin-bottom: 16px;
                }

                .apply-link, .view-link {
                    color: #4f46e5;
                    text-decoration: none;
                    font-weight: 500;
                }

                .no-data {
                    grid-column: 1 / -1;
                    text-align: center;
                    padding: 40px;
                    color: #6b7280;
                }

                .cta-section {
                    background: #1f2937;
                    color: white;
                    padding: 80px 20px;
                    text-align: center;
                }

                .cta-section h2 {
                    font-size: 32px;
                    margin: 0 0 16px 0;
                }

                .cta-section p {
                    font-size: 18px;
                    opacity: 0.9;
                    margin: 0 0 32px 0;
                }

                .cta-buttons {
                    display: flex;
                    gap: 16px;
                    justify-content: center;
                }

                @media (max-width: 768px) {
                    .hero h1 {
                        font-size: 32px;
                    }

                    .stats {
                        grid-template-columns: repeat(2, 1fr);
                    }

                    .hero-buttons, .cta-buttons {
                        flex-direction: column;
                    }
                }
            `}</style>
        </div>
    );
}
