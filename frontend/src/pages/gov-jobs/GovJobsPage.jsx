import React, { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { useAuth } from '../../hooks/useAuth';
import api from '../../lib/api';
import JobFilters from '../../components/JobFilters';
import JobCard from '../../components/JobCard';
import { Search, MapPin, ShieldCheck, ArrowRight, Loader2, Info, X } from 'lucide-react';

export const GovJobsPage = () => {
  const router = useRouter();
  const { user, isAuthenticated } = useAuth();
  const [jobs, setJobs] = useState([]);
  const [totalItems, setTotalItems] = useState(0);
  const [pagination, setPagination] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filters, setFilters] = useState({
    location: '',
    department: '',
    category: '',
  });
  const [searchQuery, setSearchQuery] = useState('');
  const [debouncedQuery, setDebouncedQuery] = useState('');

  useEffect(() => {
    const timer = setTimeout(() => setDebouncedQuery(searchQuery), 300);
    return () => clearTimeout(timer);
  }, [searchQuery]);

  const fetchGovJobs = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const queryParams = {
        location: filters.location,
        department: filters.department,
        q: debouncedQuery,
        search: debouncedQuery,
      };

      const res = await api.get('/api/v1/gov-jobs', { params: queryParams });
      const data = res.data?.data || res.data || [];
      const safeData = Array.isArray(data) ? data : [];
      setJobs(safeData);
      setTotalItems(res.data?.pagination?.total || safeData.length);
      setPagination(res.data?.pagination);
    } catch (err) {
      console.error("Fetch error:", err);
      setError("Failed to sync with live government registries. Please try again later.");
      setJobs([]);
    } finally {
      setLoading(false);
    }
  }, [filters, debouncedQuery]);

  useEffect(() => {
    fetchGovJobs();
  }, [fetchGovJobs]);

  // Sync URL with filters and search
  useEffect(() => {
    if (typeof window !== 'undefined') {
      const query = {};
      if (filters.location) query.location = filters.location;
      if (filters.department) query.department = filters.department;
      if (filters.category) query.category = filters.category;
      if (searchQuery) query.search = searchQuery;
      router.replace({ query }, undefined, { shallow: true });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filters, searchQuery]);

  const handleFilterChange = (newFilters) => {
    setFilters(prev => ({ ...prev, ...newFilters }));
  };

  const handleClearAll = () => {
    setFilters({ location: '', department: '', category: '' });
    setSearchQuery('');
  };

  const handleRemoveFilter = (key) => {
    setFilters(prev => ({ ...prev, [key]: '' }));
  };

  // Get active filters for display
  const activeFilters = Object.entries(filters).filter(([_, value]) => value !== '');

  return (
    <div className="min-h-screen bg-stone-50">
      {/* Premium Hero Section */}
      <section className="relative pt-24 pb-20 overflow-hidden bg-white border-b border-stone-200/60">
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

            <div className="flex-1 max-w-md bg-white border border-stone-200 p-10 rounded-2xl shadow-2xl shadow-stone-200/50 relative">
               <div className="absolute -top-6 -right-6 w-24 h-24 bg-stone-800 rounded-full flex items-center justify-center text-white font-black text-xs rotate-12 shadow-xl border-4 border-white">
                  LIVE<br/>AGGREGATOR
               </div>
               <h3 className="text-xl font-black text-stone-800 mb-8">Instant Search</h3>
               <div className="space-y-5">
                  <div className="relative group">
                    <Search className="absolute left-4 top-4 w-5 h-5 text-stone-400 group-focus-within:text-stone-800 transition-colors" />
                    <input
                      id="gov-job-search"
                      name="gov-job-search"
                      type="text"
                      placeholder="Job title or Department..."
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      className="w-full pl-12 pr-4 py-4 bg-stone-50 border-none rounded-2xl text-sm font-bold focus:ring-4 focus:ring-stone-700/10 transition-all outline-none"
                    />
                  </div>
                  <button
                    onClick={fetchGovJobs}
                    className="w-full py-4 bg-stone-800 text-white font-black rounded-2xl hover:bg-stone-900 transition-all shadow-xl shadow-blue-900/20 uppercase tracking-widest text-xs"
                  >
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

            {/* Active Filters Pills */}
            {(activeFilters.length > 0 || searchQuery) && (
              <div className="flex flex-wrap items-center gap-2 mb-6">
                {activeFilters.map(([key, value]) => (
                  <button
                    key={key}
                    onClick={() => handleRemoveFilter(key)}
                    className="inline-flex items-center gap-1 px-3 py-1.5 bg-stone-800 text-white text-xs font-bold rounded-full hover:bg-stone-700 transition-all"
                  >
                    {key}: {value}
                    <X className="w-3 h-3" />
                  </button>
                ))}
                {searchQuery && (
                  <button
                    onClick={() => setSearchQuery('')}
                    className="inline-flex items-center gap-1 px-3 py-1.5 bg-stone-800 text-white text-xs font-bold rounded-full hover:bg-stone-700 transition-all"
                  >
                    Search: {searchQuery}
                    <X className="w-3 h-3" />
                  </button>
                )}
                <button
                  onClick={handleClearAll}
                  className="text-xs font-bold text-stone-500 hover:text-stone-800 underline"
                >
                  Clear All
                </button>
              </div>
            )}

            <div className="flex items-center justify-between">
              <h2 className="text-2xl font-black text-stone-800 tracking-tight">
                Active Notifications <span className="text-stone-400 ml-1 font-medium">({totalItems})</span>
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
