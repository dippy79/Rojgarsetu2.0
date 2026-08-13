import React, { useEffect, useState } from 'react';

export default function PublicCounter() {
  const [stats, setStats] = useState({ total_jobs: 0, total_candidates: 0, total_companies: 0, total_placements: 0 });

  useEffect(() => {
    fetch('http://localhost:3001/api/v1/stats')
      .then(r => r.ok ? r.json() : null)
      .then(data => { if (data) setStats(data); });
  }, []);

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-6 bg-white/80 backdrop-blur-xl border border-slate-100 p-8 rounded-2xl shadow-sm">
      <div className="text-center">
        <p className="text-4xl font-extrabold font-mono text-indigo-600">{stats.total_jobs}</p>
        <p className="text-xs text-slate-500 font-semibold uppercase mt-1">Available Jobs</p>
      </div>
      <div className="text-center">
        <p className="text-4xl font-extrabold font-mono text-blue-600">{stats.total_candidates}</p>
        <p className="text-xs text-slate-500 font-semibold uppercase mt-1">Candidates</p>
      </div>
      <div className="text-center">
        <p className="text-4xl font-extrabold font-mono text-emerald-600">{stats.total_companies}</p>
        <p className="text-xs text-slate-500 font-semibold uppercase mt-1">Companies</p>
      </div>
      <div className="text-center">
        <p className="text-4xl font-extrabold font-mono text-amber-600">{stats.total_placements}</p>
        <p className="text-xs text-slate-500 font-semibold uppercase mt-1">Placements</p>
      </div>
    </div>
  );
}
