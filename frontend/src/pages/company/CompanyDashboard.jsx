import React, { useEffect, useState } from 'react';
import { Briefcase, Users, Calendar, CheckCircle2, PlusCircle, ArrowUpRight } from 'lucide-react';
import { Link } from 'react-router-dom';

export default function CompanyDashboard() {
  const [stats, setStats] = useState({ activeJobs: 0, totalApplicants: 0, interviewsThisWeek: 0, hired: 0 });
  const [recentApplicants, setRecentApplicants] = useState([]);

  useEffect(() => {
    const token = localStorage.getItem('rojgar_token');
    const headers = { Authorization: `Bearer ${token}` };

    Promise.all([
      fetch('http://localhost:3001/api/v1/companies/me', { headers }).then(r => r.ok ? r.json() : {}),
      fetch('http://localhost:3001/api/v1/companies/me/jobs', { headers }).then(r => r.ok ? r.json() : [])
    ]).then(([compData, jobsData]) => {
      if (compData.stats) setStats(compData.stats);
      if (Array.isArray(jobsData)) setRecentApplicants(jobsData.slice(0, 5));
    }).catch(err => console.error("Error fetching company stats:", err));
  }, []);

  return (
    <div className="flex min-h-screen bg-slate-50">
      <aside className="w-64 bg-white border-r border-slate-200 p-6 flex flex-col justify-between">
        <div className="space-y-6">
          <h2 className="text-xl font-bold text-slate-800">RojgarSetu <span className="text-xs bg-indigo-100 text-indigo-700 px-2 py-0.5 rounded">Recruiter</span></h2>
          <nav className="space-y-2">
            <Link to="/dashboard/company" className="flex items-center gap-3 px-4 py-3 bg-indigo-50 text-indigo-600 rounded-xl font-medium">
              <Briefcase className="w-5 h-5" /> Dashboard
            </Link>
            <Link to="/company/post-job" className="flex items-center gap-3 px-4 py-3 text-slate-600 hover:bg-slate-50 rounded-xl font-medium">
              <PlusCircle className="w-5 h-5" /> Post Job
            </Link>
            <Link to="/company/applicants" className="flex items-center gap-3 px-4 py-3 text-slate-600 hover:bg-slate-50 rounded-xl font-medium">
              <Users className="w-5 h-5" /> Applicants
            </Link>
          </nav>
        </div>
      </aside>

      <main className="flex-1 p-8">
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-2xl font-bold text-slate-900">Company Overview</h1>
            <p className="text-slate-500">Track candidates, post positions, and manage pipelines.</p>
          </div>
          <Link to="/company/post-job" className="flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white px-5 py-2.5 rounded-xl font-medium transition">
            <PlusCircle className="w-5 h-5" /> Post New Job
          </Link>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <div className="bg-white/80 backdrop-blur-xl border border-slate-100 p-6 rounded-2xl shadow-sm">
            <div className="text-indigo-600 mb-2"><Briefcase className="w-6 h-6" /></div>
            <p className="text-3xl font-bold text-slate-900">{stats.activeJobs || 0}</p>
            <p className="text-sm font-medium text-slate-500">Active Jobs</p>
          </div>
          <div className="bg-white/80 backdrop-blur-xl border border-slate-100 p-6 rounded-2xl shadow-sm">
            <div className="text-blue-600 mb-2"><Users className="w-6 h-6" /></div>
            <p className="text-3xl font-bold text-slate-900">{stats.totalApplicants || 0}</p>
            <p className="text-sm font-medium text-slate-500">Total Applicants</p>
          </div>
          <div className="bg-white/80 backdrop-blur-xl border border-slate-100 p-6 rounded-2xl shadow-sm">
            <div className="text-amber-600 mb-2"><Calendar className="w-6 h-6" /></div>
            <p className="text-3xl font-bold text-slate-900">{stats.interviewsThisWeek || 0}</p>
            <p className="text-sm font-medium text-slate-500">Interviews This Week</p>
          </div>
          <div className="bg-white/80 backdrop-blur-xl border border-slate-100 p-6 rounded-2xl shadow-sm">
            <div className="text-emerald-600 mb-2"><CheckCircle2 className="w-6 h-6" /></div>
            <p className="text-3xl font-bold text-slate-900">{stats.hired || 0}</p>
            <p className="text-sm font-medium text-slate-500">Hired Candidates</p>
          </div>
        </div>

        <div className="bg-white/80 backdrop-blur-xl border border-slate-100 rounded-2xl shadow-sm p-6">
          <h2 className="text-lg font-bold text-slate-900 mb-4">Recent Applicants</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="border-b border-slate-100 text-xs font-semibold text-slate-400 uppercase">
                  <th className="py-3 px-4">Candidate</th>
                  <th className="py-3 px-4">Applied Job</th>
                  <th className="py-3 px-4">Score</th>
                  <th className="py-3 px-4">Status</th>
                  <th className="py-3 px-4 text-right">Action</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 text-sm">
                {recentApplicants.length === 0 ? (
                  <tr><td colSpan="5" className="py-6 text-center text-slate-400">No applicants found yet.</td></tr>
                ) : (
                  recentApplicants.map((app) => (
                    <tr key={app.id || app._id} className="hover:bg-slate-50/50">
                      <td className="py-3 px-4 font-medium text-slate-800">{app.candidate_name || "Applicant"}</td>
                      <td className="py-3 px-4 text-slate-600">{app.job_title || "N/A"}</td>
                      <td className="py-3 px-4 font-semibold text-indigo-600">{app.score ? `${app.score}%` : 'N/A'}</td>
                      <td className="py-3 px-4">
                        <span className="px-2.5 py-1 text-xs font-semibold rounded-full bg-blue-50 text-blue-600 uppercase">
                          {app.status || 'Applied'}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-right">
                        <Link to="/company/applicants" className="text-indigo-600 hover:text-indigo-800 font-medium inline-flex items-center gap-1">
                          Pipeline <ArrowUpRight className="w-4 h-4" />
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
