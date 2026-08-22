// Course Card Component
import React from 'react';
// import './CourseCard.css';

const CourseCard = ({ course }) => {
  return (
    <div className="course-card">
      <div className="course-thumbnail">
        {course.thumbnail_url ? (
          <img src={course.thumbnail_url} alt={course.title} />
        ) : (
          <div className="thumbnail-placeholder">📚</div>
        )}
        {course.is_free && <span className="free-badge">Free</span>}
      </div>

      <div className="course-content">
        <h3 className="course-title">{course.title}</h3>
        
        <p className="course-provider">🏛️ {course.provider}</p>

        <div className="course-meta">
          {course.duration && <span className="meta-item">⏱️ {course.duration}</span>}
          {course.mode && <span className="meta-item">💻 {course.mode}</span>}
          {course.level && <span className="meta-item">📶 {course.level}</span>}
        </div>

        {course.rating && (
          <div className="course-rating">
            ⭐ {course.rating}
          </div>
        )}

        <div className="course-actions">
          <a 
            href={course.url} 
            target="_blank" 
            rel="noopener noreferrer" 
            className="view-btn"
          >
            View Course
          </a>
        </div>
      </div>
    </div>
  );
};

export default CourseCard;

