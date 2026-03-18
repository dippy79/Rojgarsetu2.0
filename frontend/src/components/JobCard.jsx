// Job Card Component
import React from 'react';
import './JobCard.css';

const JobCard = ({ job, type }) => {
  const isGov = type === 'government';

  return (
    <div className={`job-card ${isGov ? 'gov-job' : 'priv-job'}`}>
      <div className="job-header">
        <h3 className="job-title">{job.title}</h3>
        <span className="job-source">{job.source}</span>
      </div>

      <div className="job-details">
        {isGov ? (
          <>
            {job.department && <p className="job-department">📋 {job.department}</p>}
          </>
        ) : (
          <>
            <p className="job-company">🏢 {job.company}</p>
            {job.salary && <p className="job-salary">💰 {job.salary}</p>}
            {job.experience && <p className="job-experience">📊 {job.experience}</p>}
          </>
        )}

        {job.location && <p className="job-location">📍 {job.location}</p>}
        
        {isGov && job.last_date && (
          <p className="job-last-date">📅 Last Date: {job.last_date}</p>
        )}
        
        {job.vacancyCount && (
          <p className="job-vacancies">🎯 Vacancies: {job.vacancyCount}</p>
        )}
      </div>

      <div className="job-actions">
        <a 
          href={isGov ? job.apply_url : job.url} 
          target="_blank" 
          rel="noopener noreferrer" 
          className="apply-btn"
        >
          Apply Now
        </a>
      </div>

      <div className="job-footer">
        <span className="job-created">
          Posted: {new Date(job.created_at).toLocaleDateString()}
        </span>
      </div>
    </div>
  );
};

export default JobCard;

