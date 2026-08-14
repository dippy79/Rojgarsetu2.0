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

const AnimatedCounter = ({ value, label }) => {
  const [count, setCount] = useState(0);

  useEffect(() => {
    let start = 0;
    const end = parseInt(value, 10) || 0;
    if (start === end) {
      setCount(end);
      return;
    }
    const duration = 1500;
    const incrementTime = 30;
    const steps = duration / incrementTime;
    const stepValue = Math.ceil((end - start) / steps);

    const timer = setInterval(() => {
      start += stepValue;
      if (start >= end) {
        setCount(end);
        clearInterval(timer);
      } else {
        setCount(start);
      }
    }, incrementTime);

    return () => clearInterval(timer);
  }, [value]);

  return (
    <div className="text-center p-4 bg-white/60 backdrop-blur-md rounded-xl border border-slate-200/60">
      <div className="font-mono text-3xl font-extrabold text-slate-900">{count?.toLocaleString() ?? 0}+</div>
      <div className="text-xs font-semibold text-slate-500 uppercase tracking-wider mt-1">{label ?? 'Stat'}</div>
    </div>
  );
};

export const CandidateDashboard = () => {
  const { logout, user } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [candidateProfile, setCandidateProfile] = useState(null);
  const [recentApps, setRecentApps] = useState([]);
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [platformStats, setPlatformStats] = useState({
    totalJobs: 1250,
    companies: 450,
    candidates: 8900,
    placements: 3200,
  });

  const API_BASE = process.env.REACT_APP_BACKEND_URL || 'http://localhost:3001';

  const handleAuthError = useCallback(() => {
    localStorage.removeItem('rojgar_token');
    alert('Session expired. Please login again.');
    navigate('/login');
  }, [navigate]);

  const fetchData = useCallback(async () => {
    setLoading(true);
    setError(null);
    const token = localStorage.getItem('rojgar_token');

    try {
      const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      };

      const [profileRes, appsRes, statsRes] = await Promise.allSettled([
        fetch(`${API_BASE}/api/v1/candidates/me`, { headers }),
        fetch(`${API_BASE}/api/v1/candidates/me/applications?limit=5`, { headers }),
        fetch(`${API_BASE}/api/v1/stats`),
      ]);

      if (profileRes.status === 'fulfilled') {
        if (profileRes.value.status === 401) {
          handleAuthError();
          return;
        }
        if (profileRes.value.ok) {
          const pData = await profileRes.value.json();
          setCandidateProfile(pData?.candidate ?? pData);
        }
      }

      if (appsRes.status === 'fulfilled') {
        if (appsRes.value.status === 401) {
          handleAuthError();
          return;
        }
        if (appsRes.value.ok) {
          const aData = await appsRes.value.json();
          setRecentApps(aData?.applications ?? aData ?? []);
        }
      }

      if (statsRes.status === 'fulfilled' && statsRes.value.ok) {
        const sData = await statsRes.value.json();
        setPlatformStats({
          totalJobs: sData?.totalJobs ?? 1250,
          companies: sData?.companies ?? 450,
          candidates: sData?.candidates ?? 8900,
          placements: sData?.placements ?? 3200,
        });
      }
    } catch (err) {
      setError('Failed to load dashboard data. Please check your connection.');
    } finally {
      setLoading(false);
    }
  }, [API_BASE, handleAuthError]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const calculateCompletion = (profile) => {
    if (!profile) return 0;
    let score = 0;
    if (profile?.full_name || profile?.name) score += 20;
    if (profile?.skills?.length > 0) score += 20;
    if (profile?.experience?.length > 0) score += 20;
    if (profile?.education?.length > 0) score += 20;
    if (profile?.resume_url || profile?.resume) score += 20;
    return score;
  };

  const profileCompletion = calculateCompletion(candidateProfile);

  const getInitials = (name) => {
    if (!name) return 'CN';
    const parts = name.trim().split(' ');
    if (parts.length >= 2) {
      return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
    }
    return name.slice(0, 2).toUpperCase();
  };

  const navItems = [
    { label: 'Dashboard', path: '/dashboard/candidate', icon: LayoutDashboardIcon },
    { label: 'My Profile', path: '/candidate/profile', icon: UserIcon },
    { label: 'Applications', path: '/candidate/applications', icon: FileTextIcon },
    { label: 'Saved Jobs', path: '/candidate/saved-jobs', icon: BookmarkIcon },
    { label: 'AI Matches', path: '/candidate/ai-matches', icon: SparklesIcon },
  ];

  const getStatusBadge = (status) => {
    const s = (status || '').toUpperCase();
    switch (s) {
      case 'SHORTLISTED':
        return 'bg-amber-50 text-amber-700 border-amber-200';
      case 'INTERVIEW':
      case 'INTERVIEW_SCHEDULED':
        return 'bg-violet-50 text-violet-700 border-violet-200';
      case 'OFFERED':
      case 'SELECTED':
        return 'bg-emerald-50 text-emerald-700 border-emerald-200';
      case 'REJECTED':
        return 'bg-red-50 text-red-600 border-red-200';
      default:
        return 'bg-blue-50 text-blue-700 border-blue-200';
    }
  };

  const aiMatches = [
    { id: 1, title: 'Senior React Developer', company: 'Razorpay', match: '95%' },
    { id: 2, title: 'Full Stack Node/React', company: 'Swiggy', match: '88%' },
    { id: 3, title: 'Frontend Specialist', company: 'Postman', match: '82%' },
  ];

  return (
    <div className="flex min-h-screen bg-slate-50 font-sans">
      {/* Mobile Sidebar Overlay */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 bg-slate-900/50 z-40 md:hidden backdrop-blur-sm"
          onClick={() => setSidebarOpen(false)}
        ></div>
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

          {/* User Avatar & Info */}
          <div className="flex items-center gap-3 mb-6">
            <div className="bg-blue-600 w-12 h-12 rounded-full flex items-center justify-center font-bold text-white text-lg shadow-md shrink-0">
              {getInitials(candidateProfile?.full_name ?? user?.name)}
            </div>
            <div className="overflow-hidden">
              <h2 className="font-semibold text-white truncate text-base">
                {candidateProfile?.full_name ?? user?.name ?? 'Candidate'}
              </h2>
              <span className="inline-block bg-blue-500/20 text-blue-300 text-xs px-2 py-0.5 rounded-md font-medium mt-0.5">
                Candidate
              </span>
            </div>
          </div>

          {/* Profile Completion Bar */}
          <div className="mb-8">
            <div className="flex justify-between text-xs text-slate-400 mb-1 font-medium">
              <span>Profile Completion</span>
              <span>{profileCompletion}%</span>
            </div>
            <div className="w-full bg-slate-700 rounded-full h-2 overflow-hidden">
              <div
                className="bg-blue-500 h-2 rounded-full transition-all duration-500"
                style={{ width: `${profileCompletion}%` }}
              ></div>
            </div>
          </div>

          {/* Navigation Links */}
          <nav className="space-y-1">
            {(navItems || []).map((item) => {
              const Icon = item.icon;
              const isActive = location.pathname === item.path;
              return (
                <Link
                  key={item.path}
                  to={item.path}
                  onClick={() => setSidebarOpen(false)}
                  className={`flex items-center px-4 py-3 text-sm font-medium transition-all ${
                    isActive
                      ? 'bg-white/10 text-white rounded-lg shadow-sm'
                      : 'text-slate-400 hover:text-white hover:bg-white/5 rounded-lg'
                  }`}
                >
                  <Icon />
                  {item.label}
                </Link>
              );
            })}
          </nav>
        </div>

        {/* Logout Button */}
        <button
          onClick={handleLogout}
          className="flex items-center w-full px-4 py-3 text-sm font-medium bg-red-500/10 text-red-400 hover:bg-red-500/20 rounded-lg transition-all mt-auto"
        >
          <LogOutIcon />
          Logout
        </button>
      </aside>

      {/* Main Content Area */}
      <main className="flex-1 overflow-auto p-4 md:p-8 bg-slate-50">
        {/* Mobile Header */}
        <div className="flex items-center gap-4 mb-6 md:hidden">
          <button
            onClick={() => setSidebarOpen(true)}
            className="p-2 bg-white border border-slate-200 rounded-lg shadow-sm text-slate-600"
          >
            <MenuIcon />
          </button>
          <h1 className="text-xl font-bold text-slate-900">Dashboard</h1>
        </div>

        {/* Desktop Header */}
        <header className="mb-8 hidden md:block">
          <h1 className="text-3xl font-bold text-slate-900 tracking-tight">
            Good morning, {(candidateProfile?.full_name ?? user?.name)?.split(' ')[0] ?? 'Candidate'} 👋
          </h1>
          <p className="text-slate-500 text-sm mt-1 font-medium">
            {new Date().toLocaleDateString('en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}
          </p>
        </header>

        {/* Error Banner */}
        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-2xl flex flex-col md:flex-row items-center justify-between text-red-700 gap-4">
            <span className="text-sm font-medium">{error}</span>
            <button
              onClick={fetchData}
              className="px-4 py-2 bg-red-600 text-white text-xs font-semibold rounded-lg hover:bg-red-700 transition-all shadow-sm"
            >
              Retry
            </button>
          </div>
        )}

        {/* Loading Skeleton */}
        {loading ? (
          <div className="space-y-6 animate-pulse">
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
              {[1, 2, 3, 4].map((i) => (
                <div key={i} className="h-32 bg-slate-200 rounded-2xl"></div>
              ))}
            </div>
            <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
              <div className="lg:col-span-7 h-64 bg-slate-200 rounded-2xl"></div>
              <div className="lg:col-span-5 h-64 bg-slate-200 rounded-2xl"></div>
            </div>
            <div className="h-32 bg-slate-200 rounded-2xl"></div>
          </div>
        ) : (
          /* Bento Grid Layout */
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
            {/* Row 1: Stats Cards */}
            <div className="col-span-1 sm:col-span-2 lg:col-span-3 bg-white/80 backdrop-blur-xl border border-slate-200/60 rounded-2xl p-6 shadow-sm">
              <div className="flex items-center justify-between">
                <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Applied</span>
                <div className="w-10 h-10 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center font-bold">
                  📄
                </div>
              </div>
              <div className="font-mono text-3xl font-extrabold text-slate-900 mt-4">12</div>
              <span className="text-xs text-slate-400 mt-1 block">Total submitted</span>
            </div>

            <div className="col-span-1 sm:col-span-2 lg:col-span-3 bg-white/80 backdrop-blur-xl border border-slate-200/60 rounded-2xl p-6 shadow-sm">
              <div className="flex items-center justify-between">
                <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Shortlisted</span>
                <div className="w-10 h-10 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center font-bold">
                  ⭐
                </div>
              </div>
              <div className="font-mono text-3xl font-extrabold text-slate-900 mt-4">4</div>
              <span className="text-xs text-slate-400 mt-1 block">Under review</span>
            </div>

            <div className="col-span-1 sm:col-span-2 lg:col-span-3 bg-white/80 backdrop-blur-xl border border-slate-200/60 rounded-2xl p-6 shadow-sm">
              <div className="flex items-center justify-between">
                <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Interview</span>
                <div className="w-10 h-10 rounded-xl bg-violet-50 text-violet-600 flex items-center justify-center font-bold">
                  🎙️
                </div>
              </div>
              <div className="font-mono text-3xl font-extrabold text-slate-900 mt-4">2</div>
              <span className="text-xs text-slate-400 mt-1 block">Scheduled calls</span>
            </div>

            <div className="col-span-1 sm:col-span-2 lg:col-span-3 bg-white/80 backdrop-blur-xl border border-slate-200/60 rounded-2xl p-6 shadow-sm">
              <div className="flex items-center justify-between">
                <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Offered</span>
                <div className="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center font-bold">
                  🎉
                </div>
              </div>
              <div className="font-mono text-3xl font-extrabold text-slate-900 mt-4">1</div>
              <span className="text-xs text-slate-400 mt-1 block">Offers received</span>
            </div>

            {/* Row 2: Recent Applications (lg:col-span-7) */}
            <div className="col-span-1 lg:col-span-7 bg-white/80 backdrop-blur-xl border border-slate-200/60 rounded-2xl p-6 shadow-sm flex flex-col min-h-[320px]">
              <div>
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-lg font-bold text-slate-900">Recent Applications</h3>
                  <Link to="/candidate/applications" className="text-xs font-semibold text-blue-600 hover:text-blue-700">
                    View all →
                  </Link>
                </div>

                {(!recentApps || recentApps.length === 0) ? (
                  <div className="flex flex-col items-center justify-center py-10 text-center flex-1">
                    <div className="text-4xl mb-4">🔍</div>
                    <h4 className="text-slate-900 font-bold text-lg mb-1">No Applications Yet</h4>
                    <p className="text-slate-500 text-sm max-w-xs mb-6">Start applying to jobs to track them here</p>
                    <Link
                      to="/gov-jobs"
                      className="px-6 py-2 bg-blue-600 text-white text-sm font-semibold rounded-xl hover:bg-blue-700 transition-all shadow-md"
                    >
                      Browse Jobs →
                    </Link>
                  </div>
                ) : (
                  <div className="space-y-3">
                    {(recentApps || []).slice(0, 5).map((app, idx) => (
                      <div
                        key={app?.id ?? idx}
                        className="flex flex-col sm:flex-row sm:items-center justify-between p-4 rounded-xl bg-slate-50 border border-slate-100 hover:bg-slate-100/60 transition-all gap-3"
                      >
                        <div className="overflow-hidden">
                          <h4 className="font-semibold text-slate-900 text-sm truncate">{app?.job_title ?? 'Software Developer'}</h4>
                          <p className="text-xs text-slate-500 truncate">{app?.company ?? 'Tech Company'}</p>
                        </div>
                        <div className="flex items-center justify-between sm:justify-end gap-3 shrink-0">
                          <span className={`text-[10px] px-2.5 py-1 rounded-full border font-bold uppercase tracking-wider ${getStatusBadge(app?.status)}`}>
                            {app?.status ?? 'APPLIED'}
                          </span>
                          <span className="text-[10px] text-slate-400 font-mono">
                            {app?.created_at ? new Date(app.created_at).toLocaleDateString() : 'Recent'}
                          </span>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>

            {/* Row 2: AI Job Matches Preview (lg:col-span-5) */}
            <div className="col-span-1 lg:col-span-5 bg-white/80 backdrop-blur-xl border border-slate-200/60 rounded-2xl p-6 shadow-sm flex flex-col min-h-[320px]">
              <div>
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-lg font-bold text-slate-900 flex items-center gap-2">
                    <span>✨</span> AI Job Matches
                  </h3>
                  <Link to="/candidate/ai-matches" className="text-xs font-semibold text-blue-600 hover:text-blue-700">
                    See all matches →
                  </Link>
                </div>

                <div className="space-y-3">
                  {(aiMatches || []).map((item) => (
                    <div
                      key={item?.id ?? Math.random()}
                      className="p-4 rounded-xl bg-slate-50 border border-slate-100 flex items-center justify-between gap-3"
                    >
                      <div className="overflow-hidden">
                        <div className="flex items-center gap-2 mb-1">
                          <span className="text-[10px] font-bold text-blue-600 bg-blue-100 px-2 py-0.5 rounded uppercase">
                            {item?.match ?? '0%'}
                          </span>
                          <h4 className="font-semibold text-slate-900 text-sm truncate">{item?.title ?? 'Job Title'}</h4>
                        </div>
                        <p className="text-xs text-slate-500 truncate">{item?.company ?? 'Company'}</p>
                      </div>
                      <button className="px-3 py-1.5 bg-blue-600 text-white text-[11px] font-bold rounded-lg hover:bg-blue-700 transition-all shadow-sm shrink-0 uppercase tracking-wide">
                        Apply
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* Row 3: Platform Counter (lg:col-span-12) */}
            <div className="col-span-1 lg:col-span-12 bg-white/80 backdrop-blur-xl border border-slate-200/60 rounded-2xl p-6 shadow-sm">
              <h3 className="text-xs font-bold text-slate-500 uppercase tracking-wider mb-6 text-center">
                Platform Statistics
              </h3>
              <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 md:gap-6">
                <AnimatedCounter value={platformStats?.totalJobs ?? 0} label="Total Jobs" />
                <AnimatedCounter value={platformStats?.companies ?? 0} label="Active Companies" />
                <AnimatedCounter value={platformStats?.candidates ?? 0} label="Candidates" />
                <AnimatedCounter value={platformStats?.placements ?? 0} label="Placements" />
              </div>
            </div>
          </div>
        )}
      </main>
    </div>
  );
};

export default CandidateDashboard;