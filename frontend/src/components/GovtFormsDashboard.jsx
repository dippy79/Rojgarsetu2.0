import React, { useState, useEffect } from 'react';
import { apiUrl } from '../apiConfig';
import { FileText, Calendar, ArrowRight, ShieldAlert, Clock, CheckCircle2, Loader2, Search } from 'lucide-react';

const GovtFormsDashboard = () => {
  const [forms, setForms] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    async function load() {
      try {
        setLoading(true);
        const response = await fetch(apiUrl('/api/v1/forms'));
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        const data = await response.json();
        setForms(data.data || data.forms || []);
      } catch (err) {
        console.error(err);
        setError("Live feed synchronization interrupted. Please try again later.");
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  const daysUntil = (isoDate) => {
    const diff = new Date(isoDate).getTime() - new Date().getTime();
    return Math.ceil(diff / (1000 * 60 * 60 * 24));
  };

  const getPriorityInfo = (days) => {
    if (days <= 3) return { label: 'Ending Soon', color: 'text-rose-600', bg: 'bg-rose-50', icon: ShieldAlert };
    if (days <= 14) return { label: 'Limited Time', color: 'text-amber-600', bg: 'bg-amber-50', icon: Clock };
    return { label: 'Active', color: 'text-emerald-600', bg: 'bg-emerald-50', icon: CheckCircle2 };
  };

  const sortedForms = [...forms].sort((a, b) => new Date(a.last_date) - new Date(b.last_date));

  return (
    <div className="min-h-screen bg-[#FBFBFB] pt-24 pb-20">
      <div className="max-w-7xl mx-auto px-6">
        {/* Page Header */}
        <header className="mb-16 space-y-4">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-slate-900 rounded-2xl shadow-xl shadow-slate-900/10">
              <FileText className="w-6 h-6 text-white" />
            </div>
            <h1 className="text-4xl md:text-5xl font-black text-slate-900 tracking-tight">Official Registrations.</h1>
          </div>
          <p className="text-slate-500 font-medium text-lg max-w-2xl">
            Centralized tracking for government entrance exams, recruitment forms, and official gazette notifications.
            Sorted by deadline proximity.
          </p>
        </header>

        {/* Dynamic Search & Utility */}
        <div className="bg-white border border-slate-200 p-4 rounded-[2rem] shadow-sm flex flex-col md:flex-row items-center gap-4 mb-12">
          <div className="flex-1 relative w-full group">
            <Search className="absolute left-4 top-3.5 w-5 h-5 text-slate-400 group-focus-within:text-slate-900 transition-colors" />
            <input type="text" placeholder="Filter by examination or board..." className="w-full pl-12 pr-4 py-3 bg-slate-50 border-none rounded-xl text-sm font-bold outline-none" />
          </div>
          <div className="px-6 py-3 bg-indigo-50 text-indigo-700 rounded-xl text-[10px] font-black uppercase tracking-widest">
            {forms.length} Dynamic Forms Detected
          </div>
        </div>

        {loading ? (
          <div className="flex flex-col items-center justify-center py-24 space-y-6">
            <Loader2 className="w-12 h-12 text-slate-900 animate-spin" />
            <span className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Polling Official APIs...</span>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {sortedForms.map((form) => {
              const days = daysUntil(form.last_date);
              const p = getPriorityInfo(days);
              const StatusIcon = p.icon;

              return (
                <div key={form.id} className="group bg-white border border-slate-200 rounded-[2.5rem] p-8 hover:border-slate-900 hover:shadow-2xl transition-all flex flex-col justify-between">
                  <div>
                    <div className="flex items-center justify-between mb-8">
                      <div className={`px-4 py-1.5 ${p.bg} ${p.color} rounded-full flex items-center gap-2 text-[10px] font-black uppercase tracking-widest border border-current/10`}>
                        <StatusIcon className="w-3.5 h-3.5" />
                        {p.label}
                      </div>
                      <span className="text-[10px] font-bold text-slate-400 uppercase tracking-tighter">ID: {form.id?.slice(0, 8)}</span>
                    </div>

                    <h3 className="text-2xl font-black text-slate-900 group-hover:text-indigo-600 transition-colors mb-2 leading-tight">
                      {form.title}
                    </h3>
                    <p className="text-slate-400 text-xs font-bold uppercase tracking-widest mb-8">{form.department || 'Official Board'}</p>
                  </div>

                  <div className="space-y-6">
                    <div className="p-5 bg-slate-50 rounded-2xl flex items-center justify-between group-hover:bg-slate-900 group-hover:text-white transition-all">
                      <div className="flex items-center gap-3">
                        <Calendar className="w-5 h-5 text-slate-400" />
                        <span className="text-xs font-black uppercase tracking-widest">Closing Date</span>
                      </div>
                      <span className="text-sm font-black">{new Date(form.last_date).toLocaleDateString('en-IN', { day: 'numeric', month: 'short', year: 'numeric' })}</span>
                    </div>

                    <a
                      href={form.apply_url}
                      target="_blank"
                      rel="noreferrer"
                      className="w-full flex items-center justify-center gap-3 py-5 bg-slate-900 text-white font-black rounded-2xl hover:bg-indigo-600 transition-all shadow-xl shadow-slate-900/10 uppercase text-xs tracking-widest"
                    >
                      Official Apply <ArrowRight className="w-4 h-4" />
                    </a>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
};

export default GovtFormsDashboard;
