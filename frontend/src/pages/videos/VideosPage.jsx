// Videos Page — Tech/Career Prep default view, Government/News in separate tab (Phase 2 fix)
import React, { useState, useEffect, useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';
import VideoCard from '../../components/VideoCard';
import FilterBar from '../../components/FilterBar';
import Pagination from '../../components/Pagination';
import { apiUrl } from '../../apiConfig';
import './Videos.css';

// Tabs controlling the video feed segmentation.
export const VIDEO_TABS = {
  TECH: 'tech',
  NEWS: 'news',
};

// Channels/categories flagged as Government / News content (excluded from the
// default Tech view).
// TODO: Move this explicit filter to the Backend API (a `?category=Tech` filter
//       and/or a `?exclude=Government` query parameter) so it scales cleanly.
const GOV_CHANNELS = ['PIB India', 'DD News', 'DD India', 'PIB'];
const GOV_CATEGORY = 'Government';

const isGovVideo = (video) => {
  const cat = (video.category || '').toLowerCase();
  const channel = (video.channel || '');
  if (cat === GOV_CATEGORY.toLowerCase()) return true;
  return GOV_CHANNELS.some((g) => channel.toLowerCase().includes(g.toLowerCase()));
};

const VideosPage = () => {
  const [videos, setVideos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [pagination, setPagination] = useState({ page: 1, limit: 20, total: 0, totalPages: 0 });
  const [searchParams, setSearchParams] = useSearchParams();

  const page = parseInt(searchParams.get('page')) || 1;
  const channel = searchParams.get('channel') || '';
  const category = searchParams.get('category') || '';
  // Active tab from URL (defaults to TECH so gov/news are hidden by default).
  const activeTab = searchParams.get('tab') === VIDEO_TABS.NEWS ? VIDEO_TABS.NEWS : VIDEO_TABS.TECH;

  const fetchVideos = useCallback(async () => {
    try {
      setLoading(true);
      const params = new URLSearchParams({
        page: page.toString(),
        limit: '20',
      });
      if (channel) params.append('channel', channel);
      // When on the Tech tab, only request non-Government categories from the API
      // (the robust filter is still applied client-side below as a safeguard).
      const effectiveCategory =
        activeTab === VIDEO_TABS.TECH && category.toLowerCase() === GOV_CATEGORY.toLowerCase()
          ? ''
          : category;
      if (effectiveCategory) params.append('category', effectiveCategory);

      const response = await fetch(`${apiUrl('/api/v1/videos')}?${params}`);
      const data = await response.json();

      if (data.status === 'success') {
        const raw = data.data || [];
        // Phase 2 fix: segment the feed. Tech tab excludes Government/news videos.
        const filtered =
          activeTab === VIDEO_TABS.TECH ? raw.filter((v) => !isGovVideo(v)) : raw;
        setVideos(filtered);
        setPagination(data.pagination);
      } else {
        setError(data.error?.message || 'Failed to fetch videos');
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, [page, channel, category, activeTab]);

  useEffect(() => {
    fetchVideos();
  }, [fetchVideos]);

  const switchTab = (tab) => {
    const params = new URLSearchParams(searchParams);
    params.set('tab', tab);
    params.set('page', '1');
    setSearchParams(params);
  };

  const handleFilterChange = (filters) => {
    const params = new URLSearchParams(searchParams);
    if (filters.channel) params.set('channel', filters.channel);
    else params.delete('channel');
    if (filters.category) params.set('category', filters.category);
    else params.delete('category');
    params.set('tab', activeTab);
    params.set('page', '1');
    setSearchParams(params);
  };

  const handlePageChange = (newPage) => {
    const params = new URLSearchParams(searchParams);
    params.set('page', newPage.toString());
    setSearchParams(params);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  return (
    <div className="videos-page">
      <div className="page-header">
        <h1>Videos</h1>
        <p>Career-prep and tech learning videos — plus official government &amp; news updates</p>
      </div>

      <div className="videos-tabs" role="tablist">
        <button
          role="tab"
          aria-selected={activeTab === VIDEO_TABS.TECH}
          className={`video-tab ${activeTab === VIDEO_TABS.TECH ? 'active' : ''}`}
          onClick={() => switchTab(VIDEO_TABS.TECH)}
        >
          🎯 Tech &amp; Career Prep
        </button>
        <button
          role="tab"
          aria-selected={activeTab === VIDEO_TABS.NEWS}
          className={`video-tab ${activeTab === VIDEO_TABS.NEWS ? 'active' : ''}`}
          onClick={() => switchTab(VIDEO_TABS.NEWS)}
        >
          📰 Government &amp; News
        </button>
      </div>

      <FilterBar
        filters={{ channel, category }}
        onFilterChange={handleFilterChange}
        filterOptions={{
          categories: activeTab === VIDEO_TABS.TECH
            ? ['Jobs', 'Education', 'Skills', 'Tech', 'Interviews']
            : ['Government', 'News', 'Jobs'],
          channels: []
        }}
      />

      {loading ? (
        <div className="loading">Loading videos...</div>
      ) : error ? (
        <div className="error">{error}</div>
      ) : (
        <>
          <div className="videos-grid">
            {(videos || []).map((video) => (
              <VideoCard key={video.id} video={video} />
            ))}
          </div>

          {videos.length === 0 && (
            <div className="no-results">No videos found in this tab</div>
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
