import React, { useState, useEffect, useCallback } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';

// Icon Components
const LayoutDashboardIcon = () => (
  <svg className="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <rect x="3" y="3" width="7" height="9" rx="1" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    <rect x="14" y="3" width="7" height="5" rx="1" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    <rect x="14" y="12" width="7" height="9" rx="1" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    <rect x="3" y="16" width="7" height="5" rx="1" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
  </svg>
);

const UserIcon = () => (
  <svg className="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
  </svg>
);

const FileTextIcon = () => (
  <svg className="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
  </svg>
);

const BookmarkIcon = () => (
  <svg className="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
  </svg>
);

const SparklesIcon = () => (
  <svg className="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
  </svg>
);

const LogOutIcon = () => (
  <svg className="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
  </svg>
);

const MenuIcon = () => (
  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h16m-7 6h7" />
  </svg>
);

const XIcon = () => (
  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
  </svg>
);

export const CandidateProfile = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

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
    profile_completion: 0
  });

  const [newSkill, setNewSkill] = useState({ name: '', level: 'Intermediate' });
  const [showCertForm, setShowCertForm] = useState(false);
  const [newCert, setNewCert] = useState({ name: '', issuer: '', year: '' });
  const [showExpForm, setShowExpForm] = useState(false);
  const [showEduForm, setShowEduForm] = useState(false);
  const [newExp, setNewExp] = useState({ company: '', role: '', years: '', description: '' });
  const [newEdu, setNewEdu] = useState({ degree: '', institution: '', year: '' });
  const API_BASE = process.env.REACT_APP_BACKEND_URL || 'http://localhost:3001';

  const showToast = useCallback((message, type = 'success') => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 3500);
  }, []);

  const handleAuthError = useCallback(() => {
    localStorage.removeItem('rojgar_token');
    setToast({ message: 'Session expired. Please login again.', type: 'error' });
    setTimeout(() => navigate('/login'), 1500);
  }, [navigate]);

  const fetchProfile = useCallback(async () => {
    setLoading(true);
    setError(null);
    const token = localStorage.getItem('rojgar_token');

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
        const candidateData = data?.data?.data || data?.data || {};
        const cand = candidateData?.candidate || candidateData;
        setProfile({
          full_name: cand?.full_name || user?.name || 'Test Candidate',
          email: cand?.email || user?.email || 'test@rojgarsetu.com',
          phone: cand?.phone || '',
          title: cand?.title || 'Full Stack Developer',
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
          profile_completion: cand?.profile_completion || 0
        });
      } else if (res.status === 404) {
        // Gracefully handle missing profile - use defaults
        setProfile({
          full_name: user?.name || 'Test Candidate',
          email: user?.email || 'test@rojgarsetu.com',
          phone: '',
          title: 'Full Stack Developer',
          bio: '',
          location: '',
          targetSalary: '',
          targetRole: '',
          skills: [],
          experience: [],
          education: [],
          resume_url: '',
          profile_completion: 0
        });
      } else {
        setError('Failed to load profile data.');
      }
    } catch (err) {
      // On network error, use defaults instead of showing error
      setProfile({
        full_name: user?.name || 'Test Candidate',
        email: user?.email || 'test@rojgarsetu.com',
        phone: '',
        title: 'Full Stack Developer',
        bio: '',
        location: '',
        targetSalary: '',
        targetRole: '',
        skills: [],
        experience: [],
        education: [],
        resume_url: '',
        profile_completion: 0
      });
    } finally {
      setLoading(false);
    }
  }, [API_BASE, handleAuthError, user?.email, user?.name]);

  useEffect(() => {
    fetchProfile();
  }, [fetchProfile]);

  const calculateCompletion = (p) => {
    let score = 0;
    if (p?.full_name) score += 10;
    if (p?.skills?.length > 0) score += 20;
    if (p?.experience?.length > 0) score += 20;
    if (p?.education?.length > 0) score += 20;
    if (p?.resume_url) score += 10;
    if (p?.portfolio_links?.github || p?.portfolio_links?.linkedin) score += 10;
    if (p?.certifications?.length > 0) score += 10;
    return score;
  };

  const completionPercentage = calculateCompletion(profile);

  const handleSaveProfile = async (e) => {
    e.preventDefault();
    setSaving(true);
    const token = localStorage.getItem('rojgar_token');

    try {
      const res = await fetch(`${API_BASE}/api/v1/candidates/me`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(profile)
      });

      if (res.status === 401) {
        handleAuthError();
        return;
      }

      if (res.ok) {
        showToast('Profile updated successfully!');
      } else {
        showToast('Failed to update profile.', 'error');
      }
    } catch (err) {
      showToast('Network error during save.', 'error');
    } finally {
      setSaving(false);
    }
  };

  const handleAddSkill = (e) => {
    e.preventDefault();
    if (newSkill.name.trim() && !profile.skills.find(s => s.name === newSkill.name.trim())) {
      setProfile(prev => ({ ...prev, skills: [...prev.skills, { ...newSkill, id: Date.now() }] }));
      setNewSkill({ name: '', level: 'Intermediate' });
    }
  };

  const handleRemoveSkill = (id) => {
    setProfile(prev => ({
      ...prev,
      skills: prev.skills.filter(s => s.id !== id)
    }));
  };

  const handleAddExperience = (e) => {
    e.preventDefault();
    if (newExp.company && newExp.role) {
      setProfile(prev => ({
        ...prev,
        experience: [...(prev.experience || []), { ...newExp, id: Date.now() }]
      }));
      setNewExp({ company: '', role: '', years: '', description: '' });
      setShowExpForm(false);
    }
  };

  const handleRemoveExperience = (id) => {
    setProfile(prev => ({
      ...prev,
      experience: prev.experience.filter(exp => exp.id !== id)
    }));
  };

  const handleAddEducation = (e) => {
    e.preventDefault();
    if (newEdu.degree && newEdu.institution) {
      setProfile(prev => ({
        ...prev,
        education: [...(prev.education || []), { ...newEdu, id: Date.now() }]
      }));
      setNewEdu({ degree: '', institution: '', year: '' });
      setShowEduForm(false);
    }
  };

  const handleRemoveEducation = (id) => {
    setProfile(prev => ({
      ...prev,
      education: prev.education.filter(edu => edu.id !== id)
    }));
  };

  const handleAddCertification = (e) => {
    e.preventDefault();
    if (newCert.name && newCert.issuer) {
      setProfile(prev => ({
        ...prev,
        certifications: [...(prev.certifications || []), { ...newCert, id: Date.now() }]
      }));
      setNewCert({ name: '', issuer: '', year: '' });
      setShowCertForm(false);
    }
  };

  const handleRemoveCertification = (id) => {
    setProfile(prev => ({
      ...prev,
      certifications: prev.certifications.filter(c => c.id !== id)
    }));
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const navItems = [
    { label: 'Dashboard', path: '/dashboard/candidate', icon: LayoutDashboardIcon },
    { label: 'My Profile', path: '/candidate/profile', icon: UserIcon },
    { label: 'Applications', path: '/candidate/applications', icon: FileTextIcon },
    { label: 'Saved Jobs', path: '/candidate/saved-jobs', icon: BookmarkIcon },
    { label: 'AI Matches', path: '/candidate/ai-matches', icon: SparklesIcon },
  ];

  return (
    <div className="flex min-h-screen bg-slate-50 font-sans">
      {/* Toast Notification */}
      {toast && (
        <div className={`fixed top-6 right-6 z-[100] px-5 py-3 rounded-xl shadow-2xl border text-sm font-bold animate-in slide-in-from-right-10 ${
          toast.type === 'success' ? 'bg-emerald-50 text-emerald-800 border-emerald-200' : 'bg-red-50 text-red-800 border-red-200'
        }`}>
          {toast.message}
        </div>
      )}

      {/* Mobile Sidebar Overlay */}
      {sidebarOpen && (
        <div className="fixed inset-0 bg-slate-900/50 z-40 md:hidden backdrop-blur-sm" onClick={() => setSidebarOpen(false)}></div>
      )}

      {/* Sidebar */}
      <aside className={`
        fixed inset-y-0 left-0 z-50 w-64 bg-slate-900 text-white transform transition-transform duration-300 ease-in-out
        md:translate-x-0 md:static md:h-screen sticky top-0
        ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
        flex flex-col justify-between p-6 shadow-xl
      `}>
        <div>
          <div className="flex items-center justify-between mb-8 md:hidden">
            <span className="font-bold text-xl">RojgarSetu</span>
            <button onClick={() => setSidebarOpen(false)} className="p-2 text-slate-400 hover:text-white">
              <XIcon />
            </button>
          </div>

          <div className="flex items-center gap-3 mb-6">
            <div className="bg-blue-600 w-12 h-12 rounded-full flex items-center justify-center font-bold text-white text-lg shrink-0">
              {(profile.full_name || 'C').charAt(0)}
            </div>
            <div className="overflow-hidden">
              <h2 className="font-semibold text-white truncate text-base">{profile.full_name || 'Candidate'}</h2>
              <span className="inline-block bg-blue-500/20 text-blue-300 text-xs px-2 py-0.5 rounded-md">Candidate</span>
            </div>
          </div>

          <nav className="space-y-1">
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = location.pathname === item.path;
              return (
                <Link
                  key={item.path}
                  to={item.path}
                  onClick={() => setSidebarOpen(false)}
                  className={`flex items-center px-4 py-3 text-sm font-medium transition-all ${
                    isActive ? 'bg-white/10 text-white rounded-lg' : 'text-slate-400 hover:text-white hover:bg-white/5 rounded-lg'
                  }`}
                >
                  <Icon />
                  {item.label}
                </Link>
              );
            })}
          </nav>
        </div>

        <button onClick={handleLogout} className="flex items-center w-full px-4 py-3 text-sm font-medium bg-red-500/10 text-red-400 hover:bg-red-500/20 rounded-lg transition-all">
          <LogOutIcon />
          Logout
        </button>
      </aside>

      {/* Main Content */}
      <main className="flex-1 overflow-auto p-4 md:p-8">
        <div className="max-w-4xl mx-auto">
          {/* Header */}
          <div className="flex items-center gap-4 mb-8">
            <button onClick={() => setSidebarOpen(true)} className="p-2 bg-white border border-slate-200 rounded-lg shadow-sm text-slate-600 md:hidden">
              <MenuIcon />
            </button>
            <div>
              <h1 className="text-3xl font-bold text-slate-900">Profile Settings</h1>
              <p className="text-slate-500 text-sm">Update your professional identity and job preferences.</p>
            </div>
          </div>

          {/* Error Banner */}
          {error && (
            <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-2xl flex items-center justify-between text-red-700">
              <span className="text-sm font-medium">{error}</span>
              <button onClick={fetchProfile} className="px-4 py-1.5 bg-red-600 text-white text-xs font-bold rounded-lg hover:bg-red-700">Retry</button>
            </div>
          )}

          {loading ? (
            <div className="space-y-6 animate-pulse">
              <div className="h-64 bg-slate-200 rounded-2xl"></div>
              <div className="h-48 bg-slate-200 rounded-2xl"></div>
            </div>
          ) : (
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
              {/* Profile Completion Card */}
              <div className="lg:col-span-1 space-y-6">
                <div className="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm text-center">
                  <h3 className="text-sm font-bold text-slate-500 uppercase tracking-wider mb-4">Profile Strength</h3>
                  <div className="relative inline-flex items-center justify-center mb-4">
                    <svg className="w-32 h-32 transform -rotate-90">
                      <circle cx="64" cy="64" r="58" stroke="currentColor" strokeWidth="8" fill="transparent" className="text-slate-100" />
                      <circle cx="64" cy="64" r="58" stroke="currentColor" strokeWidth="8" fill="transparent"
                        strokeDasharray={364.4} strokeDashoffset={364.4 - (364.4 * completionPercentage) / 100}
                        className="text-blue-600 transition-all duration-1000 ease-out" strokeLinecap="round" />
                    </svg>
                    <span className="absolute text-2xl font-black text-slate-900">{completionPercentage}%</span>
                  </div>
                  <p className="text-xs text-slate-500 px-4">Complete your profile to increase your visibility to recruiters.</p>
                </div>
              </div>

              {/* Edit Form */}
              <div className="lg:col-span-2">
                <form onSubmit={handleSaveProfile} className="bg-white p-6 md:p-8 rounded-2xl border border-slate-200 shadow-sm space-y-6">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-1.5">
                      <label htmlFor="full_name" className="text-xs font-bold text-slate-700 uppercase">Full Name</label>
                      <input id="full_name" name="full_name" type="text" value={profile.full_name} onChange={(e) => setProfile({...profile, full_name: e.target.value})} className="w-full px-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-blue-500/20 outline-none" placeholder="Enter your full name" autoComplete="name" />
                    </div>
                    <div className="space-y-1.5">
                      <label htmlFor="email" className="text-xs font-bold text-slate-700 uppercase">Email Address</label>
                      <input id="email" name="email" type="email" value={profile.email} disabled className="w-full px-4 py-2.5 rounded-xl border border-slate-200 bg-slate-50 text-slate-500 cursor-not-allowed" autoComplete="email" />
                    </div>
                  </div>

                  <div className="space-y-1.5">
                    <label htmlFor="title" className="text-xs font-bold text-slate-700 uppercase">Professional Title</label>
                    <input id="title" name="title" type="text" value={profile.title} onChange={(e) => setProfile({...profile, title: e.target.value})} className="w-full px-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-blue-500/20 outline-none" placeholder="e.g. Senior Software Engineer" autoComplete="organization-title" />
                  </div>

                  <div className="space-y-1.5">
                    <label htmlFor="bio" className="text-xs font-bold text-slate-700 uppercase">Professional Summary</label>
                    <textarea id="bio" name="bio" rows="4" value={profile.bio} onChange={(e) => setProfile({...profile, bio: e.target.value})} className="w-full px-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-blue-500/20 outline-none resize-none" placeholder="Tell us about your background and experience..." autoComplete="off" />
                  </div>

                  <div className="space-y-1.5">
                    <label htmlFor="location" className="text-xs font-bold text-slate-700 uppercase">Location</label>
                    <input id="location" name="location" type="text" value={profile.location} onChange={(e) => setProfile({...profile, location: e.target.value})} className="w-full px-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-blue-500/20 outline-none" placeholder="City, State" autoComplete="address-level2" />
                  </div>

                  <div className="space-y-1.5">
                    <label htmlFor="phone" className="text-xs font-bold text-slate-700 uppercase">Phone Number</label>
                    <input id="phone" name="phone" type="tel" value={profile.phone || ''} onChange={(e) => setProfile({...profile, phone: e.target.value})} className="w-full px-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-blue-500/20 outline-none" placeholder="+91 98765 43210" autoComplete="tel" />
                  </div>

                  {/* Portfolio Links */}
                  <div className="space-y-3 p-4 bg-blue-50/50 rounded-2xl border border-blue-100">
                    <label className="text-xs font-bold text-blue-800 uppercase">Professional Links</label>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                      <div className="space-y-1">
                        <label htmlFor="github" className="text-[10px] font-bold text-slate-500 uppercase">GitHub</label>
                        <input id="github" type="url" value={profile.portfolio_links.github} onChange={(e) => setProfile({...profile, portfolio_links: {...profile.portfolio_links, github: e.target.value}})} className="w-full px-3 py-2 rounded-lg border border-slate-200 text-sm outline-none" placeholder="github.com/username" />
                      </div>
                      <div className="space-y-1">
                        <label htmlFor="linkedin" className="text-[10px] font-bold text-slate-500 uppercase">LinkedIn</label>
                        <input id="linkedin" type="url" value={profile.portfolio_links.linkedin} onChange={(e) => setProfile({...profile, portfolio_links: {...profile.portfolio_links, linkedin: e.target.value}})} className="w-full px-3 py-2 rounded-lg border border-slate-200 text-sm outline-none" placeholder="linkedin.com/in/user" />
                      </div>
                      <div className="space-y-1">
                        <label htmlFor="portfolio" className="text-[10px] font-bold text-slate-500 uppercase">Website</label>
                        <input id="portfolio" type="url" value={profile.portfolio_links.website} onChange={(e) => setProfile({...profile, portfolio_links: {...profile.portfolio_links, website: e.target.value}})} className="w-full px-3 py-2 rounded-lg border border-slate-200 text-sm outline-none" placeholder="portfolio.com" />
                      </div>
                    </div>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-1.5">
                      <label htmlFor="targetRole" className="text-xs font-bold text-slate-700 uppercase">Target Role</label>
                      <input id="targetRole" name="targetRole" type="text" value={profile.targetRole} onChange={(e) => setProfile({...profile, targetRole: e.target.value})} className="w-full px-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-blue-500/20 outline-none" placeholder="e.g. Senior Developer" autoComplete="off" />
                    </div>
                    <div className="space-y-1.5">
                      <label htmlFor="targetSalary" className="text-xs font-bold text-slate-700 uppercase">Target Salary</label>
                      <input id="targetSalary" name="targetSalary" type="text" value={profile.targetSalary} onChange={(e) => setProfile({...profile, targetSalary: e.target.value})} className="w-full px-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-blue-500/20 outline-none" placeholder="e.g. ₹15-20 LPA" autoComplete="off" />
                    </div>
                  </div>

                  {/* Skills Tag Input */}
                  <div className="space-y-3">
                    <label htmlFor="skill-input" className="text-xs font-bold text-slate-700 uppercase">Skill Graph (Expertise Level)</label>
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                      {(profile.skills || []).map((skill) => (
                        <div key={skill.id} className="p-3 bg-white rounded-xl border border-slate-200 flex justify-between items-center group hover:border-blue-400 transition-colors shadow-sm">
                          <div>
                            <p className="text-sm font-bold text-slate-900">{skill.name}</p>
                            <p className="text-[10px] text-blue-600 font-bold uppercase">{skill.level}</p>
                          </div>
                          <button type="button" onClick={() => handleRemoveSkill(skill.id)} className="text-slate-300 hover:text-red-500 opacity-0 group-hover:opacity-100 transition-opacity">×</button>
                        </div>
                      ))}
                    </div>
                    <div className="flex gap-2 p-2 bg-slate-50 rounded-xl border border-dashed border-slate-300">
                      <input id="skill-input" type="text" value={newSkill.name} onChange={(e) => setNewSkill({...newSkill, name: e.target.value})} placeholder="Skill name..." className="flex-1 px-3 py-1.5 rounded-lg border border-slate-200 text-sm outline-none" />
                      <select value={newSkill.level} onChange={(e) => setNewSkill({...newSkill, level: e.target.value})} className="px-2 py-1.5 rounded-lg border border-slate-200 text-sm outline-none">
                        <option>Beginner</option>
                        <option>Intermediate</option>
                        <option>Advanced</option>
                        <option>Expert</option>
                      </select>
                      <button type="button" onClick={handleAddSkill} className="px-4 py-1.5 bg-slate-900 text-white text-xs font-bold rounded-lg hover:bg-slate-800">Add</button>
                    </div>
                  </div>

                  {/* Experience Section */}
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <label className="text-xs font-bold text-slate-700 uppercase">Work Experience</label>
                      <button type="button" onClick={() => setShowExpForm(!showExpForm)} className="text-xs font-bold text-blue-600 hover:text-blue-700">
                        {showExpForm ? 'Cancel' : '+ Add Experience'}
                      </button>
                    </div>
                    
                    {showExpForm && (
                      <form onSubmit={handleAddExperience} className="p-4 bg-slate-50 rounded-xl border border-slate-200 space-y-3">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                          <div>
                            <label htmlFor="exp_company" className="text-xs font-bold text-slate-600">Company</label>
                            <input id="exp_company" name="exp_company" type="text" value={newExp.company} onChange={(e) => setNewExp({...newExp, company: e.target.value})} className="w-full px-3 py-2 rounded-lg border border-slate-200 text-sm focus:ring-2 focus:ring-blue-500/20 outline-none" placeholder="Company name" required autoComplete="organization" />
                          </div>
                          <div>
                            <label htmlFor="exp_role" className="text-xs font-bold text-slate-600">Role</label>
                            <input id="exp_role" name="exp_role" type="text" value={newExp.role} onChange={(e) => setNewExp({...newExp, role: e.target.value})} className="w-full px-3 py-2 rounded-lg border border-slate-200 text-sm focus:ring-2 focus:ring-blue-500/20 outline-none" placeholder="Job title" required autoComplete="organization-title" />
                          </div>
                          <div>
                            <label htmlFor="exp_years" className="text-xs font-bold text-slate-600">Years</label>
                            <input id="exp_years" name="exp_years" type="text" value={newExp.years} onChange={(e) => setNewExp({...newExp, years: e.target.value})} className="w-full px-3 py-2 rounded-lg border border-slate-200 text-sm focus:ring-2 focus:ring-blue-500/20 outline-none" placeholder="e.g. 2 years" autoComplete="off" />
                          </div>
                        </div>
                        <div>
                          <label htmlFor="exp_desc" className="text-xs font-bold text-slate-600">Description</label>
                          <textarea id="exp_desc" name="exp_desc" rows="2" value={newExp.description} onChange={(e) => setNewExp({...newExp, description: e.target.value})} className="w-full px-3 py-2 rounded-lg border border-slate-200 text-sm focus:ring-2 focus:ring-blue-500/20 outline-none resize-none" placeholder="Brief description of responsibilities" autoComplete="off" />
                        </div>
                        <div className="flex gap-2">
                          <button type="submit" className="px-4 py-2 bg-blue-600 text-white text-xs font-bold rounded-lg hover:bg-blue-700">Add Experience</button>
                          <button type="button" onClick={() => setShowExpForm(false)} className="px-4 py-2 bg-slate-200 text-slate-700 text-xs font-bold rounded-lg hover:bg-slate-300">Cancel</button>
                        </div>
                      </form>
                    )}

                    <div className="space-y-2">
                      {(profile.experience || []).map((exp) => (
                        <div key={exp.id} className="p-3 bg-slate-50 rounded-lg border border-slate-200 flex justify-between items-start">
                          <div>
                            <p className="text-sm font-bold text-slate-900">{exp.role} at {exp.company}</p>
                            <p className="text-xs text-slate-500">{exp.years} {exp.description && `• ${exp.description}`}</p>
                          </div>
                          <button type="button" onClick={() => handleRemoveExperience(exp.id)} className="text-slate-400 hover:text-red-500 text-sm">×</button>
                        </div>
                      ))}
                    </div>
                  </div>

                  {/* Education Section */}
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <label className="text-xs font-bold text-slate-700 uppercase">Education</label>
                      <button type="button" onClick={() => setShowEduForm(!showEduForm)} className="text-xs font-bold text-blue-600 hover:text-blue-700">
                        {showEduForm ? 'Cancel' : '+ Add Education'}
                      </button>
                    </div>
                    
                    {showEduForm && (
                      <form onSubmit={handleAddEducation} className="p-4 bg-slate-50 rounded-xl border border-slate-200 space-y-3">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                          <div>
                            <label htmlFor="edu_degree" className="text-xs font-bold text-slate-600">Degree</label>
                            <input id="edu_degree" name="edu_degree" type="text" value={newEdu.degree} onChange={(e) => setNewEdu({...newEdu, degree: e.target.value})} className="w-full px-3 py-2 rounded-lg border border-slate-200 text-sm focus:ring-2 focus:ring-blue-500/20 outline-none" placeholder="e.g. B.Tech Computer Science" required autoComplete="off" />
                          </div>
                          <div>
                            <label htmlFor="edu_inst" className="text-xs font-bold text-slate-600">Institution</label>
                            <input id="edu_inst" name="edu_inst" type="text" value={newEdu.institution} onChange={(e) => setNewEdu({...newEdu, institution: e.target.value})} className="w-full px-3 py-2 rounded-lg border border-slate-200 text-sm focus:ring-2 focus:ring-blue-500/20 outline-none" placeholder="University/College name" required autoComplete="organization" />
                          </div>
                          <div>
                            <label htmlFor="edu_year" className="text-xs font-bold text-slate-600">Year</label>
                            <input id="edu_year" name="edu_year" type="text" value={newEdu.year} onChange={(e) => setNewEdu({...newEdu, year: e.target.value})} className="w-full px-3 py-2 rounded-lg border border-slate-200 text-sm focus:ring-2 focus:ring-blue-500/20 outline-none" placeholder="e.g. 2020" autoComplete="off" />
                          </div>
                        </div>
                        <div className="flex gap-2">
                          <button type="submit" className="px-4 py-2 bg-blue-600 text-white text-xs font-bold rounded-lg hover:bg-blue-700">Add Education</button>
                          <button type="button" onClick={() => setShowEduForm(false)} className="px-4 py-2 bg-slate-200 text-slate-700 text-xs font-bold rounded-lg hover:bg-slate-300">Cancel</button>
                        </div>
                      </form>
                    )}

                    <div className="space-y-2">
                      {(profile.education || []).map((edu) => (
                        <div key={edu.id} className="p-3 bg-slate-50 rounded-lg border border-slate-200 flex justify-between items-start">
                          <div>
                            <p className="text-sm font-bold text-slate-900">{edu.degree}</p>
                            <p className="text-xs text-slate-500">{edu.institution} {edu.year && `• ${edu.year}`}</p>
                          </div>
                          <button type="button" onClick={() => handleRemoveEducation(edu.id)} className="text-slate-400 hover:text-red-500 text-sm">×</button>
                        </div>
                      ))}
                    </div>
                  </div>

                  {/* Certifications Section */}
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <label className="text-xs font-bold text-slate-700 uppercase">Certifications</label>
                      <button type="button" onClick={() => setShowCertForm(!showCertForm)} className="text-xs font-bold text-blue-600 hover:text-blue-700">
                        {showCertForm ? 'Cancel' : '+ Add Certification'}
                      </button>
                    </div>

                    {showCertForm && (
                      <div className="p-4 bg-slate-50 rounded-xl border border-slate-200 space-y-3">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                          <div>
                            <label className="text-xs font-bold text-slate-600">Certificate Name</label>
                            <input type="text" value={newCert.name} onChange={(e) => setNewCert({...newCert, name: e.target.value})} className="w-full px-3 py-2 rounded-lg border border-slate-200 text-sm outline-none" placeholder="e.g. AWS Solutions Architect" />
                          </div>
                          <div>
                            <label className="text-xs font-bold text-slate-600">Issuing Organization</label>
                            <input type="text" value={newCert.issuer} onChange={(e) => setNewCert({...newCert, issuer: e.target.value})} className="w-full px-3 py-2 rounded-lg border border-slate-200 text-sm outline-none" placeholder="e.g. Amazon Web Services" />
                          </div>
                        </div>
                        <div className="flex gap-2">
                          <button type="button" onClick={handleAddCertification} className="px-4 py-2 bg-blue-600 text-white text-xs font-bold rounded-lg hover:bg-blue-700">Add</button>
                          <button type="button" onClick={() => setShowCertForm(false)} className="px-4 py-2 bg-slate-200 text-slate-700 text-xs font-bold rounded-lg hover:bg-slate-300">Cancel</button>
                        </div>
                      </div>
                    )}

                    <div className="space-y-2">
                      {(profile.certifications || []).map((cert) => (
                        <div key={cert.id} className="p-3 bg-slate-50 rounded-lg border border-slate-200 flex justify-between items-start">
                          <div>
                            <p className="text-sm font-bold text-slate-900">{cert.name}</p>
                            <p className="text-xs text-slate-500">{cert.issuer} {cert.year && `• ${cert.year}`}</p>
                          </div>
                          <button type="button" onClick={() => handleRemoveCertification(cert.id)} className="text-slate-400 hover:text-red-500 text-sm">×</button>
                        </div>
                      ))}
                    </div>
                  </div>

                  {/* Resume Section */}
                  <div className="p-4 bg-slate-50 rounded-xl border border-slate-200 flex items-center justify-between">
                    <div>
                      <label htmlFor="resume-upload" className="text-xs font-bold text-slate-700">Resume / CV</label>
                      <p className="text-[10px] text-slate-500">{profile.resume_url ? 'Resume_Updated_2024.pdf' : 'No resume uploaded yet'}</p>
                    </div>
                    <input id="resume-upload" name="resume-upload" type="file" className="hidden" accept=".pdf,.doc,.docx" />
                    <button type="button" onClick={() => document.getElementById('resume-upload').click()} className="px-3 py-1.5 bg-white border border-slate-200 text-xs font-bold rounded-lg hover:bg-slate-50">Upload New</button>
                  </div>

                  <div className="flex justify-end pt-4 border-t border-slate-100">
                    <button type="submit" disabled={saving} className="px-8 py-3 bg-blue-600 text-white text-sm font-bold rounded-xl hover:bg-blue-700 shadow-lg shadow-blue-500/25 disabled:opacity-50 transition-all">
                      {saving ? 'Saving...' : 'Save Profile Changes'}
                    </button>
                  </div>
                </form>
              </div>
            </div>
          )}
        </div>
      </main>
    </div>
  );
};

export default CandidateProfile;