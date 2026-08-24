import React, { useEffect, useState } from 'react';
import { Briefcase, Users, Calendar, CheckCircle2, PlusCircle, ArrowUpRight, LayoutDashboard, Settings } from 'lucide-react';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { useAuth } from '../../hooks/useAuth';

export default function CompanyDashboard() {
  const router = useRouter();
  const { user, logout } = useAuth();
  const [stats, setStats] = useState({ activeJobs: 0, totalApplicants: 0, interviewsThisWeek: 0, hired: 0 });
  const [recentApplicants, setRecentApplicants] = useState([]);

  useEffect(() => {
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');
    const headers = { Authorization: `Bearer ${token}` };

    Promise.all([
      fetch('http://localhost:3001/api/v1/companies/me', { headers }).then(r => r.ok ? r.json() : {}),
      fetch('http://localhost:3001/api/v1/companies/me/jobs', { headers }).then(r => r.ok ? r.json() : [])
    ]).then(([compData, jobsData]) => {
      if (compData.stats) setStats(compData.stats);
      if (Array.isArray(jobsData)) setRecentApplicants(jobsData.slice(0, 5));
    }).catch(err => console.error("Error fetching company stats:", err));
  }, []);

  const navItems = [
    { label: 'Dashboard', path: '/dashboard/company', icon: LayoutDashboard },
    { label: 'Post Job', path: '/company/post-job', icon: PlusCircle },
    { label: 'Applicants', path: '/company/applicants', icon: Users },
    { label: 'Interviews', path: '/company/interviews', icon: Calendar },
    { label: 'Profile', path: '/company/profile', icon: Settings },
  ];

  return (
    <div className="flex min-h-screen bg-slate-50 font-sans">
      <aside className="w-64 bg-slate-900 text-white p-6 flex flex-col justify-between fixed h-full z-50">
        <div className="space-y-8">
          <div className="flex items-center gap-3">
            <div className="bg-indigo-600 p-2 rounded-xl">
              <Briefcase className="w-6 h-6 text-white" />
            </div>
            <h2 className="text-xl font-black tracking-tight">ROJGAR<span className="text-indigo-400">SETU</span></h2>
          </div>

          <nav className="space-y-1">
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = router.pathname === item.path;
              return (
                <Link
                  key={item.path}
                  href={item.path}
                  className={`flex items-center gap-3 px-4 py-3 rounded-xl font-medium transition-all ${
                    isActive ? 'bg-indigo-600 text-white' : 'text-slate-400 hover:text-white hover:bg-white/5'
                  }`}
                >
                  <Icon className="w-5 h-5" /> {item.label}
                </Link>
              );
            })}
          </nav>
        </div>

        <button
          onClick={() => { logout(); router.push('/login'); }}
          className="flex items-center gap-3 px-4 py-3 text-slate-400 hover:text-rose-400 hover:bg-rose-500/5 rounded-xl font-medium transition-all"
        >
          <LogOut className="w-5 h-5" /> Logout
        </button>
      </aside>

      <main className="flex-1 ml-64 p-8">
        <header className="flex justify-between items-center mb-10">
          <div>
            <h1 className="text-3xl font-black text-slate-900 tracking-tight">Hiring Intelligence.</h1>
            <p className="text-slate-500 font-medium mt-1 text-lg">Welcome back, <span className="text-indigo-600 font-bold">{user?.name || 'Recruiter'}</span></p>
          </div>
          <Link href="/company/post-job" className="flex items-center gap-2 bg-slate-900 hover:bg-slate-800 text-white px-6 py-3 rounded-2xl font-black text-xs uppercase tracking-widest transition shadow-xl shadow-slate-900/10">
            <PlusCircle className="w-4 h-4" /> New Posting
          </Link>
        </header>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mb-12">
          {[
            { label: 'Active Jobs', val: stats.activeJobs, icon: Briefcase, color: 'indigo' },
            { label: 'Applicants', val: stats.totalApplicants, icon: Users, color: 'blue' },
            { label: 'Interviews', val: stats.interviewsThisWeek, icon: Calendar, color: 'amber' },
            { label: 'Hired', val: stats.hired, icon: CheckCircle2, color: 'emerald' },
          ].map(s => (
            <div key={s.label} className="bg-white border border-slate-200 p-8 rounded-[2rem] shadow-sm hover:shadow-xl transition-all group">
              <div className={`w-12 h-12 rounded-2xl bg-${s.color}-50 text-${s.color}-600 flex items-center justify-center mb-6`}>
                <s.icon className="w-6 h-6" />
              </div>
              <p className="text-4xl font-black text-slate-900 tracking-tighter">{s.val || 0}</p>
              <p className="text-xs font-black text-slate-400 uppercase tracking-widest mt-2">{s.label}</p>
            </div>
          ))}
        </div>

        {/* Recent Applicants Bento */}
        <div className="bg-white border border-slate-200 rounded-[2.5rem] shadow-sm overflow-hidden">
          <div className="p-10 border-b border-slate-100 flex items-center justify-between bg-slate-50/50">
            <h2 className="text-2xl font-black text-slate-900 tracking-tight">Active Pipelines</h2>
            <Link href="/company/applicants" className="text-xs font-black text-indigo-600 hover:underline uppercase tracking-widest">Manage All</Link>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full text-left">
              <thead>
                <tr className="text-[10px] font-black text-slate-400 uppercase tracking-widest border-b border-slate-50">
                  <th className="px-10 py-6">Candidate</th>
                  <th className="px-10 py-6">Applied Job</th>
                  <th className="px-10 py-6">AI Fit</th>
                  <th className="px-10 py-6">Status</th>
                  <th className="px-10 py-6 text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-50 text-sm">
                {recentApplicants.length === 0 ? (
                  <tr><td colSpan="5" className="py-20 text-center text-slate-400 font-bold uppercase text-[10px] tracking-widest">No active candidates to show</td></tr>
                ) : (
                  recentApplicants.map((app) => (
                    <tr key={app.id || app._id} className="group hover:bg-slate-50 transition-colors">
                      <td className="px-10 py-8">
                        <div className="flex items-center gap-4">
                          <div className="w-10 h-10 rounded-full bg-slate-100 flex items-center justify-center font-black text-slate-400 group-hover:text-indigo-600 transition-colors uppercase">
                            {app.candidate_name?.[0] || 'A'}
                          </div>
                          <p className="font-black text-slate-900">{app.candidate_name || "Applicant"}</p>
                        </div>
                      </td>
                      <td className="px-10 py-8 text-slate-500 font-bold uppercase text-[10px]">{app.job_title || "N/A"}</td>
                      <td className="px-10 py-8">
                        <div className="px-3 py-1 bg-indigo-50 text-indigo-700 rounded-lg text-[10px] font-black uppercase tracking-tighter w-fit">
                          {app.score ? `${app.score}% Match` : 'Awaiting'}
                        </div>
                      </td>
                      <td className="px-10 py-8">
                        <div className="flex items-center gap-2">
                          <div className="w-2 h-2 rounded-full bg-blue-500"></div>
                          <span className="text-[10px] font-black uppercase tracking-widest text-slate-900">{app.status || 'Applied'}</span>
                        </div>
                      </td>
                      <td className="px-10 py-8 text-right">
                        <Link href="/company/applicants" className="inline-flex items-center gap-2 px-5 py-2.5 bg-slate-900 text-white text-[10px] font-black uppercase tracking-widest rounded-xl hover:bg-indigo-600 transition-all shadow-lg shadow-slate-900/5">
                          Review <ArrowUpRight className="w-3.5 h-3.5" />
                        </Link>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>
      </main>
    </div>
  );
}

// Stub for LogOut icon
const LogOut = ({className}) => (
  <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
  </svg>
);
