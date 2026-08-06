// Courses Page
import React, { useState, useEffect, useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';
import CourseCard from '../../components/CourseCard';
import FilterBar from '../../components/FilterBar';
import Pagination from '../../components/Pagination';
import './Courses.css';

const CoursesPage = () => {
  const [courses, setCourses] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [pagination, setPagination] = useState({ page: 1, limit: 20, total: 0, totalPages: 0 });
  const [searchParams, setSearchParams] = useSearchParams();

const page = parseInt(searchParams.get('page')) || 1;
  const provider = searchParams.get('provider') || '';
  const mode = searchParams.get('mode') || '';
  const level = searchParams.get('level') || '';
  const language = searchParams.get('language') || '';

const fetchCourses = useCallback(async () => {
    try {
      setLoading(true);
      const params = new URLSearchParams({
        page: page.toString(),
        limit: '20',
      });
      if (provider) params.append('provider', provider);
      if (mode) params.append('mode', mode);
      if (level) params.append('level', level);
      if (language) params.append('language', language);

      const response = await fetch(`${process.env.REACT_APP_BACKEND_URL}/api/v1/courses?${params}`);
      const data = await response.json();

      if (data.status === 'success') {
        setCourses(data.data);
        setPagination(data.pagination);
      } else {
        setError(data.error?.message || 'Failed to fetch courses');
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, [page, provider, mode, level, language]);

  useEffect(() => {
    fetchCourses();
  }, [fetchCourses]);

  const handleFilterChange = (filters) => {
    const params = new URLSearchParams();
    if (filters.provider) params.append('provider', filters.provider);
    if (filters.mode) params.append('mode', filters.mode);
    if (filters.level) params.append('level', filters.level);
    if (filters.language) params.append('language', filters.language);
    params.append('page', '1');
    setSearchParams(params);
  };

  const handlePageChange = (newPage) => {
    const params = new URLSearchParams(searchParams);
    params.set('page', newPage.toString());
    setSearchParams(params);
  };

  return (
    <div className="courses-page">
      <div className="page-header">
        <h1>Courses</h1>
        <p>Upskill with free and paid courses from top providers</p>
      </div>

<FilterBar
        filters={{ provider, mode, level, language }}
        onFilterChange={handleFilterChange}
        filterOptions={{
          providers: ['NPTEL', 'SWAYAM', 'NSDC', 'Coursera', 'Udemy', 'edX'],
          modes: ['online', 'offline', 'hybrid'],
          levels: ['beginner', 'intermediate', 'advanced'],
          languages: [
            { code: 'en', label: 'English' },
            { code: 'hi', label: 'Hindi' },
            { code: 'ta', label: 'Tamil' },
            { code: 'te', label: 'Telugu' },
            { code: 'bn', label: 'Bengali' },
            { code: 'mr', label: 'Marathi' },
            { code: 'gu', label: 'Gujarati' },
            { code: 'kn', label: 'Kannada' },
            { code: 'ml', label: 'Malayalam' },
            { code: 'pa', label: 'Punjabi' }
          ]
        }}
      />

      {loading ? (
        <div className="loading">Loading courses...</div>
      ) : error ? (
        <div className="error">{error}</div>
      ) : (
        <>
          <div className="courses-grid">
            {courses.map((course) => (
              <CourseCard key={course.id} course={course} />
            ))}
          </div>

          {courses.length === 0 && (
            <div className="no-results">No courses found</div>
          )}

          <Pagination
            currentPage={pagination.page}
            totalPages={pagination.totalPages}
            onPageChange={handlePageChange}
          />
        </>
      )}
    </div>
  );
};

export default CoursesPage;

