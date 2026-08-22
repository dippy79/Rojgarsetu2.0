import React, { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { useAuth } from '../../hooks/useAuth';
import JobFilters from '../../components/JobFilters';
import JobCard from '../../components/JobCard';
import { Search, MapPin, Building2, Calendar, Briefcase, ArrowRight, Sparkles } from 'lucide-react';

export const PrivateJobsPage = () => {
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

  const fetchPrivateJobs = useCallback(async () => {
    setLoading(true);
    setError(null);
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');

    try {
      // Pass type=private to separate crawler data
      const queryParams = new URLSearchParams({
        type: 'private',
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
        setJobs(SAMPLE_PRIVATE_JOBS);
      }
    } catch (err) {
      setJobs(SAMPLE_PRIVATE_JOBS);
    } finally {
      setLoading(false);
    }
  }, [filters, API_BASE]);

  useEffect(() => {
    fetchPrivateJobs();
  }, [fetchPrivateJobs]);

  const handleFilterChange = (newFilters) => {
    setFilters(newFilters);
  };

  return (
    <div className="min-h-screen bg-[#FBFBFB]">
      {/* Dynamic Hero Section */}
      <section className="relative pt-20 pb-24 overflow-hidden bg-slate-900 border-b border-slate-800">
        <div className="absolute inset-0 bg-[radial-gradient(#ffffff0a_1px,transparent_1px)] [background-size:24px_24px]"></div>
        <div className="max-w-7xl mx-auto px-6 relative z-10 text-center lg:text-left">
          <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-12">
            <div className="flex-1 space-y-8">
              <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-[10px] font-black uppercase tracking-[0.2em]">
                <Sparkles className="w-3.5 h-3.5" />
                Premium Career Opportunities
              </div>
              <h1 className="text-5xl md:text-7xl font-black text-white tracking-tight leading-tight">
                Architect Your <br />
                <span className="text-emerald-400 italic">Professional Journey.</span>
              </h1>
              <p className="text-lg text-slate-400 max-w-xl font-medium leading-relaxed">
                Aggregating the world's most innovative tech, finance, and creative roles.
                Optimized by AI for international career mobility.
              </p>
            </div>

            <div className="flex-1 max-w-md mx-auto lg:mx-0 bg-white/5 backdrop-blur-2xl border border-white/10 p-10 rounded-[3rem] shadow-2xl">
              <h3 className="text-lg font-black text-white mb-8">Refine Search</h3>
              <div className="space-y-5">
                <div className="relative group">
                  <Search className="absolute left-4 top-4 w-5 h-5 text-slate-500 group-focus-within:text-emerald-400 transition-colors" />
                  <input type="text" placeholder="Title, Skill, or Company" className="w-full pl-12 pr-4 py-4 bg-white/5 border-none rounded-2xl text-sm font-bold text-white focus:ring-4 focus:ring-emerald-500/10 transition-all placeholder:text-slate-600" />
                </div>
                <div className="relative group">
                  <MapPin className="absolute left-4 top-4 w-5 h-5 text-slate-500 group-focus-within:text-emerald-400 transition-colors" />
                  <input type="text" placeholder="Location (e.g. Dubai, Bengaluru)" className="w-full pl-12 pr-4 py-4 bg-white/5 border-none rounded-2xl text-sm font-bold text-white focus:ring-4 focus:ring-emerald-500/10 transition-all placeholder:text-slate-600" />
                </div>
                <button className="w-full py-4 bg-emerald-500 text-slate-950 font-black rounded-2xl hover:bg-emerald-400 transition-all shadow-xl shadow-emerald-500/10 uppercase tracking-widest text-xs mt-4">
                  Browse Positions
                </button>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-6 py-20">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-16">
          {/* Sidebar Filters */}
          <aside className="lg:col-span-3">
            <JobFilters onFilterChange={handleFilterChange} type="private" />
          </aside>

          {/* Job Feed */}
          <div className="lg:col-span-9 space-y-10">
            <div className="flex items-center justify-between border-b border-slate-100 pb-6">
              <h2 className="text-3xl font-black text-slate-900 tracking-tighter">
                Global Openings <span className="text-emerald-500 ml-2">({jobs.length})</span>
              </h2>
              <div className="flex items-center gap-6">
                <div className="flex items-center gap-2">
                  <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></div>
                  <span className="text-[10px] font-black text-emerald-600 uppercase tracking-widest">Live Updates</span>
                </div>
                <select className="bg-transparent border-none text-xs font-black text-slate-900 focus:ring-0 cursor-pointer">
                  <option>Top Relevance</option>
                  <option>Newest First</option>
                </select>
              </div>
            </div>

            {loading ? (
              <div className="space-y-8">
                {[1,2,3].map(i => (
                  <div key={i} className="h-56 w-full bg-slate-50 border border-slate-100 rounded-[2.5rem] animate-pulse"></div>
                ))}
              </div>
            ) : jobs.length === 0 ? (
              <div className="bg-white border border-slate-200 rounded-[3rem] p-24 text-center space-y-8">
                <div className="w-28 h-24 bg-slate-50 rounded-full flex items-center justify-center mx-auto text-5xl">🔭</div>
                <div>
                  <h3 className="text-2xl font-black text-slate-900">Expanding Horizons</h3>
                  <p className="text-slate-400 text-sm font-medium mt-2 max-w-sm mx-auto leading-relaxed">We couldn't find a direct match. Try broadening your criteria for better results.</p>
                </div>
                <button onClick={() => setFilters({location: '', jobType: [], category: ''})} className="px-10 py-3.5 bg-slate-900 text-white font-black rounded-2xl text-[10px] uppercase tracking-[0.2em] hover:bg-slate-800 transition-all">
                  Reset Filters
                </button>
              </div>
            ) : (
              <div className="grid grid-cols-1 gap-8">
                {jobs.map((job) => (
                  <div key={job.id} className="group bg-white border border-slate-200/60 rounded-[2.5rem] p-10 hover:border-emerald-400 hover:shadow-2xl hover:shadow-emerald-500/5 transition-all relative">
                    <div className="absolute top-0 right-0 p-10 opacity-0 group-hover:opacity-100 transition-all transform translate-x-4 group-hover:translate-x-0">
                      <ArrowRight className="w-8 h-8 text-emerald-500" />
                    </div>

                    <div className="flex flex-col md:flex-row md:items-start gap-10">
                      <div className="w-20 h-20 bg-slate-50 border border-slate-100 rounded-3xl flex items-center justify-center text-2xl font-black text-slate-300 shrink-0 group-hover:bg-emerald-50 group-hover:border-emerald-100 group-hover:text-emerald-600 transition-all">
                        {(job.company || job.company_name || 'C')[0]}
                      </div>

                      <div className="flex-1 space-y-6">
                        <div className="space-y-2">
                          <div className="flex items-center gap-3">
                            <span className="text-[10px] font-black uppercase tracking-[0.3em] text-emerald-600">{job.company || job.company_name || 'InnovateCorp'}</span>
                            <div className="w-1 h-1 rounded-full bg-slate-300"></div>
                            <span className="text-[10px] font-black uppercase tracking-[0.3em] text-slate-400">{job.type || 'Full Time'}</span>
                          </div>
                          <h3 className="text-3xl font-black text-slate-900 group-hover:text-emerald-600 transition-colors tracking-tighter leading-none">{job.title}</h3>
                        </div>

                        <div className="flex flex-wrap items-center gap-8 text-xs font-bold text-slate-500">
                          <div className="flex items-center gap-2.5">
                            <MapPin className="w-4 h-4 text-emerald-500" />
                            {job.location || 'Global Remote'}
                          </div>
                          <div className="flex items-center gap-2.5">
                            <Briefcase className="w-4 h-4 text-emerald-500" />
                            {job.salary || 'Competitive Package'}
                          </div>
                          <div className="px-4 py-1.5 bg-slate-900 text-white rounded-full text-[10px] font-black uppercase tracking-widest shadow-xl shadow-slate-900/10">
                            Verified Opportunity
                          </div>
                        </div>

                        <div className="flex flex-wrap gap-2 pt-2">
                          {(job.skills || ['Agile', 'Strategy', 'Scaling']).map(tag => (
                            <span key={tag} className="px-4 py-1.5 bg-slate-50 text-slate-500 text-[10px] font-black uppercase rounded-xl border border-slate-100 hover:border-emerald-200 hover:bg-emerald-50 transition-colors">
                              {tag}
                            </span>
                          ))}
                        </div>
                      </div>

                      <div className="md:w-56 space-y-4 pt-6 md:pt-4">
                        <Link href={`/jobs/${job.id}`} className="w-full block text-center py-4 bg-slate-900 text-white font-black rounded-2xl hover:bg-slate-800 transition-all text-xs uppercase tracking-widest shadow-2xl shadow-slate-900/10">
                          Review Detail
                        </Link>
                        <a href={job.url || job.apply_url} target="_blank" rel="noreferrer" className="w-full block text-center py-4 bg-white border-2 border-slate-100 text-slate-900 font-black rounded-2xl hover:border-emerald-400 hover:text-emerald-600 transition-all text-xs uppercase tracking-widest">
                          Quick Apply
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

const SAMPLE_PRIVATE_JOBS = [
  { id: 1, title: 'Lead Architecture Engineer', company: 'TechScale Global', location: 'London / Remote', salary: '£85k - 110k', type: 'Full-Time', skills: ['Next.js', 'Go', 'Kubernetes'] },
  { id: 2, title: 'Head of Growth Marketing', company: 'FinFlow UAE', location: 'Dubai, UAE', salary: 'AED 35k - 45k monthly', type: 'Contract', skills: ['Growth', 'Analytics', 'Retention'] },
  { id: 3, title: 'Senior UX Strategist', company: 'Studio 24', location: 'Remote', salary: '$120k - 150k', type: 'Full-Time', skills: ['Figma', 'Strategy', 'Design Ops'] },
  { id: 4, title: 'Backend Systems Developer', company: 'CloudNexus', location: 'Bengaluru', salary: '₹25 - 35 LPA', type: 'Full-Time', skills: ['Python', 'PostgreSQL', 'Redis'] },
];

export default PrivateJobsPage;
