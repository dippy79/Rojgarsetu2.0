import React, { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { useAuth } from '../../hooks/useAuth';
import {
  User, Mail, Phone, MapPin, Briefcase, Award,
  FileText, Globe, Github, Linkedin, Plus, X,
  Save, Loader2, LayoutDashboard, Bookmark, Sparkles, LogOut, Menu
} from 'lucide-react';

export const CandidateProfile = () => {
  const { user, logout } = useAuth();
  const router = useRouter();

  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState(null);
  const [toast, setToast] = useState(null);

  const [profile, setProfile] = useState({
    full_name: '',
    email: '',
    phone: '',
    title: '',
    bio: '',
    location: '',
    targetSalary: '',
    targetRole: '',
    skills: [],
    experience: [],
    education: [],
    portfolio_links: { github: '', linkedin: '', website: '' },
    certifications: [],
    resume_url: '',
  });

  const [newSkill, setNewSkill] = useState({ name: '', level: 'Intermediate' });
  const [showExpForm, setShowExpForm] = useState(false);
  const [newExp, setNewExp] = useState({ company: '', role: '', years: '', description: '' });

  const API_BASE = process.env.REACT_APP_BACKEND_URL || 'http://localhost:3001';

  const showToast = useCallback((message, type = 'success') => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 3500);
  }, []);

  const handleAuthError = useCallback(() => {
    localStorage.removeItem('rojgar_token');
    showToast('Session expired. Please login again.', 'error');
    setTimeout(() => router.push('/login'), 1500);
  }, [router, showToast]);

  const fetchProfile = useCallback(async () => {
    setLoading(true);
    setError(null);
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');

    try {
      const res = await fetch(`${API_BASE}/api/v1/candidates/me`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });

      if (res.status === 401) {
        handleAuthError();
        return;
      }

      if (res.ok) {
        const data = await res.json();
        const cand = data?.candidate || data?.data || {};
        setProfile({
          full_name: cand?.full_name || user?.name || '',
          email: cand?.email || user?.email || '',
          phone: cand?.phone || '',
          title: cand?.title || '',
          bio: cand?.bio || '',
          location: cand?.location || '',
          targetSalary: cand?.targetSalary || '',
          targetRole: cand?.targetRole || '',
          skills: Array.isArray(cand?.skills) ? cand.skills : [],
          experience: Array.isArray(cand?.experience) ? cand.experience : [],
          education: Array.isArray(cand?.education) ? cand.education : [],
          portfolio_links: cand?.portfolio_links || { github: '', linkedin: '', website: '' },
          certifications: Array.isArray(cand?.certifications) ? cand.certifications : [],
          resume_url: cand?.resume_url || '',
        });
      }
    } catch (err) {
      setError('Failed to synchronize profile with the cloud node.');
    } finally {
      setLoading(false);
    }
  }, [API_BASE, handleAuthError, user]);

  useEffect(() => {
    fetchProfile();
  }, [fetchProfile]);

  const calculateCompletion = (p) => {
    let score = 0;
    if (p?.full_name) score += 20;
    if (p?.skills?.length > 0) score += 20;
    if (p?.experience?.length > 0) score += 20;
    if (p?.education?.length > 0) score += 20;
    if (p?.resume_url) score += 20;
    return score;
  };

  const completionPercentage = calculateCompletion(profile);

  const handleSaveProfile = async (e) => {
    e.preventDefault();
    setSaving(true);
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');

    try {
      const res = await fetch(`${API_BASE}/api/v1/candidates/me`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(profile)
      });

      if (res.ok) {
        showToast('Neural profile updated successfully!');
      } else {
        showToast('Failed to commit changes to database.', 'error');
      }
    } catch (err) {
      showToast('Network latency detected. Save failed.', 'error');
    } finally {
      setSaving(false);
    }
  };

  const handleAddSkill = (e) => {
    if (e.key === 'Enter') {
       e.preventDefault();
       if (newSkill.name.trim()) {
          setProfile(prev => ({ ...prev, skills: [...prev.skills, { ...newSkill, id: Date.now() }] }));
          setNewSkill({ name: '', level: 'Intermediate' });
       }
    }
  };

  return (
    <div className="flex min-h-screen bg-[#FBFBFB] font-sans">
      {/* Sidebar - Consistent Premium Style */}
      <aside className={`
        fixed inset-y-0 left-0 z-50 w-72 bg-slate-950 text-white transform transition-transform duration-500
        md:translate-x-0 md:static md:h-screen sticky top-0
        ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
        flex flex-col p-8
      `}>
        <div className="flex items-center justify-between mb-12">
          <div className="flex items-center gap-3">
            <div className="p-2.5 bg-blue-600 rounded-2xl">
              <Briefcase className="w-6 h-6 text-white" />
            </div>
            <span className="font-black text-2xl tracking-tighter uppercase">Rojgar<span className="text-blue-500">Setu</span></span>
          </div>
        </div>

        <nav className="space-y-2 flex-1">
          {[
            { label: 'Dashboard', path: '/dashboard/candidate', icon: LayoutDashboard },
            { label: 'My Profile', path: '/candidate/profile', icon: User },
            { label: 'Applications', path: '/candidate/applications', icon: FileText },
            { label: 'Saved Jobs', path: '/candidate/saved-jobs', icon: Bookmark },
            { label: 'AI Matches', path: '/candidate/ai-matches', icon: Sparkles },
          ].map((item) => (
            <Link key={item.path} href={item.path} className={`flex items-center gap-4 px-6 py-4 rounded-2xl text-sm font-bold transition-all ${router.pathname === item.path ? 'bg-blue-600 text-white' : 'text-slate-400 hover:text-white hover:bg-white/5'}`}>
              <item.icon className="w-5 h-5" /> {item.label}
            </Link>
          ))}
        </nav>

        <button onClick={() => { logout(); router.push('/login'); }} className="flex items-center gap-4 w-full px-6 py-4 text-sm font-bold text-rose-400 hover:bg-rose-500/10 rounded-2xl transition-all mt-auto border border-rose-500/10">
          <LogOut className="w-5 h-5" /> Logout
        </button>
      </aside>

      <main className="flex-1 overflow-auto p-6 md:p-12 lg:p-16">
        <header className="flex flex-col md:flex-row md:items-center justify-between gap-8 mb-16">
           <div className="space-y-2">
              <h1 className="text-4xl font-black text-slate-900 tracking-tight">Identity Settings.</h1>
              <p className="text-slate-400 text-lg font-medium">Refining your professional metadata for <span className="text-blue-600 font-bold">AI Ingestion.</span></p>
           </div>

           <div className="flex items-center gap-8">
              <div className="text-right">
                 <p className="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Profile Strength</p>
                 <div className="flex items-center gap-3">
                    <div className="w-48 h-2.5 bg-slate-100 rounded-full overflow-hidden border border-slate-200">
                       <div className="h-full bg-blue-600 transition-all duration-1000 shadow-[0_0_10px_rgba(37,99,235,0.4)]" style={{width: `${completionPercentage}%`}}></div>
                    </div>
                    <span className="text-sm font-black text-slate-900">{completionPercentage}%</span>
                 </div>
              </div>
           </div>
        </header>

        {loading ? (
          <div className="space-y-12 animate-pulse">
             <div className="h-96 bg-white border border-slate-200 rounded-[3rem]"></div>
             <div className="h-64 bg-white border border-slate-200 rounded-[3rem]"></div>
          </div>
        ) : (
          <form onSubmit={handleSaveProfile} className="grid grid-cols-1 lg:grid-cols-12 gap-12">
             <div className="lg:col-span-8 space-y-8">
                {/* Basic Intel Section */}
                <div className="bg-white border border-slate-200 rounded-[3.5rem] p-12 shadow-sm space-y-10">
                   <h3 className="text-xl font-black text-slate-900 tracking-tight flex items-center gap-3">
                      <div className="w-2 h-2 rounded-full bg-blue-600"></div> Core Meta
                   </h3>

                   <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                      <div className="space-y-3">
                         <label htmlFor="full_name" className="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Legal Name</label>
                         <input id="full_name" name="full_name" type="text" value={profile.full_name} onChange={e => setProfile({...profile, full_name: e.target.value})} className="w-full px-6 py-4 bg-slate-50 border-none rounded-2xl text-sm font-bold focus:ring-4 focus:ring-blue-500/10 transition-all outline-none" placeholder="Simranjeet Singh" autocomplete="name" />
                      </div>
                      <div className="space-y-3">
                         <label htmlFor="email" className="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Sync Address</label>
                         <input id="email" name="email" type="email" value={profile.email} disabled className="w-full px-6 py-4 bg-slate-100 border-none rounded-2xl text-sm font-bold text-slate-400 cursor-not-allowed outline-none" autocomplete="email" />
                      </div>
                   </div>

                   <div className="space-y-3">
                      <label htmlFor="title" className="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Professional Header</label>
                      <input id="title" name="title" type="text" value={profile.title} onChange={e => setProfile({...profile, title: e.target.value})} className="w-full px-6 py-4 bg-slate-50 border-none rounded-2xl text-sm font-bold focus:ring-4 focus:ring-blue-500/10 transition-all outline-none" placeholder="Senior DevOps Architect" autocomplete="organization-title" />
                   </div>

                   <div className="space-y-3">
                      <label htmlFor="bio" className="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Executive Summary</label>
                      <textarea id="bio" name="bio" rows={5} value={profile.bio} onChange={e => setProfile({...profile, bio: e.target.value})} className="w-full px-6 py-4 bg-slate-50 border-none rounded-[2rem] text-sm font-bold focus:ring-4 focus:ring-blue-500/10 transition-all outline-none resize-none" placeholder="Deep-dive into your trajectory..." autocomplete="off" />
                   </div>
                </div>

                {/* Skill Graph Section */}
                <div className="bg-white border border-slate-200 rounded-[3.5rem] p-12 shadow-sm space-y-10">
                   <h3 className="text-xl font-black text-slate-900 tracking-tight flex items-center gap-3">
                      <div className="w-2 h-2 rounded-full bg-emerald-500"></div> Neural Skills
                   </h3>

                   <div className="space-y-6">
                      <div className="flex flex-wrap gap-3">
                         {profile.skills.map(s => (
                           <div key={s.id} className="flex items-center gap-3 px-5 py-2.5 bg-slate-900 text-white rounded-2xl text-[10px] font-black uppercase tracking-widest border border-slate-800 shadow-xl shadow-slate-900/10 transition-all hover:scale-105">
                              {s.name}
                              <button onClick={() => setProfile({...profile, skills: profile.skills.filter(sk => sk.id !== s.id)})} className="hover:text-rose-400 transition-colors"><X className="w-3 h-3" /></button>
                           </div>
                         ))}
                      </div>

                      <div className="relative group">
                         <Plus className="absolute left-5 top-4.5 w-5 h-5 text-slate-400" />
                         <input
                           type="text"
                           placeholder="Type skill and press ENTER..."
                           value={newSkill.name}
                           onChange={e => setNewSkill({...newSkill, name: e.target.value})}
                           onKeyDown={handleAddSkill}
                           className="w-full pl-14 pr-6 py-5 bg-slate-50 border-none rounded-3xl text-sm font-bold focus:ring-4 focus:ring-emerald-500/10 transition-all outline-none"
                         />
                      </div>
                   </div>
                </div>
             </div>

             <div className="lg:col-span-4 space-y-8">
                {/* Logistics Bento */}
                <div className="bg-slate-900 rounded-[3rem] p-10 text-white space-y-8 shadow-2xl shadow-blue-900/20 relative overflow-hidden">
                   <div className="absolute top-0 right-0 w-32 h-32 bg-blue-600/10 rounded-full -mr-16 -mt-16 blur-3xl"></div>

                   <div className="space-y-6 relative z-10">
                      <div className="space-y-3">
                         <label htmlFor="location" className="text-[10px] font-black text-slate-500 uppercase tracking-widest ml-1">Current Hub</label>
                         <div className="relative group">
                            <MapPin className="absolute left-4 top-4 w-4 h-4 text-slate-600 group-focus-within:text-blue-400 transition-colors" />
                            <input id="location" name="location" type="text" value={profile.location} onChange={e => setProfile({...profile, location: e.target.value})} className="w-full pl-12 pr-6 py-4 bg-white/5 border border-white/10 rounded-2xl text-xs font-bold outline-none focus:border-blue-500 transition-all" placeholder="Dubai, UAE" autocomplete="address-level2" />
                         </div>
                      </div>

                      <div className="space-y-3">
                         <label htmlFor="phone" className="text-[10px] font-black text-slate-500 uppercase tracking-widest ml-1">Direct Line</label>
                         <div className="relative">
                            <Phone className="absolute left-4 top-4 w-4 h-4 text-slate-600" />
                            <input id="phone" name="phone" type="tel" value={profile.phone} onChange={e => setProfile({...profile, phone: e.target.value})} className="w-full pl-12 pr-6 py-4 bg-white/5 border border-white/10 rounded-2xl text-xs font-bold outline-none focus:border-blue-500 transition-all" placeholder="+971 50 123 4567" autocomplete="tel" />
                         </div>
                      </div>
                   </div>

                   <button
                     type="submit"
                     disabled={saving}
                     className="w-full flex items-center justify-center gap-3 py-6 bg-blue-600 text-white font-black rounded-[2rem] hover:bg-blue-700 transition-all shadow-2xl shadow-blue-600/30 uppercase text-[10px] tracking-[0.3em]"
                   >
                     {saving ? <Loader2 className="w-4 h-4 animate-spin" /> : <><Save className="w-4 h-4" /> Commit Profile</>}
                   </button>
                </div>

                {/* Social Nodes */}
                <div className="bg-white border border-slate-200 rounded-[3rem] p-10 space-y-6">
                   <h4 className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Digital Footprint</h4>
                   <div className="space-y-4">
                      <div className="flex items-center gap-4 p-4 bg-slate-50 rounded-2xl border border-slate-100 group transition-all hover:border-slate-900">
                         <Github className="w-5 h-5 text-slate-400 group-hover:text-slate-900 transition-colors" />
                         <input type="url" value={profile.portfolio_links.github} onChange={e => setProfile({...profile, portfolio_links: {...profile.portfolio_links, github: e.target.value}})} className="bg-transparent border-none text-xs font-bold outline-none w-full" placeholder="github.com/username" autocomplete="off" />
                      </div>
                      <div className="flex items-center gap-4 p-4 bg-slate-50 rounded-2xl border border-slate-100 group transition-all hover:border-blue-600">
                         <Linkedin className="w-5 h-5 text-slate-400 group-hover:text-blue-600 transition-colors" />
                         <input type="url" value={profile.portfolio_links.linkedin} onChange={e => setProfile({...profile, portfolio_links: {...profile.portfolio_links, linkedin: e.target.value}})} className="bg-transparent border-none text-xs font-bold outline-none w-full" placeholder="linkedin.com/in/user" autocomplete="off" />
                      </div>
                   </div>
                </div>
             </div>
          </form>
        )}
      </main>
    </div>
  );
};

export default CandidateProfile;
