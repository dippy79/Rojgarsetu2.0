import React, { useState, useEffect } from 'react';
import {
  MapPin, Briefcase, LayoutGrid, ChevronDown, X, Globe,
  Building2, Terminal, Stethoscope, Landmark, Plane,
  ShieldCheck, Search, Filter, SlidersHorizontal
} from 'lucide-react';

const JobFilters = ({ onFilterChange, type = 'government' }) => {
  const [filters, setFilters] = useState({
    location: '',
    jobType: '',
    category: '',
    department: '',
    company: '',
    source: '',
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

  const industryCategories = [
    { id: 'technology', label: 'Technology', icon: Terminal },
    { id: 'healthcare', label: 'Healthcare', icon: Stethoscope },
    { id: 'banking', label: 'Banking', icon: Landmark },
    { id: 'defence', label: 'Defence', icon: ShieldCheck },
    { id: 'railway', label: 'Railways', icon: Building2 },
    { id: 'aviation', label: 'Aviation', icon: Plane },
  ];

  const handleFilterChange = (key, value) => {
    const newFilters = { ...filters, [key]: value };
    setFilters(newFilters);
    onFilterChange(newFilters);
  };

  const clearFilters = () => {
    const cleared = {
      location: '',
      jobType: '',
      category: '',
      department: '',
      company: '',
      source: ''
    };
    setFilters(cleared);
    onFilterChange(cleared);
  };

  return (
    <div className="bg-white/90 backdrop-blur-xl border border-slate-200/60 rounded-[2.5rem] p-8 shadow-2xl shadow-slate-200/30 sticky top-24 space-y-8 transition-all duration-500">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-black text-slate-900 flex items-center gap-3">
          <SlidersHorizontal className="w-5 h-5 text-indigo-600" />
          Refine Feed
        </h2>
        <button
          onClick={clearFilters}
          className="text-[10px] font-black text-slate-400 hover:text-rose-500 transition-colors uppercase tracking-[0.2em]"
        >
          Reset All
        </button>
      </div>

      {/* Location Selector */}
      <div className="space-y-3">
        <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Region / State</label>
        <div className="relative group">
          <MapPin className="absolute left-4 top-4 w-4 h-4 text-slate-400 group-focus-within:text-indigo-600 transition-colors z-10" />
          <select
            value={filters.location}
            onChange={(e) => handleFilterChange('location', e.target.value)}
            className="w-full pl-12 pr-10 py-4 bg-slate-50 border-2 border-transparent focus:border-indigo-500/20 rounded-2xl text-sm font-bold text-slate-900 appearance-none cursor-pointer hover:bg-slate-100 transition-all outline-none"
          >
            {INDIAN_STATES.map(state => (
              <option key={state} value={state === 'All India' ? '' : state}>{state}</option>
            ))}
          </select>
          <ChevronDown className="absolute right-4 top-4 w-4 h-4 text-slate-400 pointer-events-none" />
        </div>
      </div>

      {/* Job Type for Private Jobs */}
      {type === 'private' && (
        <div className="space-y-4">
          <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Job Commitment</label>
          <div className="flex flex-wrap gap-2">
            {jobTypes.map((t) => (
              <button
                key={t.id}
                onClick={() => handleFilterChange('jobType', filters.jobType === t.id ? '' : t.id)}
                className={`px-4 py-2.5 rounded-xl border-2 text-xs font-bold transition-all ${
                  filters.jobType === t.id
                    ? 'bg-slate-900 border-slate-900 text-white shadow-lg'
                    : 'bg-white border-slate-100 text-slate-600 hover:border-indigo-200'
                }`}
              >
                {t.label}
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Department/Company Specific Search */}
      <div className="space-y-3">
        <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">
          {type === 'government' ? 'Department / Body' : 'Company Name'}
        </label>
        <div className="relative group">
          <Building2 className="absolute left-4 top-4 w-4 h-4 text-slate-400 group-focus-within:text-indigo-600 transition-colors z-10" />
          <input
            type="text"
            placeholder={type === 'government' ? 'e.g. SSC, UPSC...' : 'e.g. Google, Razorpay...'}
            value={type === 'government' ? filters.department : filters.company}
            onChange={(e) => handleFilterChange(type === 'government' ? 'department' : 'company', e.target.value)}
            className="w-full pl-12 pr-4 py-4 bg-slate-50 border-2 border-transparent focus:border-indigo-500/20 rounded-2xl text-sm font-bold text-slate-900 outline-none transition-all placeholder:text-slate-300"
          />
        </div>
      </div>

      {/* Category Grid */}
      <div className="space-y-4">
        <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Industry Focus</label>
        <div className="grid grid-cols-2 gap-3">
          {industryCategories.map((cat) => {
            const Icon = cat.icon;
            const active = filters.category === cat.id;
            return (
              <button
                key={cat.id}
                onClick={() => handleFilterChange('category', active ? '' : cat.id)}
                className={`flex flex-col items-center gap-3 p-4 rounded-3xl border-2 transition-all ${
                  active
                    ? 'bg-indigo-50 border-indigo-600 shadow-inner'
                    : 'bg-white border-slate-100 hover:border-slate-200 hover:bg-slate-50'
                }`}
              >
                <Icon className={`w-5 h-5 ${active ? 'text-indigo-600' : 'text-slate-400'}`} />
                <span className={`text-[10px] font-black uppercase tracking-tighter ${active ? 'text-indigo-700' : 'text-slate-500'}`}>
                  {cat.label}
                </span>
              </button>
            );
          })}
        </div>
      </div>

      {/* Premium International Toggle */}
      <div className="pt-6 border-t border-slate-100">
        <div className="flex items-center justify-between p-4 bg-emerald-50 rounded-2xl border border-emerald-100 group cursor-pointer hover:bg-emerald-100 transition-colors">
          <div className="flex items-center gap-3">
            <Globe className="w-5 h-5 text-emerald-600" />
            <div>
              <p className="text-xs font-black text-emerald-900">International</p>
              <p className="text-[9px] text-emerald-600 font-bold uppercase">UAE • Europe • US</p>
            </div>
          </div>
          <div className="w-10 h-6 bg-slate-200 rounded-full relative">
            <div className="absolute left-1 top-1 w-4 h-4 bg-white rounded-full shadow-sm transition-transform"></div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default JobFilters;
