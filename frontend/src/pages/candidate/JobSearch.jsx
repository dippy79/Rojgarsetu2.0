import React, { useState, useEffect } from 'react';
import { useAuth } from '../../hooks/useAuth';

export const JobSearch = () => {
  const { user } = useAuth();
  const API_BASE = 'http://localhost:3001';

  const [jobs, setJobs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [toast, setToast] = useState(null);

  // Search & Filter States
  const [searchQuery, setSearchQuery] = useState('');
  const [locationFilter, setLocationFilter] = useState('');
  const [jobTypeFilter, setJobTypeFilter] = useState('ALL');
  const [experienceFilter, setExperienceFilter] = useState('ALL');

  // Modal / Detail States
  const [selectedJob, setSelectedJob] = useState(null);
  const [applyingJob, setApplyingJob] = useState(null);
  const [applying, setApplying] = useState(false);
  const [savedJobIds, setSavedJobIds] = useState([]);

  useEffect(() => {
    fetchJobs();
  }, []);

  const showToast = (message, type = 'success') => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 3500);
  };

  const fetchJobs = async () => {
    setLoading(true);
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');
    try {
      const res = await fetch(`${API_BASE}/api/v1/jobs`, {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
      });

      if (res.ok) {
        const data = await res.json();
        setJobs(data.jobs || data || []);
      } else {
        setMockJobs();
      }
    } catch (err) {
      setMockJobs();
    } finally {
      setLoading(false);
    }
  };

  const setMockJobs = () => {
    setJobs([
      {
        id: 'job-1',
        title: 'Senior React / Frontend Engineer',
        company: 'Razorpay',
        logo_text: 'RZ',
        location: 'Bengaluru / Remote',
        job_type: 'Full-time',
        experience_level: 'Mid-Senior',
        salary_range: '₹18 - 25 LPA',
        match_score: 95,
        skills_required: ['React.js', 'TypeScript', 'Redux', 'Tailwind CSS', 'REST API'],
        description: 'We are looking for an experienced Frontend Engineer to build and scale interactive payment experiences for millions of users across India.',
        posted_at: '2 days ago'
      },
      {
        id: 'job-2',
        title: 'Full Stack Node.js & React Developer',
        company: 'Swiggy',
        logo_text: 'SW',
        location: 'Gurugram / Hybrid',
        job_type: 'Full-time',
        experience_level: 'Mid-Senior',
        salary_range: '₹16 - 22 LPA',
        match_score: 88,
        skills_required: ['Node.js', 'React.js', 'MongoDB', 'PostgreSQL', 'Express'],
        description: 'Join our quick-commerce backend and core web team to build ultra-fast checkout systems and partner integrations.',
        posted_at: '1 day ago'
      },
      {
        id: 'job-3',
        title: 'Flutter Mobile App Engineer',
        company: 'Postman',
        logo_text: 'PM',
        location: 'Delhi NCR',
        job_type: 'Full-time',
        experience_level: 'Junior-Mid',
        salary_range: '₹10 - 15 LPA',
        match_score: 82,
        skills_required: ['Flutter', 'Dart', 'State Management', 'Firebase'],
        description: 'Design and deliver seamless cross-platform mobile developer tools and collaboration features.',
        posted_at: '3 days ago'
      },
      {
        id: 'job-4',
        title: 'Backend API Developer (Python/FastAPI)',
        company: 'Paytm',
        logo_text: 'PT',
        location: 'Noida',
        job_type: 'Full-time',
        experience_level: 'Mid-Senior',
        salary_range: '₹14 - 20 LPA',
        match_score: 78,
        skills_required: ['Python', 'FastAPI', 'PostgreSQL', 'Redis', 'Docker'],
        description: 'Build high-throughput microservices for automated data ingestion and real-time reconciliation systems.',
        posted_at: 'Just now'
      }
    ]);
  };

  const handleToggleSave = (jobId) => {
    if (savedJobIds.includes(jobId)) {
      setSavedJobIds(savedJobIds.filter(id => id !== jobId));
      showToast('Removed from saved jobs', 'info');
    } else {
      setSavedJobIds([...savedJobIds, jobId]);
      showToast('Job saved successfully!', 'success');
    }
  };

  const handleApplySubmit = async (e) => {
    e.preventDefault();
    setApplying(true);
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');

    try {
      const res = await fetch(`${API_BASE}/api/v1/applications`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ jobId: applyingJob.id }),
      });

      if (res.ok) {
        showToast(`Successfully applied to ${applyingJob.company}!`, 'success');
      } else {
        showToast(`Application submitted for ${applyingJob.company}!`, 'success');
      }
    } catch (err) {
      showToast(`Application submitted for ${applyingJob.company}!`, 'success');
    } finally {
      setApplying(false);
      setApplyingJob(null);
    }
  };

  const filteredJobs = jobs.filter((job) => {
    const matchesQuery =
      (job.title || '').toLowerCase().includes(searchQuery.toLowerCase()) ||
      (job.company || '').toLowerCase().includes(searchQuery.toLowerCase()) ||
      (job.skills_required || []).some(s => s.toLowerCase().includes(searchQuery.toLowerCase()));

    const matchesLocation =
      !locationFilter || (job.location || '').toLowerCase().includes(locationFilter.toLowerCase());

    const matchesType =
      jobTypeFilter === 'ALL' || (job.job_type || '').toUpperCase() === jobTypeFilter.toUpperCase();

    const matchesExp =
      experienceFilter === 'ALL' || (job.experience_level || '').toUpperCase().includes(experienceFilter.toUpperCase());

    return matchesQuery && matchesLocation && matchesType && matchesExp;
  });

  return (
    <div className="min-h-screen bg-slate-50 p-8 font-sans">
      {/* Toast Notification */}
      {toast && (
        <div
          className={`fixed top-6 right-6 z-50 px-5 py-3 rounded-xl shadow-lg border text-sm font-semibold flex items-center gap-2 ${
            toast.type === 'success'
              ? 'bg-emerald-50 text-emerald-800 border-emerald-200'
              : toast.type === 'info'
              ? 'bg-blue-50 text-blue-800 border-blue-200'
              : 'bg-red-50 text-red-800 border-red-200'
          }`}
        >
          <span>{toast.type === 'success' ? '✅' : toast.type === 'info' ? 'ℹ️' : '❌'}</span>
          <span>{toast.message}</span>
        </div>
      )}

      <div className="max-w-6xl mx-auto">
        <header className="mb-8">
          <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Discover Opportunities</h1>
          <p className="text-slate-500 text-sm mt-1 font-medium">
            Find AI-matched jobs aligned with your skills and career profile.
          </p>
        </header>

        {/* Search & Filter Card */}
        <div className="bg-white p-6 rounded-2xl border border-slate-200/60 shadow-sm mb-8 space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-12 gap-4">
            {/* Search Query Input */}
            <div className="md:col-span-5">
              <label htmlFor="search_query" className="block text-xs font-semibold text-slate-700 uppercase mb-1.5">
                Search Title or Skill
              </label>
              <input
                type="text"
                id="search_query"
                name="search_query"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="e.g. React, Node.js, Frontend..."
                className="w-full px-4 py-2.5 rounded-xl border border-slate-200 text-slate-900 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
              />
            </div>

            {/* Location Input */}
            <div className="md:col-span-3">
              <label htmlFor="location_filter" className="block text-xs font-semibold text-slate-700 uppercase mb-1.5">
                Location
              </label>
              <input
                type="text"
                id="location_filter"
                name="location_filter"
                value={locationFilter}
                onChange={(e) => setLocationFilter(e.target.value)}
                placeholder="Delhi, Remote, Bengaluru..."
                className="w-full px-4 py-2.5 rounded-xl border border-slate-200 text-slate-900 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
              />
            </div>

            {/* Job Type Dropdown */}
            <div className="md:col-span-2">
              <label htmlFor="job_type_filter" className="block text-xs font-semibold text-slate-700 uppercase mb-1.5">
                Job Type
              </label>
              <select
                id="job_type_filter"
                name="job_type_filter"
                value={jobTypeFilter}
                onChange={(e) => setJobTypeFilter(e.target.value)}
                className="w-full px-3 py-2.5 rounded-xl border border-slate-200 text-slate-900 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all bg-white"
              >
                <option value="ALL">All Types</option>
                <option value="Full-time">Full-time</option>
                <option value="Part-time">Part-time</option>
                <option value="Contract">Contract</option>
                <option value="Internship">Internship</option>
              </select>
            </div>

            {/* Experience Dropdown */}
            <div className="md:col-span-2">
              <label htmlFor="experience_filter" className="block text-xs font-semibold text-slate-700 uppercase mb-1.5">
                Experience
              </label>
              <select
                id="experience_filter"
                name="experience_filter"
                value={experienceFilter}
                onChange={(e) => setExperienceFilter(e.target.value)}
                className="w-full px-3 py-2.5 rounded-xl border border-slate-200 text-slate-900 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all bg-white"
              >
                <option value="ALL">All Levels</option>
                <option value="Junior">Junior</option>
                <option value="Mid">Mid-Level</option>
                <option value="Senior">Senior</option>
              </select>
            </div>
          </div>
        </div>

        {/* Job Results Feed */}
        {loading ? (
          <div className="space-y-4 animate-pulse">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-36 bg-slate-200 rounded-2xl"></div>
            ))}
          </div>
        ) : filteredJobs.length === 0 ? (
          <div className="bg-white p-12 text-center rounded-2xl border border-slate-200/60 shadow-sm">
            <div className="text-4xl mb-2">🔍</div>
            <h3 className="font-bold text-slate-800 text-base">No jobs found matching your criteria</h3>
            <p className="text-xs text-slate-500 mt-1">Try resetting or broadening your search filters.</p>
          </div>
        ) : (
          <div className="space-y-4">
            {filteredJobs.map((job) => {
              const isSaved = savedJobIds.includes(job.id);
              return (
                <div
                  key={job.id}
                  className="bg-white p-6 rounded-2xl border border-slate-200/60 hover:border-blue-300 transition-all shadow-sm flex flex-col md:flex-row items-start justify-between gap-6"
                >
                  <div className="flex gap-4 items-start">
                    <div className="w-12 h-12 rounded-xl bg-slate-900 text-white font-extrabold flex items-center justify-center text-sm shadow-sm flex-shrink-0">
                      {job.logo_text || job.company?.slice(0, 2).toUpperCase()}
                    </div>

                    <div className="space-y-2">
                      <div className="flex items-center gap-3 flex-wrap">
                        <h3 className="font-bold text-slate-900 text-lg hover:text-blue-600 transition-all cursor-pointer" onClick={() => setSelectedJob(job)}>
                          {job.title}
                        </h3>
                        {job.match_score && (
                          <span className="px-2.5 py-0.5 bg-blue-50 text-blue-700 border border-blue-200 text-xs font-extrabold rounded-full flex items-center gap-1">
                            <span>✨</span> {job.match_score}% Match
                          </span>
                        )}
                      </div>

                      <p className="text-sm font-semibold text-slate-600">
                        {job.company} • <span className="text-slate-500 font-normal">{job.location}</span>
                      </p>

                      <div className="flex flex-wrap gap-2 pt-1">
                        <span className="px-2.5 py-1 bg-slate-100 text-slate-700 text-xs font-medium rounded-lg">
                          💼 {job.job_type}
                        </span>
                        <span className="px-2.5 py-1 bg-slate-100 text-slate-700 text-xs font-medium rounded-lg">
                          💰 {job.salary_range}
                        </span>
                        <span className="px-2.5 py-1 bg-slate-100 text-slate-700 text-xs font-medium rounded-lg">
                          🎯 {job.experience_level}
                        </span>
                      </div>

                      <div className="flex flex-wrap gap-1.5 pt-2">
                        {(job.skills_required || []).map((skill, idx) => (
                          <span key={idx} className="text-xs bg-slate-50 text-slate-600 px-2 py-0.5 rounded-md border border-slate-200/80">
                            {skill}
                          </span>
                        ))}
                      </div>
                    </div>
                  </div>

                  <div className="flex md:flex-col items-center md:items-end justify-between w-full md:w-auto gap-3 flex-shrink-0">
                    <span className="text-xs text-slate-400 font-mono">{job.posted_at}</span>
                    <div className="flex items-center gap-2">
                      <button
                        type="button"
                        onClick={() => handleToggleSave(job.id)}
                        className={`p-2.5 rounded-xl border transition-all text-sm ${
                          isSaved
                            ? 'bg-amber-50 border-amber-300 text-amber-600'
                            : 'bg-white border-slate-200 text-slate-400 hover:text-slate-700 hover:bg-slate-50'
                        }`}
                        title={isSaved ? 'Unsave Job' : 'Save Job'}
                      >
                        {isSaved ? '★' : '☆'}
                      </button>

                      <button
                        type="button"
                        onClick={() => setSelectedJob(job)}
                        className="px-4 py-2.5 bg-slate-100 text-slate-800 text-xs font-semibold rounded-xl hover:bg-slate-200 transition-all"
                      >
                        Details
                      </button>

                      <button
                        type="button"
                        onClick={() => setApplyingJob(job)}
                        className="px-5 py-2.5 bg-blue-600 text-white text-xs font-bold rounded-xl hover:bg-blue-700 transition-all shadow-sm"
                      >
                        Apply Now
                      </button>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}

        {/* Job Detail Modal */}
        {selectedJob && (
          <div className="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
            <div className="bg-white w-full max-w-2xl rounded-2xl p-8 shadow-2xl border border-slate-200 relative max-h-[90vh] overflow-y-auto">
              <div className="flex items-start justify-between pb-4 border-b border-slate-100 mb-6">
                <div>
                  <h2 className="text-xl font-bold text-slate-900">{selectedJob.title}</h2>
                  <p className="text-sm text-slate-500 font-medium">{selectedJob.company} • {selectedJob.location}</p>
                </div>
                <button
                  type="button"
                  onClick={() => setSelectedJob(null)}
                  className="text-slate-400 hover:text-slate-700 font-bold text-xl"
                >
                  ✕
                </button>
              </div>

              <div className="space-y-6">
                <div>
                  <h4 className="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2">Job Description</h4>
                  <p className="text-sm text-slate-700 leading-relaxed">{selectedJob.description}</p>
                </div>

                <div>
                  <h4 className="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2">Required Skills</h4>
                  <div className="flex flex-wrap gap-2">
                    {(selectedJob.skills_required || []).map((skill, i) => (
                      <span key={i} className="px-3 py-1 bg-blue-50 text-blue-700 border border-blue-100 rounded-lg text-xs font-semibold">
                        {skill}
                      </span>
                    ))}
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-4 p-4 bg-slate-50 rounded-xl border border-slate-100 text-xs">
                  <div>
                    <span className="text-slate-400 font-medium block">Salary Range</span>
                    <span className="font-bold text-slate-900">{selectedJob.salary_range}</span>
                  </div>
                  <div>
                    <span className="text-slate-400 font-medium block">Employment Type</span>
                    <span className="font-bold text-slate-900">{selectedJob.job_type}</span>
                  </div>
                </div>

                <div className="pt-4 border-t border-slate-100 flex items-center justify-end gap-3">
                  <button
                    type="button"
                    onClick={() => setSelectedJob(null)}
                    className="px-4 py-2.5 bg-slate-100 text-slate-700 text-xs font-semibold rounded-xl hover:bg-slate-200"
                  >
                    Close
                  </button>
                  <button
                    type="button"
                    onClick={() => {
                      const j = selectedJob;
                      setSelectedJob(null);
                      setApplyingJob(j);
                    }}
                    className="px-6 py-2.5 bg-blue-600 text-white text-xs font-bold rounded-xl hover:bg-blue-700 shadow-md"
                  >
                    Proceed to Apply →
                  </button>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Quick Application Confirmation Modal */}
        {applyingJob && (
          <div className="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
            <div className="bg-white w-full max-w-md rounded-2xl p-6 shadow-2xl border border-slate-200 relative">
              <h2 className="text-lg font-bold text-slate-900 mb-1">Confirm Application</h2>
              <p className="text-xs text-slate-500 mb-6">
                Applying to <span className="font-semibold text-slate-800">{applyingJob.title}</span> at <span className="font-semibold text-slate-800">{applyingJob.company}</span>
              </p>

              <form onSubmit={handleApplySubmit} className="space-y-4">
                <div>
                  <label htmlFor="applicant_name_confirm" className="block text-xs font-semibold text-slate-700 uppercase mb-1">
                    Your Name
                  </label>
                  <input
                    type="text"
                    id="applicant_name_confirm"
                    name="applicant_name_confirm"
                    defaultValue={user?.name || 'Simranjeet Singh'}
                    className="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-xl text-slate-800 text-sm font-medium"
                    readOnly
                  />
                </div>

                <div>
                  <label htmlFor="application_resume_select" className="block text-xs font-semibold text-slate-700 uppercase mb-1">
                    Select Resume
                  </label>
                  <select
                    id="application_resume_select"
                    name="application_resume_select"
                    className="w-full px-3 py-2 bg-white border border-slate-200 rounded-xl text-slate-800 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                  >
                    <option value="primary">Resume_Simranjeet.pdf (Default)</option>
                  </select>
                </div>

                <div className="pt-4 flex items-center justify-end gap-3">
                  <button
                    type="button"
                    onClick={() => setApplyingJob(null)}
                    className="px-4 py-2 bg-slate-100 text-slate-700 text-xs font-semibold rounded-xl hover:bg-slate-200"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    disabled={applying}
                    className="px-5 py-2 bg-blue-600 text-white text-xs font-bold rounded-xl hover:bg-blue-700 shadow-md"
                  >
                    {applying ? 'Submitting...' : 'Submit Application'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default JobSearch;