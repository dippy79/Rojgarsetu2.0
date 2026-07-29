// Videos Page
import React, { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import VideoCard from '../../components/VideoCard';
import FilterBar from '../../components/FilterBar';
import Pagination from '../../components/Pagination';
import './Videos.css';

const VideosPage = () => {
  const [videos, setVideos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [pagination, setPagination] = useState({ page: 1, limit: 20, total: 0, totalPages: 0 });
  const [searchParams, setSearchParams] = useSearchParams();

  const page = parseInt(searchParams.get('page')) || 1;
  const channel = searchParams.get('channel') || '';
  const category = searchParams.get('category') || '';

  useEffect(() => {
    fetchVideos();
  }, [page, channel, category]);

  const fetchVideos = async () => {
    try {
      setLoading(true);
      const params = new URLSearchParams({
        page: page.toString(),
        limit: '20',
      });
      if (channel) params.append('channel', channel);
      if (category) params.append('category', category);

      const response = await fetch(`${process.env.REACT_APP_BACKEND_URL}/api/v1/videos?${params}`);
      const data = await response.json();

      if (data.status === 'success') {
        setVideos(data.data);
        setPagination(data.pagination);
      } else {
        setError(data.error?.message || 'Failed to fetch videos');
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleFilterChange = (filters) => {
    const params = new URLSearchParams();
    if (filters.channel) params.append('channel', filters.channel);
    if (filters.category) params.append('category', filters.category);
    params.append('page', '1');
    setSearchParams(params);
  };

  const handlePageChange = (newPage) => {
    const params = new URLSearchParams(searchParams);
    params.set('page', newPage.toString());
    setSearchParams(params);
  };

  return (
    <div className="videos-page">
      <div className="page-header">
        <h1>Videos</h1>
        <p>Watch educational and job-related videos from official channels</p>
      </div>

      <FilterBar
        filters={{ channel, category }}
        onFilterChange={handleFilterChange}
        filterOptions={{
          channels: ['Naukri', 'LinkedIn', 'Study IQ', 'Unacademy', 'Byju\'s', 'Skill India'],
          categories: ['Jobs', 'Government', 'Education', 'Skills', 'UPSC', 'SSC', 'Railways']
        }}
      />

      {loading ? (
        <div className="loading">Loading videos...</div>
      ) : error ? (
        <div className="error">{error}</div>
      ) : (
        <>
          <div className="videos-grid">
            {videos.map((video) => (
              <VideoCard key={video.id} video={video} />
            ))}
          </div>

          {videos.length === 0 && (
            <div className="no-results">No videos found</div>
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

export default VideosPage;

