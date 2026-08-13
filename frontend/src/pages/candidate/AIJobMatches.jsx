import React, { useState, useEffect } from 'react';
import { Sparkles, Briefcase, MapPin, Building } from 'lucide-react';

export default function AIJobMatches() {
  const [matches, setMatches] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('rojgar_token');
    fetch('http://localhost:3001/api/v1/ai/recommend/jobs', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify({ skills: ['React', 'Node.js'], experience_years: 2 })
    })
      .then(r => r.ok ? r.json() : [])
      .then(data => {
        setMatches(Array.isArray(data) ? data : []);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, []);

  return (
    <div className="p-8 max-w-5xl mx-auto">
      <div className="flex items-center gap-3 mb-6">
        <Sparkles className="w-8 h-8 text-indigo-600" />
        <h1 className="text-2xl font-bold text-slate-900">AI Recommendation Engine</h1>
      </div>

      {loading ? (
        <p className="text-slate-500">Analyzing your profile and matching opportunities...</p>
      ) : matches.length === 0 ? (
        <div className="text-center py-12 bg-white rounded-2xl border border-slate-100">
          <p className="text-slate-600 font-medium">Complete your profile for AI matches</p>
        </div>
      ) : (
        <div className="grid gap-4">
          {matches.map(m => (
            <div key={m.job_id} className="bg-white/80 backdrop-blur-xl border border-slate-100 p-6 rounded-2xl shadow-sm flex items-center justify-between">
              <div className="space-y-2">
                <div className="flex items-center gap-3">
                  <span className="text-lg font-bold text-slate-900">{m.title}</span>
                  <span className="px-3 py-1 bg-indigo-50 text-indigo-600 font-bold rounded-full text-xs">
                    {m.score}% Match
                  </span>
                </div>
                <div className="flex gap-2 flex-wrap">
                  {m.matching_skills?.map(s => (
                    <span key={s} className="bg-emerald-50 text-emerald-700 text-xs px-2.5 py-0.5 rounded-md font-medium">
                      ✓ {s}
                    </span>
                  ))}
                </div>
              </div>
              <button className="bg-indigo-600 hover:bg-indigo-700 text-white px-5 py-2.5 rounded-xl font-medium text-sm transition">
                Apply Now
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
