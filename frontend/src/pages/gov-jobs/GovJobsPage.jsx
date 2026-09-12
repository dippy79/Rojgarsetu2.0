import React, { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { useAuth } from '../../hooks/useAuth';
import api from '../../lib/api';
import JobFilters from '../../components/JobFilters';
import JobCard from '../../components/JobCard';
import { Search, MapPin, ShieldCheck, ArrowRight, Loader2, Info } from 'lucide-react';

export const GovJobsPage = () => {
  const { user, isAuthenticated } = useAuth();
  const [jobs, setJobs] = useState([]);
  const [pagination, setPagination] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filters, setFilters] = useState({
    location: '',
    department: '',
    category: '',
  });

  const fetchGovJobs = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const queryParams = {
        location: filters.location,
        department: filters.department,
      };

      const res = await api.get('/api/v1/gov-jobs', { params: queryParams });
      const data = res.data?.data || res.data || [];
      const safeData = Array.isArray(data) ? data : [];
      setJobs(safeData);
      setPagination(res.data?.pagination);
    } catch (err) {
      console.error("Fetch error:", err);
      setError("Failed to sync with live government registries. Please try again later.");
      setJobs([]);
    } finally {
      setLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    fetchGovJobs();
  }, [fetchGovJobs]);

  const handleFilterChange = (newFilters) => {
    setFilters(prev => ({ ...prev, ...newFilters }));
  };

  return (
    <div className="min-h-screen bg-stone-50">
      {/* Premium Hero Section */}
      <section className="relative pt-24 pb-20 overflow-hidden bg-white border-b border-slate-200/60">
        <div className="absolute inset-0 bg-[radial-gradient(#e2e8f0_1px,transparent_1px)] [background-size:24px_24px] opacity-40"></div>
        <div className="max-w-7xl mx-auto px-6 relative z-10">
          <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-16">
            <div className="flex-1 space-y-8">
              <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-stone-100 border border-stone-200 text-stone-600 text-[10px] font-black uppercase tracking-[0.2em]">
                <ShieldCheck className="w-3.5 h-3.5" />
                Verified Public Sector Intelligence
              </div>
              <h1 className="text-5xl md:text-7xl font-black text-stone-800 tracking-tight leading-[1.1]">
                Your Career in <br />
                <span className="text-stone-600">Public Service</span> Starts Here.
              </h1>
              <p className="text-lg text-stone-500 max-w-xl font-medium leading-relaxed">
                Directly connected to SSC, UPSC, and Railway Recruitment Boards.
                Real-time synchronization with official gazettes and notifications.
              </p>

              {!isAuthenticated && (
                <div className="flex items-center gap-4 pt-4">
                  <Link href="/login" className="px-8 py-4 bg-stone-800 text-white font-black rounded-2xl hover:bg-stone-900 transition-all shadow-xl shadow-blue-900/20 text-xs uppercase tracking-widest">
                    Unlock All Notifications
                  </Link>
                </div>
              )}
            </div>

            <div className="flex-1 max-w-md bg-white border border-slate-200 p-10 rounded-2xl shadow-2xl shadow-slate-200/50 relative">
               <div className="absolute -top-6 -right-6 w-24 h-24 bg-stone-800 rounded-full flex items-center justify-center text-white font-black text-xs rotate-12 shadow-xl border-4 border-white">
                  LIVE<br/>AGGREGATOR
               </div>
               <h3 className="text-xl font-black text-stone-800 mb-8">Instant Search</h3>
               <div className="space-y-5">
                  <div className="relative group">
                    <Search className="absolute left-4 top-4 w-5 h-5 text-slate-400 group-focus-within:text-stone-800 transition-colors" />
                    <input type="text" placeholder="Job title or Department..." className="w-full pl-12 pr-4 py-4 bg-slate-50 border-none rounded-2xl text-sm font-bold focus:ring-4 focus:ring-stone-700/10 transition-all outline-none" />
                  </div>
                  <button className="w-full py-4 bg-stone-800 text-white font-black rounded-2xl hover:bg-stone-900 transition-all shadow-xl shadow-blue-900/20 uppercase tracking-widest text-xs">
                    Fetch Opportunities
                  </button>
               </div>
            </div>
          </div>
        </div>
      </section>

      {/* Main Grid Layout */}
      <main className="max-w-7xl mx-auto px-6 py-16">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-12">
          {/* Sidebar Filters */}
          <aside className="lg:col-span-3">
            <JobFilters onFilterChange={handleFilterChange} type="government" />
          </aside>

          {/* Job Feed */}
          <div className="lg:col-span-9 space-y-10">
            {error && (
              <div className="p-4 bg-stone-100 border border-stone-200 rounded-2xl flex items-center gap-3 text-stone-900 text-sm font-bold">
                <Info className="w-5 h-5" />
                {error}
              </div>
            )}

            <div className="flex items-center justify-between">
              <h2 className="text-2xl font-black text-stone-800 tracking-tight">
                Active Notifications <span className="text-stone-800 ml-1">({jobs.length})</span>
              </h2>
            </div>

            {loading ? (
              <div className="flex flex-col items-center justify-center py-20 space-y-4">
                <Loader2 className="w-12 h-12 text-stone-800 animate-spin" />
                <p className="text-stone-500 font-bold uppercase tracking-widest text-[10px]">Synchronizing with Govt Portals...</p>
              </div>
            ) : (
              <>
                {Array.isArray(jobs) && jobs.length === 0 ? (
                  <div className="flex flex-col items-center justify-center py-20 text-center">
                    <div className="w-16 h-16 bg-stone-100 rounded-full flex items-center justify-center mb-4">
                      <ShieldCheck className="w-8 h-8 text-stone-400" />
                    </div>
                    <h3 className="text-xl font-bold text-stone-900 mb-2">No government jobs found</h3>
                    <p className="text-stone-500 text-sm">Check back later for new notifications</p>
                  </div>
                ) : (
                  <div className="grid grid-cols-1 gap-6">
                    {jobs.map((job) => (
                      <JobCard key={job.id} job={job} type="government" />
                    ))}
                  </div>
                )}
              </>
            )}
          </div>
        </div>
      </main>
    </div>
  );
};

export default GovJobsPage;
