// Government Jobs Page
import React, { useState, useEffect, useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';

// Fixed 2-level relative paths to components
import JobCard from '../../components/JobCard';
import FilterBar from '../../components/FilterBar';
import Pagination from '../../components/Pagination';
import { apiUrl } from '../../apiConfig';
import './GovJobs.css';

const GovJobsPage = () => {
  const [jobs, setJobs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 20,
    total: 0,
    totalPages: 1,
  });

  const [searchParams, setSearchParams] = useSearchParams();

  // Extract filter parameters safely from URL
  const search = searchParams.get('search') || '';
  const department = searchParams.get('department') || '';
  const location = searchParams.get('location') || '';
  const source = searchParams.get('source') || '';
  const category = searchParams.get('category') || '';
  const region = searchParams.get('region') || '';
  const language = searchParams.get('language') || '';
  const page = parseInt(searchParams.get('page'), 10) || 1;

  // Fetch jobs with AbortController for clean lifecycle & race-condition prevention
  const fetchJobs = useCallback(async (signal) => {
    try {
      setLoading(true);
      setError(null);

      const params = new URLSearchParams({
        page: page.toString(),
        limit: '20',
      });

if (search) params.append('search', search);
      if (department) params.append('department', department);
      if (location) params.append('location', location);
      if (source) params.append('source', source);
      if (category) params.append('category', category);
      if (region) params.append('region', region);
      if (language) params.append('language', language);

      const endpoint = `${apiUrl('/api/v1/gov-jobs')}?${params.toString()}`;

      const response = await fetch(endpoint, { signal });

      if (!response.ok) {
        throw new Error(`Server responded with status ${response.status}`);
      }

      const data = await response.json();

      if (data.status === 'success' || data.success) {
        const fetchedJobs = data.data || data.jobs || [];
        setJobs(fetchedJobs);

        setPagination({
          page: data.pagination?.page || page,
          limit: data.pagination?.limit || 20,
          total: data.pagination?.total || fetchedJobs.length,
          totalPages:
            data.pagination?.totalPages ||
            Math.ceil((data.pagination?.total || fetchedJobs.length) / 20) ||
            1,
        });
      } else {
        setError(data.error?.message || data.message || 'Failed to fetch government jobs');
      }
    } catch (err) {
      if (err.name !== 'AbortError') {
        setError(err.message || 'Unable to load government jobs. Please check your backend connection.');
      }
    } finally {
      setLoading(false);
    }
}, [page, search, department, location, source, category, region, language]);

  useEffect(() => {
    const controller = new AbortController();
    fetchJobs(controller.signal);

    return () => controller.abort();
  }, [fetchJobs]);

  const handleFilterChange = (filters) => {
    const params = new URLSearchParams();

    if (filters.search) params.append('search', filters.search);
    if (filters.department) params.append('department', filters.department);
    if (filters.location) params.append('location', filters.location);
    if (filters.source) params.append('source', filters.source);
if (filters.category) params.append('category', filters.category);
    if (filters.region) params.append('region', filters.region);
    if (filters.language) params.append('language', filters.language);

    params.append('page', '1'); // Reset to page 1 on filter update
    setSearchParams(params);
  };

  const handlePageChange = (newPage) => {
    const params = new URLSearchParams(searchParams);
    params.set('page', newPage.toString());
    setSearchParams(params);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleResetFilters = () => {
    setSearchParams({});
  };

const activeFilterCount = [search, department, location, source, category, region, language].filter(Boolean).length;

  return (
    <div className="gov-jobs-page">
      <div className="page-header">
        <div className="header-title-section">
          <h1>Government Jobs</h1>
          <p>Verified notifications from UPSC, SSC, Railways, State PSCs & Central Agencies</p>
        </div>
        {!loading && !error && (
          <div className="stats-badge">
            <span>{pagination.total}</span> Jobs Available
          </div>
        )}
      </div>

<FilterBar
        filters={{ search, department, location, source, category, region, language }}
        onFilterChange={handleFilterChange}
        filterOptions={{
          departments: ['UPSC', 'SSC', 'Railway', 'State PSC', 'Banking', 'Teaching', 'Defense', 'Medical'],
          locations: ['All India', 'Delhi', 'Mumbai', 'Bangalore', 'Chennai', 'Kolkata', 'Hyderabad', 'Pune', 'Uttar Pradesh', 'Bihar'],
          sources: ['ncs', 'ssc', 'upsc', 'rrb', 'employment_news'],
          categories: ['10th/12th Pass', 'Graduate', 'Post Graduate', 'Diploma', 'Engineering'],
          regions: ['India', 'Overseas', 'Global Remote'],
          languages: ['English', 'Hindi']
        }}
      />

      {activeFilterCount > 0 && (
        <div className="active-filters-bar">
          <span>Showing results for active filters ({activeFilterCount})</span>
          <button className="reset-filter-btn" onClick={handleResetFilters}>
            Clear All Filters
          </button>
        </div>
      )}

      {loading ? (
        <div className="loading-container">
          <div className="spinner"></div>
          <p>Fetching latest government notifications...</p>
        </div>
      ) : error ? (
        <div className="error-container">
          <p className="error-message">{error}</p>
          <button className="retry-btn" onClick={() => fetchJobs()}>
            Try Again
          </button>
        </div>
      ) : (
        <>
          {jobs.length > 0 ? (
            <div className="jobs-grid">
              {jobs.map((job, index) => (
                <JobCard key={job.id || job._id || index} job={job} type="government" />
              ))}
            </div>
          ) : (
            <div className="no-results-card">
              <h3>No Government Jobs Found</h3>
              <p>Try adjusting your search terms, department, or location filters.</p>
              {activeFilterCount > 0 && (
                <button className="reset-filter-btn" onClick={handleResetFilters}>
                  Reset All Filters
                </button>
              )}
            </div>
          )}

          {pagination.totalPages > 1 && (
            <Pagination
              currentPage={pagination.page}
              totalPages={pagination.totalPages}
              onPageChange={handlePageChange}
            />
          )}
        </>
      )}
    </div>
  );
};

export default GovJobsPage;
