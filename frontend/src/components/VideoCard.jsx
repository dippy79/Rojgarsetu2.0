// Video Card Component
import React from 'react';
// import './VideoCard.css';

const VideoCard = ({ video }) => {
  return (
    <div className="video-card">
      <div className="video-thumbnail">
        {video.thumbnail ? (
          <img src={video.thumbnail} alt={video.title} />
        ) : (
          <div className="thumbnail-placeholder">🎬</div>
        )}
        {video.duration && <span className="duration-badge">{video.duration}</span>}
      </div>

      <div className="video-content">
        <h3 className="video-title">{video.title}</h3>
        
        <p className="video-channel">📺 {video.channel}</p>

        <div className="video-meta">
          {video.viewCount && (
            <span className="meta-item">👁️ {formatViewCount(video.viewCount)}</span>
          )}
          {video.published_at && (
            <span className="meta-item">📅 {new Date(video.published_at).toLocaleDateString()}</span>
          )}
        </div>

        <div className="video-actions">
          <a 
            href={video.url} 
            target="_blank" 
            rel="noopener noreferrer" 
            className="watch-btn"
          >
            Watch Video
          </a>
        </div>
      </div>
    </div>
  );
};

function formatViewCount(count) {
  if (count >= 1000000) {
    return (count / 1000000).toFixed(1) + 'M';
  }
  if (count >= 1000) {
    return (count / 1000).toFixed(1) + 'K';
  }
  return count.toString();
}

export default VideoCard;

