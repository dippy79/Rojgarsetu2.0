import React, { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { useAuth } from '../../hooks/useAuth';
import JobFilters from '../../components/JobFilters';
import JobCard from '../../components/JobCard';
import { Search, MapPin, Sparkles, Filter } from 'lucide-react';

export const JobSearch = () => {
  const { user } = useAuth();
  const router = useRouter();
  const [jobs, setJobs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filters, setFilters] = useState({
    location: '',
    jobType: [],
    category: '',
  });

  const API_BASE = process.env.REACT_APP_BACKEND_URL || 'http://localhost:3001';

  const fetchJobs = useCallback(async () => {
    setLoading(true);
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');

    try {
      const queryParams = new URLSearchParams({
        location: filters.location,
        category: filters.category,
        jobType: filters.jobType.join(','),
      });

      const res = await fetch(`${API_BASE}/api/v1/jobs?${queryParams.toString()}`, {
        headers: {
          'Authorization': token ? `Bearer ${token}` : '',
          'Content-Type': 'application/json'
        }
      });

      if (res.ok) {
        const data = await res.json();
        setJobs(data?.jobs ?? data ?? []);
      } else {
        setJobs(SAMPLE_JOBS);
      }
    } catch (err) {
      setJobs(SAMPLE_JOBS);
    } finally {
      setLoading(false);
    }
  }, [filters, API_BASE]);

  useEffect(() => {
    fetchJobs();
  }, [fetchJobs]);

  const handleFilterChange = (newFilters) => {
    setFilters(newFilters);
  };

  return (
    <div className="min-h-screen bg-[#FBFBFB]">
      {/* Search Header */}
      <section className="bg-white border-b border-slate-100 pt-16 pb-12">
        <div className="max-w-7xl mx-auto px-6">
          <div className="flex flex-col md:flex-row md:items-end justify-between gap-8">
            <div className="space-y-4">
              <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-50 text-blue-600 text-[10px] font-black uppercase tracking-widest">
                <Sparkles className="w-3 h-3" /> AI Discovery Engine
              </div>
              <h1 className="text-4xl md:text-5xl font-black text-slate-900 tracking-tight">Explore Opportunities.</h1>
              <p className="text-slate-400 font-medium max-w-lg">Universal search across government registries and premium private boards.</p>
            </div>

            <div className="flex-1 max-w-2xl w-full">
              <div className="relative group">
                <Search className="absolute left-5 top-5 w-5 h-5 text-slate-400 group-focus-within:text-blue-600 transition-colors" />
                <input
                  type="text"
                  placeholder="Search by role, company, or keyword..."
                  className="w-full pl-14 pr-6 py-5 bg-slate-50 border-2 border-transparent focus:border-blue-500/20 focus:bg-white rounded-[2rem] text-sm font-bold shadow-sm transition-all outline-none"
                />
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Main Layout */}
      <main className="max-w-7xl mx-auto px-6 py-12">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-12">
          {/* Filters Sidebar */}
          <aside className="lg:col-span-3">
            <div className="lg:hidden mb-6">
              <button className="w-full py-4 bg-slate-900 text-white rounded-2xl font-black flex items-center justify-center gap-3">
                <Filter className="w-4 h-4" /> Show Filters
              </button>
            </div>
            <div className="hidden lg:block">
              <JobFilters onFilterChange={handleFilterChange} type="all" />
            </div>
          </aside>

          {/* Results Area */}
          <div className="lg:col-span-9 space-y-8">
            <div className="flex items-center justify-between">
              <h2 className="text-xl font-black text-slate-900 tracking-tight uppercase tracking-[0.1em]">
                Available Positions <span className="text-blue-600 ml-2">({jobs.length})</span>
              </h2>
            </div>

            {loading ? (
              <div className="space-y-6">
                {[1,2,3,4].map(i => (
                  <div key={i} className="h-44 bg-slate-100 rounded-[2rem] animate-pulse"></div>
                ))}
              </div>
            ) : (
              <div className="grid grid-cols-1 gap-6">
                {jobs.map(job => (
                  <JobCard key={job.id} job={job} type={job.type === 'government' ? 'government' : 'private'} />
                ))}
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  );
};

const SAMPLE_JOBS = [
  { id: 's1', title: 'Senior Software Architect', company: 'TechFlow Global', location: 'London / Remote', salary: '£90k - 120k', type: 'Full-Time', tags: ['High Fit', 'Cloud'] },
  { id: 's2', title: 'Executive Officer 2026', dept: 'Staff Selection Commission', location: 'New Delhi', last_date: '2026-08-30', type: 'government', tags: ['Central Govt'] },
  { id: 's3', title: 'Product Design Lead', company: 'CreativePulse', location: 'Dubai, UAE', salary: 'AED 30k+', type: 'Full-Time', tags: ['Design', 'Relocation'] },
];

export default JobSearch;
