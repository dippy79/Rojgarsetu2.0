// FilterBar Component
import React from 'react';
import './FilterBar.css';

const FilterBar = ({ filters, onFilterChange, filterOptions }) => {
  const handleChange = (key, value) => {
    onFilterChange({ ...filters, [key]: value });
  };

  return (
    <div className="filter-bar">
      {filterOptions.departments && (
        <div className="filter-group">
          <label>Department</label>
          <select 
            value={filters.department || ''} 
            onChange={(e) => handleChange('department', e.target.value)}
          >
            <option value="">All Departments</option>
            {filterOptions.departments.map(dept => (
              <option key={dept} value={dept}>{dept}</option>
            ))}
          </select>
        </div>
      )}

      {filterOptions.companies && (
        <div className="filter-group">
          <label>Company</label>
          <select 
            value={filters.company || ''} 
            onChange={(e) => handleChange('company', e.target.value)}
          >
            <option value="">All Companies</option>
            {filterOptions.companies.map(company => (
              <option key={company} value={company}>{company}</option>
            ))}
          </select>
        </div>
      )}

      {filterOptions.locations && (
        <div className="filter-group">
          <label>Location</label>
          <select 
            value={filters.location || ''} 
            onChange={(e) => handleChange('location', e.target.value)}
          >
            <option value="">All Locations</option>
            {filterOptions.locations.map(loc => (
              <option key={loc} value={loc}>{loc}</option>
            ))}
          </select>
        </div>
      )}

      {filterOptions.providers && (
        <div className="filter-group">
          <label>Provider</label>
          <select 
            value={filters.provider || ''} 
            onChange={(e) => handleChange('provider', e.target.value)}
          >
            <option value="">All Providers</option>
            {filterOptions.providers.map(provider => (
              <option key={provider} value={provider}>{provider}</option>
            ))}
          </select>
        </div>
      )}

      {filterOptions.modes && (
        <div className="filter-group">
          <label>Mode</label>
          <select 
            value={filters.mode || ''} 
            onChange={(e) => handleChange('mode', e.target.value)}
          >
            <option value="">All Modes</option>
            {filterOptions.modes.map(mode => (
              <option key={mode} value={mode}>{mode}</option>
            ))}
          </select>
        </div>
      )}

      {filterOptions.levels && (
        <div className="filter-group">
          <label>Level</label>
          <select 
            value={filters.level || ''} 
            onChange={(e) => handleChange('level', e.target.value)}
          >
            <option value="">All Levels</option>
            {filterOptions.levels.map(level => (
              <option key={level} value={level}>{level}</option>
            ))}
          </select>
        </div>
      )}

      {filterOptions.channels && (
        <div className="filter-group">
          <label>Channel</label>
          <select 
            value={filters.channel || ''} 
            onChange={(e) => handleChange('channel', e.target.value)}
          >
            <option value="">All Channels</option>
            {filterOptions.channels.map(channel => (
              <option key={channel} value={channel}>{channel}</option>
            ))}
          </select>
        </div>
      )}

      {filterOptions.categories && (
        <div className="filter-group">
          <label>Category</label>
          <select 
            value={filters.category || ''} 
            onChange={(e) => handleChange('category', e.target.value)}
          >
            <option value="">All Categories</option>
            {filterOptions.categories.map(cat => (
              <option key={cat} value={cat}>{cat}</option>
            ))}
          </select>
        </div>
      )}

      {filterOptions.sources && (
        <div className="filter-group">
          <label>Source</label>
          <select 
            value={filters.source || ''} 
            onChange={(e) => handleChange('source', e.target.value)}
          >
            <option value="">All Sources</option>
            {filterOptions.sources.map(source => (
              <option key={source} value={source}>{source}</option>
            ))}
          </select>
        </div>
      )}

{filterOptions.jobTypes && (
        <div className="filter-group">
          <label>Job Type</label>
          <select 
            value={filters.jobType || ''} 
            onChange={(e) => handleChange('jobType', e.target.value)}
          >
            <option value="">All Job Types</option>
            {filterOptions.jobTypes.map(type => (
              <option key={type} value={type}>{type}</option>
            ))}
          </select>
        </div>
      )}

      {filterOptions.languages && (
        <div className="filter-group">
          <label>Language</label>
          <select 
            value={filters.language || ''} 
            onChange={(e) => handleChange('language', e.target.value)}
          >
            <option value="">All Languages</option>
            {filterOptions.languages.map(lang => (
              <option 
                key={typeof lang === 'string' ? lang : lang.code} 
                value={typeof lang === 'string' ? lang : lang.code}
              >
                {typeof lang === 'string' ? lang : lang.label}
              </option>
            ))}
          </select>
        </div>
      )}

      <button
        className="clear-filters"
        onClick={() => onFilterChange({})}
      >
        Clear Filters
      </button>
    </div>
  );
};

export default FilterBar;

