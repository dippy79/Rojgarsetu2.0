import React, { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useAuth } from '../../hooks/useAuth';
import JobFilters from '../../components/JobFilters';
import JobCard from '../../components/JobCard';
import { Search, Sparkles, Loader2, Info } from 'lucide-react';
import { apiUrl } from '../../apiConfig';

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
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');

    try {
      const queryParams = new URLSearchParams({
        location: filters.location,
        job_type: filters.jobType,
        company: filters.company,
      });

      const response = await fetch(apiUrl(`/api/v1/priv-jobs?${queryParams.toString()}`), {
        headers: {
          'Authorization': token ? `Bearer ${token}` : '',
          'Content-Type': 'application/json'
        }
      });

      if (!response.ok) throw new Error(`HTTP Error: ${response.status}`);

      const result = await response.json();
      setJobs(result.data || []);
    } catch (err) {
      console.error("Fetch error:", err);
      setError("AI-aggregator offline. Displaying premium archive.");
      setJobs(SAMPLE_PRIVATE_JOBS);
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
    <div className="min-h-screen bg-[#FBFBFB]">
      {/* Global Standard Hero */}
      <section className="relative pt-24 pb-20 overflow-hidden bg-slate-950 border-b border-slate-800">
        <div className="absolute inset-0 bg-[radial-gradient(#ffffff0a_1px,transparent_1px)] [background-size:32px_32px]"></div>
        <div className="max-w-7xl mx-auto px-6 relative z-10 text-center lg:text-left">
          <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-16">
            <div className="flex-1 space-y-8">
              <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-[10px] font-black uppercase tracking-[0.2em]">
                <Sparkles className="w-3.5 h-3.5" />
                World-Class Career Portals
              </div>
              <h1 className="text-5xl md:text-8xl font-black text-white tracking-tight leading-[0.9]">
                Elevate Your <br />
                <span className="text-emerald-400 italic">Trajectory.</span>
              </h1>
              <p className="text-lg text-slate-400 max-w-xl font-medium leading-relaxed">
                We aggregate high-impact roles from top-tier tech startups to global conglomerates.
                Sourced from LinkedIn, Indeed, and direct company career pages.
              </p>
            </div>

            <div className="flex-1 max-w-md mx-auto lg:mx-0 bg-white/5 backdrop-blur-3xl border border-white/10 p-12 rounded-[4rem] shadow-2xl relative group">
               <div className="absolute inset-0 bg-emerald-500/5 blur-[80px] rounded-full opacity-0 group-hover:opacity-100 transition-opacity"></div>
               <h3 className="text-xl font-black text-white mb-10 relative">Refine Search</h3>
               <div className="space-y-6 relative">
                  <div className="relative">
                    <Search className="absolute left-4 top-4 w-5 h-5 text-slate-500" />
                    <input type="text" placeholder="Title, Skill, or Firm" className="w-full pl-12 pr-4 py-4 bg-white/5 border-none rounded-2xl text-sm font-bold text-white focus:ring-4 focus:ring-emerald-500/10 transition-all outline-none" />
                  </div>
                  <button className="w-full py-5 bg-emerald-500 text-slate-950 font-black rounded-2xl hover:bg-emerald-400 transition-all shadow-xl shadow-emerald-500/20 uppercase tracking-widest text-[11px]">
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
              <div className="p-4 bg-emerald-50 text-emerald-700 rounded-2xl border border-emerald-100 text-xs font-black uppercase tracking-widest">
                {error}
              </div>
            )}

            <div className="flex items-center justify-between border-b border-slate-100 pb-8">
              <h2 className="text-3xl font-black text-slate-900 tracking-tighter">
                Global Openings <span className="text-emerald-500 ml-2">({jobs.length})</span>
              </h2>
            </div>

            {loading ? (
              <div className="flex flex-col items-center justify-center py-20 space-y-4">
                <Loader2 className="w-12 h-12 text-emerald-500 animate-spin" />
                <p className="text-slate-400 font-bold uppercase tracking-widest text-[10px]">Filtering Global Boards...</p>
              </div>
            ) : (
              <div className="grid grid-cols-1 gap-8">
                {jobs.map((job) => (
                  <JobCard key={job.id} job={job} type="private" />
                ))}
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  );
};

const SAMPLE_PRIVATE_JOBS = [
  { id: 'p1', title: 'Senior AI Engineer', company: 'TechFlow Global', location: 'Dubai, UAE / Remote', salary: '$120k - $160k', type: 'Full-Time', is_verified: true, tags: ['AI', 'Next.js'] },
  { id: 'p2', title: 'Product Strategist', company: 'NexusScale', location: 'London, UK', salary: '£80k+', type: 'Contract', is_verified: true, tags: ['UX', 'Growth'] },
];

export default PrivateJobsPage;
