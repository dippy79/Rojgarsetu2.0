import React from 'react';
import { MapPin, Building2, Calendar, Briefcase, ArrowUpRight, ShieldCheck, Zap, Bookmark } from 'lucide-react';
import Link from 'next/link';

const JobCard = ({ job, type = 'government' }) => {
  const isGov = type === 'government';
  const displayCompany = isGov ? (job?.dept || job?.department) : (job?.company || job?.company_name);

  const handleApplyClick = (e) => {
    const applyUrl = isGov
      ? (job?.apply_url || job?.apply_link || job?.url)
      : (job?.url || job?.apply_url || job?.apply_link);

    if (applyUrl) {
      window.open(applyUrl, '_blank', 'noopener,noreferrer');
    } else {
      console.warn("No apply URL found for job:", job?.id);
    }
  };

  return (
    <div className="group bg-white border border-stone-200 rounded-xl p-8 hover:border-stone-400 hover:shadow-lg transition-all relative overflow-hidden">
      <div className="flex flex-col md:flex-row md:items-start gap-8 relative z-10">
        {/* Logo/Icon Area */}
        <div className={`w-16 h-16 rounded-lg flex items-center justify-center text-xl font-black shrink-0 transition-all duration-300 ${
          isGov
            ? 'bg-stone-100 text-stone-800'
            : 'bg-stone-800 text-white'
        }`}>
          {(displayCompany || 'J')[0]}
        </div>
        
        {/* Info Area */}
        <div className="flex-1 space-y-4">
          <div className="space-y-1">
            <div className="flex items-center gap-2">
              <span className="text-[10px] font-bold uppercase tracking-widest text-stone-400">
                {displayCompany || (isGov ? 'Government Dept' : 'Private Sector')}
              </span>
              {job?.is_verified && (
                <ShieldCheck className="w-3 h-3 text-stone-600" />
              )}
            </div>
            <h3 className="text-2xl font-black text-stone-800 transition-colors tracking-tight leading-tight">
              {job?.title || 'Job Opening'}
            </h3>
          </div>

          <div className="flex flex-wrap items-center gap-6 text-xs font-bold text-stone-500">
            <div className="flex items-center gap-2">
              <MapPin className="w-4 h-4 text-stone-300" />
              {job?.location || 'Location Independent'}
            </div>
            {isGov ? (
              <div className="flex items-center gap-2">
                <Calendar className="w-4 h-4 text-stone-300" />
                Due: {job?.last_date ? new Date(job.last_date).toLocaleDateString() : 'Rolling'}
              </div>
            ) : (
              <div className="flex items-center gap-2">
                <Briefcase className="w-4 h-4 text-stone-300" />
                {job?.salary || 'Market Rate'}
              </div>
            )}
            <div className="px-3 py-1 bg-stone-50 rounded text-stone-600 border border-stone-100 uppercase text-[10px] tracking-tighter">
              {job?.type || job?.job_type || 'Full Time'}
            </div>
          </div>

          {/* Tags */}
          <div className="pt-2 flex flex-wrap gap-2">
            {(job?.skills || job?.tags || (isGov ? ['Public Sector'] : ['Private Sector'])).slice(0, 4).map(tag => (
              <span key={tag} className="px-2 py-1 bg-white border border-stone-100 text-stone-400 text-[9px] font-bold uppercase rounded shadow-sm">
                {tag}
              </span>
            ))}
          </div>
        </div>
        
        {/* Actions Area */}
        <div className="md:w-48 space-y-3 pt-4 md:pt-0 shrink-0">
          <Link
            href={`/jobs/${job?.id}`}
            className="w-full block text-center py-3.5 bg-stone-100 text-stone-800 hover:bg-stone-200 font-bold rounded-lg transition-all text-[10px] uppercase tracking-widest"
          >
            Review Details
          </Link>

          <div className="flex gap-2">
            <button
              onClick={handleApplyClick}
              className="flex-1 block text-center py-3.5 bg-stone-800 text-white hover:bg-stone-900 font-bold rounded-lg transition-all text-[10px] uppercase tracking-widest shadow-lg"
            >
              Direct Apply
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default JobCard;
