import React, { useState, useEffect } from 'react';
import { Building2, Globe, Mail, Phone, MapPin, BadgeCheck, Settings, Users, Link as LinkIcon, Image as ImageIcon } from 'lucide-react';
import Link from 'next/link';
import { useRouter } from 'next/router';

export default function CompanyProfile() {
  const router = useRouter();
  const [profile, setProfile] = useState({
    name: '',
    industry: '',
    website: '',
    description: '',
    headquarters: '',
    company_size: '',
    logo_url: '',
    employer_badge: 'verified',
    social_links: { linkedin: '', twitter: '' },
    ats_integration_meta: { status: 'disconnected', provider: '' }
  });

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    // Mock fetch for now, will connect to /api/v1/companies/me
    setTimeout(() => {
      setProfile({
        name: 'TechCorp Solutions',
        industry: 'Software Development',
        website: 'https://techcorp.com',
        description: 'Building the future of enterprise software.',
        headquarters: 'Bangalore, India',
        company_size: '50-200',
        logo_url: '',
        employer_badge: 'gold',
        social_links: { linkedin: 'linkedin.com/company/techcorp', twitter: 'twitter.com/techcorp' },
        ats_integration_meta: { status: 'connected', provider: 'Greenhouse' }
      });
      setLoading(false);
    }, 800);
  }, []);

  const handleSave = (e) => {
    e.preventDefault();
    setSaving(true);
    setTimeout(() => {
      setSaving(false);
      alert('Profile updated!');
    }, 1000);
  };

  if (loading) return <div className="flex h-screen items-center justify-center">Loading Profile...</div>;

  return (
    <div className="flex min-h-screen bg-slate-50">
      <aside className="w-64 bg-white border-r border-slate-200 p-6">
        <h2 className="text-xl font-bold text-slate-800 mb-8">RojgarSetu <span className="text-xs bg-indigo-100 text-indigo-700 px-2 py-0.5 rounded">Recruiter</span></h2>
        <nav className="space-y-2">
          <Link href="/dashboard/company" className="flex items-center gap-3 px-4 py-3 text-slate-600 hover:bg-slate-50 rounded-xl font-medium">
            <Building2 className="w-5 h-5" /> Dashboard
          </Link>
          <Link href="/company/profile" className="flex items-center gap-3 px-4 py-3 bg-indigo-50 text-indigo-600 rounded-xl font-medium">
            <Settings className="w-5 h-5" /> Profile Settings
          </Link>
          <Link href="/company/applicants" className="flex items-center gap-3 px-4 py-3 text-slate-600 hover:bg-slate-50 rounded-xl font-medium">
            <Users className="w-5 h-5" /> Applicants
          </Link>
        </nav>
      </aside>

      <main className="flex-1 p-8">
        <div className="max-w-4xl mx-auto">
          <header className="flex justify-between items-start mb-8">
            <div>
              <div className="flex items-center gap-3 mb-1">
                <h1 className="text-3xl font-bold text-slate-900">{profile.name}</h1>
                {profile.employer_badge === 'gold' && (
                  <span className="flex items-center gap-1 bg-amber-100 text-amber-700 px-3 py-1 rounded-full text-xs font-bold border border-amber-200">
                    <BadgeCheck className="w-3.5 h-3.5" /> GOLD PARTNER
                  </span>
                )}
              </div>
              <p className="text-slate-500">Manage your company presence and recruiter branding.</p>
            </div>
            <button form="profile-form" disabled={saving} className="bg-indigo-600 hover:bg-indigo-700 text-white px-8 py-2.5 rounded-xl font-bold transition shadow-lg shadow-indigo-500/20 disabled:opacity-50">
              {saving ? 'Saving...' : 'Save Changes'}
            </button>
          </header>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-2 space-y-6">
              <form id="profile-form" onSubmit={handleSave} className="bg-white p-8 rounded-2xl border border-slate-200 shadow-sm space-y-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-1.5">
                    <label className="text-xs font-bold text-slate-700 uppercase">Industry</label>
                    <input type="text" value={profile.industry} onChange={e => setProfile({...profile, industry: e.target.value})} className="w-full px-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-indigo-500/20 outline-none" />
                  </div>
                  <div className="space-y-1.5">
                    <label className="text-xs font-bold text-slate-700 uppercase">Website</label>
                    <div className="relative">
                      <Globe className="absolute left-3 top-3 w-4 h-4 text-slate-400" />
                      <input type="url" value={profile.website} onChange={e => setProfile({...profile, website: e.target.value})} className="w-full pl-10 pr-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-indigo-500/20 outline-none" />
                    </div>
                  </div>
                </div>

                <div className="space-y-1.5">
                  <label className="text-xs font-bold text-slate-700 uppercase">About Company</label>
                  <textarea rows="4" value={profile.description} onChange={e => setProfile({...profile, description: e.target.value})} className="w-full px-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-indigo-500/20 outline-none resize-none" />
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-1.5">
                    <label className="text-xs font-bold text-slate-700 uppercase">Headquarters</label>
                    <div className="relative">
                      <MapPin className="absolute left-3 top-3 w-4 h-4 text-slate-400" />
                      <input type="text" value={profile.headquarters} onChange={e => setProfile({...profile, headquarters: e.target.value})} className="w-full pl-10 pr-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-indigo-500/20 outline-none" />
                    </div>
                  </div>
                  <div className="space-y-1.5">
                    <label className="text-xs font-bold text-slate-700 uppercase">Company Size</label>
                    <select value={profile.company_size} onChange={e => setProfile({...profile, company_size: e.target.value})} className="w-full px-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-indigo-500/20 outline-none appearance-none bg-white">
                      <option value="1-10">1-10 employees</option>
                      <option value="11-50">11-50 employees</option>
                      <option value="51-200">51-200 employees</option>
                      <option value="201-500">201-500 employees</option>
                      <option value="500+">500+ employees</option>
                    </select>
                  </div>
                </div>

                <div className="pt-6 border-t border-slate-100">
                  <h3 className="text-sm font-bold text-slate-800 mb-4 flex items-center gap-2">
                    <ImageIcon className="w-4 h-4" /> Company Gallery
                  </h3>
                  <div className="flex gap-4">
                    <div className="w-24 h-24 rounded-xl bg-slate-100 border-2 border-dashed border-slate-300 flex items-center justify-center text-slate-400 cursor-pointer hover:bg-slate-200 transition">
                      <ImageIcon className="w-6 h-6" />
                    </div>
                    <p className="text-xs text-slate-500 mt-2">Upload workspace photos to attract talent.</p>
                  </div>
                </div>
              </form>
            </div>

            <div className="space-y-6">
              <div className="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm">
                <h3 className="text-sm font-bold text-slate-800 mb-4 flex items-center gap-2">
                  <LinkIcon className="w-4 h-4" /> ATS Integration
                </h3>
                <div className={`p-4 rounded-xl border ${profile.ats_integration_meta.status === 'connected' ? 'bg-emerald-50 border-emerald-100' : 'bg-slate-50 border-slate-100'}`}>
                  <p className="text-xs font-bold text-slate-500 uppercase mb-1">Status</p>
                  <p className={`text-sm font-bold ${profile.ats_integration_meta.status === 'connected' ? 'text-emerald-700' : 'text-slate-700'}`}>
                    {profile.ats_integration_meta.status === 'connected' ? `Connected to ${profile.ats_integration_meta.provider}` : 'Disconnected'}
                  </p>
                  <button className="mt-4 w-full bg-white border border-slate-200 py-2 rounded-lg text-xs font-bold text-slate-700 hover:bg-slate-50">
                    Manage Integration
                  </button>
                </div>
              </div>

              <div className="bg-indigo-900 text-white p-6 rounded-2xl shadow-xl shadow-indigo-900/20">
                <h3 className="text-sm font-bold mb-2">Recruiter Pro</h3>
                <p className="text-xs text-indigo-200 mb-4">Unlock advanced analytics, AI candidate matching, and premium badges.</p>
                <button className="w-full bg-indigo-500 hover:bg-indigo-400 text-white py-2 rounded-lg text-xs font-bold transition">
                  Upgrade Plan
                </button>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
