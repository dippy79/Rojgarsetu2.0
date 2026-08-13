import React, { useState } from 'react';
import { useAuth } from '../../hooks/useAuth';

export const CandidateProfile = () => {
  const { user } = useAuth();
  const [toast, setToast] = useState(null);

  // Candidate Profile State
  const [profile, setProfile] = useState({
    name: user?.name || 'Simranjeet Singh',
    email: user?.email || 'simranjeet@example.com',
    title: 'Full Stack & AI Engineer',
    bio: 'Passionate full-stack developer with expertise in React, Node.js, Flutter, and AI-driven platforms.',
    location: 'Delhi NCR, India',
    targetSalary: '₹20 - 28 LPA',
    targetRole: 'Senior Full Stack Engineer',
    skills: ['React.js', 'Node.js', 'Flutter', 'Python', 'Tailwind CSS', 'PostgreSQL', 'Docker'],
    resumeName: 'Resume_Simranjeet_2026.pdf',
    lastUpdated: '2026-03-01'
  });

  const [newSkill, setNewSkill] = useState('');
  const [saving, setSaving] = useState(false);

  // AI Skill Gap State
  const skillGapAnalysis = {
    targetRole: 'Lead AI & Full Stack Architect',
    readinessScore: 82,
    matchingSkills: ['React.js', 'Node.js', 'Python', 'PostgreSQL', 'Docker'],
    missingSkills: ['System Design (LLD/HLD)', 'Kubernetes', 'GraphQL', 'Vector Databases (Qdrant/Pinecone)'],
    recommendedRoadmap: [
      { skill: 'Vector Databases', course: 'Pinecone & Qdrant Essentials', effort: '2 weeks' },
      { skill: 'System Design', course: 'High-Scale Distributed Systems Architecture', effort: '4 weeks' },
      { skill: 'Kubernetes', course: 'K8s for Developers & Microservices', effort: '2 weeks' }
    ]
  };

  const showToast = (message, type = 'success') => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 3500);
  };

  const handleAddSkill = (e) => {
    e.preventDefault();
    if (newSkill.trim() && !profile.skills.includes(newSkill.trim())) {
      setProfile({ ...profile, skills: [...profile.skills, newSkill.trim()] });
      setNewSkill('');
      showToast('Skill added');
    }
  };

  const handleRemoveSkill = (skillToRemove) => {
    setProfile({
      ...profile,
      skills: profile.skills.filter((s) => s !== skillToRemove)
    });
  };

  const handleSaveProfile = (e) => {
    e.preventDefault();
    setSaving(true);
    setTimeout(() => {
      setSaving(false);
      showToast('Profile & AI preferences saved successfully!');
    }, 800);
  };

  return (
    <div className="min-h-screen bg-slate-50 p-8 font-sans">
      {toast && (
        <div className={`fixed top-6 right-6 z-50 px-5 py-3 rounded-xl shadow-lg border text-sm font-semibold flex items-center gap-2 ${
          toast.type === 'success' ? 'bg-emerald-50 text-emerald-800 border-emerald-200' : 'bg-blue-50 text-blue-800 border-blue-200'
        }`}>
          <span>{toast.message}</span>
        </div>
      )}

      <div className="max-w-6xl mx-auto space-y-8">
        <header>
          <h1 className="text-3xl font-bold text-slate-900 tracking-tight">Candidate Profile & AI Insights</h1>
          <p className="text-slate-500 text-sm mt-1 font-medium">
            Manage your professional identity and analyze skill gaps for your dream target roles.
          </p>
        </header>

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
          {/* Left Column: Edit Profile */}
          <div className="lg:col-span-7 space-y-6">
            <form onSubmit={handleSaveProfile} className="bg-white p-6 rounded-2xl border border-slate-200/60 shadow-sm space-y-5">
              <h2 className="text-lg font-bold text-slate-900 border-b border-slate-100 pb-3">Personal Details</h2>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label htmlFor="cand_name" className="block text-xs font-semibold text-slate-700 uppercase mb-1">Full Name</label>
                  <input
                    type="text"
                    id="cand_name"
                    value={profile.name}
                    onChange={(e) => setProfile({ ...profile, name: e.target.value })}
                    className="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-slate-900 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                  />
                </div>

                <div>
                  <label htmlFor="cand_email" className="block text-xs font-semibold text-slate-700 uppercase mb-1">Email Address</label>
                  <input
                    type="email"
                    id="cand_email"
                    value={profile.email}
                    disabled
                    className="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 bg-slate-50 text-slate-500 text-sm cursor-not-allowed"
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label htmlFor="cand_title" className="block text-xs font-semibold text-slate-700 uppercase mb-1">Professional Title</label>
                  <input
                    type="text"
                    id="cand_title"
                    value={profile.title}
                    onChange={(e) => setProfile({ ...profile, title: e.target.value })}
                    className="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-slate-900 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                  />
                </div>

                <div>
                  <label htmlFor="cand_location" className="block text-xs font-semibold text-slate-700 uppercase mb-1">Location</label>
                  <input
                    type="text"
                    id="cand_location"
                    value={profile.location}
                    onChange={(e) => setProfile({ ...profile, location: e.target.value })}
                    className="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-slate-900 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                  />
                </div>
              </div>

              <div>
                <label htmlFor="cand_bio" className="block text-xs font-semibold text-slate-700 uppercase mb-1">Professional Summary</label>
                <textarea
                  id="cand_bio"
                  rows={3}
                  value={profile.bio}
                  onChange={(e) => setProfile({ ...profile, bio: e.target.value })}
                  className="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-slate-900 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                />
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-2">
                <div>
                  <label htmlFor="cand_target_role" className="block text-xs font-semibold text-slate-700 uppercase mb-1">Target Role</label>
                  <input
                    type="text"
                    id="cand_target_role"
                    value={profile.targetRole}
                    onChange={(e) => setProfile({ ...profile, targetRole: e.target.value })}
                    className="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-slate-900 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                  />
                </div>

                <div>
                  <label htmlFor="cand_target_salary" className="block text-xs font-semibold text-slate-700 uppercase mb-1">Expected Salary</label>
                  <input
                    type="text"
                    id="cand_target_salary"
                    value={profile.targetSalary}
                    onChange={(e) => setProfile({ ...profile, targetSalary: e.target.value })}
                    className="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-slate-900 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                  />
                </div>
              </div>

              {/* Skills Section */}
              <div className="pt-3">
                <label className="block text-xs font-semibold text-slate-700 uppercase mb-2">Skills & Expertise</label>
                <div className="flex flex-wrap gap-2 mb-3">
                  {profile.skills.map((skill) => (
                    <span key={skill} className="px-3 py-1 bg-slate-100 text-slate-800 text-xs font-semibold rounded-lg flex items-center gap-1.5 border border-slate-200">
                      {skill}
                      <button
                        type="button"
                        onClick={() => handleRemoveSkill(skill)}
                        className="text-slate-400 hover:text-red-500 font-bold ml-1"
                      >
                        ×
                      </button>
                    </span>
                  ))}
                </div>

                <div className="flex gap-2">
                  <input
                    type="text"
                    placeholder="Add a new skill (e.g. Redis, GraphQL)"
                    value={newSkill}
                    onChange={(e) => setNewSkill(e.target.value)}
                    className="flex-1 px-3.5 py-2 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                  />
                  <button
                    type="button"
                    onClick={handleAddSkill}
                    className="px-4 py-2 bg-slate-800 text-white text-xs font-bold rounded-xl hover:bg-slate-900"
                  >
                    Add Skill
                  </button>
                </div>
              </div>

              {/* Resume Box */}
              <div className="p-4 bg-slate-50 rounded-xl border border-slate-200/80 flex items-center justify-between">
                <div>
                  <p className="text-xs font-bold text-slate-700">Active Resume</p>
                  <p className="text-xs text-slate-500">{profile.resumeName}</p>
                </div>
                <button type="button" className="px-3 py-1.5 bg-white border border-slate-200 text-slate-700 text-xs font-semibold rounded-lg hover:bg-slate-100">
                  Replace PDF
                </button>
              </div>

              <div className="pt-2 flex justify-end">
                <button
                  type="submit"
                  disabled={saving}
                  className="px-6 py-2.5 bg-blue-600 text-white text-xs font-bold rounded-xl hover:bg-blue-700 shadow-md transition-all"
                >
                  {saving ? 'Saving...' : 'Save Profile Changes'}
                </button>
              </div>
            </form>
          </div>

          {/* Right Column: AI Skill Gap Analysis */}
          <div className="lg:col-span-5 space-y-6">
            <div className="bg-gradient-to-br from-slate-900 to-slate-800 text-white p-6 rounded-2xl shadow-md border border-slate-800 space-y-5">
              <div className="flex items-center justify-between">
                <div>
                  <span className="text-[10px] uppercase font-mono tracking-wider text-blue-400 font-bold">AI Skill Gap Engine</span>
                  <h3 className="text-lg font-bold mt-0.5">{skillGapAnalysis.targetRole}</h3>
                </div>
                <div className="text-right">
                  <div className="text-2xl font-black text-emerald-400">{skillGapAnalysis.readinessScore}%</div>
                  <span className="text-[10px] text-slate-400 font-medium">Readiness Index</span>
                </div>
              </div>

              {/* Readiness Progress Bar */}
              <div className="w-full bg-slate-700/60 rounded-full h-2 overflow-hidden">
                <div className="bg-emerald-400 h-full rounded-full" style={{ width: `${skillGapAnalysis.readinessScore}%` }}></div>
              </div>

              <div>
                <h4 className="text-xs font-bold uppercase text-slate-400 tracking-wider mb-2">Identified Skill Gaps</h4>
                <div className="space-y-2">
                  {skillGapAnalysis.missingSkills.map((gap, i) => (
                    <div key={i} className="flex items-center gap-2 text-xs bg-slate-800/80 px-3 py-2 rounded-xl border border-slate-700/60 text-slate-200">
                      <span className="text-amber-400 font-bold">⚠️</span>
                      <span>{gap}</span>
                    </div>
                  ))}
                </div>
              </div>

              <div>
                <h4 className="text-xs font-bold uppercase text-slate-400 tracking-wider mb-2">Suggested Upskilling Plan</h4>
                <div className="space-y-2.5">
                  {skillGapAnalysis.recommendedRoadmap.map((item, idx) => (
                    <div key={idx} className="p-3 bg-slate-800/40 rounded-xl border border-slate-700/40 text-xs space-y-1">
                      <div className="flex justify-between items-center font-bold text-white">
                        <span>{item.skill}</span>
                        <span className="text-[10px] font-mono bg-blue-500/20 text-blue-300 px-2 py-0.5 rounded">
                          Est. {item.effort}
                        </span>
                      </div>
                      <p className="text-slate-400 text-[11px]">{item.course}</p>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CandidateProfile;