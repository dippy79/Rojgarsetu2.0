import React, { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { useAuth } from '../../hooks/useAuth';
import {
  LayoutDashboard, User, FileText, Bookmark, Sparkles,
  LogOut, Menu, X, Briefcase, MapPin, Calendar,
  Search, ArrowRight, Loader2, Filter, AlertCircle, Award
} from 'lucide-react';

export const MyApplications = () => {
  const { user, logout } = useAuth();
  const router = useRouter();

  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [applications, setApplications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('ALL');
  const [selectedApp, setSelectedApp] = useState(null);

  const API_BASE = process.env.REACT_APP_BACKEND_URL || 'http://localhost:3001';

  const fetchApplications = useCallback(async () => {
    setLoading(true);
    setError(null);
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');

    try {
      const res = await fetch(`${API_BASE}/api/v1/applications/me`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': token ? `Bearer ${token}` : '',
        },
      });

      if (res.status === 401) {
        localStorage.removeItem('rojgar_token');
        router.push('/login');
        return;
      }

      if (res.ok) {
        const data = await res.json();
        setApplications(data?.applications ?? data?.data ?? data ?? []);
      }
    } catch (err) {
      setError('Neural link synchronization failed. Displaying cached applications.');
    } finally {
      setLoading(false);
    }
  }, [API_BASE, router]);

  useEffect(() => {
    fetchApplications();
  }, [fetchApplications]);

  const filteredApps = (applications || []).filter((app) => {
    const matchesSearch =
      (app?.job_title || '').toLowerCase().includes(searchTerm.toLowerCase()) ||
      (app?.company_name || '').toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus =
      statusFilter === 'ALL' ? true : (app?.status || '').toUpperCase() === statusFilter.toUpperCase();
    return matchesSearch && matchesStatus;
  });

  return (
    <div className="flex min-h-screen bg-[#FBFBFB] font-sans">
      {/* Sidebar - Consistent Premium Style */}
      <aside className={`
        fixed inset-y-0 left-0 z-50 w-72 bg-slate-950 text-white transform transition-transform duration-500
        md:translate-x-0 md:static md:h-screen sticky top-0
        ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
        flex flex-col p-8 shadow-[20px_0_60px_-15px_rgba(0,0,0,0.3)]
      `}>
        <div className="flex items-center justify-between mb-12">
          <div className="flex items-center gap-3">
            <div className="p-2.5 bg-blue-600 rounded-2xl">
              <Briefcase className="w-6 h-6 text-white" />
            </div>
            <span className="font-black text-2xl tracking-tighter uppercase">Rojgar<span className="text-blue-500">Setu</span></span>
          </div>
          <button onClick={() => setSidebarOpen(false)} className="p-2 text-slate-400 hover:text-white md:hidden">
            <X className="w-6 h-6" />
          </button>
        </div>

        <nav className="space-y-2 flex-1">
          {[
            { label: 'Dashboard', path: '/dashboard/candidate', icon: LayoutDashboard },
            { label: 'My Profile', path: '/candidate/profile', icon: User },
            { label: 'Applications', path: '/candidate/applications', icon: FileText },
            { label: 'Saved Jobs', path: '/candidate/saved-jobs', icon: Bookmark },
            { label: 'AI Matches', path: '/candidate/ai-matches', icon: Sparkles },
          ].map((item) => (
            <Link key={item.path} href={item.path} className={`flex items-center gap-4 px-6 py-4 rounded-2xl text-sm font-bold transition-all ${router.pathname === item.path ? 'bg-blue-600 text-white' : 'text-slate-400 hover:text-white hover:bg-white/5'}`}>
              <item.icon className="w-5 h-5" /> {item.label}
            </Link>
          ))}
        </nav>

        <button onClick={() => { logout(); router.push('/login'); }} className="flex items-center gap-4 w-full px-6 py-4 text-sm font-bold text-rose-400 hover:bg-rose-500/10 rounded-2xl transition-all mt-auto border border-rose-500/10">
          <LogOut className="w-5 h-5" /> Logout
        </button>
      </aside>

      {/* Main Content */}
      <main className="flex-1 overflow-auto p-6 md:p-12 lg:p-16">
        {error && (
          <div className="mb-10 p-5 bg-amber-50 border border-amber-100 rounded-3xl flex flex-col md:flex-row items-center justify-between gap-4 text-amber-700">
             <div className="flex items-center gap-4">
                <AlertCircle className="w-6 h-6 shrink-0" />
                <p className="text-sm font-bold">{error}</p>
             </div>
             <button onClick={fetchApplications} className="px-6 py-2.5 bg-amber-600 text-white text-[10px] font-black uppercase tracking-widest rounded-xl hover:bg-amber-700 transition-all shadow-lg shadow-amber-600/20">Retry Sync</button>
          </div>
        )}

        <header className="flex flex-col md:flex-row md:items-center justify-between gap-8 mb-16">
          <div className="space-y-2">
            <div className="flex items-center gap-4">
              <button onClick={() => setSidebarOpen(true)} className="p-3 bg-white border border-slate-200 rounded-2xl md:hidden shadow-sm">
                <Menu className="w-5 h-5" />
              </button>
              <h1 className="text-4xl md:text-5xl font-black text-slate-900 tracking-tight tracking-tight">Active Pipelines.</h1>
            </div>
            <p className="text-slate-400 text-lg font-medium">Tracking {applications.length} deployments through the <span className="text-blue-600 font-bold">Hiring Cycle.</span></p>
          </div>
        </header>

        {/* Dynamic Controls */}
        <div className="bg-white border border-slate-200 p-6 rounded-[2.5rem] shadow-sm flex flex-col lg:flex-row items-center gap-8 mb-12">
           <div className="flex-1 relative w-full group">
              <Search className="absolute left-5 top-4.5 w-5 h-5 text-slate-400 group-focus-within:text-blue-600 transition-colors" />
              <input
                type="text"
                placeholder="Search by job title or company..."
                value={searchTerm}
                onChange={e => setSearchTerm(e.target.value)}
                className="w-full pl-14 pr-6 py-4.5 bg-slate-50 border-none rounded-2xl text-sm font-bold focus:ring-4 focus:ring-blue-500/10 transition-all outline-none"
              />
           </div>

           <div className="flex bg-slate-100 p-1.5 rounded-2xl border border-slate-200 gap-1 shrink-0 overflow-x-auto max-w-full">
              {['ALL', 'APPLIED', 'SHORTLISTED', 'INTERVIEW', 'REJECTED'].map(f => (
                <button
                  key={f}
                  onClick={() => setStatusFilter(f)}
                  className={`px-5 py-2 rounded-xl text-[9px] font-black uppercase tracking-widest transition-all ${statusFilter === f ? 'bg-white text-blue-600 shadow-sm border border-slate-200' : 'text-slate-400 hover:text-slate-600'}`}
                >
                  {f}
                </button>
              ))}
           </div>
        </div>

        {loading ? (
          <div className="space-y-6">
             {[1,2,3].map(i => (
               <div key={i} className="h-44 bg-white border border-slate-200 rounded-[2.5rem] animate-pulse"></div>
             ))}
          </div>
        ) : filteredApps.length === 0 ? (
          <div className="bg-white border border-slate-200 rounded-[3rem] p-24 text-center space-y-8 shadow-sm">
             <div className="w-28 h-28 bg-slate-50 rounded-full flex items-center justify-center mx-auto text-5xl shadow-inner">🛰️</div>
             <div>
                <h3 className="text-2xl font-black text-slate-900 tracking-tight">No Signals Detected</h3>
                <p className="text-slate-400 text-sm font-medium mt-1">Start broadcasting your profile to top firms to see pipelines here.</p>
             </div>
             <Link href="/gov-jobs" className="px-10 py-4 bg-slate-900 text-white font-black rounded-[2rem] text-xs uppercase tracking-[0.2em] hover:bg-blue-600 transition-all shadow-2xl shadow-slate-900/10 inline-block">Explore Roles →</Link>
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-6">
            {filteredApps.map((app) => (
              <div key={app.id} className="group bg-white border border-slate-200 rounded-[2.5rem] p-10 hover:border-blue-400 hover:shadow-2xl hover:shadow-blue-500/5 transition-all relative overflow-hidden flex flex-col lg:flex-row items-center justify-between gap-10">
                <div className="absolute top-0 right-0 w-32 h-32 bg-blue-500/5 rounded-full -mr-16 -mt-16 blur-3xl group-hover:scale-150 transition-transform duration-700"></div>

                <div className="flex-1 flex flex-col md:flex-row items-center gap-8 relative z-10">
                   <div className="w-20 h-20 bg-slate-50 border border-slate-100 rounded-3xl flex items-center justify-center text-2xl font-black text-slate-300 group-hover:bg-blue-50 group-hover:border-blue-100 group-hover:text-blue-600 transition-all uppercase">
                      {(app?.company_name || app?.company || 'C')[0]}
                   </div>

                   <div className="space-y-4 text-center md:text-left">
                      <div className="space-y-1">
                         <div className="flex flex-wrap items-center justify-center md:justify-start gap-3">
                            <h3 className="text-3xl font-black text-slate-900 group-hover:text-blue-600 transition-colors tracking-tighter leading-tight">{app?.job_title}</h3>
                            <div className={`px-4 py-1.5 rounded-full text-[9px] font-black uppercase tracking-widest border ${getStatusColor(app?.status)} shadow-sm`}>
                               {app?.status}
                            </div>
                         </div>
                         <p className="text-sm font-bold text-slate-400 uppercase tracking-widest">{app?.company_name}</p>
                      </div>

                      <div className="flex flex-wrap items-center justify-center md:justify-start gap-6 text-[10px] font-black text-slate-400 uppercase tracking-widest">
                         <div className="flex items-center gap-2">
                            <MapPin className="w-3.5 h-3.5" />
                            {app?.location || 'Remote'}
                         </div>
                         <div className="flex items-center gap-2">
                            <Calendar className="w-3.5 h-3.5" />
                            Applied: {app?.applied_date ? new Date(app.applied_date).toLocaleDateString() : 'Active'}
                         </div>
                         {app?.match_score && (
                           <div className="px-3 py-1 bg-indigo-50 text-indigo-600 rounded-lg border border-indigo-100 flex items-center gap-2">
                              <Sparkles className="w-3 h-3" />
                              {app.match_score}% Fit
                           </div>
                         )}
                      </div>
                   </div>
                </div>

                <div className="flex flex-col md:flex-row lg:flex-col gap-4 relative z-10 shrink-0 w-full md:w-auto">
                   <button onClick={() => setSelectedApp(app)} className="px-8 py-4 bg-slate-900 text-white font-black rounded-2xl hover:bg-blue-600 transition-all shadow-xl shadow-slate-900/10 uppercase text-[10px] tracking-widest flex items-center justify-center gap-3 group">
                      Timeline Data <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
                   </button>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Timeline Overlay */}
        {selectedApp && (
          <div className="fixed inset-0 bg-slate-950/60 backdrop-blur-md z-[100] flex items-center justify-center p-6 animate-in fade-in duration-300" onClick={() => setSelectedApp(null)}>
            <div className="bg-white w-full max-w-2xl rounded-[3rem] p-12 shadow-2xl relative overflow-hidden" onClick={e => e.stopPropagation()}>
               <div className="absolute top-0 left-0 w-full h-2 bg-gradient-to-r from-blue-600 via-indigo-600 to-purple-600"></div>

               <button onClick={() => setSelectedApp(null)} className="absolute top-8 right-8 p-3 bg-slate-50 rounded-2xl hover:bg-slate-100 transition-all"><X className="w-5 h-5 text-slate-900" /></button>

               <div className="mb-12">
                  <h2 className="text-3xl font-black text-slate-900 tracking-tight">{selectedApp.job_title}</h2>
                  <p className="text-slate-400 font-bold uppercase tracking-widest mt-2">{selectedApp.company_name}</p>
               </div>

               <div className="space-y-8 relative">
                  <div className="absolute left-4 top-2 bottom-2 w-px bg-slate-100"></div>

                  {[
                    { status: 'Applied', date: selectedApp.applied_date, completed: true, icon: FileText },
                    { status: 'Profile Review', date: 'Processing', completed: selectedApp.status !== 'Applied', icon: User },
                    { status: 'Assessment', date: 'Awaiting', completed: false, icon: Award },
                  ].map((step, idx) => (
                    <div key={idx} className="flex gap-8 items-start relative z-10">
                       <div className={`w-8 h-8 rounded-full flex items-center justify-center shrink-0 shadow-lg ${step.completed ? 'bg-blue-600 text-white shadow-blue-500/20' : 'bg-white border-2 border-slate-100 text-slate-300'}`}>
                          <step.icon className="w-3.5 h-3.5" />
                       </div>
                       <div className="pt-1">
                          <p className={`text-sm font-black uppercase tracking-widest ${step.completed ? 'text-slate-900' : 'text-slate-300'}`}>{step.status}</p>
                          <p className="text-[10px] font-bold text-slate-400 mt-1 uppercase tracking-tighter">{step.date}</p>
                       </div>
                    </div>
                  ))}
               </div>

               <div className="mt-12 pt-8 border-t border-slate-50">
                  <button onClick={() => setSelectedApp(null)} className="w-full py-4 bg-slate-50 text-slate-900 font-black rounded-2xl hover:bg-slate-100 transition-all uppercase text-[10px] tracking-widest">Acknowledge Interface</button>
               </div>
            </div>
          </div>
        )}
      </main>
    </div>
  );
};

const getStatusColor = (s) => {
  switch((s || '').toUpperCase()) {
    case 'SHORTLISTED': return 'bg-indigo-50 text-indigo-600 border-indigo-100';
    case 'INTERVIEW': return 'bg-amber-50 text-amber-600 border-amber-100';
    case 'REJECTED': return 'bg-rose-50 text-rose-600 border-rose-100';
    default: return 'bg-blue-50 text-blue-600 border-blue-100';
  }
};

export default MyApplications;
