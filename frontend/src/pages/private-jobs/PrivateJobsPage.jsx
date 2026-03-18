// Private Jobs Page
import React, { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import JobCard from '../components/JobCard';
import FilterBar from '../components/FilterBar';
import Pagination from '../components/Pagination';
import './PrivateJobs.css';

const PrivateJobsPage = () => {
  const [jobs, setJobs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [pagination, setPagination] = useState({ page: 1, limit: 20, total: 0, totalPages: 0 });
  const [searchParams, setSearchParams] = useSearchParams();

  const page = parseInt(searchParams.get('page')) || 1;
  const company = searchParams.get('company') || '';
  const location = searchParams.get('location') || '';
  const jobType = searchParams.get('job_type') || '';
  const source = searchParams.get('source') || '';

  useEffect(() => {
    fetchJobs();
  }, [page, company, location, jobType, source]);

  const fetchJobs = async () => {
    try {
      setLoading(true);
      const params = new URLSearchParams({
        page: page.toString(),
        limit: '20',
      });
      if (company) params.append('company', company);
      if (location) params.append('location', location);
      if (jobType) params.append('job_type', jobType);
      if (source) params.append('source', source);

      const response = await fetch(`${process.env.REACT_APP_BACKEND_URL}/api/v1/private-jobs?${params}`);
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
    if (filters.company) params.append('company', filters.company);
    if (filters.location) params.append('location', filters.location);
    if (filters.jobType) params.append('job_type', filters.jobType);
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
    <div className="private-jobs-page">
      <div className="page-header">
        <h1>Private Jobs</h1>
        <p>Find the best private sector job opportunities</p>
      </div>

      <FilterBar
        filters={{ company, location, jobType, source }}
        onFilterChange={handleFilterChange}
        filterOptions={{
          companies: ['TCS', 'Infosys', 'Wipro', 'Amazon', 'Google', 'Microsoft', 'Flipkart'],
          locations: ['Delhi', 'Mumbai', 'Bangalore', 'Chennai', 'Kolkata', 'Hyderabad', 'Pune', 'Remote'],
          jobTypes: ['full-time', 'part-time', 'contract', 'internship', 'remote'],
          sources: ['linkedin', 'indeed', 'google_jobs', 'company_pages']
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
              <JobCard key={job.id} job={job} type="private" />
            ))}
          </div>

          {jobs.length === 0 && (
            <div className="no-results">No private jobs found</div>
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

export default PrivateJobsPage;

