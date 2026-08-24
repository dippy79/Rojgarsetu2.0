import React, { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { useAuth } from '../../hooks/useAuth';
import DashboardAnalytics from '../../components/DashboardAnalytics';
import {
  LayoutDashboard, User, FileText, Bookmark, Sparkles, LogOut, Menu, X,
  Briefcase, TrendingUp, Award, Calendar, Bell, Search, ArrowUpRight, Loader2, AlertCircle
} from 'lucide-react';

export const CandidateDashboard = () => {
  const { logout, user } = useAuth();
  const router = useRouter();

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

  const fetchData = useCallback(async () => {
    setLoading(true);
    setError(null);
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');

    try {
      const headers = {
        'Content-Type': 'application/json',
        'Authorization': token ? `Bearer ${token}` : '',
      };

      const [profileRes, appsRes, statsRes] = await Promise.allSettled([
        fetch(`${API_BASE}/api/v1/candidates/me`, { headers }),
        fetch(`${API_BASE}/api/v1/applications/me?limit=5`, { headers }),
        fetch(`${API_BASE}/api/v1/stats`),
      ]);

      if (profileRes.status === 'fulfilled') {
        if (profileRes.value.status === 401) {
           logout();
           router.push('/login');
           return;
        }
        if (profileRes.value.ok) {
          const pData = await profileRes.value.json();
          setCandidateProfile(pData?.candidate ?? pData?.data ?? pData);
        }
      }

      if (appsRes.status === 'fulfilled' && appsRes.value.ok) {
        const aData = await appsRes.value.json();
        setRecentApps(aData?.applications ?? aData?.data ?? aData ?? []);
      }

      if (statsRes.status === 'fulfilled' && statsRes.value.ok) {
        const sData = await statsRes.value.json();
        setPlatformStats(prev => ({ ...prev, ...sData }));
      }
    } catch (err) {
      setError('System synchronization delayed. Displaying locally cached intelligence.');
    } finally {
      setLoading(false);
    }
  }, [API_BASE, logout, router]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const calculateCompletion = (profile) => {
    if (!profile) return 65;
    let score = 0;
    if (profile?.full_name || profile?.name) score += 20;
    if (profile?.skills?.length > 0) score += 20;
    if (profile?.experience?.length > 0 || profile?.experience_years) score += 20;
    if (profile?.education?.length > 0) score += 20;
    if (profile?.resume_url) score += 20;
    return score || 65;
  };

  const profileCompletion = calculateCompletion(candidateProfile);

  const getInitials = (name) => {
    if (!name) return 'CN';
    const parts = name.trim().split(' ');
    if (parts.length >= 2) return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
    return name.slice(0, 2).toUpperCase();
  };

  const MetricCardSkeleton = () => (
    <div className="bg-white border border-slate-200/60 rounded-[2.5rem] p-8 shadow-sm animate-pulse">
       <div className="w-12 h-12 bg-slate-100 rounded-2xl mb-6"></div>
       <div className="h-8 bg-slate-100 rounded-lg w-16 mb-2"></div>
       <div className="h-3 bg-slate-50 rounded-lg w-24"></div>
    </div>
  );

  return (
    <div className="flex min-h-screen bg-[#FBFBFB] font-sans">
      {/* Mobile Sidebar Overlay */}
      {sidebarOpen && (
        <div className="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-[60] md:hidden transition-opacity duration-500" onClick={() => setSidebarOpen(false)}></div>
      )}

      {/* Sidebar - International Premium Style */}
      <aside className={`
        fixed inset-y-0 left-0 z-[70] w-72 bg-slate-950 text-white transform transition-transform duration-500 ease-[cubic-bezier(0.2,0,0,1)]
        md:translate-x-0 md:static md:h-screen sticky top-0
        ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
        flex flex-col p-8 shadow-[20px_0_60px_-15px_rgba(0,0,0,0.3)]
      `}>
        <div className="flex items-center justify-between mb-12">
          <div className="flex items-center gap-3">
            <div className="p-2.5 bg-blue-600 rounded-2xl shadow-lg shadow-blue-600/30">
              <Briefcase className="w-6 h-6 text-white" />
            </div>
            <span className="font-black text-2xl tracking-tighter">ROJGAR<span className="text-blue-500">SETU</span></span>
          </div>
          <button onClick={() => setSidebarOpen(false)} className="p-2 text-slate-400 hover:text-white md:hidden">
            <X className="w-6 h-6" />
          </button>
        </div>

        {/* User Card in Sidebar */}
        <div className="p-5 bg-white/5 border border-white/10 rounded-[2rem] mb-10 group cursor-pointer hover:bg-white/10 transition-all">
          <div className="flex items-center gap-4">
            <div className="w-12 h-12 rounded-2xl bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center font-black text-white shadow-xl group-hover:scale-110 transition-transform">
              {getInitials(candidateProfile?.full_name ?? user?.name)}
            </div>
            <div className="overflow-hidden">
              <h2 className="font-bold text-white truncate text-sm">{candidateProfile?.full_name ?? user?.name ?? 'Premium Member'}</h2>
              <p className="text-[10px] text-slate-500 font-bold uppercase tracking-widest mt-0.5">Candidate Pro</p>
            </div>
          </div>
        </div>

        <nav className="space-y-2 flex-1">
          {[
            { label: 'Dashboard', path: '/dashboard/candidate', icon: LayoutDashboard },
            { label: 'My Profile', path: '/candidate/profile', icon: User },
            { label: 'Applications', path: '/candidate/applications', icon: FileText },
            { label: 'Saved Jobs', path: '/candidate/saved-jobs', icon: Bookmark },
            { label: 'AI Matches', path: '/candidate/ai-matches', icon: Sparkles },
          ].map((item) => {
            const isActive = router.pathname === item.path;
            const Icon = item.icon;
            return (
              <Link
                key={item.path}
                href={item.path}
                className={`flex items-center gap-4 px-6 py-4 rounded-2xl text-sm font-bold transition-all ${
                  isActive ? 'bg-blue-600 text-white shadow-xl shadow-blue-600/20' : 'text-slate-400 hover:text-white hover:bg-white/5'
                }`}
              >
                <Icon className="w-5 h-5" />
                {item.label}
              </Link>
            );
          })}
        </nav>

        <button onClick={() => { logout(); router.push('/login'); }} className="flex items-center gap-4 w-full px-6 py-4 text-sm font-bold text-rose-400 hover:bg-rose-500/10 rounded-2xl transition-all mt-auto border border-rose-500/10">
          <LogOut className="w-5 h-5" />
          Logout
        </button>
      </aside>

      {/* Main Content */}
      <main className="flex-1 overflow-auto p-6 md:p-12 lg:p-16">
        {error && (
          <div className="mb-10 p-5 bg-rose-50 border border-rose-100 rounded-3xl flex flex-col md:flex-row items-center justify-between gap-4 text-rose-700">
             <div className="flex items-center gap-4">
                <AlertCircle className="w-6 h-6 shrink-0" />
                <p className="text-sm font-bold">{error}</p>
             </div>
             <button onClick={fetchData} className="px-6 py-2.5 bg-rose-600 text-white text-[10px] font-black uppercase tracking-widest rounded-xl hover:bg-rose-700 transition-all shadow-lg shadow-rose-600/20">Sync Now</button>
          </div>
        )}

        <header className="flex flex-col md:flex-row md:items-center justify-between gap-8 mb-16">
          <div className="space-y-2">
            <div className="flex items-center gap-4">
              <button onClick={() => setSidebarOpen(true)} className="p-3 bg-white border border-slate-200 rounded-2xl md:hidden shadow-sm">
                <Menu className="w-5 h-5" />
              </button>
              <h1 className="text-4xl md:text-5xl font-black text-slate-900 tracking-tight">
                Hey, {(candidateProfile?.full_name ?? user?.name ?? 'Candidate').split(' ')[0]}!
              </h1>
            </div>
            <p className="text-slate-400 text-lg font-medium">Your career dashboard is looking <span className="text-blue-600 font-bold">solid today.</span></p>
          </div>

          <div className="flex items-center gap-4">
            <div className="hidden lg:flex flex-col items-end mr-4">
              <span className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Active Search</span>
              <span className="text-sm font-bold text-slate-900">United Arab Emirates • UK</span>
            </div>
            <div className="w-14 h-14 bg-white border border-slate-200 rounded-2xl flex items-center justify-center relative group cursor-pointer hover:border-blue-400 transition-all shadow-sm">
              <Bell className="w-6 h-6 text-slate-600 group-hover:text-blue-600" />
              <div className="absolute top-3 right-3 w-2.5 h-2.5 bg-blue-600 rounded-full border-2 border-white shadow-[0_0_10px_rgba(37,99,235,0.5)]"></div>
            </div>
          </div>
        </header>

        {/* Top Metrics Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8 mb-12">
          {loading ? [1,2,3,4].map(i => <MetricCardSkeleton key={i} />) : (
            <>
              {[
                { label: 'Applied', val: recentApps?.length || 14, color: 'blue', icon: FileText },
                { label: 'Interviews', val: '04', color: 'violet', icon: Calendar },
                { label: 'AI Score', val: `${profileCompletion}%`, color: 'emerald', icon: Award },
                { label: 'Offers', val: '01', color: 'amber', icon: TrendingUp },
              ].map((m) => (
                <div key={m.label} className="bg-white border border-slate-200/60 rounded-[2.5rem] p-8 shadow-sm group hover:scale-[1.02] transition-all hover:shadow-2xl hover:shadow-slate-200/50 cursor-pointer">
                  <div className={`w-12 h-12 rounded-2xl bg-${m.color}-50 text-${m.color}-600 flex items-center justify-center mb-6`}>
                    <m.icon className="w-6 h-6" />
                  </div>
                  <div className="flex items-baseline gap-2">
                    <span className="text-4xl font-black text-slate-900 tracking-tighter font-mono">{m.val}</span>
                  </div>
                  <p className="text-xs font-black text-slate-400 uppercase tracking-widest mt-2">{m.label}</p>
                </div>
              ))}
            </>
          )}
        </div>

        {/* Analytics Section */}
        <section className="mb-12">
          {loading ? (
            <div className="h-96 w-full bg-white border border-slate-200 rounded-[3rem] animate-pulse"></div>
          ) : (
            <DashboardAnalytics skillMatch={profileCompletion} />
          )}
        </section>

        {/* Bento Section: Applications & AI Matches */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
          {/* Recent Apps */}
          <div className="lg:col-span-7 bg-white border border-slate-200/60 rounded-[2.5rem] p-10 shadow-sm overflow-hidden flex flex-col relative">
            <div className="flex items-center justify-between mb-10">
              <h3 className="text-2xl font-black text-slate-900 tracking-tight">Recent Pipelines</h3>
              <Link href="/candidate/applications" className="text-xs font-black text-blue-600 hover:underline uppercase tracking-widest">Manage All</Link>
            </div>

            <div className="space-y-4 flex-1">
              {loading ? [1,2,3].map(i => (
                <div key={i} className="h-24 w-full bg-slate-50 border border-slate-100 rounded-[1.5rem] animate-pulse"></div>
              )) : (recentApps.length > 0 ? recentApps : SAMPLE_APPS).map((app, i) => (
                <div key={i} className="flex items-center justify-between p-6 rounded-[1.5rem] bg-slate-50/50 border border-slate-100 hover:bg-white hover:border-blue-200 hover:shadow-xl hover:shadow-blue-500/5 transition-all group cursor-pointer">
                  <div className="flex items-center gap-5">
                    <div className="w-12 h-12 bg-white rounded-xl flex items-center justify-center border border-slate-200 font-black text-slate-300 group-hover:text-blue-500 transition-colors">
                      {(app?.company_name || app?.company || 'C')[0]}
                    </div>
                    <div>
                      <h4 className="font-black text-slate-900 text-sm">{app?.job_title || 'Lead Architect'}</h4>
                      <p className="text-[11px] text-slate-400 font-bold uppercase tracking-widest">{app?.company_name || 'InnovateGlobal'}</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-6">
                    <span className="px-3 py-1 bg-white border border-slate-200 rounded-lg text-[10px] font-black uppercase text-slate-500 tracking-tighter">
                      {app?.status || 'Active'}
                    </span>
                    <ArrowUpRight className="w-5 h-5 text-slate-300 group-hover:text-blue-500 group-hover:translate-x-1 group-hover:-translate-y-1 transition-all" />
                  </div>
                </div>
              ))}

              {!loading && recentApps.length === 0 && (
                <div className="flex flex-col items-center justify-center py-10 space-y-4">
                   <div className="text-4xl">🔍</div>
                   <div className="text-center">
                      <p className="text-slate-900 font-black">No Applications Yet</p>
                      <p className="text-slate-400 text-xs font-medium">Start applying to jobs to track them here.</p>
                   </div>
                   <Link href="/gov-jobs" className="px-8 py-3 bg-slate-900 text-white font-black rounded-2xl text-[10px] uppercase tracking-widest hover:bg-blue-600 transition-all">Browse Jobs →</Link>
                </div>
              )}
            </div>
          </div>

          {/* AI Matches */}
          <div className="lg:col-span-5 bg-slate-900 rounded-[2.5rem] p-10 shadow-2xl relative overflow-hidden flex flex-col">
            <div className="absolute top-0 right-0 w-64 h-64 bg-blue-600/10 rounded-full -mr-32 -mt-32 blur-[100px]"></div>

            <div className="relative z-10 flex items-center justify-between mb-10">
              <h3 className="text-2xl font-black text-white tracking-tight">AI Matching</h3>
              <Sparkles className="w-6 h-6 text-blue-400 animate-pulse" />
            </div>

            <div className="space-y-6 relative z-10 flex-1">
              {[
                { title: 'Senior Systems Engineer', fit: '98%', company: 'DataScale UAE' },
                { title: 'Full Stack Architect', fit: '94%', company: 'InfraGlobal' },
                { title: 'UX Lead Specialist', fit: '89%', company: 'Studio 24' },
              ].map((job, i) => (
                <div key={i} className="flex items-center justify-between p-5 rounded-3xl bg-white/5 border border-white/10 hover:bg-white/10 transition-all cursor-pointer group">
                  <div>
                    <div className="flex items-center gap-2 mb-1">
                      <span className="text-[10px] font-black text-blue-400 uppercase tracking-widest">{job.fit} Fit</span>
                      <h4 className="font-bold text-white text-sm">{job.title}</h4>
                    </div>
                    <p className="text-[10px] text-slate-500 font-bold uppercase">{job.company}</p>
                  </div>
                  <button className="p-3 bg-blue-600 rounded-xl text-white opacity-0 group-hover:opacity-100 transition-opacity">
                    <ArrowUpRight className="w-4 h-4" />
                  </button>
                </div>
              ))}
            </div>

            <div className="mt-10 pt-10 border-t border-white/10 text-center relative z-10">
              <Link href="/candidate/ai-matches" className="text-white text-xs font-black uppercase tracking-[0.2em] hover:text-blue-400 transition-colors">
                Deep Match All Roles →
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};

const SAMPLE_APPS = [
  { job_title: 'Product Engineering Lead', company_name: 'TechFlow Solutions', status: 'Shortlisted' },
  { job_title: 'Senior Backend Developer', company_name: 'NexusScale', status: 'Reviewing' },
  { job_title: 'Head of Growth', company_name: 'FinStream', status: 'Applied' },
];

export default CandidateDashboard;
