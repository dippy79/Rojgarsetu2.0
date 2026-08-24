import React, { useState, useEffect, useCallback } from 'react';
import { Star, MoreVertical, Search, Filter, Loader2, ArrowLeft, Calendar, User, MessageSquare } from 'lucide-react';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { useAuth } from '../../hooks/useAuth';

const COLUMNS = ['Applied', 'Shortlisted', 'Interview', 'Hired', 'Rejected'];

export default function JobApplicants() {
  const router = useRouter();
  const { user } = useAuth();
  const [applications, setApplications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchQuery] = useState('');

  const fetchApplications = useCallback(async () => {
    setLoading(true);
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');
    try {
      const res = await fetch('http://localhost:3001/api/v1/company/applications', {
        headers: { Authorization: `Bearer ${token}` }
      });
      if (res.ok) {
        const data = await res.json();
        setApplications(Array.isArray(data) ? data : data.data || []);
      } else {
        // Fallback to sample for demo
        setApplications(SAMPLE_APPLICANTS);
      }
    } catch (err) {
      setApplications(SAMPLE_APPLICANTS);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { fetchApplications(); }, [fetchApplications]);

  const updateStatus = async (id, status) => {
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');
    try {
      await fetch(`http://localhost:3001/api/v1/applications/${id}/status`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
        body: JSON.stringify({ status, recruiter_notes: `Status changed to ${status}` })
      });
      fetchApplications();
    } catch (err) {
      console.error("Status update failed:", err);
    }
  };

  const filteredApps = applications.filter(app =>
    (app.candidate_name || '').toLowerCase().includes(searchTerm.toLowerCase()) ||
    (app.job_title || '').toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="min-h-screen bg-[#FBFBFB] font-sans">
      {/* Premium Header */}
      <header className="bg-white border-b border-slate-200 pt-12 pb-8 px-10 sticky top-0 z-40">
        <div className="max-w-[1600px] mx-auto space-y-8">
           <div className="flex items-center justify-between">
              <div className="flex items-center gap-4">
                 <Link href="/dashboard/company" className="p-3 bg-slate-50 border border-slate-200 rounded-2xl hover:bg-slate-100 transition-colors">
                    <ArrowLeft className="w-5 h-5 text-slate-600" />
                 </Link>
                 <div>
                    <h1 className="text-3xl font-black text-slate-900 tracking-tight">Hiring Pipeline.</h1>
                    <p className="text-slate-400 text-sm font-bold uppercase tracking-widest mt-1">Manage {applications.length} Active Applicants</p>
                 </div>
              </div>

              <div className="flex items-center gap-4">
                 <div className="relative group">
                    <Search className="absolute left-4 top-3.5 w-4 h-4 text-slate-400 group-focus-within:text-indigo-600 transition-colors" />
                    <input
                      type="text"
                      placeholder="Search candidates..."
                      value={searchTerm}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      className="pl-12 pr-6 py-3.5 bg-slate-50 border-none rounded-2xl text-sm font-bold w-72 focus:ring-4 focus:ring-indigo-500/10 transition-all outline-none"
                    />
                 </div>
                 <button className="p-3.5 bg-slate-900 text-white rounded-2xl hover:bg-slate-800 transition-all shadow-xl shadow-slate-900/10">
                    <Filter className="w-5 h-5" />
                 </button>
              </div>
           </div>
        </div>
      </header>

      <main className="p-10">
        <div className="max-w-[1600px] mx-auto">
          {loading ? (
            <div className="flex flex-col items-center justify-center py-40">
               <Loader2 className="w-12 h-12 text-indigo-600 animate-spin mb-4" />
               <p className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Hydrating Recruitment Data...</p>
            </div>
          ) : (
            <div className="flex gap-8 overflow-x-auto pb-10 custom-scrollbar">
              {COLUMNS.map(col => {
                const colApps = filteredApps.filter(a => (a.status || 'Applied').toLowerCase() === col.toLowerCase());
                return (
                  <div key={col} className="min-w-[340px] flex flex-col gap-6">
                    <div className="flex items-center justify-between px-2">
                       <div className="flex items-center gap-3">
                          <div className={`w-2 h-2 rounded-full ${getStatusColor(col)}`}></div>
                          <h2 className="font-black text-slate-900 text-sm uppercase tracking-widest">{col}</h2>
                       </div>
                       <span className="px-3 py-1 bg-slate-100 border border-slate-200 text-slate-400 rounded-lg text-[10px] font-black">
                          {colApps.length}
                       </span>
                    </div>

                    <div className="flex-1 space-y-4 min-h-[600px]">
                      {colApps.map(app => (
                        <div key={app.id} className="group bg-white border border-slate-200 rounded-[2rem] p-6 shadow-sm hover:shadow-2xl hover:border-indigo-400 transition-all cursor-pointer relative overflow-hidden">
                           <div className="absolute top-0 right-0 w-24 h-24 bg-indigo-500/5 rounded-full -mr-12 -mt-12 blur-2xl group-hover:scale-150 transition-transform"></div>

                           <div className="relative z-10 space-y-4">
                              <div className="flex justify-between items-start">
                                 <div>
                                    <h3 className="font-black text-slate-900 text-lg group-hover:text-indigo-600 transition-colors leading-tight">{app.candidate_name || "Applicant"}</h3>
                                    <p className="text-[10px] font-black text-slate-400 uppercase tracking-tighter mt-1">{app.job_title || 'General Application'}</p>
                                 </div>
                                 <button className="p-2 hover:bg-slate-50 rounded-lg text-slate-300 hover:text-slate-900 transition-all"><MoreVertical className="w-4 h-4" /></button>
                              </div>

                              <div className="flex items-center gap-1 text-amber-400">
                                 {[...Array(5)].map((_, i) => (
                                    <Star key={i} className={`w-3.5 h-3.5 ${i < (app.star_rating || 3) ? 'fill-current' : 'text-slate-100 fill-slate-100'}`} />
                                 ))}
                              </div>

                              <div className="flex flex-wrap gap-1.5">
                                 {(app.skills || ['Agile', 'React']).slice(0, 3).map(s => (
                                    <span key={s} className="px-2 py-0.5 bg-slate-50 border border-slate-100 text-slate-500 text-[8px] font-black uppercase rounded-md">
                                       {s}
                                    </span>
                                 ))}
                              </div>

                              <div className="flex items-center justify-between pt-4 border-t border-slate-50 text-[9px] font-black text-slate-400 uppercase tracking-widest">
                                 <div className="flex items-center gap-1.5">
                                    <Calendar className="w-3 h-3" />
                                    {app.applied_date ? new Date(app.applied_date).toLocaleDateString() : 'Active'}
                                 </div>
                                 <div className="flex items-center gap-1.5 text-indigo-600">
                                    <User className="w-3 h-3" />
                                    Review Profile
                                 </div>
                              </div>

                              <div className="grid grid-cols-2 gap-2 pt-2 opacity-0 group-hover:opacity-100 transition-opacity">
                                 {col !== 'Shortlisted' && <button onClick={() => updateStatus(app.id, 'Shortlisted')} className="py-2 bg-indigo-50 text-indigo-600 rounded-xl text-[9px] font-black uppercase tracking-tighter hover:bg-indigo-600 hover:text-white transition-all">Shortlist</button>}
                                 {col !== 'Interview' && <button onClick={() => updateStatus(app.id, 'Interview')} className="py-2 bg-amber-50 text-amber-600 rounded-xl text-[9px] font-black uppercase tracking-tighter hover:bg-amber-600 hover:text-white transition-all">Interview</button>}
                                 {col !== 'Hired' && <button onClick={() => updateStatus(app.id, 'Hired')} className="py-2 bg-emerald-50 text-emerald-600 rounded-xl text-[9px] font-black uppercase tracking-tighter hover:bg-emerald-600 hover:text-white transition-all">Hire</button>}
                                 {col !== 'Rejected' && <button onClick={() => updateStatus(app.id, 'Rejected')} className="py-2 bg-rose-50 text-rose-600 rounded-xl text-[9px] font-black uppercase tracking-tighter hover:bg-rose-600 hover:text-white transition-all">Reject</button>}
                              </div>
                           </div>
                        </div>
                      ))}
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </main>
    </div>
  );
}

const getStatusColor = (col) => {
  switch(col) {
    case 'Applied': return 'bg-blue-500';
    case 'Shortlisted': return 'bg-indigo-500';
    case 'Interview': return 'bg-amber-500';
    case 'Hired': return 'bg-emerald-500';
    case 'Rejected': return 'bg-rose-500';
    default: return 'bg-slate-500';
  }
};

const SAMPLE_APPLICANTS = [
  { id: 'a1', candidate_name: 'Aditya Sharma', job_title: 'Full Stack Engineer', status: 'Applied', star_rating: 4, skills: ['Node.js', 'React', 'AWS'], applied_date: '2026-08-22' },
  { id: 'a2', candidate_name: 'Priya Verma', job_title: 'UX Architect', status: 'Shortlisted', star_rating: 5, skills: ['Figma', 'Strategy', 'Design System'], applied_date: '2026-08-21' },
  { id: 'a3', candidate_name: 'Vikram Singh', job_title: 'Backend Lead', status: 'Interview', star_rating: 3, skills: ['Python', 'PostgreSQL', 'Redis'], applied_date: '2026-08-20' },
];
