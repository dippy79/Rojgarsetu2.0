import React, { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useAuth } from '../../hooks/useAuth';
import { apiUrl } from '../../apiConfig';
import { Search, GraduationCap, Clock, BookOpen, Filter, ArrowRight, Loader2 } from 'lucide-react';

export const CoursesPage = () => {
  const { isAuthenticated } = useAuth();
  const [courses, setCourses] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filters, setFilters] = useState({
    provider: '',
    mode: '',
    level: '',
  });

  const fetchCourses = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const queryParams = new URLSearchParams(filters);
      const response = await fetch(apiUrl(`/api/v1/courses?${queryParams.toString()}`));
      if (!response.ok) throw new Error(`HTTP ${response.status}`);
      const result = await response.json();
      setCourses(result.data || []);
    } catch (err) {
      console.error(err);
      setError("Unable to sync courses. Please try again later.");
    } finally {
      setLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    fetchCourses();
  }, [fetchCourses]);

  return (
    <div className="min-h-screen bg-[#FAFAFA]">
      {/* Academy Hero */}
      <section className="bg-white border-b border-slate-200 pt-24 pb-16">
        <div className="max-w-7xl mx-auto px-6">
          <div className="flex flex-col md:flex-row md:items-end justify-between gap-10">
            <div className="space-y-6">
              <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-purple-50 text-purple-600 text-[10px] font-black uppercase tracking-widest">
                <GraduationCap className="w-4 h-4" /> Rojgar Academy
              </div>
              <h1 className="text-4xl md:text-6xl font-black text-slate-900 tracking-tight leading-none">
                Master the <br /> <span className="text-purple-600">New Economy.</span>
              </h1>
              <p className="text-slate-500 font-medium max-w-lg">Certified learning pathways designed by industry leaders to bridge the skill gap.</p>
            </div>

            <div className="flex-1 max-w-xl">
              <div className="relative group">
                <Search className="absolute left-5 top-5 w-5 h-5 text-slate-400 group-focus-within:text-purple-600 transition-colors" />
                <input
                  type="text"
                  placeholder="Search skills (e.g. AI, Management, Coding)..."
                  className="w-full pl-14 pr-6 py-5 bg-slate-50 border-2 border-transparent focus:border-purple-500/20 focus:bg-white rounded-[2rem] text-sm font-bold shadow-sm transition-all outline-none"
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
            <div className="bg-white border border-slate-200 rounded-[2.5rem] p-8 shadow-sm space-y-8 sticky top-24">
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-black uppercase tracking-widest text-slate-900">Filters</h3>
                <Filter className="w-4 h-4 text-slate-400" />
              </div>

              <div className="space-y-4">
                <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Learning Mode</label>
                <div className="grid grid-cols-1 gap-2">
                  {['Online', 'Hybrid', 'Self-Paced'].map(m => (
                    <button key={m} className="w-full text-left px-4 py-3 rounded-xl border border-slate-100 text-xs font-bold text-slate-600 hover:border-purple-200 transition-all">
                      {m}
                    </button>
                  ))}
                </div>
              </div>

              <div className="space-y-4">
                <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Expertise Level</label>
                <div className="grid grid-cols-1 gap-2">
                  {['Beginner', 'Intermediate', 'Expert'].map(l => (
                    <button key={l} className="w-full text-left px-4 py-3 rounded-xl border border-slate-100 text-xs font-bold text-slate-600 hover:border-purple-200 transition-all">
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
                <Loader2 className="w-10 h-10 text-purple-600 animate-spin mb-4" />
                <span className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Loading Academic Catalog...</span>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                {courses.map((course) => (
                  <div key={course.id} className="group bg-white border border-slate-200 rounded-[2.5rem] p-8 hover:border-purple-400 hover:shadow-2xl hover:shadow-purple-500/5 transition-all">
                    <div className="flex items-center justify-between mb-6">
                      <span className="px-3 py-1 bg-purple-50 text-purple-600 text-[10px] font-black uppercase tracking-widest rounded-lg border border-purple-100">
                        {course.provider_name || 'Premium Academy'}
                      </span>
                      <BookOpen className="w-5 h-5 text-slate-300 group-hover:text-purple-500 transition-colors" />
                    </div>

                    <h3 className="text-2xl font-black text-slate-900 group-hover:text-purple-600 transition-colors mb-4 leading-tight">
                      {course.name || course.title}
                    </h3>

                    <div className="flex flex-wrap gap-4 mb-8">
                      <div className="flex items-center gap-2 text-[10px] font-bold text-slate-400 uppercase tracking-widest">
                        <Clock className="w-3.5 h-3.5" /> {course.duration || '8 Weeks'}
                      </div>
                      <div className="flex items-center gap-2 text-[10px] font-bold text-slate-400 uppercase tracking-widest">
                        <GraduationCap className="w-3.5 h-3.5" /> {course.level || 'Intermediate'}
                      </div>
                    </div>

                    <div className="pt-6 border-t border-slate-100 flex items-center justify-between">
                      <span className="text-xl font-black text-slate-900">
                        {course.fees_amount === 0 ? 'FREE' : `₹${course.fees_amount?.toLocaleString()}`}
                      </span>
                      <button className="flex items-center gap-2 px-6 py-3 bg-slate-900 text-white text-[10px] font-black uppercase tracking-widest rounded-xl hover:bg-purple-600 transition-all shadow-xl shadow-slate-900/10 group">
                        Explore Syllabus <ArrowRight className="w-3 h-3 group-hover:translate-x-1 transition-transform" />
                      </button>
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

export default CoursesPage;
