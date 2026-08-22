import React, { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { useAuth } from '../../hooks/useAuth';
import JobFilters from '../../components/JobFilters';
import JobCard from '../../components/JobCard';
import { Search, MapPin, Building2, Calendar, ShieldCheck, ArrowRight } from 'lucide-react';

export const GovJobsPage = () => {
  const { user } = useAuth();
  const router = useRouter();
  const [jobs, setJobs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filters, setFilters] = useState({
    location: '',
    jobType: [],
    category: '',
  });

  const API_BASE = process.env.REACT_APP_BACKEND_URL || 'http://localhost:3001';

  const fetchGovJobs = useCallback(async () => {
    setLoading(true);
    setError(null);
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');

    try {
      // Pass type=government to separate crawler data
      const queryParams = new URLSearchParams({
        type: 'government',
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
        // Fallback to sample data for demo if API fails
        setJobs(SAMPLE_GOV_JOBS);
      }
    } catch (err) {
      setJobs(SAMPLE_GOV_JOBS);
    } finally {
      setLoading(false);
    }
  }, [filters, API_BASE]);

  useEffect(() => {
    fetchGovJobs();
  }, [fetchGovJobs]);

  const handleFilterChange = (newFilters) => {
    setFilters(newFilters);
  };

  return (
    <div className="min-h-screen bg-[#FBFBFB]">
      {/* Dynamic Hero Section */}
      <section className="relative pt-20 pb-24 overflow-hidden bg-white border-b border-slate-100">
        <div className="absolute top-0 right-0 w-1/3 h-full bg-blue-50/50 skew-x-12 transform translate-x-1/2"></div>
        <div className="max-w-7xl mx-auto px-6 relative z-10">
          <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-12">
            <div className="flex-1 space-y-6">
              <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-50 border border-blue-100 text-blue-600 text-xs font-black uppercase tracking-widest">
                <ShieldCheck className="w-3.5 h-3.5" />
                Verified Government Portals
              </div>
              <h1 className="text-5xl md:text-6xl font-black text-slate-900 tracking-tight leading-tight">
                Secure Your <span className="text-blue-600">Future</span> <br />
                in Public Service.
              </h1>
              <p className="text-lg text-slate-500 max-w-xl font-medium leading-relaxed">
                Direct aggregation from SSC, UPSC, Railways, and State PSC portals.
                Zero spam, 100% official notifications.
              </p>

              <div className="flex items-center gap-6 pt-4">
                <div className="flex -space-x-3">
                  {[1,2,3,4].map(i => (
                    <div key={i} className="w-10 h-10 rounded-full border-2 border-white bg-slate-200"></div>
                  ))}
                </div>
                <p className="text-xs font-bold text-slate-400">
                  <span className="text-slate-900">10k+</span> Candidates apply every month
                </p>
              </div>
            </div>

            <div className="flex-1 max-w-md bg-white border border-slate-200 p-8 rounded-[2.5rem] shadow-2xl shadow-slate-200/50">
              <h3 className="text-lg font-black text-slate-900 mb-6">Quick Search</h3>
              <div className="space-y-4">
                <div className="relative">
                  <Search className="absolute left-4 top-3.5 w-5 h-5 text-slate-400" />
                  <input type="text" placeholder="Job title, department..." className="w-full pl-12 pr-4 py-3.5 bg-slate-50 border-none rounded-2xl text-sm font-bold focus:ring-4 focus:ring-blue-500/10 transition-all" />
                </div>
                <div className="relative">
                  <MapPin className="absolute left-4 top-3.5 w-5 h-5 text-slate-400" />
                  <input type="text" placeholder="State / City" className="w-full pl-12 pr-4 py-3.5 bg-slate-50 border-none rounded-2xl text-sm font-bold focus:ring-4 focus:ring-blue-500/10 transition-all" />
                </div>
                <button className="w-full py-4 bg-slate-900 text-white font-black rounded-2xl hover:bg-slate-800 transition-all shadow-xl shadow-slate-900/20 uppercase tracking-widest text-xs">
                  Find Opportunities
                </button>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-6 py-16">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-12">
          {/* Sidebar Filters */}
          <aside className="lg:col-span-3">
            <JobFilters onFilterChange={handleFilterChange} />
          </aside>

          {/* Job Feed */}
          <div className="lg:col-span-9 space-y-8">
            <div className="flex items-center justify-between">
              <h2 className="text-2xl font-black text-slate-900 tracking-tight">
                Latest Notifications <span className="text-blue-600 ml-1">({jobs.length})</span>
              </h2>
              <div className="flex items-center gap-3">
                <span className="text-xs font-bold text-slate-400">Sort by:</span>
                <select className="bg-transparent border-none text-xs font-black text-slate-900 focus:ring-0 cursor-pointer">
                  <option>Newest First</option>
                  <option>Closing Soon</option>
                </select>
              </div>
            </div>

            {loading ? (
              <div className="space-y-6">
                {[1,2,3].map(i => (
                  <div key={i} className="h-48 w-full bg-slate-100 rounded-3xl animate-pulse"></div>
                ))}
              </div>
            ) : jobs.length === 0 ? (
              <div className="bg-white border border-slate-200 rounded-[2.5rem] p-20 text-center space-y-6">
                <div className="w-24 h-24 bg-slate-50 rounded-full flex items-center justify-center mx-auto text-4xl">🔍</div>
                <div>
                  <h3 className="text-xl font-bold text-slate-900">No matches found</h3>
                  <p className="text-slate-400 text-sm font-medium mt-1">Try adjusting your filters or search terms.</p>
                </div>
                <button onClick={() => setFilters({location: '', jobType: [], category: ''})} className="px-8 py-3 bg-slate-900 text-white font-black rounded-2xl text-xs uppercase tracking-widest">
                  Reset Filters
                </button>
              </div>
            ) : (
              <div className="grid grid-cols-1 gap-6">
                {jobs.map((job) => (
                  <div key={job.id} className="group bg-white border border-slate-200/60 rounded-[2rem] p-8 hover:border-blue-400 hover:shadow-2xl hover:shadow-blue-500/5 transition-all relative overflow-hidden">
                    <div className="absolute top-0 right-0 p-8 opacity-0 group-hover:opacity-100 transition-opacity">
                      <ArrowRight className="w-6 h-6 text-blue-600" />
                    </div>

                    <div className="flex flex-col md:flex-row md:items-start gap-8">
                      <div className="w-16 h-16 bg-slate-50 border border-slate-100 rounded-2xl flex items-center justify-center text-xl font-black text-slate-300 shrink-0 group-hover:bg-blue-50 group-hover:border-blue-100 group-hover:text-blue-600 transition-colors">
                        {(job.company_name || job.dept || 'G')[0]}
                      </div>

                      <div className="flex-1 space-y-4">
                        <div className="space-y-1">
                          <span className="text-[10px] font-black uppercase tracking-[0.2em] text-blue-600/60">{job.dept || 'Government Dept'}</span>
                          <h3 className="text-2xl font-black text-slate-900 group-hover:text-blue-600 transition-colors tracking-tight">{job.title}</h3>
                        </div>

                        <div className="flex flex-wrap items-center gap-6 text-xs font-bold text-slate-400">
                          <div className="flex items-center gap-2">
                            <MapPin className="w-4 h-4 text-slate-300" />
                            {job.location || 'All India'}
                          </div>
                          <div className="flex items-center gap-2">
                            <Calendar className="w-4 h-4 text-slate-300" />
                            Apply by {job.last_date || 'N/A'}
                          </div>
                          <div className="flex items-center gap-2 px-3 py-1 bg-slate-50 rounded-lg text-slate-500">
                            💼 {job.type || 'Full Time'}
                          </div>
                        </div>

                        <div className="pt-2 flex flex-wrap gap-2">
                          {(job.tags || ['Public Sector', 'Central Govt']).map(tag => (
                            <span key={tag} className="px-3 py-1 bg-slate-100 text-slate-500 text-[10px] font-black uppercase rounded-lg">
                              {tag}
                            </span>
                          ))}
                        </div>
                      </div>

                      <div className="md:w-48 space-y-3 pt-4 md:pt-0">
                        <Link href={`/jobs/${job.id}`} className="w-full block text-center py-3 bg-slate-900 text-white font-black rounded-xl hover:bg-slate-800 transition-all text-[10px] uppercase tracking-widest shadow-lg shadow-slate-900/10">
                          View Details
                        </Link>
                        <a href={job.apply_url} target="_blank" rel="noreferrer" className="w-full block text-center py-3 bg-white border border-slate-200 text-slate-900 font-black rounded-xl hover:bg-slate-50 transition-all text-[10px] uppercase tracking-widest">
                          Direct Apply
                        </a>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  );
};

const SAMPLE_GOV_JOBS = [
  { id: 1, title: 'Executive Officer 2026', dept: 'Staff Selection Commission', location: 'New Delhi', last_date: '2026-08-30', type: 'Full-Time', tags: ['Central Govt', 'Grade B'] },
  { id: 2, title: 'Civil Services Prelims', dept: 'UPSC', location: 'All India', last_date: '2026-09-15', type: 'Full-Time', tags: ['Gazetted', 'IAS/IPS'] },
  { id: 3, title: 'Senior Section Engineer', dept: 'Indian Railways', location: 'Multiple Zones', last_date: '2026-08-25', type: 'Full-Time', tags: ['Technical', 'Permanent'] },
  { id: 4, title: 'Bank Probationary Officer', dept: 'IBPS', location: 'Pan India', last_date: '2026-09-01', type: 'Full-Time', tags: ['Banking', 'Pan India'] },
];

export default GovJobsPage;
