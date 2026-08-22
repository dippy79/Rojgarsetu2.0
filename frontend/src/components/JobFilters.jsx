import React, { useState, useEffect } from 'react';
import { MapPin, Briefcase, LayoutGrid, ChevronDown, X, Globe, Building2, Terminal, Stethoscope, Landmark, Plane } from 'lucide-react';

const JobFilters = ({ onFilterChange, type = 'government' }) => {
  const [isOpen, setIsOpen] = useState(true);
  const [filters, setFilters] = useState({
    location: '',
    jobType: [],
    category: '',
    experience: '',
  });

  const INDIAN_STATES = [
    'All India', 'Andhra Pradesh', 'Arunachal Pradesh', 'Assam', 'Bihar', 'Chhattisgarh', 'Delhi', 'Goa', 'Gujarat',
    'Haryana', 'Himachal Pradesh', 'Jammu & Kashmir', 'Jharkhand', 'Karnataka', 'Kerala', 'Madhya Pradesh',
    'Maharashtra', 'Manipur', 'Meghalaya', 'Mizoram', 'Nagaland', 'Odisha', 'Punjab', 'Rajasthan', 'Sikkim',
    'Tamil Nadu', 'Telangana', 'Tripura', 'Uttar Pradesh', 'Uttarakhand', 'West Bengal'
  ];

  const jobTypes = [
    { id: 'full-time', label: 'Full-Time' },
    { id: 'part-time', label: 'Part-Time' },
    { id: 'internship', label: 'Internship' },
    { id: 'contract', label: 'Contract' },
    { id: 'remote', label: 'Remote' },
  ];

  const categories = [
    { id: 'tech', label: 'Technology', icon: Terminal, color: 'text-blue-500' },
    { id: 'medical', label: 'Healthcare', icon: Stethoscope, color: 'text-rose-500' },
    { id: 'finance', label: 'Banking & Finance', icon: Landmark, color: 'text-emerald-500' },
    { id: 'defence', label: 'Defence', icon: ShieldCheck, color: 'text-slate-700' },
    { id: 'railway', label: 'Railways', icon: Building2, color: 'text-amber-600' },
    { id: 'aviation', label: 'Aviation', icon: Plane, color: 'text-cyan-500' },
  ];

  const handleTypeToggle = (typeId) => {
    const updatedTypes = filters.jobType.includes(typeId)
      ? filters.jobType.filter(t => t !== typeId)
      : [...filters.jobType, typeId];

    const newFilters = { ...filters, jobType: updatedTypes };
    setFilters(newFilters);
    onFilterChange(newFilters);
  };

  const handleLocationChange = (val) => {
    const newFilters = { ...filters, location: val === 'All India' ? '' : val };
    setFilters(newFilters);
    onFilterChange(newFilters);
  };

  const handleCategoryChange = (catId) => {
    const newFilters = { ...filters, category: catId === filters.category ? '' : catId };
    setFilters(newFilters);
    onFilterChange(newFilters);
  };

  const clearFilters = () => {
    const cleared = { location: '', jobType: [], category: '', experience: '' };
    setFilters(cleared);
    onFilterChange(cleared);
  };

  return (
    <div className="bg-white/80 backdrop-blur-2xl border border-slate-200/50 rounded-[2.5rem] p-8 shadow-2xl shadow-slate-200/40 sticky top-24 space-y-10 overflow-hidden">
      <div className="absolute top-0 right-0 w-32 h-32 bg-blue-500/5 rounded-full -mr-16 -mt-16 blur-3xl"></div>

      <div className="flex items-center justify-between relative z-10">
        <h2 className="text-xl font-black text-slate-900 flex items-center gap-3">
          <div className="p-2 bg-slate-900 rounded-xl">
            <LayoutGrid className="w-5 h-5 text-blue-400" />
          </div>
          Refine Search
        </h2>
        <button
          onClick={clearFilters}
          className="text-xs font-black text-blue-600 hover:text-red-500 transition-colors uppercase tracking-tighter"
        >
          Reset
        </button>
      </div>

      {/* State / Location Dropdown */}
      <div className="space-y-4">
        <label className="text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] ml-1">Preferred Location</label>
        <div className="relative group">
          <MapPin className="absolute left-4 top-4 w-4 h-4 text-slate-400 group-focus-within:text-blue-600 transition-colors z-10" />
          <select
            value={filters.location || 'All India'}
            onChange={(e) => handleLocationChange(e.target.value)}
            className="w-full pl-12 pr-10 py-4 bg-slate-50 border-2 border-transparent focus:border-blue-500/20 rounded-[1.25rem] text-sm font-bold text-slate-900 appearance-none cursor-pointer hover:bg-slate-100 transition-all outline-none"
          >
            {INDIAN_STATES.map(state => (
              <option key={state} value={state}>{state}</option>
            ))}
          </select>
          <ChevronDown className="absolute right-4 top-4 w-4 h-4 text-slate-400 pointer-events-none" />
        </div>
      </div>

      {/* Employment Type */}
      <div className="space-y-4">
        <label className="text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] ml-1">Job Commitment</label>
        <div className="grid grid-cols-1 gap-2.5">
          {jobTypes.map((t) => {
            const active = filters.jobType.includes(t.id);
            return (
              <button
                key={t.id}
                onClick={() => handleTypeToggle(t.id)}
                className={`flex items-center justify-between px-5 py-4 rounded-2xl border-2 transition-all group ${
                  active
                    ? 'bg-slate-900 border-slate-900 text-white shadow-xl shadow-slate-900/20'
                    : 'bg-white border-slate-100 text-slate-600 hover:border-blue-200'
                }`}
              >
                <span className="flex items-center gap-3 text-sm font-bold">
                  <div className={`w-2 h-2 rounded-full transition-all ${active ? 'bg-blue-400 scale-125' : 'bg-slate-200 group-hover:bg-slate-300'}`}></div>
                  {t.label}
                </span>
                {active && <div className="w-5 h-5 bg-white/10 rounded-lg flex items-center justify-center"><X className="w-3 h-3" /></div>}
              </button>
            );
          })}
        </div>
      </div>

      {/* Modern Category Selector */}
      <div className="space-y-4">
        <label className="text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] ml-1">Industry Focus</label>
        <div className="grid grid-cols-2 gap-3">
          {categories.map((cat) => {
            const active = filters.category === cat.id;
            const Icon = cat.icon;
            return (
              <button
                key={cat.id}
                onClick={() => handleCategoryChange(cat.id)}
                className={`flex flex-col items-center gap-3 p-4 rounded-3xl border-2 transition-all ${
                  active
                    ? 'bg-blue-50 border-blue-600 shadow-inner'
                    : 'bg-white border-slate-100 hover:border-slate-200 hover:bg-slate-50'
                }`}
              >
                <Icon className={`w-6 h-6 ${active ? 'text-blue-600' : 'text-slate-400'}`} />
                <span className={`text-[10px] font-black uppercase tracking-tighter ${active ? 'text-blue-700' : 'text-slate-500'}`}>
                  {cat.label}
                </span>
              </button>
            );
          })}
        </div>
      </div>

      {/* Global Opportunities Toggle (International Standard) */}
      <div className="pt-6 border-t border-slate-100 space-y-4">
        <div className="flex items-center justify-between p-4 bg-emerald-50 rounded-2xl border border-emerald-100 group cursor-pointer hover:bg-emerald-100 transition-colors">
          <div className="flex items-center gap-3">
            <Globe className="w-5 h-5 text-emerald-600" />
            <div>
              <p className="text-xs font-black text-emerald-900 tracking-tight">International Roles</p>
              <p className="text-[10px] text-emerald-600 font-bold uppercase">UAE • Europe • US</p>
            </div>
          </div>
          <div className="w-10 h-6 bg-slate-200 rounded-full relative transition-colors group-hover:bg-emerald-300">
            <div className="absolute left-1 top-1 w-4 h-4 bg-white rounded-full shadow-sm transition-transform"></div>
          </div>
        </div>
      </div>
    </div>
  );
};

// Helper component for Icon
const ShieldCheck = ({className}) => (
  <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
  </svg>
);

export default JobFilters;
