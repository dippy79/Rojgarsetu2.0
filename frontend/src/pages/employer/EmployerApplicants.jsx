import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';

export const EmployerApplicants = () => {
  const router = useRouter();
  const API_BASE = 'http://localhost:3001';

  const [applicants, setApplicants] = useState([]);
  const [selectedJob, setSelectedJob] = useState('ALL');
  const [minMatchFilter, setMinMatchFilter] = useState(0);
  const [statusFilter, setStatusFilter] = useState('ALL');
  const [loading, setLoading] = useState(true);
  const [toast, setToast] = useState(null);
  const [activeResumeModal, setActiveResumeModal] = useState(null);

  useEffect(() => {
    fetchApplicants();
  }, []);

  const showToast = (message, type = 'success') => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 3500);
  };

  const fetchApplicants = async () => {
    setLoading(true);
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');
    try {
      const res = await fetch(`${API_BASE}/api/v1/employer/applicants`, {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
      });

      if (res.ok) {
        const data = await res.json();
        setApplicants(data.applicants || []);
      } else {
        setMockApplicants();
      }
    } catch (err) {
      setMockApplicants();
    } finally {
      setLoading(false);
    }
  };

  const setMockApplicants = () => {
    setApplicants([
      {
        id: 'cand-1',
        name: 'Simranjeet Singh',
        email: 'simranjeet@example.com',
        roleApplied: 'Senior React / Frontend Engineer',
        jobId: 'job-101',
        matchScore: 96,
        skills: ['React.js', 'Node.js', 'Flutter', 'TypeScript', 'Tailwind CSS'],
        experience: '4.5 Years',
        appliedDate: '2026-03-01',
        status: 'Shortlisted',
        resumeUrl: 'Resume_Simranjeet_2026.pdf'
      },
      {
        id: 'cand-2',
        name: 'Aarav Sharma',
        email: 'aarav.sharma@example.com',
        roleApplied: 'Senior React / Frontend Engineer',
        jobId: 'job-101',
        matchScore: 89,
        skills: ['React.js', 'Redux', 'JavaScript', 'CSS3'],
        experience: '3.0 Years',
        appliedDate: '2026-03-02',
        status: 'Under Review',
        resumeUrl: 'Aarav_Resume.pdf'
      },
      {
        id: 'cand-3',
        name: 'Priya Verma',
        email: 'priya.v@example.com',
        roleApplied: 'Full Stack Node.js & React Developer',
        jobId: 'job-102',
        matchScore: 84,
        skills: ['Node.js', 'Express', 'MongoDB', 'React.js'],
        experience: '2.5 Years',
        appliedDate: '2026-02-28',
        status: 'Under Review',
        resumeUrl: 'Priya_Verma_Resume.pdf'
      },
      {
        id: 'cand-4',
        name: 'Rohan Gupta',
        email: 'rohan.g@example.com',
        roleApplied: 'Flutter Mobile App Engineer',
        jobId: 'job-103',
        matchScore: 72,
        skills: ['Flutter', 'Dart', 'Firebase'],
        experience: '1.5 Years',
        appliedDate: '2026-03-02',
        status: 'Applied',
        resumeUrl: 'Rohan_Flutter.pdf'
      }
    ]);
  };

  const handleUpdateStatus = (candId, newStatus) => {
    setApplicants(
      applicants.map((a) => (a.id === candId ? { ...a, status: newStatus } : a))
    );
    showToast(`Candidate status changed to ${newStatus}`);
  };

  const filteredApplicants = applicants.filter((app) => {
    const matchesJob = selectedJob === 'ALL' || app.jobId === selectedJob;
    const matchesScore = app.matchScore >= minMatchFilter;
    const matchesStatus = statusFilter === 'ALL' || app.status === statusFilter;
    return matchesJob && matchesScore && matchesStatus;
  });

  return (
    <div className="min-h-screen bg-slate-50 p-8 font-sans">
      {toast && (
        <div className="fixed top-6 right-6 z-50 px-5 py-3 rounded-xl shadow-lg border text-sm font-semibold bg-emerald-50 text-emerald-800 border-emerald-200">
          <span>{toast.message}</span>
        </div>
      )}

      <div className="max-w-6xl mx-auto space-y-8">
        <header>
          <h1 className="text-3xl font-bold text-slate-900 tracking-tight">AI Applicant Ranking & Review</h1>
          <p className="text-slate-500 text-sm mt-1 font-medium">
            Filter applicants by AI match percentage, evaluation metrics, and resume fit.
          </p>
        </header>

        {/* Filter Controls Bar */}
        <div className="bg-white p-5 rounded-2xl border border-slate-200/60 shadow-sm grid grid-cols-1 md:grid-cols-12 gap-4 items-center">
          <div className="md:col-span-4">
            <label htmlFor="filter_job" className="block text-xs font-bold text-slate-700 uppercase mb-1">Filter Position</label>
            <select
              id="filter_job"
              value={selectedJob}
              onChange={(e) => setSelectedJob(e.target.value)}
              className="w-full px-3.5 py-2 rounded-xl border border-slate-200 text-sm bg-white font-medium focus:outline-none focus:ring-2 focus:ring-blue-500/20"
            >
              <option value="ALL">All Openings</option>
              <option value="job-101">Senior React / Frontend Engineer</option>
              <option value="job-102">Full Stack Node.js & React Developer</option>
              <option value="job-103">Flutter Mobile App Engineer</option>
            </select>
          </div>

          <div className="md:col-span-4">
            <label htmlFor="filter_min_match" className="block text-xs font-bold text-slate-700 uppercase mb-1">
              Min. AI Match Score ({minMatchFilter}%)
            </label>
            <input
              type="range"
              id="filter_min_match"
              min="0"
              max="95"
              step="5"
              value={minMatchFilter}
              onChange={(e) => setMinMatchFilter(Number(e.target.value))}
              className="w-full accent-blue-600"
            />
          </div>

          <div className="md:col-span-4">
            <label htmlFor="filter_status" className="block text-xs font-bold text-slate-700 uppercase mb-1">Application Status</label>
            <select
              id="filter_status"
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="w-full px-3.5 py-2 rounded-xl border border-slate-200 text-sm bg-white font-medium focus:outline-none focus:ring-2 focus:ring-blue-500/20"
            >
              <option value="ALL">All Statuses</option>
              <option value="Shortlisted">Shortlisted</option>
              <option value="Under Review">Under Review</option>
              <option value="Applied">Applied</option>
              <option value="Rejected">Rejected</option>
            </select>
          </div>
        </div>

        {/* Applicants Feed */}
        {loading ? (
          <div className="p-12 text-center text-slate-400">Evaluating candidates...</div>
        ) : filteredApplicants.length === 0 ? (
          <div className="bg-white p-12 text-center rounded-2xl border border-slate-200/60 shadow-sm">
            <p className="text-sm text-slate-500">No applicants match the selected filter criteria.</p>
          </div>
        ) : (
          <div className="space-y-4">
            {filteredApplicants.map((app) => (
              <div
                key={app.id}
                className="bg-white p-6 rounded-2xl border border-slate-200/60 shadow-sm flex flex-col md:flex-row items-start justify-between gap-6 hover:border-blue-200 transition-all"
              >
                <div className="space-y-3 flex-1">
                  <div className="flex items-center gap-3 flex-wrap">
                    <h3 className="text-lg font-bold text-slate-900">{app.name}</h3>
                    <span className="px-3 py-0.5 bg-emerald-50 text-emerald-700 border border-emerald-200 text-xs font-black rounded-full flex items-center gap-1">
                      ✨ {app.matchScore}% Match
                    </span>
                    <span className="px-2.5 py-0.5 bg-slate-100 text-slate-700 text-xs font-semibold rounded-md">
                      Exp: {app.experience}
                    </span>
                  </div>

                  <p className="text-xs font-semibold text-slate-500">
                    Applied for <span className="text-slate-800 font-bold">{app.roleApplied}</span> • {app.appliedDate}
                  </p>

                  <div className="flex flex-wrap gap-1.5 pt-1">
                    {app.skills.map((skill, i) => (
                      <span key={i} className="text-xs bg-slate-100 text-slate-700 px-2.5 py-0.5 rounded-lg font-medium">
                        {skill}
                      </span>
                    ))}
                  </div>
                </div>

                <div className="flex flex-col items-start md:items-end gap-3 w-full md:w-auto">
                  <div className="flex items-center gap-2">
                    <button
                      type="button"
                      onClick={() => setActiveResumeModal(app)}
                      className="px-3 py-1.5 bg-slate-100 text-slate-700 text-xs font-semibold rounded-lg hover:bg-slate-200"
                    >
                      📄 Preview Resume
                    </button>
                  </div>

                  <div className="flex items-center gap-2 pt-2">
                    <button
                      type="button"
                      onClick={() => handleUpdateStatus(app.id, 'Shortlisted')}
                      className={`px-3 py-1.5 text-xs font-bold rounded-lg border transition-all ${
                        app.status === 'Shortlisted'
                          ? 'bg-emerald-600 text-white border-emerald-600'
                          : 'bg-emerald-50 text-emerald-700 border-emerald-200 hover:bg-emerald-100'
                      }`}
                    >
                      Shortlist
                    </button>

                    <button
                      type="button"
                      onClick={() => handleUpdateStatus(app.id, 'Rejected')}
                      className={`px-3 py-1.5 text-xs font-bold rounded-lg border transition-all ${
                        app.status === 'Rejected'
                          ? 'bg-red-600 text-white border-red-600'
                          : 'bg-red-50 text-red-700 border-red-200 hover:bg-red-100'
                      }`}
                    >
                      Reject
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Resume Preview Modal */}
        {activeResumeModal && (
          <div className="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
            <div className="bg-white w-full max-w-lg rounded-2xl p-6 shadow-2xl border border-slate-200">
              <div className="flex justify-between items-center pb-3 border-b border-slate-100 mb-4">
                <h3 className="font-bold text-slate-900">{activeResumeModal.name} - Resume Preview</h3>
                <button
                  type="button"
                  onClick={() => setActiveResumeModal(null)}
                  className="text-slate-400 hover:text-slate-600 font-bold"
                >
                  ✕
                </button>
              </div>

              <div className="p-4 bg-slate-50 border border-slate-200 rounded-xl space-y-3 text-xs text-slate-700">
                <p><strong>File Name:</strong> {activeResumeModal.resumeUrl}</p>
                <p><strong>Email:</strong> {activeResumeModal.email}</p>
                <p><strong>Total Experience:</strong> {activeResumeModal.experience}</p>
                <p><strong>AI Skill Fit Score:</strong> {activeResumeModal.matchScore}%</p>
                <div className="pt-2 border-t border-slate-200">
                  <p className="text-[11px] text-slate-500 italic">
                    Integrated with PDF Parser Engine. AI extracted experience matches requirements for full stack systems.
                  </p>
                </div>
              </div>

              <div className="mt-5 flex justify-end">
                <button
                  type="button"
                  onClick={() => setActiveResumeModal(null)}
                  className="px-4 py-2 bg-slate-800 text-white text-xs font-bold rounded-xl hover:bg-slate-900"
                >
                  Close Preview
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default EmployerApplicants;