import React, { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { useAuth } from '../../hooks/useAuth';
import DashboardAnalytics from '../../components/DashboardAnalytics';
import {
  Plus, Users, Briefcase, Zap, Search, Bell,
  Menu, X, LayoutDashboard, Settings, FileText,
  ChevronRight, ArrowUpRight, TrendingUp, Filter, Loader2
} from 'lucide-react';
import { apiUrl } from '../../apiConfig';

export const EmployerDashboard = () => {
  const { user, logout } = useAuth();
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState({
    activeJobs: 4,
    totalApplicants: 128,
    shortlisted: 18,
    aiMatched: 42
  });
  const [recentJobs, setRecentJobs] = useState([]);
  const [sidebarOpen, setSidebarOpen] = useState(false);

  const fetchData = useCallback(async () => {
    setLoading(true);
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');
    try {
      const res = await fetch(apiUrl('/api/v1/employer/dashboard'), {
        headers: {
          'Authorization': token ? `Bearer ${token}` : '',
          'Content-Type': 'application/json'
        },
      });

      if (res.ok) {
        const data = await res.json();
        setStats(data.stats || stats);
        setRecentJobs(data.jobs || []);
      } else {
        setRecentJobs(SAMPLE_EMPLOYER_JOBS);
      }
    } catch (err) {
      setRecentJobs(SAMPLE_EMPLOYER_JOBS);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return (
    <div className="flex min-h-screen bg-[#F8FAFC]">
      {/* Sidebar - Pro Recruitment Interface */}
      <aside className={`
        fixed inset-y-0 left-0 z-50 w-72 bg-slate-950 text-white transform transition-transform duration-500 ease-[cubic-bezier(0.2,0,0,1)]
        md:translate-x-0 md:static md:h-screen sticky top-0
        ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
        flex flex-col p-8
      `}>
        <div className="flex items-center gap-3 mb-16">
          <div className="p-2 bg-emerald-500 rounded-xl">
            <Briefcase className="w-6 h-6 text-slate-950" />
          </div>
          <span className="font-black text-2xl tracking-tighter">ROJGAR<span className="text-emerald-500 uppercase">Pro</span></span>
        </div>

        <nav className="space-y-2 flex-1">
          {[
            { label: 'Overview', path: '/dashboard/company', icon: LayoutDashboard },
            { label: 'Job Postings', path: '/employer/jobs', icon: FileText },
            { label: 'Talent Pool', path: '/employer/candidates', icon: Users },
            { label: 'Settings', path: '/employer/settings', icon: Settings },
          ].map((item) => {
            const isActive = router.pathname.includes(item.path) || (item.path === '/dashboard/company' && router.pathname.includes('company'));
            const Icon = item.icon;
            return (
              <Link
                key={item.path}
                href={item.path}
                className={`flex items-center gap-4 px-6 py-4 rounded-2xl text-sm font-bold transition-all ${
                  isActive ? 'bg-emerald-500 text-slate-950 shadow-xl shadow-emerald-500/20' : 'text-slate-500 hover:text-white hover:bg-white/5'
                }`}
              >
                <Icon className="w-5 h-5" />
                {item.label}
              </Link>
            );
          })}
        </nav>

        <div className="mt-auto p-6 bg-white/5 rounded-[2rem] border border-white/10">
          <p className="text-[10px] font-black text-slate-500 uppercase tracking-widest mb-3">Enterprise Plan</p>
          <div className="flex items-center gap-4">
            <div className="w-10 h-10 rounded-full bg-slate-800 flex items-center justify-center font-black text-emerald-500">
              {(user?.name || 'C')[0]}
            </div>
            <div className="overflow-hidden">
              <p className="text-sm font-bold truncate">{user?.name || 'Recruiter'}</p>
              <button onClick={() => logout()} className="text-[10px] font-black text-rose-500 uppercase hover:text-rose-400 transition-colors">Sign Out</button>
            </div>
          </div>
        </div>
      </aside>

      <main className="flex-1 overflow-auto p-6 md:p-12 lg:p-16">
        <header className="flex flex-col md:flex-row md:items-center justify-between gap-8 mb-16">
          <div className="space-y-2">
             <div className="flex items-center gap-4">
                <button onClick={() => setSidebarOpen(true)} className="p-3 bg-white border border-slate-200 rounded-2xl md:hidden">
                  <Menu className="w-5 h-5" />
                </button>
                <h1 className="text-4xl md:text-5xl font-black text-slate-900 tracking-tight">Hiring Intelligence.</h1>
             </div>
             <p className="text-slate-400 text-lg font-medium">Monitoring <span className="text-emerald-600 font-bold">Talent Pipelines</span> for your organization.</p>
          </div>

          <div className="flex items-center gap-4">
            <Link href="/employer/post-job" className="px-8 py-4 bg-slate-900 text-white font-black rounded-2xl flex items-center gap-3 hover:bg-slate-800 transition-all shadow-xl shadow-slate-900/10 uppercase text-[10px] tracking-widest">
              <Plus className="w-4 h-4" /> New Posting
            </Link>
          </div>
        </header>

        {/* Vital Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-12">
          {[
            { label: 'Active Roles', val: stats.activeJobs, color: 'emerald', icon: Briefcase },
            { label: 'Candidate Fit', val: `${stats.aiMatched}%`, color: 'blue', icon: Zap },
            { label: 'Unread Apps', val: stats.totalApplicants, color: 'indigo', icon: Users },
          ].map(s => (
            <div key={s.label} className="bg-white border border-slate-200 p-8 rounded-[2.5rem] shadow-sm hover:shadow-2xl transition-all group">
              <div className={`w-12 h-12 rounded-2xl bg-${s.color}-50 text-${s.color}-600 flex items-center justify-center mb-6`}>
                <s.icon className="w-6 h-6" />
              </div>
              <p className="text-4xl font-black text-slate-900 tracking-tighter">{s.val}</p>
              <p className="text-xs font-black text-slate-400 uppercase tracking-widest mt-2">{s.label}</p>
            </div>
          ))}
        </div>

        {/* Analytics Section */}
        <section className="mb-12">
           <DashboardAnalytics />
        </section>

        {/* Talent Management Section */}
        <div className="bg-white border border-slate-200 rounded-[3rem] overflow-hidden shadow-sm">
           <div className="p-10 border-b border-slate-100 flex items-center justify-between bg-slate-50/50">
              <h3 className="text-2xl font-black text-slate-900 tracking-tight">Active Pipelines</h3>
              <div className="flex items-center gap-4">
                <div className="relative">
                  <Search className="absolute left-3 top-2.5 w-4 h-4 text-slate-400" />
                  <input type="text" placeholder="Search postings..." className="pl-10 pr-4 py-2 bg-white border border-slate-200 rounded-xl text-xs font-bold outline-none focus:ring-2 focus:ring-emerald-500/20" />
                </div>
                <button className="p-2.5 bg-white border border-slate-200 rounded-xl hover:bg-slate-50"><Filter className="w-4 h-4 text-slate-500" /></button>
              </div>
           </div>

           {loading ? (
             <div className="p-20 text-center flex flex-col items-center">
                <Loader2 className="w-10 h-10 text-emerald-500 animate-spin mb-4" />
                <span className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Compiling Pipeline Data...</span>
             </div>
           ) : (
             <div className="overflow-x-auto">
                <table className="w-full text-left">
                   <thead>
                      <tr className="text-[10px] font-black text-slate-400 uppercase tracking-widest border-b border-slate-100">
                         <th className="px-10 py-6">Position</th>
                         <th className="px-10 py-6">Intelligence</th>
                         <th className="px-10 py-6">Throughput</th>
                         <th className="px-10 py-6">Status</th>
                         <th className="px-10 py-6 text-right">Actions</th>
                      </tr>
                   </thead>
                   <tbody className="divide-y divide-slate-50">
                      {recentJobs.map(job => (
                        <tr key={job.id} className="group hover:bg-slate-50/50 transition-colors">
                           <td className="px-10 py-8">
                              <p className="font-black text-slate-900 group-hover:text-emerald-600 transition-colors">{job.title}</p>
                              <p className="text-[10px] font-bold text-slate-400 uppercase tracking-tighter mt-1">{job.location} • {job.type}</p>
                           </td>
                           <td className="px-10 py-8">
                              <div className="flex items-center gap-2">
                                 <div className="px-3 py-1 bg-emerald-50 text-emerald-600 rounded-lg text-[10px] font-black uppercase tracking-tighter">
                                    {job.aiMatchedCount || 12} Top Matches
                                 </div>
                              </div>
                           </td>
                           <td className="px-10 py-8">
                              <p className="text-sm font-black text-slate-900">{job.applicantsCount || 45}</p>
                              <p className="text-[10px] font-bold text-slate-400 uppercase">Candidates</p>
                           </td>
                           <td className="px-10 py-8">
                              <div className="flex items-center gap-2">
                                 <div className="w-2 h-2 rounded-full bg-blue-500"></div>
                                 <span className="text-[10px] font-black uppercase tracking-widest text-slate-900">Active</span>
                              </div>
                           </td>
                           <td className="px-10 py-8 text-right">
                              <Link href={`/employer/applicants?jobId=${job.id}`} className="inline-flex items-center gap-2 px-5 py-2.5 bg-slate-900 text-white text-[10px] font-black uppercase tracking-widest rounded-xl hover:bg-emerald-600 transition-all shadow-lg shadow-slate-900/5">
                                 Review <ArrowUpRight className="w-3.5 h-3.5" />
                              </Link>
                           </td>
                        </tr>
                      ))}
                   </tbody>
                </table>
             </div>
           )}
        </div>
      </main>
    </div>
  );
};

const SAMPLE_EMPLOYER_JOBS = [
  { id: 'ej1', title: 'Senior Backend Architect', location: 'Remote', type: 'Full-Time', applicantsCount: 84, aiMatchedCount: 22 },
  { id: 'ej2', title: 'Lead UX Strategist', location: 'Dubai, UAE', type: 'Full-Time', applicantsCount: 36, aiMatchedCount: 14 },
];

export default EmployerDashboard;
