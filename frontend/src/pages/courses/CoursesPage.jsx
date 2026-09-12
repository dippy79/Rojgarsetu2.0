import React, { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useAuth } from '../../hooks/useAuth';
import api from '../../lib/api';
import { Search, GraduationCap, Clock, BookOpen, Filter, ArrowRight, Loader2 } from 'lucide-react';

export const CoursesPage = () => {
  const { isAuthenticated } = useAuth();
  const [courses, setCourses] = useState([]);
  const [providers, setProviders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filters, setFilters] = useState({
    provider: '',
    mode: '',
    level: '',
  });

  const fetchInitialData = useCallback(async () => {
    try {
      const provRes = await api.get('/api/v1/courses/providers');
      setProviders(provRes.data.data || []);
    } catch (err) {
      console.error("Failed to fetch providers:", err);
    }
  }, []);

  const fetchCourses = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await api.get('/api/v1/courses', { params: filters });
      const data = res.data?.data || res.data || [];
      const safeData = Array.isArray(data) ? data : [];
      setCourses(safeData);
    } catch (err) {
      console.error(err);
      setError("Unable to sync courses. Please try again later.");
    } finally {
      setLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    fetchInitialData();
  }, [fetchInitialData]);

  useEffect(() => {
    fetchCourses();
  }, [fetchCourses]);

  const dynamicProviders = React.useMemo(() => {
    const fromApi = providers.map(p => p.name || p);
    const fromCourses = courses.map(c => c.provider_name || c.provider || c.source || c.platform);
    return [...new Set([...fromApi, ...fromCourses])].filter(Boolean).sort();
  }, [providers, courses]);

  return (
    <div className="min-h-screen bg-stone-50">
      {/* Academy Hero */}
      <section className="bg-white border-b border-stone-200 pt-24 pb-16">
        <div className="max-w-7xl mx-auto px-6">
          <div className="flex flex-col md:flex-row md:items-end justify-between gap-10">
            <div className="space-y-6">
              <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-stone-100 text-stone-600 text-[10px] font-black uppercase tracking-widest">
                <GraduationCap className="w-4 h-4" /> Rojgar Academy
              </div>
              <h1 className="text-4xl md:text-6xl font-black text-stone-800 tracking-tight leading-none">
                Master the <br /> <span className="text-stone-950">New Economy.</span>
              </h1>
              <p className="text-stone-500 font-medium max-w-lg">Certified learning pathways designed by industry leaders to bridge the skill gap.</p>
            </div>

            <div className="flex-1 max-w-xl">
              <div className="relative group">
                <Search className="absolute left-5 top-5 w-5 h-5 text-stone-400 group-focus-within:text-stone-800 transition-colors" />
                <input
                  type="text"
                  placeholder="Search skills (e.g. AI, Management, Coding)..."
                  className="w-full pl-14 pr-6 py-5 bg-stone-50 border-2 border-transparent focus:border-stone-800/10 focus:bg-white rounded-xl text-sm font-bold shadow-sm transition-all outline-none"
                />
              </div>
            </div>
          </div>
        </div>
      </section>

      <main className="max-w-7xl mx-auto px-6 py-16">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-12">
          {/* Filtering Sidebar */}
          <aside className="lg:col-span-3 space-y-8">
            <div className="bg-white border border-stone-200 rounded-2xl p-8 shadow-sm space-y-8 sticky top-24">
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-black uppercase tracking-widest text-stone-900">Filters</h3>
                <Filter className="w-4 h-4 text-stone-400" />
              </div>

              <div className="space-y-4">
                <label className="text-[10px] font-black text-stone-400 uppercase tracking-widest">Platform</label>
                <div className="grid grid-cols-1 gap-2 max-h-48 overflow-y-auto pr-2 custom-scrollbar">
                  <button
                    onClick={() => setFilters(f => ({ ...f, provider: '' }))}
                    className={`w-full text-left px-4 py-3 rounded-xl border text-xs font-bold transition-all ${!filters.provider ? 'bg-stone-800 text-white border-stone-800 shadow-md' : 'bg-white border-stone-100 text-stone-600 hover:border-stone-200'}`}
                  >
                    All Platforms
                  </button>
                  {dynamicProviders.map(provider => (
                    <button
                      key={provider}
                      onClick={() => setFilters(f => ({ ...f, provider: provider }))}
                      className={`w-full text-left px-4 py-3 rounded-xl border text-xs font-bold transition-all ${filters.provider === provider ? 'bg-stone-800 text-white border-stone-800 shadow-md' : 'bg-white border-stone-100 text-stone-600 hover:border-stone-200'}`}
                    >
                      {provider}
                    </button>
                  ))}
                </div>
              </div>

              <div className="space-y-4">
                <label className="text-[10px] font-black text-stone-400 uppercase tracking-widest">Learning Mode</label>
                <div className="grid grid-cols-1 gap-2">
                  {['Online', 'Hybrid', 'Self-Paced'].map(m => (
                    <button key={m} className="w-full text-left px-4 py-3 rounded-xl border border-stone-100 text-xs font-bold text-stone-600 hover:border-stone-200 transition-all">
                      {m}
                    </button>
                  ))}
                </div>
              </div>

              <div className="space-y-4">
                <label className="text-[10px] font-black text-stone-400 uppercase tracking-widest">Expertise Level</label>
                <div className="grid grid-cols-1 gap-2">
                  {['Beginner', 'Intermediate', 'Expert'].map(l => (
                    <button key={l} className="w-full text-left px-4 py-3 rounded-xl border border-stone-100 text-xs font-bold text-stone-600 hover:border-stone-200 transition-all">
                      {l}
                    </button>
                  ))}
                </div>
              </div>
            </div>
          </aside>

          {/* Courses Feed */}
          <div className="lg:col-span-9 space-y-10">
            {loading ? (
              <div className="flex flex-col items-center justify-center py-20">
                <Loader2 className="w-10 h-10 text-stone-800 animate-spin mb-4" />
                <span className="text-[10px] font-black text-stone-400 uppercase tracking-widest">Loading Academic Catalog...</span>
              </div>
            ) : (
              <>
                {Array.isArray(courses) && courses.length === 0 ? (
                  <div className="flex flex-col items-center justify-center py-20 text-center">
                    <div className="w-16 h-16 bg-stone-100 rounded-full flex items-center justify-center mb-4">
                      <GraduationCap className="w-8 h-8 text-stone-400" />
                    </div>
                    <h3 className="text-xl font-bold text-stone-900 mb-2">No courses found</h3>
                    <p className="text-stone-500 text-sm">Check back later for new learning opportunities</p>
                  </div>
                ) : (
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                    {courses.map((course) => (
                      <div key={course.id} className="group bg-white border border-stone-200 rounded-2xl p-8 hover:border-stone-400 hover:shadow-2xl transition-all">
                        <div className="flex items-center justify-between mb-6">
                          <span className="px-3 py-1 bg-stone-50 text-stone-600 text-[10px] font-black uppercase tracking-widest rounded-lg border border-stone-100">
                            {course.provider_name || 'Premium Academy'}
                          </span>
                          <BookOpen className="w-5 h-5 text-stone-300 group-hover:text-stone-800 transition-colors" />
                        </div>

                        <h3 className="text-2xl font-black text-stone-900 group-hover:text-stone-800 transition-colors mb-4 leading-tight">
                          {course.name || course.title}
                        </h3>

                        <div className="flex flex-wrap gap-4 mb-8">
                          <div className="flex items-center gap-2 text-[10px] font-bold text-stone-400 uppercase tracking-widest">
                            <Clock className="w-3.5 h-3.5" /> {course.duration || '8 Weeks'}
                          </div>
                          <div className="flex items-center gap-2 text-[10px] font-bold text-stone-400 uppercase tracking-widest">
                            <GraduationCap className="w-3.5 h-3.5" /> {course.level || 'Intermediate'}
                          </div>
                        </div>

                        <div className="pt-6 border-t border-stone-100 flex items-center justify-between">
                          <span className="text-xl font-black text-stone-900">
                            {course.price === '0' || course.is_free ? 'FREE' : `₹${course.price || course.fees_amount || 'Free'}`}
                          </span>
                          <button
                            onClick={() => {
                              const url = course.url || course.apply_link || course.course_url;
                              if (url) window.open(url, '_blank', 'noopener,noreferrer');
                            }}
                            className="flex items-center gap-2 px-6 py-3 bg-stone-900 text-white text-[10px] font-black uppercase tracking-widest rounded-xl hover:bg-stone-950 transition-all shadow-xl shadow-stone-900/10 group"
                          >
                            Explore Syllabus <ArrowRight className="w-3 h-3 group-hover:translate-x-1 transition-transform" />
                          </button>
                        </div>
                      </div>
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

export default CoursesPage;
