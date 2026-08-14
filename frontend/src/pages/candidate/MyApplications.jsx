import React, { useState, useEffect, useCallback } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';

// Icon Components
const LayoutDashboardIcon = () => (
  <svg className="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <rect x="3" y="3" width="7" height="9" rx="1" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    <rect x="14" y="3" width="7" height="5" rx="1" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    <rect x="14" y="12" width="7" height="9" rx="1" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    <rect x="3" y="16" width="7" height="5" rx="1" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
  </svg>
);

const UserIcon = () => (
  <svg className="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
  </svg>
);

const FileTextIcon = () => (
  <svg className="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
  </svg>
);

const BookmarkIcon = () => (
  <svg className="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
  </svg>
);

const SparklesIcon = () => (
  <svg className="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
  </svg>
);

const LogOutIcon = () => (
  <svg className="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
  </svg>
);

const MenuIcon = () => (
  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h16m-7 6h7" />
  </svg>
);

const XIcon = () => (
  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
  </svg>
);

export const MyApplications = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [applications, setApplications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('ALL');
  const [selectedApp, setSelectedApp] = useState(null);

  const API_BASE = process.env.REACT_APP_BACKEND_URL || 'http://localhost:3001';

  const handleAuthError = useCallback(() => {
    localStorage.removeItem('rojgar_token');
    alert('Session expired. Please login again.');
    navigate('/login');
  }, [navigate]);

  const fetchApplications = useCallback(async () => {
    setLoading(true);
    setError(null);
    const token = localStorage.getItem('rojgar_token');

    try {
      const res = await fetch(`${API_BASE}/api/v1/applications/me`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      });

      if (res.status === 401) {
        handleAuthError();
        return;
      }

      if (res.ok) {
        const data = await res.json();
        setApplications(data?.applications ?? data ?? []);
      } else {
        setError('Failed to fetch applications.');
      }
    } catch (err) {
      setError('Connection error. Please try again later.');
    } finally {
      setLoading(false);
    }
  }, [API_BASE, handleAuthError]);

  useEffect(() => {
    fetchApplications();
  }, [fetchApplications]);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const filteredApps = (applications || []).filter((app) => {
    const matchesSearch =
      (app?.job_title || '').toLowerCase().includes(searchTerm.toLowerCase()) ||
      (app?.company_name || '').toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus =
      statusFilter === 'ALL' ? true : (app?.status || '').toUpperCase() === statusFilter.toUpperCase();
    return matchesSearch && matchesStatus;
  });

  const getStatusBadge = (status) => {
    switch ((status || '').toUpperCase()) {
      case 'SHORTLISTED':
        return <span className="px-3 py-1 bg-emerald-50 text-emerald-700 border border-emerald-200 text-[10px] font-bold uppercase rounded-full tracking-wider">Shortlisted</span>;
      case 'INTERVIEW':
        return <span className="px-3 py-1 bg-blue-50 text-blue-700 border border-blue-200 text-[10px] font-bold uppercase rounded-full tracking-wider">Interview</span>;
      case 'OFFER':
      case 'SELECTED':
        return <span className="px-3 py-1 bg-purple-50 text-purple-700 border border-purple-200 text-[10px] font-bold uppercase rounded-full tracking-wider">Offered</span>;
      case 'REJECTED':
        return <span className="px-3 py-1 bg-rose-50 text-rose-700 border border-rose-200 text-[10px] font-bold uppercase rounded-full tracking-wider">Rejected</span>;
      default:
        return <span className="px-3 py-1 bg-amber-50 text-amber-700 border border-amber-200 text-[10px] font-bold uppercase rounded-full tracking-wider">Applied</span>;
    }
  };

  const navItems = [
    { label: 'Dashboard', path: '/dashboard/candidate', icon: LayoutDashboardIcon },
    { label: 'My Profile', path: '/candidate/profile', icon: UserIcon },
    { label: 'Applications', path: '/candidate/applications', icon: FileTextIcon },
    { label: 'Saved Jobs', path: '/candidate/saved-jobs', icon: BookmarkIcon },
    { label: 'AI Matches', path: '/candidate/ai-matches', icon: SparklesIcon },
  ];

  return (
    <div className="flex min-h-screen bg-slate-50 font-sans">
      {/* Mobile Sidebar Overlay */}
      {sidebarOpen && (
        <div className="fixed inset-0 bg-slate-900/50 z-40 md:hidden backdrop-blur-sm" onClick={() => setSidebarOpen(false)}></div>
      )}

      {/* Sidebar */}
      <aside className={`
        fixed inset-y-0 left-0 z-50 w-64 bg-slate-900 text-white transform transition-transform duration-300 ease-in-out
        md:translate-x-0 md:static md:h-screen sticky top-0
        ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
        flex flex-col justify-between p-6 shadow-xl
      `}>
        <div>
          <div className="flex items-center justify-between mb-8 md:hidden">
            <span className="font-bold text-xl">RojgarSetu</span>
            <button onClick={() => setSidebarOpen(false)} className="p-2 text-slate-400 hover:text-white">
              <XIcon />
            </button>
          </div>

          <div className="flex items-center gap-3 mb-6">
            <div className="bg-blue-600 w-12 h-12 rounded-full flex items-center justify-center font-bold text-white text-lg shrink-0">
              {(user?.name || 'C').charAt(0)}
            </div>
            <div className="overflow-hidden">
              <h2 className="font-semibold text-white truncate text-base">{user?.name || 'Candidate'}</h2>
              <span className="inline-block bg-blue-500/20 text-blue-300 text-xs px-2 py-0.5 rounded-md">Candidate</span>
            </div>
          </div>

          <nav className="space-y-1">
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = location.pathname === item.path;
              return (
                <Link
                  key={item.path}
                  to={item.path}
                  onClick={() => setSidebarOpen(false)}
                  className={`flex items-center px-4 py-3 text-sm font-medium transition-all ${
                    isActive ? 'bg-white/10 text-white rounded-lg shadow-sm' : 'text-slate-400 hover:text-white hover:bg-white/5 rounded-lg'
                  }`}
                >
                  <Icon />
                  {item.label}
                </Link>
              );
            })}
          </nav>
        </div>

        <button onClick={handleLogout} className="flex items-center w-full px-4 py-3 text-sm font-medium bg-red-500/10 text-red-400 hover:bg-red-500/20 rounded-lg transition-all">
          <LogOutIcon />
          Logout
        </button>
      </aside>

      {/* Main Content */}
      <main className="flex-1 overflow-auto p-4 md:p-8">
        <div className="max-w-5xl mx-auto">
          {/* Header */}
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
            <div className="flex items-center gap-4">
              <button onClick={() => setSidebarOpen(true)} className="p-2 bg-white border border-slate-200 rounded-lg shadow-sm text-slate-600 md:hidden">
                <MenuIcon />
              </button>
              <div>
                <h1 className="text-3xl font-bold text-slate-900 tracking-tight">My Applications</h1>
                <p className="text-slate-500 text-sm font-medium">Manage and track your job application progress.</p>
              </div>
            </div>
          </div>

          {/* Error Banner */}
          {error && (
            <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-2xl flex flex-col sm:flex-row items-center justify-between text-red-700 gap-4">
              <span className="text-sm font-medium">{error}</span>
              <button onClick={fetchApplications} className="px-4 py-2 bg-red-600 text-white text-xs font-bold rounded-lg hover:bg-red-700 shadow-md">Retry</button>
            </div>
          )}

          {/* Search & Filter Bar */}
          <div className="bg-white p-4 rounded-2xl shadow-sm border border-slate-200/60 mb-6 space-y-4">
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
              <div className="relative flex-1 max-w-md">
                <input
                  type="text"
                  placeholder="Search by job or company..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="w-full pl-4 pr-10 py-2.5 rounded-xl border border-slate-200 text-sm focus:ring-2 focus:ring-blue-500/20 outline-none"
                />
                <span className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400">🔍</span>
              </div>
              <div className="flex flex-wrap items-center gap-2">
                {['ALL', 'APPLIED', 'SHORTLISTED', 'INTERVIEW', 'REJECTED'].map((st) => (
                  <button
                    key={st}
                    onClick={() => setStatusFilter(st)}
                    className={`px-4 py-2 text-xs font-bold rounded-xl transition-all ${
                      statusFilter === st ? 'bg-slate-900 text-white shadow-md' : 'bg-slate-50 text-slate-600 hover:bg-slate-100'
                    }`}
                  >
                    {st === 'ALL' ? 'All' : st}
                  </button>
                ))}
              </div>
            </div>
          </div>

          {/* Applications List */}
          {loading ? (
            <div className="space-y-4 animate-pulse">
              {[1, 2, 3].map((n) => (
                <div key={n} className="h-32 bg-slate-200 rounded-2xl border border-slate-200"></div>
              ))}
            </div>
          ) : filteredApps.length === 0 ? (
            <div className="bg-white p-16 text-center rounded-3xl border border-slate-200/60 shadow-sm flex flex-col items-center">
              <div className="text-5xl mb-4">🔍</div>
              <h3 className="text-xl font-bold text-slate-900 mb-2">No Applications Found</h3>
              <p className="text-slate-500 max-w-sm mb-8 font-medium">
                {searchTerm ? 'Try adjusting your search filters to find what you are looking for.' : 'Start applying to jobs to track your progress here.'}
              </p>
              {!searchTerm && (
                <Link to="/gov-jobs" className="px-6 py-3 bg-blue-600 text-white font-bold rounded-xl hover:bg-blue-700 shadow-lg shadow-blue-500/25 transition-all">
                  Browse Jobs →
                </Link>
              )}
            </div>
          ) : (
            <div className="space-y-4">
              {filteredApps.map((app) => (
                <div
                  key={app?.id ?? Math.random()}
                  className="bg-white p-6 rounded-2xl border border-slate-200 hover:border-blue-300 transition-all shadow-sm group"
                >
                  <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
                    <div className="space-y-1 flex-1 overflow-hidden">
                      <div className="flex flex-wrap items-center gap-3">
                        <h3 className="text-lg font-bold text-slate-900 truncate">{app?.job_title ?? 'Untitled Position'}</h3>
                        {getStatusBadge(app?.status)}
                      </div>
                      <p className="text-sm font-semibold text-slate-600">
                        {app?.company_name ?? 'Unknown Company'} • <span className="text-slate-400 font-medium">{app?.location ?? 'Not specified'}</span>
                      </p>
                      <div className="flex flex-wrap items-center gap-4 text-xs text-slate-400 pt-2 font-medium">
                        <span className="flex items-center gap-1.5">📅 Applied: {app?.applied_date ?? 'Recent'}</span>
                        <span className="flex items-center gap-1.5">💰 {app?.salary ?? 'Competitive'}</span>
                        {app?.match_score && (
                          <span className="px-2 py-0.5 bg-blue-50 text-blue-600 rounded-md font-bold">✨ AI Match: {app.match_score}%</span>
                        )}
                      </div>
                    </div>
                    <button
                      onClick={() => setSelectedApp(app)}
                      className="w-full md:w-auto px-5 py-2.5 bg-slate-50 text-slate-800 text-xs font-bold rounded-xl border border-slate-200 hover:bg-slate-900 hover:text-white transition-all shadow-sm"
                    >
                      Track Timeline →
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}

          {/* Timeline Modal */}
          {selectedApp && (
            <div className="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4" onClick={() => setSelectedApp(null)}>
              <div className="bg-white w-full max-w-lg rounded-3xl p-8 shadow-2xl animate-in zoom-in-95" onClick={e => e.stopPropagation()}>
                <div className="flex justify-between items-start mb-6">
                  <div>
                    <h2 className="text-xl font-black text-slate-900">{selectedApp?.job_title}</h2>
                    <p className="text-sm text-slate-500 font-bold">{selectedApp?.company_name}</p>
                  </div>
                  <button onClick={() => setSelectedApp(null)} className="text-slate-400 hover:text-slate-900 text-2xl font-bold">×</button>
                </div>

                <div className="space-y-6 py-4">
                  {(selectedApp?.timeline || [
                    { status: 'Applied', date: selectedApp?.applied_date ?? 'Today', completed: true },
                    { status: 'Under Review', date: 'Processing', completed: false }
                  ]).map((step, idx, arr) => (
                    <div key={idx} className="flex gap-4">
                      <div className="flex flex-col items-center">
                        <div className={`w-8 h-8 rounded-full flex items-center justify-center text-xs font-black shrink-0 ${
                          step.completed ? 'bg-emerald-500 text-white' : 'bg-slate-100 text-slate-400'
                        }`}>
                          {step.completed ? '✓' : idx + 1}
                        </div>
                        {idx !== arr.length - 1 && <div className={`w-0.5 h-full ${step.completed ? 'bg-emerald-500' : 'bg-slate-100'}`}></div>}
                      </div>
                      <div className="pb-6">
                        <p className={`text-sm font-bold ${step.completed ? 'text-slate-900' : 'text-slate-400'}`}>{step.status}</p>
                        <p className="text-xs text-slate-400 font-medium">{step.date}</p>
                      </div>
                    </div>
                  ))}
                </div>

                <div className="pt-6 border-t border-slate-100 flex justify-end">
                  <button onClick={() => setSelectedApp(null)} className="px-6 py-2.5 bg-slate-900 text-white text-xs font-bold rounded-xl hover:bg-slate-800 transition-all">Close Tracker</button>
                </div>
              </div>
            </div>
          )}
        </div>
      </main>
    </div>
  );
};

export default MyApplications;