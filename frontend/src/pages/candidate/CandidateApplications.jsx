import React, { useState, useEffect } from 'react';

export const CandidateApplications = () => {
  const API_BASE = 'http://localhost:3001';

  const [activeTab, setActiveTab] = useState('applications'); // 'applications' | 'saved'
  const [applications, setApplications] = useState([]);
  const [savedJobs, setSavedJobs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [toast, setToast] = useState(null);

  useEffect(() => {
    fetchData();
  }, []);

  const showToast = (message, type = 'success') => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 3500);
  };

  const fetchData = async () => {
    setLoading(true);
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');

    try {
      const res = await fetch(`${API_BASE}/api/v1/candidate/applications`, {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
      });

      if (res.ok) {
        const data = await res.json();
        setApplications(data.applications || []);
        setSavedJobs(data.savedJobs || []);
      } else {
        setMockData();
      }
    } catch (err) {
      setMockData();
    } finally {
      setLoading(false);
    }
  };

  const setMockData = () => {
    setApplications([
      {
        id: 'app-101',
        job_title: 'Senior React / Frontend Engineer',
        company: 'Razorpay',
        logo_text: 'RZ',
        applied_date: '2026-03-01',
        status: 'Shortlisted',
        next_step: 'Technical Interview on March 12, 2:00 PM IST',
        salary_range: '₹18 - 25 LPA',
        location: 'Bengaluru / Remote'
      },
      {
        id: 'app-102',
        job_title: 'Full Stack Node.js & React Developer',
        company: 'Swiggy',
        logo_text: 'SW',
        applied_date: '2026-02-27',
        status: 'Under Review',
        next_step: 'Resume reviewed by Hiring Manager',
        salary_range: '₹16 - 22 LPA',
        location: 'Gurugram / Hybrid'
      },
      {
        id: 'app-103',
        job_title: 'Flutter Mobile App Engineer',
        company: 'Postman',
        logo_text: 'PM',
        applied_date: '2026-02-20',
        status: 'Applied',
        next_step: 'Awaiting recruiter evaluation',
        salary_range: '₹10 - 15 LPA',
        location: 'Delhi NCR'
      }
    ]);

    setSavedJobs([
      {
        id: 'job-4',
        title: 'Backend API Developer (Python/FastAPI)',
        company: 'Paytm',
        logo_text: 'PT',
        location: 'Noida',
        salary_range: '₹14 - 20 LPA',
        saved_date: '2026-03-02'
      }
    ]);
    setLoading(false);
  };

  const handleWithdraw = (appId, company) => {
    setApplications(applications.filter((a) => a.id !== appId));
    showToast(`Withdrew application for ${company}`, 'info');
  };

  const handleRemoveSaved = (jobId) => {
    setSavedJobs(savedJobs.filter((j) => j.id !== jobId));
    showToast('Removed from saved jobs', 'info');
  };

  const getStatusBadge = (status) => {
    switch (status) {
      case 'Shortlisted':
        return 'bg-emerald-50 text-emerald-700 border-emerald-200';
      case 'Under Review':
        return 'bg-amber-50 text-amber-700 border-amber-200';
      case 'Rejected':
        return 'bg-red-50 text-red-700 border-red-200';
      default:
        return 'bg-blue-50 text-blue-700 border-blue-200';
    }
  };

  return (
    <div className="min-h-screen bg-slate-50 p-8 font-sans">
      {/* Toast */}
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
          <span>{toast.message}</span>
        </div>
      )}

      <div className="max-w-5xl mx-auto">
        <header className="mb-8">
          <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Applications & Saved Jobs</h1>
          <p className="text-slate-500 text-sm mt-1 font-medium">
            Track active job applications, interview schedules, and bookmarked roles.
          </p>
        </header>

        {/* Tab Switcher */}
        <div className="flex border-b border-slate-200 mb-8">
          <button
            type="button"
            onClick={() => setActiveTab('applications')}
            className={`pb-4 px-6 text-sm font-bold border-b-2 transition-all ${
              activeTab === 'applications'
                ? 'border-blue-600 text-blue-600'
                : 'border-transparent text-slate-500 hover:text-slate-800'
            }`}
          >
            My Applications ({applications.length})
          </button>
          <button
            type="button"
            onClick={() => setActiveTab('saved')}
            className={`pb-4 px-6 text-sm font-bold border-b-2 transition-all ${
              activeTab === 'saved'
                ? 'border-blue-600 text-blue-600'
                : 'border-transparent text-slate-500 hover:text-slate-800'
            }`}
          >
            Saved Jobs ({savedJobs.length})
          </button>
        </div>

        {/* Applications View */}
        {activeTab === 'applications' && (
          <div className="space-y-4">
            {loading ? (
              <div className="space-y-3 animate-pulse">
                {[1, 2].map((i) => (
                  <div key={i} className="h-28 bg-slate-200 rounded-2xl"></div>
                ))}
              </div>
            ) : applications.length === 0 ? (
              <div className="bg-white p-12 text-center rounded-2xl border border-slate-200/60 shadow-sm">
                <p className="text-sm text-slate-500">You haven't applied to any positions yet.</p>
              </div>
            ) : (
              applications.map((app) => (
                <div
                  key={app.id}
                  className="bg-white p-6 rounded-2xl border border-slate-200/60 shadow-sm flex flex-col md:flex-row items-start justify-between gap-6"
                >
                  <div className="flex gap-4 items-start">
                    <div className="w-12 h-12 rounded-xl bg-slate-900 text-white font-extrabold flex items-center justify-center text-sm shadow-sm flex-shrink-0">
                      {app.logo_text || app.company.slice(0, 2).toUpperCase()}
                    </div>

                    <div className="space-y-2">
                      <div className="flex items-center gap-3 flex-wrap">
                        <h3 className="font-bold text-slate-900 text-lg">{app.job_title}</h3>
                        <span
                          className={`px-3 py-0.5 border text-xs font-bold rounded-full ${getStatusBadge(
                            app.status
                          )}`}
                        >
                          {app.status}
                        </span>
                      </div>

                      <p className="text-xs font-semibold text-slate-600">
                        {app.company} • <span className="text-slate-400 font-normal">{app.location}</span>
                      </p>

                      <div className="bg-slate-50 border border-slate-100 p-3 rounded-xl text-xs text-slate-700 font-medium mt-2">
                        <span className="text-slate-400 uppercase text-[10px] font-bold block mb-0.5">
                          Latest Update
                        </span>
                        {app.next_step}
                      </div>
                    </div>
                  </div>

                  <div className="flex flex-col items-start md:items-end justify-between w-full md:w-auto h-full gap-4">
                    <span className="text-xs text-slate-400 font-mono">Applied on {app.applied_date}</span>
                    <button
                      type="button"
                      onClick={() => handleWithdraw(app.id, app.company)}
                      className="px-3 py-1.5 bg-slate-100 text-slate-600 text-xs font-semibold rounded-lg hover:bg-red-50 hover:text-red-600 transition-all"
                    >
                      Withdraw
                    </button>
                  </div>
                </div>
              ))
            )}
          </div>
        )}

        {/* Saved Jobs View */}
        {activeTab === 'saved' && (
          <div className="space-y-4">
            {savedJobs.length === 0 ? (
              <div className="bg-white p-12 text-center rounded-2xl border border-slate-200/60 shadow-sm">
                <p className="text-sm text-slate-500">No bookmarked jobs found.</p>
              </div>
            ) : (
              savedJobs.map((job) => (
                <div
                  key={job.id}
                  className="bg-white p-6 rounded-2xl border border-slate-200/60 shadow-sm flex items-center justify-between gap-4"
                >
                  <div className="flex gap-4 items-center">
                    <div className="w-10 h-10 rounded-xl bg-slate-900 text-white font-extrabold flex items-center justify-center text-xs shadow-sm flex-shrink-0">
                      {job.logo_text}
                    </div>
                    <div>
                      <h4 className="font-bold text-slate-900 text-base">{job.title}</h4>
                      <p className="text-xs text-slate-500 font-medium">
                        {job.company} • {job.location} • <span className="text-slate-700">{job.salary_range}</span>
                      </p>
                    </div>
                  </div>

                  <div className="flex items-center gap-2">
                    <button
                      type="button"
                      onClick={() => handleRemoveSaved(job.id)}
                      className="px-3 py-2 bg-slate-100 text-slate-600 text-xs font-semibold rounded-xl hover:bg-slate-200"
                    >
                      Remove
                    </button>
                  </div>
                </div>
              ))
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default CandidateApplications;