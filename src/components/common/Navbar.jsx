import React, { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

export const Navbar = () => {
  const location = useLocation();
  const { user, logout } = useAuth();
  
  // Default mode toggle: candidate | employer
  const [role, setRole] = useState(
    location.pathname.startsWith('/employer') ? 'employer' : 'candidate'
  );

  const isActive = (path) => location.pathname === path;

  return (
    <header className="bg-slate-900 text-white border-b border-slate-800 sticky top-0 z-40">
      <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
        
        {/* Brand Logo */}
        <div className="flex items-center gap-8">
          <Link to="/" className="flex items-center gap-2 font-black text-xl tracking-tight text-white">
            <span className="w-8 h-8 rounded-lg bg-blue-600 flex items-center justify-center text-white text-sm font-extrabold shadow-md">
              R
            </span>
            <span>Rojgar<span className="text-blue-500">Setu</span></span>
          </Link>

          {/* Role Mode Switcher Toggle */}
          <div className="bg-slate-800/80 p-1 rounded-xl border border-slate-700/60 flex items-center gap-1">
            <button
              type="button"
              onClick={() => setRole('candidate')}
              className={`px-3 py-1 rounded-lg text-xs font-bold transition-all ${
                role === 'candidate'
                  ? 'bg-blue-600 text-white shadow'
                  : 'text-slate-400 hover:text-white'
              }`}
            >
              Candidate View
            </button>
            <button
              type="button"
              onClick={() => setRole('employer')}
              className={`px-3 py-1 rounded-lg text-xs font-bold transition-all ${
                role === 'employer'
                  ? 'bg-purple-600 text-white shadow'
                  : 'text-slate-400 hover:text-white'
              }`}
            >
              Employer Portal
            </button>
          </div>
        </div>

        {/* Dynamic Navigation Links */}
        <nav className="hidden md:flex items-center gap-1 text-xs font-semibold">
          {role === 'candidate' ? (
            <>
              <Link
                to="/jobs"
                className={`px-3.5 py-2 rounded-xl transition-all ${
                  isActive('/jobs')
                    ? 'bg-slate-800 text-white font-bold'
                    : 'text-slate-300 hover:text-white hover:bg-slate-800/50'
                }`}
              >
                🔍 Find Jobs
              </Link>
              <Link
                to="/candidate/applications"
                className={`px-3.5 py-2 rounded-xl transition-all ${
                  isActive('/candidate/applications')
                    ? 'bg-slate-800 text-white font-bold'
                    : 'text-slate-300 hover:text-white hover:bg-slate-800/50'
                }`}
              >
                📁 Applications Tracker
              </Link>
              <Link
                to="/candidate/profile"
                className={`px-3.5 py-2 rounded-xl transition-all ${
                  isActive('/candidate/profile')
                    ? 'bg-slate-800 text-white font-bold'
                    : 'text-slate-300 hover:text-white hover:bg-slate-800/50'
                }`}
              >
                ✨ AI Skill Gap & Profile
              </Link>
            </>
          ) : (
            <>
              <Link
                to="/employer/dashboard"
                className={`px-3.5 py-2 rounded-xl transition-all ${
                  isActive('/employer/dashboard')
                    ? 'bg-slate-800 text-white font-bold'
                    : 'text-slate-300 hover:text-white hover:bg-slate-800/50'
                }`}
              >
                📊 Dashboard
              </Link>
              <Link
                to="/employer/post-job"
                className={`px-3.5 py-2 rounded-xl transition-all ${
                  isActive('/employer/post-job')
                    ? 'bg-slate-800 text-white font-bold'
                    : 'text-slate-300 hover:text-white hover:bg-slate-800/50'
                }`}
              >
                ➕ Post New Job
              </Link>
              <Link
                to="/employer/applicants"
                className={`px-3.5 py-2 rounded-xl transition-all ${
                  isActive('/employer/applicants')
                    ? 'bg-slate-800 text-white font-bold'
                    : 'text-slate-300 hover:text-white hover:bg-slate-800/50'
                }`}
              >
                👥 AI Applicant Ranker
              </Link>
            </>
          )}
        </nav>

        {/* User Info / Profile Pill */}
        <div className="flex items-center gap-3">
          <div className="text-right hidden sm:block">
            <p className="text-xs font-bold text-white leading-tight">
              {user?.name || 'Simranjeet Singh'}
            </p>
            <p className="text-[10px] text-slate-400 capitalize">{role} Account</p>
          </div>
          <div className="w-8 h-8 rounded-full bg-slate-700 border border-slate-600 flex items-center justify-center text-xs font-bold text-white">
            {user?.name ? user.name[0] : 'S'}
          </div>
        </div>

      </div>
    </header>
  );
};

export default Navbar;