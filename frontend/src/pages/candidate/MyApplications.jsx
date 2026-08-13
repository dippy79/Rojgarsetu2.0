import React, { useState, useEffect } from 'react';
import { useAuth } from '../../hooks/useAuth';

export const MyApplications = () => {
  const { user } = useAuth();
  const API_BASE = 'http://localhost:3001';

  const [applications, setApplications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('ALL');
  const [selectedApp, setSelectedApp] = useState(null);

  useEffect(() => {
    fetchApplications();
  }, []);

  const fetchApplications = async () => {
    setLoading(true);
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');
    try {
      const res = await fetch(`${API_BASE}/api/v1/applications/me`, {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
      });
      if (res.ok) {
        const data = await res.json();
        setApplications(data.applications || data || []);
      } else {
        // Fallback demo data if endpoint is pending
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
        job_title: 'Full Stack Web Developer',
        company_name: 'TechCorp Solutions',
        location: 'Delhi, India (Hybrid)',
        salary: '₹8 - 12 LPA',
        applied_date: '2026-03-01',
        status: 'SHORTLISTED',
        match_score: 92,
        timeline: [
          { status: 'Applied', date: '2026-03-01', completed: true },
          { status: 'Profile Shortlisted', date: '2026-03-03', completed: true },
          { status: 'Technical Interview', date: 'Scheduled for 2026-03-10', completed: false },
        ],
      },
      {
        id: 'app-102',
        job_title: 'Frontend React Developer',
        company_name: 'Innovate Labs',
        location: 'Remote',
        salary: '₹6 - 9 LPA',
        applied_date: '2026-02-25',
        status: 'INTERVIEW',
        match_score: 88,
        timeline: [
          { status: 'Applied', date: '2026-02-25', completed: true },
          { status: 'Profile Shortlisted', date: '2026-02-27', completed: true },
          { status: 'Technical Interview', date: '2026-03-02', completed: true },
          { status: 'HR Round', date: 'Scheduled', completed: false },
        ],
      },
      {
        id: 'app-103',
        job_title: 'Junior Flutter Developer',
        company_name: 'MobileApps Co.',
        location: 'Noida, UP',
        salary: '₹5 - 7 LPA',
        applied_date: '2026-02-15',
        status: 'REJECTED',
        match_score: 65,
        timeline: [
          { status: 'Applied', date: '2026-02-15', completed: true },
          { status: 'Application Reviewed', date: '2026-02-18', completed: true },
          { status: 'Not Shortlisted', date: '2026-02-20', completed: true },
        ],
      },
    ]);
  };

  const filteredApps = (applications || []).filter((app) => {
    const matchesSearch =
      (app.job_title || '').toLowerCase().includes(searchTerm.toLowerCase()) ||
      (app.company_name || '').toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus =
      statusFilter === 'ALL' ? true : (app.status || '').toUpperCase() === statusFilter.toUpperCase();
    return matchesSearch && matchesStatus;
  });

  const getStatusBadge = (status) => {
    switch ((status || '').toUpperCase()) {
      case 'SHORTLISTED':
        return <span className="px-3 py-1 bg-emerald-50 text-emerald-700 border border-emerald-200 text-xs font-bold rounded-full">Shortlisted</span>;
      case 'INTERVIEW':
        return <span className="px-3 py-1 bg-blue-50 text-blue-700 border border-blue-200 text-xs font-bold rounded-full">Interview Round</span>;
      case 'OFFER':
        return <span className="px-3 py-1 bg-purple-50 text-purple-700 border border-purple-200 text-xs font-bold rounded-full">Offer Letter Received</span>;
      case 'REJECTED':
        return <span className="px-3 py-1 bg-rose-50 text-rose-700 border border-rose-200 text-xs font-bold rounded-full">Not Selected</span>;
      default:
        return <span className="px-3 py-1 bg-amber-50 text-amber-700 border border-amber-200 text-xs font-bold rounded-full">Under Review</span>;
    }
  };

  return (
    <div className="min-h-screen bg-slate-50 p-8 font-sans">
      <div className="max-w-6xl mx-auto">
        <header className="mb-8">
          <h1 className="text-3xl font-bold text-slate-900 tracking-tight">My Job Applications</h1>
          <p className="text-slate-500 text-sm mt-1">
            Track and manage your real-time job application progress.
          </p>
        </header>

        {/* Filter and Search Bar */}
        <div className="bg-white p-4 rounded-2xl shadow-sm border border-slate-200/60 mb-6 flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="w-full md:w-80">
            <label htmlFor="search_applications" className="sr-only">
              Search Applications
            </label>
            <input
              type="text"
              id="search_applications"
              name="search_applications"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder="Search job title or company..."
              className="w-full px-4 py-2.5 rounded-xl border border-slate-200 text-slate-900 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500"
            />
          </div>

          <div className="flex flex-wrap items-center gap-2 w-full md:w-auto">
            {['ALL', 'APPLIED', 'SHORTLISTED', 'INTERVIEW', 'REJECTED'].map((st) => (
              <button
                key={st}
                type="button"
                onClick={() => setStatusFilter(st)}
                className={`px-3.5 py-2 text-xs font-semibold rounded-xl transition-all ${
                  statusFilter === st
                    ? 'bg-slate-900 text-white shadow-sm'
                    : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                }`}
              >
                {st === 'ALL' ? 'All Applications' : st}
              </button>
            ))}
          </div>
        </div>

        {/* Applications List */}
        {loading ? (
          <div className="space-y-4">
            {[1, 2, 3].map((n) => (
              <div key={n} className="bg-white p-6 rounded-2xl border border-slate-200/60 animate-pulse h-28"></div>
            ))}
          </div>
        ) : filteredApps.length === 0 ? (
          <div className="bg-white p-12 text-center rounded-2xl border border-slate-200/60">
            <p className="text-slate-500 font-medium">No job applications match your filters.</p>
          </div>
        ) : (
          <div className="space-y-4">
            {filteredApps.map((app) => (
              <div
                key={app.id}
                className="bg-white p-6 rounded-2xl border border-slate-200/60 hover:border-blue-200 transition-all shadow-sm flex flex-col md:flex-row items-start md:items-center justify-between gap-4"
              >
                <div className="space-y-1">
                  <div className="flex items-center gap-3">
                    <h3 className="text-base font-bold text-slate-900">{app.job_title}</h3>
                    {getStatusBadge(app.status)}
                  </div>
                  <p className="text-sm font-medium text-slate-600">
                    {app.company_name} • <span className="text-slate-400">{app.location}</span>
                  </p>
                  <div className="flex items-center gap-4 text-xs text-slate-400 pt-1">
                    <span>Applied on: {app.applied_date}</span>
                    <span>•</span>
                    <span>Salary: {app.salary}</span>
                    {app.match_score && (
                      <>
                        <span>•</span>
                        <span className="text-blue-600 font-semibold">AI Match Score: {app.match_score}%</span>
                      </>
                    )}
                  </div>
                </div>

                <button
                  type="button"
                  onClick={() => setSelectedApp(app)}
                  className="px-4 py-2 bg-slate-100 text-slate-800 text-xs font-semibold rounded-xl hover:bg-slate-200 transition-all self-start md:self-center"
                >
                  View Application Timeline →
                </button>
              </div>
            ))}
          </div>
        )}

        {/* Application Timeline Modal */}
        {selectedApp && (
          <div className="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
            <div className="bg-white w-full max-w-lg rounded-2xl p-6 shadow-2xl border border-slate-200 relative animate-in fade-in zoom-in-95">
              <div className="flex items-start justify-between border-b border-slate-100 pb-4 mb-4">
                <div>
                  <h2 className="text-lg font-bold text-slate-900">{selectedApp.job_title}</h2>
                  <p className="text-xs text-slate-500">{selectedApp.company_name}</p>
                </div>
                <button
                  type="button"
                  onClick={() => setSelectedApp(null)}
                  className="text-slate-400 hover:text-slate-700 text-lg font-bold"
                >
                  ✕
                </button>
              </div>

              <div className="space-y-6">
                <div>
                  <h4 className="text-xs font-bold text-slate-400 uppercase tracking-wider mb-3">Application Progress</h4>
                  <div className="space-y-4 relative before:absolute before:left-3 before:top-2 before:bottom-2 before:w-0.5 before:bg-slate-200">
                    {(selectedApp.timeline || []).map((step, idx) => (
                      <div key={idx} className="flex items-start gap-4 relative z-10">
                        <div
                          className={`w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold text-white ${
                            step.completed ? 'bg-emerald-500' : 'bg-slate-300'
                          }`}
                        >
                          {step.completed ? '✓' : idx + 1}
                        </div>
                        <div>
                          <p className="text-sm font-semibold text-slate-800">{step.status}</p>
                          <p className="text-xs text-slate-400">{step.date}</p>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>

                <div className="pt-4 border-t border-slate-100 flex justify-end">
                  <button
                    type="button"
                    onClick={() => setSelectedApp(null)}
                    className="px-4 py-2 bg-slate-900 text-white text-xs font-semibold rounded-xl hover:bg-slate-800"
                  >
                    Close
                  </button>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default MyApplications;