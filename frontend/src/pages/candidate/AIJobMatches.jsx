import React, { useState, useEffect, useCallback } from 'react';
import { Sparkles, Briefcase, MapPin, Building, Zap, Loader2, ArrowUpRight, Award } from 'lucide-react';
import { useAuth } from '../../hooks/useAuth';
import api from '../../lib/api';

export default function AIJobMatches() {
  const { user } = useAuth();
  const [matches, setMatches] = useState([]);
  const [loading, setLoading] = useState(true);
  const [profile, setProfile] = useState(null);

  const fetchRecommendations = useCallback(async (candidateProfile) => {
    setLoading(true);
    try {
      const res = await api.post('/api/v1/ai/recommend/jobs', {
        user_skills: candidateProfile?.skills || [],
        experience_years: candidateProfile?.experience_years || 0,
        preferred_locations: candidateProfile?.preferred_location || []
      });

      setMatches(res.data.recommendations || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  const fetchProfile = useCallback(async () => {
    try {
      const res = await api.get('/api/v1/candidates/me');
      const data = res.data;
      const p = data.candidate || data.data || data;
      setProfile(p);
      fetchRecommendations(p);
    } catch (err) {
      console.error(err);
      setLoading(false);
    }
  }, [fetchRecommendations]);

  useEffect(() => {
    fetchProfile();
  }, [fetchProfile]);

  return (
    <div className="min-h-screen bg-[#FBFBFB] py-20 px-6 font-sans">
      <div className="max-w-5xl mx-auto space-y-12">
        <header className="space-y-4">
           <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-indigo-50 text-indigo-600 text-[10px] font-black uppercase tracking-widest">
              <Sparkles className="w-3.5 h-3.5" /> Neural Engine v2.0
           </div>
           <h1 className="text-4xl md:text-6xl font-black text-slate-900 tracking-tight leading-none">AI Recommendations.</h1>
           <p className="text-slate-500 font-medium text-lg max-w-2xl">
             Hyper-personalized job matches based on your unique skill graph and professional trajectory.
           </p>
        </header>

        {loading ? (
          <div className="flex flex-col items-center justify-center py-32 space-y-4">
             <Loader2 className="w-12 h-12 text-indigo-600 animate-spin" />
             <p className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Running Match Algorithms...</p>
          </div>
        ) : matches.length === 0 ? (
          <div className="bg-white border border-slate-200 rounded-[3rem] p-20 text-center space-y-6 shadow-sm">
             <div className="w-24 h-24 bg-slate-50 rounded-full flex items-center justify-center mx-auto text-4xl shadow-inner">🧠</div>
             <div>
                <h3 className="text-xl font-bold text-slate-900">Neural Gap Detected</h3>
                <p className="text-slate-400 text-sm font-medium mt-1">Complete your profile skills to enable deep matching.</p>
             </div>
             <button className="px-8 py-3 bg-slate-900 text-white font-black rounded-2xl text-[10px] uppercase tracking-widest hover:bg-indigo-600 transition-all shadow-xl shadow-slate-900/10">Update Profile</button>
          </div>
        ) : (
          <div className="grid gap-6">
            {matches.map((m, idx) => (
              <div key={m.job_id || idx} className="group bg-white border border-slate-200 rounded-[2.5rem] p-8 hover:border-indigo-400 hover:shadow-2xl hover:shadow-indigo-500/5 transition-all relative overflow-hidden flex flex-col md:flex-row items-center justify-between gap-8">
                <div className="absolute top-0 right-0 w-32 h-32 bg-indigo-500/5 rounded-full -mr-16 -mt-16 blur-3xl group-hover:scale-150 transition-transform duration-700"></div>

                <div className="space-y-4 relative z-10 flex-1">
                  <div className="flex items-center gap-3">
                    <div className="px-4 py-1.5 bg-indigo-600 text-white rounded-full flex items-center gap-2 text-[10px] font-black uppercase tracking-widest shadow-lg shadow-indigo-600/20">
                      <Zap className="w-3 h-3 fill-current" />
                      {Math.floor(m.match_score * 100)}% Match
                    </div>
                    <span className="text-[10px] font-black text-slate-400 uppercase tracking-widest">{m.source_table || 'Direct'}</span>
                  </div>

                  <div>
                    <h3 className="text-2xl font-black text-slate-900 group-hover:text-indigo-600 transition-colors tracking-tight">{m.title}</h3>
                    <p className="text-sm font-bold text-slate-400 uppercase tracking-widest mt-1 flex items-center gap-2">
                       <Building className="w-3.5 h-3.5" /> {m.company || 'InnovateCorp'}
                    </p>
                  </div>

                  <div className="flex flex-wrap gap-2">
                    {(m.matched_skills || []).map(s => (
                      <span key={s} className="bg-emerald-50 text-emerald-700 text-[9px] font-black uppercase px-3 py-1 rounded-lg border border-emerald-100 flex items-center gap-1.5">
                        <Award className="w-3 h-3" /> {s}
                      </span>
                    ))}
                  </div>
                </div>

                <div className="flex flex-col items-center gap-4 relative z-10 shrink-0">
                   <div className="flex items-center gap-2 text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">
                      <MapPin className="w-3.5 h-3.5" /> {m.location || 'Global'}
                   </div>
                   <button className="w-full md:w-48 py-4 bg-slate-900 text-white font-black rounded-2xl hover:bg-indigo-600 transition-all shadow-xl shadow-slate-900/10 uppercase text-[10px] tracking-widest flex items-center justify-center gap-3 group">
                     Apply Now <ArrowUpRight className="w-4 h-4 group-hover:translate-x-1 group-hover:-translate-y-1 transition-transform" />
                   </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
