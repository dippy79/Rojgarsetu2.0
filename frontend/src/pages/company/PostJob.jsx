import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

export default function PostJob() {
  const navigate = useNavigate();
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
    const token = localStorage.getItem('rojgar_token');
    const payload = {
      ...formData,
      skills_required: formData.skills_required.split(',').map(s => s.trim()).filter(Boolean),
      salary_min: Number(formData.salary_min),
      salary_max: Number(formData.salary_max),
      experience_years: Number(formData.experience_years)
    };

    const res = await fetch('http://localhost:3001/api/v1/company-jobs', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify(payload)
    });

    if (res.ok) navigate('/dashboard/company');
  };

  return (
    <div className="max-w-3xl mx-auto my-10 p-8 bg-white/80 backdrop-blur-xl border border-slate-100 rounded-2xl shadow-sm">
      <h1 className="text-2xl font-bold text-slate-900 mb-6">Create Job Opening</h1>
      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label htmlFor="title" className="block text-sm font-medium text-slate-700 mb-1">Job Title</label>
          <input id="title" name="title" value={formData.title} onChange={handleChange} required className="w-full border rounded-xl p-3 text-sm border-slate-200 outline-none focus:ring-2 focus:ring-indigo-500" />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label htmlFor="job_type" className="block text-sm font-medium text-slate-700 mb-1">Job Type</label>
            <select id="job_type" name="job_type" value={formData.job_type} onChange={handleChange} className="w-full border rounded-xl p-3 text-sm border-slate-200 outline-none">
              <option value="full-time">Full-time</option>
              <option value="part-time">Part-time</option>
              <option value="contract">Contract</option>
              <option value="internship">Internship</option>
            </select>
          </div>
          <div>
            <label htmlFor="location" className="block text-sm font-medium text-slate-700 mb-1">Location</label>
            <input id="location" name="location" value={formData.location} onChange={handleChange} required className="w-full border rounded-xl p-3 text-sm border-slate-200 outline-none" />
          </div>
        </div>

        <div className="grid grid-cols-3 gap-4">
          <div>
            <label htmlFor="salary_min" className="block text-sm font-medium text-slate-700 mb-1">Min Salary</label>
            <input id="salary_min" name="salary_min" type="number" value={formData.salary_min} onChange={handleChange} className="w-full border rounded-xl p-3 text-sm border-slate-200" />
          </div>
          <div>
            <label htmlFor="salary_max" className="block text-sm font-medium text-slate-700 mb-1">Max Salary</label>
            <input id="salary_max" name="salary_max" type="number" value={formData.salary_max} onChange={handleChange} className="w-full border rounded-xl p-3 text-sm border-slate-200" />
          </div>
          <div>
            <label htmlFor="currency_code" className="block text-sm font-medium text-slate-700 mb-1">Currency</label>
            <select id="currency_code" name="currency_code" value={formData.currency_code} onChange={handleChange} className="w-full border rounded-xl p-3 text-sm border-slate-200">
              <option value="INR">INR</option>
              <option value="USD">USD</option>
              <option value="GBP">GBP</option>
            </select>
          </div>
        </div>

        <div>
          <label htmlFor="experience_years" className="block text-sm font-medium text-slate-700 mb-1">Experience (Years)</label>
          <input id="experience_years" name="experience_years" type="number" value={formData.experience_years} onChange={handleChange} className="w-full border rounded-xl p-3 text-sm border-slate-200" />
        </div>

        <div>
          <label htmlFor="skills_required" className="block text-sm font-medium text-slate-700 mb-1">Skills (comma separated)</label>
          <input id="skills_required" name="skills_required" value={formData.skills_required} onChange={handleChange} placeholder="React, Node.js, PostgreSQL" className="w-full border rounded-xl p-3 text-sm border-slate-200" />
        </div>

        <div>
          <label htmlFor="description" className="block text-sm font-medium text-slate-700 mb-1">Description</label>
          <textarea id="description" name="description" rows={4} value={formData.description} onChange={handleChange} className="w-full border rounded-xl p-3 text-sm border-slate-200" />
        </div>

        <div>
          <label htmlFor="requirements" className="block text-sm font-medium text-slate-700 mb-1">Requirements</label>
          <textarea id="requirements" name="requirements" rows={3} value={formData.requirements} onChange={handleChange} className="w-full border rounded-xl p-3 text-sm border-slate-200" />
        </div>

        <div className="grid grid-cols-2 gap-4 items-center">
          <div className="flex items-center gap-2">
            <input id="is_remote" name="is_remote" type="checkbox" checked={formData.is_remote} onChange={handleChange} className="w-4 h-4 text-indigo-600 rounded" />
            <label htmlFor="is_remote" className="text-sm text-slate-700 font-medium">Remote Position</label>
          </div>
          <div>
            <label htmlFor="expires_at" className="block text-sm font-medium text-slate-700 mb-1">Expiration Date</label>
            <input id="expires_at" name="expires_at" type="date" value={formData.expires_at} onChange={handleChange} className="w-full border rounded-xl p-3 text-sm border-slate-200" />
          </div>
        </div>

        <button type="submit" className="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-3 rounded-xl transition">
          Publish Job Posting
        </button>
      </form>
    </div>
  );
}
