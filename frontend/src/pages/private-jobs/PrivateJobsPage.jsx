import React, { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useAuth } from '../../hooks/useAuth';
import api from '../../lib/api';
import JobFilters from '../../components/JobFilters';
import JobCard from '../../components/JobCard';
import { Search, Sparkles, Loader2, Info } from 'lucide-react';

export const PrivateJobsPage = () => {
  const { isAuthenticated } = useAuth();
  const [jobs, setJobs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filters, setFilters] = useState({
    location: '',
    jobType: '',
    company: '',
  });

  const fetchPrivateJobs = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const queryParams = {
        location: filters.location,
        job_type: filters.jobType,
        company: filters.company,
      };

      const res = await api.get('/api/v1/priv-jobs', { params: queryParams });
      const data = res.data?.data || res.data || [];
      const safeData = Array.isArray(data) ? data : [];
      setJobs(safeData);
    } catch (err) {
      console.error("Fetch error:", err);
      setError("AI-aggregator offline. Please try again later.");
    } finally {
      setLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    fetchPrivateJobs();
  }, [fetchPrivateJobs]);

  const handleFilterChange = (newFilters) => {
    setFilters(prev => ({ ...prev, ...newFilters }));
  };

  return (
    <div className="min-h-screen bg-stone-50">
      {/* Global Standard Hero */}
      <section className="relative pt-24 pb-20 overflow-hidden bg-stone-900 border-b border-stone-800">
        <div className="absolute inset-0 bg-[radial-gradient(#ffffff0a_1px,transparent_1px)] [background-size:32px_32px]"></div>
        <div className="max-w-7xl mx-auto px-6 relative z-10 text-center lg:text-left">
          <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-16">
            <div className="flex-1 space-y-8">
              <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-stone-800 border border-stone-700 text-stone-400 text-[10px] font-black uppercase tracking-[0.2em]">
                <Sparkles className="w-3.5 h-3.5" />
                World-Class Career Portals
              </div>
              <h1 className="text-5xl md:text-8xl font-black text-white tracking-tight leading-[0.9]">
                Elevate Your <br />
                <span className="text-stone-400 italic">Trajectory.</span>
              </h1>
              <p className="text-lg text-stone-500 max-w-xl font-medium leading-relaxed">
                We aggregate high-impact roles from top-tier tech startups to global conglomerates.
                Sourced from LinkedIn, Indeed, and direct company career pages.
              </p>
            </div>

            <div className="flex-1 max-w-md mx-auto lg:mx-0 bg-white border border-stone-200 p-12 rounded-2xl shadow-2xl relative group">
               <h3 className="text-xl font-black text-stone-900 mb-10 relative tracking-tight">Refine Search</h3>
               <div className="space-y-6 relative">
                  <div className="relative">
                    <Search className="absolute left-4 top-4 w-5 h-5 text-stone-400" />
                    <input type="text" placeholder="Title, Skill, or Firm" className="w-full pl-12 pr-4 py-4 bg-stone-50 border-none rounded-xl text-sm font-bold text-stone-800 focus:ring-4 focus:ring-stone-800/5 transition-all outline-none" />
                  </div>
                  <button className="w-full py-5 bg-stone-800 text-white font-black rounded-xl hover:bg-stone-900 transition-all shadow-xl shadow-stone-800/20 uppercase tracking-widest text-[11px]">
                    Analyze Opportunities
                  </button>
               </div>
            </div>
          </div>
        </div>
      </section>

      {/* Grid Layout */}
      <main className="max-w-7xl mx-auto px-6 py-20">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-16">
          <aside className="lg:col-span-3">
            <JobFilters onFilterChange={handleFilterChange} type="private" />
          </aside>

          <div className="lg:col-span-9 space-y-12">
            {error && (
              <div className="p-4 bg-stone-100 text-stone-700 rounded-xl border border-stone-200 text-xs font-black uppercase tracking-widest">
                {error}
              </div>
            )}

            <div className="flex items-center justify-between border-b border-stone-100 pb-8">
              <h2 className="text-3xl font-black text-stone-800 tracking-tighter">
                Global Openings <span className="text-stone-400 ml-2 font-medium">({jobs.length})</span>
              </h2>
            </div>

            {loading ? (
              <div className="flex flex-col items-center justify-center py-20 space-y-4">
                <Loader2 className="w-12 h-12 text-stone-800 animate-spin" />
                <p className="text-stone-500 font-bold uppercase tracking-widest text-[10px]">Filtering Global Boards...</p>
              </div>
            ) : (
              <>
                {Array.isArray(jobs) && jobs.length === 0 ? (
                  <div className="flex flex-col items-center justify-center py-20 text-center">
                    <div className="w-16 h-16 bg-stone-100 rounded-full flex items-center justify-center mb-4">
                      <Sparkles className="w-8 h-8 text-stone-400" />
                    </div>
                    <h3 className="text-xl font-bold text-stone-900 mb-2">No private jobs found</h3>
                    <p className="text-stone-500 text-sm">Adjust your filters or check back later</p>
                  </div>
                ) : (
                  <div className="grid grid-cols-1 gap-8">
                    {jobs.map((job) => (
                      <JobCard key={job.id} job={job} type="private" />
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

export default PrivateJobsPage;
