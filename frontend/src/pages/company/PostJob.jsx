import React, { useState } from 'react';
import { useRouter } from 'next/router';
import { ArrowLeft, Briefcase, MapPin, DollarSign, Calendar, Sparkles, Send, Loader2, Globe, Laptop } from 'lucide-react';
import Link from 'next/link';

export default function PostJob() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState({
    title: '', description: '', requirements: '', location: '',
    job_type: 'full-time', salary_min: '', salary_max: '', currency_code: 'INR',
    experience_years: 0, skills_required: '', is_remote: false, expires_at: ''
  });

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({ ...prev, [name]: type === 'checkbox' ? checked : value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');
    const payload = {
      ...formData,
      skills_required: formData.skills_required.split(',').map(s => s.trim()).filter(Boolean),
      salary_min: Number(formData.salary_min),
      salary_max: Number(formData.salary_max),
      experience_years: Number(formData.experience_years)
    };

    try {
      const res = await fetch('http://localhost:3001/api/v1/company-jobs', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
        body: JSON.stringify(payload)
      });

      if (res.ok) {
        router.push('/dashboard/company');
      } else {
        alert("Failed to publish job. Please try again.");
      }
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#FBFBFB] py-20 px-6 font-sans">
      <div className="max-w-4xl mx-auto space-y-12">
        {/* Navigation & Header */}
        <div className="flex items-center justify-between">
           <div className="flex items-center gap-6">
              <Link href="/dashboard/company" className="p-4 bg-white border border-slate-200 rounded-3xl hover:bg-slate-50 transition-all shadow-sm">
                 <ArrowLeft className="w-6 h-6 text-slate-900" />
              </Link>
              <div>
                 <h1 className="text-4xl font-black text-slate-900 tracking-tight">Post Opening.</h1>
                 <p className="text-slate-400 font-bold uppercase tracking-widest text-[10px] mt-1">Broadcast to 10k+ Certified Candidates</p>
              </div>
           </div>

           <div className="hidden md:flex items-center gap-2 px-4 py-2 bg-indigo-50 text-indigo-600 rounded-2xl border border-indigo-100 text-[10px] font-black uppercase tracking-widest">
              <Sparkles className="w-3.5 h-3.5" /> AI Ranking Enabled
           </div>
        </div>

        <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-12 gap-8">
           {/* Left Section: Core Details */}
           <div className="md:col-span-8 space-y-8">
              <div className="bg-white border border-slate-200 rounded-[2.5rem] p-10 shadow-sm space-y-8">
                 <div className="space-y-4">
                    <label htmlFor="title" className="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Position Title</label>
                    <div className="relative group">
                       <Briefcase className="absolute left-4 top-4 w-5 h-5 text-slate-300 group-focus-within:text-indigo-600 transition-colors" />
                       <input id="title" name="title" value={formData.title} onChange={handleChange} required placeholder="e.g. Senior Systems Architect" className="w-full pl-12 pr-6 py-4 bg-slate-50 border-none rounded-2xl text-sm font-bold outline-none focus:ring-4 focus:ring-indigo-500/10 transition-all" />
                    </div>
                 </div>

                 <div className="space-y-4">
                    <label htmlFor="description" className="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Role Description</label>
                    <textarea id="description" name="description" rows={6} value={formData.description} onChange={handleChange} required placeholder="Detail the core responsibilities and team culture..." className="w-full px-6 py-4 bg-slate-50 border-none rounded-3xl text-sm font-bold outline-none focus:ring-4 focus:ring-indigo-500/10 transition-all resize-none" />
                 </div>

                 <div className="space-y-4">
                    <label htmlFor="requirements" className="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Core Requirements</label>
                    <textarea id="requirements" name="requirements" rows={4} value={formData.requirements} onChange={handleChange} placeholder="Technical stacks, education, and specific traits..." className="w-full px-6 py-4 bg-slate-50 border-none rounded-3xl text-sm font-bold outline-none focus:ring-4 focus:ring-indigo-500/10 transition-all resize-none" />
                 </div>
              </div>

              <div className="bg-white border border-slate-200 rounded-[2.5rem] p-10 shadow-sm space-y-8">
                 <h3 className="text-sm font-black text-slate-900 uppercase tracking-widest flex items-center gap-2">
                    <Sparkles className="w-4 h-4 text-indigo-600" /> Professional Fit
                 </h3>

                 <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                    <div className="space-y-4">
                       <label htmlFor="skills_required" className="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Skills (Comma Separated)</label>
                       <input id="skills_required" name="skills_required" value={formData.skills_required} onChange={handleChange} placeholder="Next.js, Python, AWS" className="w-full px-6 py-4 bg-slate-50 border-none rounded-2xl text-sm font-bold outline-none focus:ring-4 focus:ring-indigo-500/10 transition-all" />
                    </div>
                    <div className="space-y-4">
                       <label htmlFor="experience_years" className="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Min. Experience (Years)</label>
                       <input id="experience_years" name="experience_years" type="number" value={formData.experience_years} onChange={handleChange} className="w-full px-6 py-4 bg-slate-50 border-none rounded-2xl text-sm font-bold outline-none focus:ring-4 focus:ring-indigo-500/10 transition-all" />
                    </div>
                 </div>
              </div>
           </div>

           {/* Right Section: Metadata */}
           <div className="md:col-span-4 space-y-8">
              <div className="bg-slate-900 rounded-[2.5rem] p-8 text-white space-y-8 shadow-2xl shadow-indigo-900/20">
                 <h3 className="text-xs font-black uppercase tracking-[0.2em] text-indigo-400">Logistics</h3>

                 <div className="space-y-6">
                    <div className="space-y-3">
                       <label htmlFor="location" className="text-[10px] font-black text-slate-500 uppercase tracking-widest">Office Hub</label>
                       <div className="relative group">
                          <MapPin className="absolute left-3 top-3 w-4 h-4 text-slate-600" />
                          <input id="location" name="location" value={formData.location} onChange={handleChange} required placeholder="e.g. Dubai, UAE" className="w-full pl-10 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl text-xs font-bold outline-none focus:border-indigo-500 transition-all" />
                       </div>
                    </div>

                    <div className="space-y-3">
                       <label htmlFor="job_type" className="text-[10px] font-black text-slate-500 uppercase tracking-widest">Commitment</label>
                       <select id="job_type" name="job_type" value={formData.job_type} onChange={handleChange} className="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-xs font-bold outline-none appearance-none cursor-pointer">
                          <option value="full-time" className="bg-slate-900 text-white">Full-time</option>
                          <option value="part-time" className="bg-slate-900 text-white">Part-time</option>
                          <option value="contract" className="bg-slate-900 text-white">Contract</option>
                          <option value="internship" className="bg-slate-900 text-white">Internship</option>
                       </select>
                    </div>

                    <div className="flex items-center justify-between p-4 bg-white/5 rounded-2xl border border-white/5">
                       <div className="flex items-center gap-3">
                          <Laptop className="w-4 h-4 text-indigo-400" />
                          <span className="text-[10px] font-black uppercase tracking-widest">Remote Friendly</span>
                       </div>
                       <input id="is_remote" name="is_remote" type="checkbox" checked={formData.is_remote} onChange={handleChange} className="w-5 h-5 rounded-lg border-none bg-white/10 text-indigo-600 focus:ring-0 cursor-pointer" />
                    </div>
                 </div>

                 <div className="h-px bg-white/10"></div>

                 <div className="space-y-6">
                    <div className="space-y-3">
                       <label className="text-[10px] font-black text-slate-500 uppercase tracking-widest">Compensation Range</label>
                       <div className="flex gap-2">
                          <input name="salary_min" type="number" value={formData.salary_min} onChange={handleChange} placeholder="Min" className="w-1/2 px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-xs font-bold outline-none" />
                          <input name="salary_max" type="number" value={formData.salary_max} onChange={handleChange} placeholder="Max" className="w-1/2 px-4 py-3 bg-white/5 border border-white/10 rounded-xl text-xs font-bold outline-none" />
                       </div>
                    </div>

                    <div className="space-y-3">
                       <label htmlFor="expires_at" className="text-[10px] font-black text-slate-500 uppercase tracking-widest">Application Deadline</label>
                       <div className="relative">
                          <Calendar className="absolute left-3 top-3 w-4 h-4 text-slate-600" />
                          <input id="expires_at" name="expires_at" type="date" value={formData.expires_at} onChange={handleChange} className="w-full pl-10 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl text-xs font-bold outline-none" />
                       </div>
                    </div>
                 </div>
              </div>

              <button
                type="submit"
                disabled={loading}
                className="w-full flex items-center justify-center gap-3 py-6 bg-indigo-600 text-white font-black rounded-[2rem] hover:bg-indigo-700 transition-all shadow-2xl shadow-indigo-600/30 uppercase text-xs tracking-[0.2em] group"
              >
                {loading ? <Loader2 className="w-5 h-5 animate-spin" /> : <>Publish Opening <Send className="w-4 h-4 group-hover:translate-x-1 group-hover:-translate-y-1 transition-transform" /></>}
              </button>
           </div>
        </form>
      </div>
    </div>
  );
}
