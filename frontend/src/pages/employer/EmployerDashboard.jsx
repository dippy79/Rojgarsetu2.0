import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { useAuth } from '../../hooks/useAuth';
import { useRouter } from 'next/router';
import { useAuth } from '../../hooks/useAuth';

export const EmployerDashboard = () => {
  const { user } = useAuth();
  const router = useRouter();
  const API_BASE = 'http://localhost:3001';

  const [stats, setStats] = useState({
    activeJobs: 4,
    totalApplicants: 128,
    shortlisted: 18,
    aiMatched: 42
  });

  const [recentJobs, setRecentJobs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [toast, setToast] = useState(null);

  useEffect(() => {
    fetchDashboardData();
  }, []);

  const showToast = (message, type = 'success') => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 3500);
  };

  const fetchDashboardData = async () => {
    setLoading(true);
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');
    try {
      const res = await fetch(`${API_BASE}/api/v1/employer/dashboard`, {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
      });

      if (res.ok) {
        const data = await res.json();
        setStats(data.stats || stats);
        setRecentJobs(data.jobs || []);
      } else {
        setMockJobs();
      }
    } catch (err) {
      setMockJobs();
    } finally {
      setLoading(false);
    }
  };

  const setMockJobs = () => {
    setRecentJobs([
      {
        id: 'job-101',
        title: 'Senior React / Frontend Engineer',
        location: 'Bengaluru / Remote',
        type: 'Full-time',
        applicantsCount: 45,
        aiMatchedCount: 14,
        status: 'Active',
        postedDate: '2026-02-28'
      },
      {
        id: 'job-102',
        title: 'Full Stack Node.js & React Developer',
        location: 'Gurugram / Hybrid',
        type: 'Full-time',
        applicantsCount: 32,
        aiMatchedCount: 11,
        status: 'Active',
        postedDate: '2026-03-01'
      },
      {
        id: 'job-103',
        title: 'Flutter Mobile App Engineer',
        location: 'Delhi NCR',
        type: 'Full-time',
        applicantsCount: 29,
        aiMatchedCount: 9,
        status: 'Active',
        postedDate: '2026-02-22'
      },
      {
        id: 'job-104',
        title: 'Backend API Developer (Python/FastAPI)',
        location: 'Noida',
        type: 'Full-time',
        applicantsCount: 22,
        aiMatchedCount: 8,
        status: 'Draft',
        postedDate: '2026-03-02'
      }
    ]);
  };

  const handleToggleJobStatus = (jobId, currentStatus) => {
    const nextStatus = currentStatus === 'Active' ? 'Closed' : 'Active';
    setRecentJobs(
      recentJobs.map((j) => (j.id === jobId ? { ...j, status: nextStatus } : j))
    );
    showToast(`Job status updated to ${nextStatus}`);
  };

  return (
    <div className="min-h-screen bg-slate-50 p-8 font-sans">
      {toast && (
        <div className="fixed top-6 right-6 z-50 px-5 py-3 rounded-xl shadow-lg border text-sm font-semibold bg-emerald-50 text-emerald-800 border-emerald-200">
          <span>{toast.message}</span>
        </div>
      )}

      <div className="max-w-6xl mx-auto space-y-8">
        {/* Header */}
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Employer Portal</h1>
            <p className="text-slate-500 text-sm mt-1 font-medium">
              Welcome back, {user?.name || 'Recruiter'}. Manage postings and review AI-ranked candidates.
            </p>
          </div>
          <Link
            href="/employer/post-job"
            className="px-5 py-2.5 bg-blue-600 text-white text-xs font-bold rounded-xl hover:bg-blue-700 shadow-md transition-all self-start md:self-auto flex items-center gap-2"
          >
            <span>+</span> Post New Opening
          </Link>
        </div>

        {/* Metrics Overview Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
          <div className="bg-white p-5 rounded-2xl border border-slate-200/60 shadow-sm">
            <span className="text-xs font-bold text-slate-400 uppercase tracking-wider">Active Openings</span>
            <div className="text-3xl font-black text-slate-900 mt-2">{stats.activeJobs}</div>
            <p className="text-[11px] text-emerald-600 font-semibold mt-1">4 roles receiving applications</p>
          </div>

          <div className="bg-white p-5 rounded-2xl border border-slate-200/60 shadow-sm">
            <span className="text-xs font-bold text-slate-400 uppercase tracking-wider">Total Applicants</span>
            <div className="text-3xl font-black text-slate-900 mt-2">{stats.totalApplicants}</div>
            <p className="text-[11px] text-blue-600 font-semibold mt-1">+12 this week</p>
          </div>

          <div className="bg-white p-5 rounded-2xl border border-slate-200/60 shadow-sm">
            <span className="text-xs font-bold text-slate-400 uppercase tracking-wider">Shortlisted</span>
            <div className="text-3xl font-black text-slate-900 mt-2">{stats.shortlisted}</div>
            <p className="text-[11px] text-purple-600 font-semibold mt-1">Ready for interviews</p>
          </div>

          <div className="bg-white p-5 rounded-2xl border border-slate-200/60 shadow-sm">
            <span className="text-xs font-bold text-slate-400 uppercase tracking-wider">High AI Match Score</span>
            <div className="text-3xl font-black text-emerald-600 mt-2">{stats.aiMatched}</div>
            <p className="text-[11px] text-slate-500 font-semibold mt-1">&gt;85% Skill Fit</p>
          </div>
        </div>

        {/* Recent Jobs Table */}
        <div className="bg-white rounded-2xl border border-slate-200/60 shadow-sm overflow-hidden">
          <div className="p-6 border-b border-slate-100 flex items-center justify-between">
            <h2 className="text-lg font-bold text-slate-900">Manage Active Job Listings</h2>
            <span className="text-xs font-medium text-slate-500">Showing {recentJobs.length} positions</span>
          </div>

          {loading ? (
            <div className="p-8 text-center text-slate-400">Loading active job postings...</div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm text-slate-600">
                <thead className="bg-slate-50 text-[11px] font-bold text-slate-400 uppercase tracking-wider border-b border-slate-100">
                  <tr>
                    <th className="py-3.5 px-6">Job Title</th>
                    <th className="py-3.5 px-6">Location</th>
                    <th className="py-3.5 px-6">Applicants</th>
                    <th className="py-3.5 px-6">AI Fit</th>
                    <th className="py-3.5 px-6">Status</th>
                    <th className="py-3.5 px-6 text-right">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {recentJobs.map((job) => (
                    <tr key={job.id} className="hover:bg-slate-50/80 transition-colors">
                      <td className="py-4 px-6 font-bold text-slate-900">
                        {job.title}
                        <span className="block text-[11px] font-normal text-slate-400">Posted on {job.postedDate}</span>
                      </td>
                      <td className="py-4 px-6 text-xs font-medium">{job.location}</td>
                      <td className="py-4 px-6 font-bold text-slate-800">{job.applicantsCount}</td>
                      <td className="py-4 px-6">
                        <span className="px-2.5 py-1 bg-emerald-50 text-emerald-700 text-xs font-extrabold rounded-full border border-emerald-100">
                          {job.aiMatchedCount} Top Matches
                        </span>
                      </td>
                      <td className="py-4 px-6">
                        <span
                          className={`px-2.5 py-1 text-xs font-bold rounded-lg ${
                            job.status === 'Active'
                              ? 'bg-blue-50 text-blue-700 border border-blue-200'
                              : 'bg-slate-100 text-slate-600'
                          }`}
                        >
                          {job.status}
                        </span>
                      </td>
                      <td className="py-4 px-6 text-right space-x-2">
                        <button
                          type="button"
                          onClick={() => handleToggleJobStatus(job.id, job.status)}
                          className="px-3 py-1.5 bg-slate-100 text-slate-700 text-xs font-semibold rounded-lg hover:bg-slate-200"
                        >
                          {job.status === 'Active' ? 'Close' : 'Activate'}
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default EmployerDashboard;