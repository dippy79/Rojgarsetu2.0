import React, { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import api from '../../lib/api';
import { Play, Search, Youtube, TrendingUp, Clock, Filter, Loader2, ArrowUpRight } from 'lucide-react';

export const VideosPage = () => {
  const [videos, setVideos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filters, setFilters] = useState({
    channel: '',
    category: '',
  });

  const fetchVideos = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await api.get('/api/v1/videos', { params: filters });
      setVideos(res.data.data || []);
    } catch (err) {
      console.error(err);
      setError("Video feed sync failed. Please try again later.");
    } finally {
      setLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    fetchVideos();
  }, [fetchVideos]);

  return (
    <div className="min-h-screen bg-[#F0F2F5]">
      {/* Cinematic Hero */}
      <section className="relative h-[60vh] bg-slate-900 overflow-hidden flex items-center">
        <div className="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1516321318423-f06f85e504b3?q=80&w=2070')] bg-cover bg-center opacity-30 blur-sm"></div>
        <div className="absolute inset-0 bg-gradient-to-r from-slate-950 via-slate-950/80 to-transparent"></div>

        <div className="max-w-7xl mx-auto px-6 relative z-10">
          <div className="space-y-6 max-w-2xl">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-red-500/10 border border-red-500/20 text-red-400 text-[10px] font-black uppercase tracking-widest">
              <Youtube className="w-4 h-4" /> Live Content Stream
            </div>
            <h1 className="text-5xl md:text-7xl font-black text-white tracking-tighter leading-[0.9]">
              Visual Learning <br /> <span className="text-red-500 italic">Accelerated.</span>
            </h1>
            <p className="text-slate-400 text-lg font-medium leading-relaxed">
              Curated masterclasses from the world's leading industry experts.
              Interview strategies, technical deep-dives, and career guidance.
            </p>
            <div className="pt-4 flex items-center gap-6">
              <button className="px-8 py-4 bg-white text-slate-950 font-black rounded-2xl flex items-center gap-3 hover:scale-105 transition-all shadow-xl shadow-white/5 uppercase text-xs tracking-widest">
                <Play className="w-4 h-4 fill-current" /> Start Watching
              </button>
              <div className="flex items-center gap-2 text-white/60 text-xs font-bold uppercase tracking-widest">
                <TrendingUp className="w-4 h-4 text-emerald-400" /> 120+ New Videos This Week
              </div>
            </div>
          </div>
        </div>
      </section>

      <main className="max-w-7xl mx-auto px-6 -mt-16 relative z-20 pb-24">
        {/* Filtering & Search Bar */}
        <div className="bg-white/80 backdrop-blur-xl border border-white/20 p-6 rounded-[2.5rem] shadow-2xl flex flex-col md:flex-row items-center gap-6 mb-16">
          <div className="flex-1 relative w-full group">
            <Search className="absolute left-4 top-4 w-5 h-5 text-slate-400 group-focus-within:text-red-500 transition-colors" />
            <input type="text" placeholder="Search masterclasses..." className="w-full pl-12 pr-4 py-4 bg-slate-50 border-none rounded-2xl text-sm font-bold outline-none focus:ring-4 focus:ring-red-500/10 transition-all" />
          </div>
          <div className="flex items-center gap-4 w-full md:w-auto">
            <select className="flex-1 md:w-48 bg-slate-50 border-none px-6 py-4 rounded-2xl text-xs font-black uppercase tracking-widest outline-none">
              <option>All Channels</option>
              <option>Interview Prep</option>
              <option>Coding Deep-dives</option>
            </select>
            <button className="p-4 bg-slate-950 text-white rounded-2xl hover:bg-red-600 transition-all">
              <Filter className="w-5 h-5" />
            </button>
          </div>
        </div>

        {/* Video Grid */}
        <div className="space-y-12">
          <div className="flex items-center justify-between">
            <h2 className="text-2xl font-black text-slate-900 tracking-tight uppercase tracking-widest">Featured Series</h2>
            <div className="h-px flex-1 mx-8 bg-slate-200 hidden md:block"></div>
            <Link href="#" className="text-xs font-black text-red-600 hover:underline uppercase tracking-[0.2em]">View All Stream →</Link>
          </div>

          {loading ? (
            <div className="flex flex-col items-center justify-center py-20">
              <Loader2 className="w-12 h-12 text-red-600 animate-spin mb-4" />
              <p className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Buffering Global Feed...</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-10">
              {videos.map((video) => (
                <div key={video.id} className="group cursor-pointer">
                  <div className="relative aspect-video rounded-[2rem] overflow-hidden shadow-xl mb-6">
                    <img src={video.thumbnail_url || 'https://images.unsplash.com/photo-1611162617474-5b21e879e113?q=80&w=1974'} className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-700" alt={video.title} />
                    <div className="absolute inset-0 bg-slate-900/40 group-hover:bg-slate-900/20 transition-all"></div>
                    <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-all transform scale-150 group-hover:scale-100">
                      <div className="w-16 h-16 bg-white rounded-full flex items-center justify-center shadow-2xl">
                        <Play className="w-6 h-6 fill-slate-900 ml-1" />
                      </div>
                    </div>
                    <div className="absolute bottom-4 right-4 px-2 py-1 bg-black/80 backdrop-blur-md rounded-lg text-[10px] font-bold text-white uppercase tracking-widest">
                      {video.duration || '12:45'}
                    </div>
                  </div>

                  <div className="space-y-3 px-2">
                    <div className="flex items-center gap-2">
                      <span className="px-2 py-0.5 bg-red-50 text-red-600 text-[9px] font-black uppercase rounded-md tracking-tighter">
                        {video.channel_name || 'Industry Lead'}
                      </span>
                      <div className="w-1 h-1 rounded-full bg-slate-300"></div>
                      <span className="text-[9px] font-black text-slate-400 uppercase tracking-tighter">{video.views || '12k'} Views</span>
                    </div>
                    <h3 className="text-xl font-black text-slate-900 group-hover:text-red-600 transition-colors leading-tight">
                      {video.title}
                    </h3>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </main>
    </div>
  );
};

export default VideosPage;
