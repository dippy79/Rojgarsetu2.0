import React, { useState, useEffect } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';

// Icon Components (Inline SVG to avoid missing package dependencies)
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
      <div className="font-mono text-3xl font-extrabold text-slate-900">{count.toLocaleString()}+</div>
      <div className="text-xs font-semibold text-slate-500 uppercase tracking-wider mt-1">{label}</div>
    </div>
  );
};

export const CandidateDashboard = () => {
  const { logout, user } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const [candidateProfile, setCandidateProfile] = useState(null);
  const [recentApps, setRecentApps] = useState([]);
  const [platformStats, setPlatformStats] = useState({
    totalJobs: 1250,
    companies: 450,
    candidates: 8900,
    placements: 3200,
  });

  const API_BASE = 'http://localhost:3001';

  const fetchData = async () => {
    setLoading(true);
    setError(false);
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');

    try {
      const headers = {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      };

      const [profileRes, appsRes, statsRes] = await Promise.allSettled([
        fetch(`${API_BASE}/api/v1/candidates/me`, { headers }),
        fetch(`${API_BASE}/api/v1/candidates/me/applications?limit=5`, { headers }),
        fetch(`${API_BASE}/api/v1/stats`),
      ]);

      if (profileRes.status === 'fulfilled' && profileRes.value.ok) {
        const pData = await profileRes.value.json();
        setCandidateProfile(pData.candidate || pData);
      } else {
        setCandidateProfile({
          full_name: user?.name || user?.email?.split('@')[0] || 'Candidate User',
          profile_completion: 75,
          skills: ['React.js', 'Node.js', 'Tailwind CSS', 'JavaScript'],
        });
      }

      if (appsRes.status === 'fulfilled' && appsRes.value.ok) {
        const aData = await appsRes.value.json();
        setRecentApps(aData.applications || aData || []);
      } else {
        setRecentApps([
          { id: '1', job_title: 'Frontend Engineer', company: 'TechCorp India', status: 'SHORTLISTED', created_at: '2026-08-10' },
          { id: '2', job_title: 'Full Stack Developer', company: 'Innovate Labs', status: 'APPLIED', created_at: '2026-08-11' },
          { id: '3', job_title: 'React Native Developer', company: 'MobileFirst Systems', status: 'INTERVIEW', created_at: '2026-08-08' },
        ]);
      }

      if (statsRes.status === 'fulfilled' && statsRes.value.ok) {
        const sData = await statsRes.value.json();
        setPlatformStats({
          totalJobs: sData.totalJobs || 1250,
          companies: sData.companies || 450,
          candidates: sData.candidates || 8900,
          placements: sData.placements || 3200,
        });
      }
    } catch (err) {
      setError(true);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

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
      {/* Sidebar */}
      <aside className="w-64 bg-slate-900 text-white h-screen sticky top-0 flex flex-col justify-between p-6 shadow-xl z-20">
        <div>
          {/* User Avatar & Info */}
          <div className="flex items-center gap-3 mb-6">
            <div className="bg-blue-600 w-12 h-12 rounded-full flex items-center justify-center font-bold text-white text-lg shadow-md">
              {getInitials(candidateProfile?.full_name)}
            </div>
            <div className="overflow-hidden">
              <h2 className="font-semibold text-white truncate text-base">
                {candidateProfile?.full_name || 'Candidate'}
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
              <span>{candidateProfile?.profile_completion || 75}%</span>
            </div>
            <div className="w-full bg-slate-700 rounded-full h-2 overflow-hidden">
              <div
                className="bg-blue-500 h-2 rounded-full transition-all duration-500"
                style={{ width: `${candidateProfile?.profile_completion || 75}%` }}
              ></div>
            </div>
          </div>

          {/* Navigation Links */}
          <nav className="space-y-1">
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = location.pathname === item.path;
              return (
                <Link
                  key={item.path}
                  to={item.path}
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
          className="flex items-center w-full px-4 py-3 text-sm font-medium bg-red-500/10 text-red-400 hover:bg-red-500/20 rounded-lg transition-all"
        >
          <LogOutIcon />
          Logout
        </button>
      </aside>

      {/* Main Content Area */}
      <main className="flex-1 overflow-auto p-8 bg-slate-50">
        {/* Header */}
        <header className="mb-8">
          <h1 className="text-3xl font-bold text-slate-900 tracking-tight">
            Good morning, {candidateProfile?.full_name?.split(' ')[0] || 'Candidate'} 👋
          </h1>
          <p className="text-slate-500 text-sm mt-1 font-medium">
            {new Date().toLocaleDateString('en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}
          </p>
        </header>

        {/* Error Banner */}
        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-2xl flex items-center justify-between text-red-700">
            <span className="text-sm font-medium">Failed to load live candidate data — showing default preview.</span>
            <button
              onClick={fetchData}
              className="px-3 py-1 bg-red-600 text-white text-xs font-semibold rounded-lg hover:bg-red-700 transition-all"
            >
              Retry
            </button>
          </div>
        )}

        {/* Loading Skeleton */}
        {loading ? (
          <div className="space-y-6 animate-pulse">
            <div className="grid grid-cols-12 gap-6">
              {[1, 2, 3, 4].map((i) => (
                <div key={i} className="col-span-3 h-32 bg-slate-200 rounded-2xl"></div>
              ))}
            </div>
            <div className="grid grid-cols-12 gap-6">
              <div className="col-span-7 h-64 bg-slate-200 rounded-2xl"></div>
              <div className="col-span-5 h-64 bg-slate-200 rounded-2xl"></div>
            </div>
            <div className="h-32 bg-slate-200 rounded-2xl"></div>
          </div>
        ) : (
          /* Bento Grid Layout */
          <div className="grid grid-cols-12 gap-6">
            {/* Row 1: Stats Cards */}
            <div className="col-span-3 bg-white/80 backdrop-blur-xl border border-slate-200/60 rounded-2xl p-6 shadow-sm">
              <div className="flex items-center justify-between">
                <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Applied</span>
                <div className="w-10 h-10 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center font-bold">
                  📄
                </div>
              </div>
              <div className="font-mono text-3xl font-extrabold text-slate-900 mt-4">12</div>
              <span className="text-xs text-slate-400 mt-1 block">Total submitted</span>
            </div>

            <div className="col-span-3 bg-white/80 backdrop-blur-xl border border-slate-200/60 rounded-2xl p-6 shadow-sm">
              <div className="flex items-center justify-between">
                <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Shortlisted</span>
                <div className="w-10 h-10 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center font-bold">
                  ⭐
                </div>
              </div>
              <div className="font-mono text-3xl font-extrabold text-slate-900 mt-4">4</div>
              <span className="text-xs text-slate-400 mt-1 block">Under review</span>
            </div>

            <div className="col-span-3 bg-white/80 backdrop-blur-xl border border-slate-200/60 rounded-2xl p-6 shadow-sm">
              <div className="flex items-center justify-between">
                <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Interview</span>
                <div className="w-10 h-10 rounded-xl bg-violet-50 text-violet-600 flex items-center justify-center font-bold">
                  🎙️
                </div>
              </div>
              <div className="font-mono text-3xl font-extrabold text-slate-900 mt-4">2</div>
              <span className="text-xs text-slate-400 mt-1 block">Scheduled calls</span>
            </div>

            <div className="col-span-3 bg-white/80 backdrop-blur-xl border border-slate-200/60 rounded-2xl p-6 shadow-sm">
              <div className="flex items-center justify-between">
                <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Offered</span>
                <div className="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center font-bold">
                  🎉
                </div>
              </div>
              <div className="font-mono text-3xl font-extrabold text-slate-900 mt-4">1</div>
              <span className="text-xs text-slate-400 mt-1 block">Offers received</span>
            </div>

            {/* Row 2: Recent Applications (col-span-7) */}
            <div className="col-span-7 bg-white/80 backdrop-blur-xl border border-slate-200/60 rounded-2xl p-6 shadow-sm flex flex-col justify-between">
              <div>
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-lg font-bold text-slate-900">Recent Applications</h3>
                  <Link to="/candidate/applications" className="text-xs font-semibold text-blue-600 hover:text-blue-700">
                    View all →
                  </Link>
                </div>

                {(!recentApps || recentApps.length === 0) ? (
                  <div className="text-center py-8">
                    <div className="text-4xl mb-2">📥</div>
                    <p className="text-slate-500 text-sm font-medium">No recent applications found</p>
                  </div>
                ) : (
                  <div className="space-y-3">
                    {(recentApps || []).slice(0, 5).map((app, idx) => (
                      <div
                        key={app.id || idx}
                        className="flex items-center justify-between p-3 rounded-xl bg-slate-50 border border-slate-100 hover:bg-slate-100/60 transition-all"
                      >
                        <div>
                          <h4 className="font-semibold text-slate-900 text-sm">{app.job_title || 'Software Developer'}</h4>
                          <p className="text-xs text-slate-500">{app.company || 'Tech Company'}</p>
                        </div>
                        <div className="flex items-center gap-3">
                          <span className={`text-xs px-2.5 py-1 rounded-full border font-semibold ${getStatusBadge(app.status)}`}>
                            {app.status || 'APPLIED'}
                          </span>
                          <span className="text-xs text-slate-400 font-mono">
                            {app.created_at ? new Date(app.created_at).toLocaleDateString() : 'Recent'}
                          </span>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>

            {/* Row 2: AI Job Matches Preview (col-span-5) */}
            <div className="col-span-5 bg-white/80 backdrop-blur-xl border border-slate-200/60 rounded-2xl p-6 shadow-sm flex flex-col justify-between">
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
                  {aiMatches.map((item) => (
                    <div
                      key={item.id}
                      className="p-3 rounded-xl bg-slate-50 border border-slate-100 flex items-center justify-between"
                    >
                      <div>
                        <div className="flex items-center gap-2">
                          <span className="text-xs font-bold text-blue-600 bg-blue-100 px-2 py-0.5 rounded-md">
                            {item.match}
                          </span>
                          <h4 className="font-semibold text-slate-900 text-sm truncate max-w-[140px]">{item.title}</h4>
                        </div>
                        <p className="text-xs text-slate-500 mt-1">{item.company}</p>
                      </div>
                      <button className="px-3 py-1.5 bg-blue-600 text-white text-xs font-semibold rounded-lg hover:bg-blue-700 transition-all shadow-sm">
                        Apply
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* Row 3: Platform Counter (col-span-12) */}
            <div className="col-span-12 bg-white/80 backdrop-blur-xl border border-slate-200/60 rounded-2xl p-6 shadow-sm">
              <h3 className="text-sm font-bold text-slate-500 uppercase tracking-wider mb-4 text-center">
                Platform Statistics
              </h3>
              <div className="grid grid-cols-4 gap-4">
                <AnimatedCounter value={platformStats.totalJobs} label="Total Jobs" />
                <AnimatedCounter value={platformStats.companies} label="Active Companies" />
                <AnimatedCounter value={platformStats.candidates} label="Registered Candidates" />
                <AnimatedCounter value={platformStats.placements} label="Successful Placements" />
              </div>
            </div>
          </div>
        )}
      </main>
    </div>
  );
};

export default CandidateDashboard;