// Government Jobs Page
import React, { useState, useEffect } from 'react';
import { useParams, useSearchParams } from 'react-router-dom';
import JobCard from '../components/JobCard';
import FilterBar from '../components/FilterBar';
import Pagination from '../components/Pagination';
import './GovJobs.css';

const GovJobsPage = () => {
  const [jobs, setJobs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [pagination, setPagination] = useState({ page: 1, limit: 20, total: 0, totalPages: 0 });
  const [searchParams, setSearchParams] = useSearchParams();

  const page = parseInt(searchParams.get('page')) || 1;
  const department = searchParams.get('department') || '';
  const location = searchParams.get('location') || '';
  const source = searchParams.get('source') || '';

  useEffect(() => {
    fetchJobs();
  }, [page, department, location, source]);

  const fetchJobs = async () => {
    try {
      setLoading(true);
      const params = new URLSearchParams({
        page: page.toString(),
        limit: '20',
      });
      if (department) params.append('department', department);
      if (location) params.append('location', location);
      if (source) params.append('source', source);

      const response = await fetch(`${process.env.REACT_APP_BACKEND_URL}/api/v1/gov-jobs?${params}`);
      const data = await response.json();

      if (data.status === 'success') {
        setJobs(data.data);
        setPagination(data.pagination);
      } else {
        setError(data.error?.message || 'Failed to fetch jobs');
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleFilterChange = (filters) => {
    const params = new URLSearchParams();
    if (filters.department) params.append('department', filters.department);
    if (filters.location) params.append('location', filters.location);
    if (filters.source) params.append('source', filters.source);
    params.append('page', '1');
    setSearchParams(params);
  };

  const handlePageChange = (newPage) => {
    const params = new URLSearchParams(searchParams);
    params.set('page', newPage.toString());
    setSearchParams(params);
  };

  return (
    <div className="gov-jobs-page">
      <div className="page-header">
        <h1>Government Jobs</h1>
        <p>Find the latest government job openings in India</p>
      </div>

      <FilterBar
        filters={{ department, location, source }}
        onFilterChange={handleFilterChange}
        filterOptions={{
          departments: ['UPSC', 'SSC', 'Railway', 'State PSC', 'Banking', 'Teaching', 'Medical'],
          locations: ['Delhi', 'Mumbai', 'Bangalore', 'Chennai', 'Kolkata', 'Hyderabad', 'Pune', 'All India'],
          sources: ['ncs', 'ssc', 'upsc', 'rrb', 'employment_news']
        }}
      />

      {loading ? (
        <div className="loading">Loading jobs...</div>
      ) : error ? (
        <div className="error">{error}</div>
      ) : (
        <>
          <div className="jobs-grid">
            {jobs.map((job) => (
              <JobCard key={job.id} job={job} type="government" />
            ))}
          </div>

          {jobs.length === 0 && (
            <div className="no-results">No government jobs found</div>
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

export default GovJobsPage;

