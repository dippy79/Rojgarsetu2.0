import React, { useState, useEffect } from 'react';
import { Briefcase, Users, Building2, Zap } from 'lucide-react';

const PublicCounter = () => {
  const [stats, setStats] = useState({
    total_jobs: 0,
    total_candidates: 0,
    total_companies: 0,
    total_placements: 0
  });

  useEffect(() => {
    fetch('http://localhost:3001/api/v1/stats')
      .then(res => res.ok ? res.json() : {})
      .then(data => {
        if (data) setStats(data);
      })
      .catch(err => console.error("Stats fetch error:", err));
  }, []);

  return (
    <div className="grid grid-cols-2 lg:grid-cols-4 gap-8 py-12 px-6 bg-slate-900 rounded-[3rem] shadow-2xl shadow-indigo-900/20 relative overflow-hidden">
      <div className="absolute top-0 left-0 w-full h-full bg-[radial-gradient(#ffffff0a_1px,transparent_1px)] [background-size:24px_24px]"></div>

      <StatItem icon={<Briefcase className="w-6 h-6 text-blue-400" />} label="Live Jobs" value={stats.total_jobs} color="blue" />
      <StatItem icon={<Users className="w-6 h-6 text-emerald-400" />} label="Candidates" value={stats.total_candidates} color="emerald" />
      <StatItem icon={<Building2 className="w-6 h-6 text-indigo-400" />} label="Companies" value={stats.total_companies} color="indigo" />
      <StatItem icon={<Zap className="w-6 h-6 text-amber-400" />} label="Placements" value={stats.total_placements} color="amber" />
    </div>
  );
};

const StatItem = ({ icon, label, value, color }) => {
  const [displayValue, setDisplayValue] = useState(0);

  useEffect(() => {
    let start = 0;
    const end = parseInt(value, 10) || 0;
    if (end === 0) return;

    const duration = 2000;
    const increment = end / (duration / 16);

    const timer = setInterval(() => {
      start += increment;
      if (start >= end) {
        setDisplayValue(end);
        clearInterval(timer);
      } else {
        setDisplayValue(Math.floor(start));
      }
    }, 16);

    return () => clearInterval(timer);
  }, [value]);

  return (
    <div className="flex flex-col items-center text-center space-y-4 relative z-10">
      <div className={`p-4 bg-white/5 rounded-2xl border border-white/10 shadow-xl`}>
        {icon}
      </div>
      <div>
        <p className="text-4xl font-black text-white tracking-tighter font-mono">{displayValue.toLocaleString()}+</p>
        <p className="text-[10px] font-black text-slate-500 uppercase tracking-[0.2em] mt-2">{label}</p>
      </div>
    </div>
  );
};

export default PublicCounter;
