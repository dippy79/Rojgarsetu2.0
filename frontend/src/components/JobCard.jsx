import React from 'react';
import { MapPin, Building2, Calendar, Briefcase, ArrowUpRight, ShieldCheck, Zap } from 'lucide-react';
import Link from 'next/link';

const JobCard = ({ job, type = 'government' }) => {
  const isGov = type === 'government';
  const displayCompany = isGov ? job.dept || job.department : job.company || job.company_name;

  return (
    <div className="group bg-white border border-slate-200/60 rounded-[2rem] p-8 hover:border-blue-400 hover:shadow-2xl hover:shadow-blue-500/5 transition-all relative overflow-hidden">
      {/* Background Accent */}
      <div className={`absolute top-0 right-0 w-32 h-32 ${isGov ? 'bg-blue-500/5' : 'bg-emerald-500/5'} rounded-full -mr-16 -mt-16 blur-3xl group-hover:scale-150 transition-transform duration-700`}></div>

      <div className="flex flex-col md:flex-row md:items-start gap-8 relative z-10">
        {/* Logo/Icon Area */}
        <div className={`w-16 h-16 rounded-2xl flex items-center justify-center text-xl font-black shrink-0 transition-all duration-300 ${
          isGov
            ? 'bg-blue-50 border border-blue-100 text-blue-600 group-hover:bg-blue-600 group-hover:text-white'
            : 'bg-emerald-50 border border-emerald-100 text-emerald-600 group-hover:bg-emerald-600 group-hover:text-white'
        }`}>
          {(displayCompany || 'G')[0]}
        </div>
        
        {/* Info Area */}
        <div className="flex-1 space-y-4">
          <div className="space-y-1">
            <div className="flex items-center gap-2">
              <span className={`text-[10px] font-black uppercase tracking-[0.2em] ${isGov ? 'text-blue-600/60' : 'text-emerald-600/60'}`}>
                {displayCompany || (isGov ? 'Government Dept' : 'Private Sector')}
              </span>
              {job.is_verified && (
                <ShieldCheck className="w-3 h-3 text-blue-500" />
              )}
              {job.scam_score < 0.2 && isGov && (
                <div className="flex items-center gap-1 px-1.5 py-0.5 bg-emerald-50 text-emerald-600 rounded text-[8px] font-black uppercase tracking-tighter">
                  <Zap className="w-2 h-2" /> Genuine
                </div>
              )}
            </div>
            <h3 className={`text-2xl font-black text-slate-900 group-hover:${isGov ? 'text-blue-600' : 'text-emerald-600'} transition-colors tracking-tight leading-tight`}>
              {job.title}
            </h3>
          </div>

          <div className="flex flex-wrap items-center gap-6 text-xs font-bold text-slate-400">
            <div className="flex items-center gap-2">
              <MapPin className="w-4 h-4 text-slate-300" />
              {job.location || 'All India'}
            </div>
            {isGov ? (
              <div className="flex items-center gap-2">
                <Calendar className="w-4 h-4 text-slate-300" />
                Apply by {job.last_date ? new Date(job.last_date).toLocaleDateString() : 'N/A'}
              </div>
            ) : (
              <div className="flex items-center gap-2">
                <Briefcase className="w-4 h-4 text-slate-300" />
                {job.salary || 'Competitive'}
              </div>
            )}
            <div className="flex items-center gap-2 px-3 py-1 bg-slate-50 rounded-lg text-slate-500 border border-slate-100">
              💼 {job.type || job.job_type || 'Full Time'}
            </div>
          </div>

          {/* Tags */}
          <div className="pt-2 flex flex-wrap gap-2">
            {(job.skills || job.tags || (isGov ? ['Public Sector', 'Central Govt'] : ['Tech', 'Innovation'])).slice(0, 4).map(tag => (
              <span key={tag} className="px-3 py-1 bg-white border border-slate-100 text-slate-500 text-[10px] font-black uppercase rounded-lg shadow-sm group-hover:border-slate-200 transition-colors">
                {tag}
              </span>
            ))}
          </div>
        </div>
        
        {/* Actions Area */}
        <div className="md:w-48 space-y-3 pt-4 md:pt-0 shrink-0">
          <Link
            href={`/jobs/${job.id}`}
            className={`w-full block text-center py-3.5 text-white font-black rounded-xl transition-all text-[10px] uppercase tracking-widest shadow-xl ${
              isGov
                ? 'bg-slate-900 hover:bg-blue-600 shadow-slate-900/10'
                : 'bg-slate-900 hover:bg-emerald-600 shadow-slate-900/10'
            }`}
          >
            Review Details
          </Link>

          <div className="flex gap-2">
            <a
              href={isGov ? job.apply_url : job.url}
              target="_blank"
              rel="noreferrer"
              className="flex-1 block text-center py-3.5 bg-white border border-slate-200 text-slate-900 font-black rounded-xl hover:bg-slate-50 hover:border-slate-300 transition-all text-[10px] uppercase tracking-widest"
            >
              Direct Apply
            </a>
            {isGov && (
              <button
                title="Track Results"
                className="px-4 bg-indigo-50 border border-indigo-100 text-indigo-600 rounded-xl hover:bg-indigo-600 hover:text-white transition-all shadow-sm group/track"
              >
                <Bookmark className="w-4 h-4 group-hover/track:fill-current" />
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default JobCard;
