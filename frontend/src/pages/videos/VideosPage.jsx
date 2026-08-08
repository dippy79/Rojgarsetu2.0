// Videos Page — Tech/Career Prep default view, Government/News in separate tab.
// Channels/categories are now loaded from dedicated backend endpoints
// (/api/v1/videos/channels, /api/v1/videos/categories) instead of hardcoded
// lists, and the Tech tab applies server-side exclusion (?exclude=Government)
// so both the grid and pagination totals are correct.
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

// Category excluded from the default Tech view (matches the backend's
// ?exclude=Government handling).
const GOV_CATEGORY = 'Government';

const VideosPage = () => {
  const [videos, setVideos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [allChannels, setAllChannels] = useState([]);
  const [allCategories, setAllCategories] = useState([]);
  const [pagination, setPagination] = useState({ page: 1, limit: 20, total: 0, totalPages: 0 });
  const [searchParams, setSearchParams] = useSearchParams();

  const page = parseInt(searchParams.get('page')) || 1;
  const channel = searchParams.get('channel') || '';
  const category = searchParams.get('category') || '';
  // Active tab from URL (defaults to TECH so gov/news are hidden by default).
  const activeTab = searchParams.get('tab') === VIDEO_TABS.NEWS ? VIDEO_TABS.NEWS : VIDEO_TABS.TECH;

  // Load distinct channels + categories for the filter dropdowns.
  const fetchFilterOptions = useCallback(async () => {
    try {
      const [chResp, catResp] = await Promise.all([
        fetch(apiUrl('/api/v1/videos/channels')),
        fetch(apiUrl('/api/v1/videos/categories')),
      ]);
      const chData = await chResp.json();
      const catData = await catResp.json();
      if (chData.status === 'success' && Array.isArray(chData.data)) {
        setAllChannels(chData.data.map((c) => c.channel).filter(Boolean));
      }
      if (catData.status === 'success' && Array.isArray(catData.data)) {
        setAllCategories(catData.data.map((c) => c.category).filter(Boolean));
      }
    } catch (err) {
      console.warn('Failed to load video filter options', err);
    }
  }, []);

  useEffect(() => {
    fetchFilterOptions();
  }, [fetchFilterOptions]);

  const fetchVideos = useCallback(async () => {
    try {
      setLoading(true);
      const params = new URLSearchParams({
        page: page.toString(),
        limit: '20',
      });
      if (channel) params.append('channel', channel);
      if (category) params.append('category', category);
      // Tech tab: exclude Government content server-side so the video list AND
      // the pagination total both reflect the post-exclusion set.
      if (activeTab === VIDEO_TABS.TECH) params.append('exclude', GOV_CATEGORY);

      const response = await fetch(`${apiUrl('/api/v1/videos')}?${params}`);
      const data = await response.json();

      if (data.status === 'success') {
        setVideos(data.data || []);
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
          categories: allCategories,
          channels: allChannels
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
