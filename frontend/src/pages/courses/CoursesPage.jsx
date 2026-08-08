// Courses Page — Dynamic Providers (Phase 1 fix)
import React, { useState, useEffect, useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';
import CourseCard from '../../components/CourseCard';
import FilterBar from '../../components/FilterBar';
import Pagination from '../../components/Pagination';
import { apiUrl } from '../../apiConfig';
import './Courses.css';

// Known provider order for consistent, sequential display.
const PROVIDER_ORDER = ['Coursera', 'TutorialsPoint', 'GeeksforGeeks', 'Udemy', 'W3Schools', 'NPTEL', 'SWAYAM', 'NSDC'];

const CoursesPage = () => {
  const [courses, setCourses] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [allProviders, setAllProviders] = useState([]);
  const [pagination, setPagination] = useState({ page: 1, limit: 20, total: 0, totalPages: 0 });
  const [searchParams, setSearchParams] = useSearchParams();

  const page = parseInt(searchParams.get('page')) || 1;
  const provider = searchParams.get('provider') || '';
  const mode = searchParams.get('mode') || '';
  const level = searchParams.get('level') || '';

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

      const response = await fetch(`${apiUrl('/api/v1/courses')}?${params}`);
      const data = await response.json();

      if (data.status === 'success') {
        setCourses(data.data || []);
        setPagination(data.pagination);
      } else {
        setError(data.error?.message || 'Failed to fetch courses');
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, [page, provider, mode, level]);

  const fetchProviders = useCallback(async () => {
    try {
      const response = await fetch(apiUrl('/api/v1/courses/providers'));
      const data = await response.json();
      if (data.status === 'success' && Array.isArray(data.data)) {
        const ordered = PROVIDER_ORDER.filter((p) => data.data.some((item) => item.provider === p));
        const extras = data.data
          .map((item) => item.provider)
          .filter((providerName) => providerName && !ordered.includes(providerName));
        setAllProviders([...ordered, ...extras]);
      }
    } catch (err) {
      console.warn('Failed to load providers', err);
    }
  }, []);

  useEffect(() => {
    fetchProviders();
  }, [fetchProviders]);

  useEffect(() => {
    fetchCourses();
  }, [fetchCourses]);

  const handleFilterChange = (filters) => {
    const params = new URLSearchParams();
    if (filters.provider) params.append('provider', filters.provider);
    if (filters.mode) params.append('mode', filters.mode);
    if (filters.level) params.append('level', filters.level);
    params.append('page', '1');
    setSearchParams(params);
  };

  const handlePageChange = (newPage) => {
    const params = new URLSearchParams(searchParams);
    params.set('page', newPage.toString());
    setSearchParams(params);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  return (
    <div className="courses-page">
      <div className="page-header">
        <h1>Courses</h1>
        <p>Upskill with free and paid courses from top providers</p>
      </div>

      {!loading && allProviders.length > 0 && (
        <div className="provider-strip" aria-label="Course providers">
          <span className="provider-strip-label">Providers:</span>
          {allProviders.map((p) => (
            <button
              key={p}
              className={`provider-chip ${provider.toLowerCase() === p.toLowerCase() ? 'active' : ''}`}
              onClick={() => handleFilterChange({ provider: provider.toLowerCase() === p.toLowerCase() ? '' : p, mode, level })}
            >
              {p}
            </button>
          ))}
        </div>
      )}

      <FilterBar
        filters={{ provider, mode, level }}
        onFilterChange={handleFilterChange}
        filterOptions={{
          providers: allProviders.length > 0 ? allProviders : PROVIDER_ORDER,
          modes: ['online', 'offline', 'hybrid'],
          levels: ['beginner', 'intermediate', 'advanced']
        }}
      />

      {loading ? (
        <div className="loading">Loading courses...</div>
      ) : error ? (
        <div className="error">{error}</div>
      ) : (
        <>
          <div className="courses-grid">
            {(courses || []).map((course) => (
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
