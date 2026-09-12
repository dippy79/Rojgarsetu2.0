import React, { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import api from '../../lib/api';
import { MonitorPlay, Search, Youtube, TrendingUp, Clock, Filter, Loader2, PlayCircle, ArrowUpRight } from 'lucide-react';

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
      const data = res.data?.data || res.data || [];
      const safeData = Array.isArray(data) ? data : [];
      setVideos(safeData);
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
    <div className="min-h-screen bg-stone-50">
      {/* Cinematic Hero */}
      <section className="relative h-[60vh] bg-stone-900 overflow-hidden flex items-center">
        <div className="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1516321318423-f06f85e504b3?q=80&w=2070')] bg-cover bg-center opacity-30 blur-sm"></div>
        <div className="absolute inset-0 bg-gradient-to-r from-stone-950 via-stone-950/80 to-transparent"></div>

        <div className="max-w-7xl mx-auto px-6 relative z-10">
          <div className="space-y-6 max-w-2xl">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-stone-800 text-stone-300 text-[10px] font-black uppercase tracking-widest">
              <Youtube className="w-4 h-4" /> Live Content Stream
            </div>
            <h1 className="text-5xl md:text-7xl font-black text-white tracking-tighter leading-[0.9]">
              Visual Learning <br /> <span className="text-stone-300 italic">Accelerated.</span>
            </h1>
            <p className="text-stone-400 text-lg font-medium leading-relaxed">
              Curated masterclasses from the world's leading industry experts.
              Interview strategies, technical deep-dives, and career guidance.
            </p>
            <div className="pt-4 flex items-center gap-6">
              <button className="px-8 py-4 bg-white text-stone-950 font-black rounded-xl flex items-center gap-3 hover:scale-105 transition-all shadow-xl shadow-white/5 uppercase text-xs tracking-widest">
                <PlayCircle className="w-4 h-4 fill-current" /> Start Watching
              </button>
              <div className="flex items-center gap-2 text-white/60 text-xs font-bold uppercase tracking-widest">
                <TrendingUp className="w-4 h-4 text-stone-400" /> 120+ New Videos This Week
              </div>
            </div>
          </div>
        </div>
      </section>

      <main className="max-w-7xl mx-auto px-6 -mt-16 relative z-20 pb-24">
        {/* Filtering & Search Bar */}
        <div className="bg-white border border-stone-200 p-6 rounded-2xl shadow-2xl flex flex-col md:flex-row items-center gap-6 mb-16">
          <div className="flex-1 relative w-full group">
            <Search className="absolute left-4 top-4 w-5 h-5 text-stone-400 group-focus-within:text-stone-800 transition-colors" />
            <input type="text" placeholder="Search masterclasses..." className="w-full pl-12 pr-4 py-4 bg-stone-50 border-none rounded-xl text-sm font-bold outline-none focus:ring-4 focus:ring-stone-800/5 transition-all" />
          </div>
          <div className="flex items-center gap-4 w-full md:w-auto">
            <select className="flex-1 md:w-48 bg-stone-50 border-none px-6 py-4 rounded-xl text-xs font-black uppercase tracking-widest outline-none text-stone-900">
              <option>All Channels</option>
              <option>Interview Prep</option>
              <option>Coding Deep-dives</option>
            </select>
            <button className="p-4 bg-stone-900 text-white rounded-xl hover:bg-stone-950 transition-all">
              <Filter className="w-5 h-5" />
            </button>
          </div>
        </div>

        {/* Video Grid */}
        <div className="space-y-12">
          <div className="flex items-center justify-between">
            <h2 className="text-2xl font-black text-stone-900 tracking-tight uppercase tracking-widest">Featured Series</h2>
            <div className="h-px flex-1 mx-8 bg-stone-200 hidden md:block"></div>
            <Link href="#" className="text-xs font-black text-stone-600 hover:underline uppercase tracking-[0.2em]">View All Stream →</Link>
          </div>

          {loading ? (
            <div className="flex flex-col items-center justify-center py-20">
              <Loader2 className="w-12 h-12 text-stone-800 animate-spin mb-4" />
              <p className="text-[10px] font-black text-stone-400 uppercase tracking-widest">Buffering Global Feed...</p>
            </div>
          ) : (
            <>
              {Array.isArray(videos) && videos.length === 0 ? (
                <div className="flex flex-col items-center justify-center py-20 text-center">
                  <div className="w-16 h-16 bg-stone-100 rounded-full flex items-center justify-center mb-4">
                    <MonitorPlay className="w-8 h-8 text-stone-400" />
                  </div>
                  <h3 className="text-xl font-bold text-stone-900 mb-2">No videos found</h3>
                  <p className="text-stone-500 text-sm">Check back later for new content</p>
                </div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-10">
                  {Array.isArray(videos) && videos.map((video) => (
                    <div key={video.id} className="bg-white border border-stone-200 rounded-2xl overflow-hidden hover:shadow-md transition-all duration-200 group">
                      <div className="relative aspect-video bg-stone-100">
                        <img
                          src={video.thumbnail_url || video.thumbnail || `https://i.ytimg.com/vi/${video.youtube_id || ''}/mqdefault.jpg`}
                          className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                          alt={video.title}
                        />
                        <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity bg-stone-900/20">
                           <PlayCircle className="w-12 h-12 text-white fill-stone-900/20" />
                        </div>
                      </div>
                      <div className="p-5">
                        <h3 className="font-bold text-stone-900 line-clamp-2 mb-1 leading-tight">{video.title}</h3>
                        <p className="text-[10px] text-stone-400 mb-4 uppercase tracking-widest font-black">{video.channel_name || video.channel || 'Industry Insight'}</p>
                        <button
                          onClick={() => {
                            const url = video.url || video.video_url || `https://youtube.com/watch?v=${video.youtube_id}`;
                            window.open(url, "_blank", "noopener,noreferrer");
                          }}
                          className="inline-flex items-center gap-2 bg-stone-800 text-white text-xs font-bold px-4 py-3 rounded-xl hover:bg-stone-950 w-full justify-center transition-all shadow-lg shadow-stone-900/10"
                        >
                          <MonitorPlay className="w-4 h-4" /> Watch Masterclass
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </>
          )}
        </div>
      </main>
    </div>
  );
};

export default VideosPage;

