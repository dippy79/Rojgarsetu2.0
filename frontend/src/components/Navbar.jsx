import React, { useState, useRef, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { User, LogOut, LayoutDashboard, ChevronDown, Briefcase } from 'lucide-react';
import { useAuth } from '../context/AuthContext';

import NotificationBell from './NotificationBell';

const Navbar = () => {
  const router = useRouter();
  const { user, logout, isAuthenticated } = useAuth();
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const dropdownRef = useRef(null);

  const isLoggedIn = Boolean(isAuthenticated || user);

  const currentUser = user || {
    name: 'User',
    role: 'candidate',
  };

  // Close dropdown on outside click
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setDropdownOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleLogout = () => {
    setDropdownOpen(false);
    if (logout) logout();
    router.push('/login');
  };

  const getInitials = (name) => {
    if (!name) return 'U';
    return name
      .split(' ')
      .filter(Boolean)
      .map((n) => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  };

  const navLinkStyle = (path) => {
    const isActive = router.pathname === path;
    return isActive
      ? 'text-blue-600 font-semibold border-b-2 border-blue-600 pb-1 transition-all'
      : 'text-slate-600 hover:text-slate-900 font-medium transition-all';
  }

  const userRole = (currentUser.role || 'candidate').toLowerCase();
  const dashboardPath = `/dashboard/${userRole}`;
  const profilePath = userRole === 'company' || userRole === 'employer' ? '/company/profile' : '/candidate/profile';

  useEffect(() => {
    let ws;
    if (isLoggedIn && typeof window !== 'undefined') {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const host = process.env.NEXT_PUBLIC_WS_HOST || 'localhost:3001';
      ws = new WebSocket(`${protocol}//${host}/api/v1/ws`);

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          if (data.type === 'notification') {
             // Logic to show notification
             console.log("Real-time signal received:", data);
          }
        } catch (e) { console.error("WS parse error:", e); }
      };
    }
    return () => ws?.close();
  }, [isLoggedIn]);

  return (
    <nav className="sticky top-0 z-50 bg-white/80 backdrop-blur-xl border-b border-slate-200 h-16 flex items-center justify-between px-6 shadow-sm">
      {/* Brand Logo */}
      <Link href="/" className="flex items-center gap-2 text-xl font-bold text-slate-900">
        <div className="bg-slate-900 text-white p-2 rounded-xl">
          <Briefcase className="w-5 h-5 text-blue-400" />
        </div>
        <span>Rojgar<span className="text-blue-600">Setu</span></span>
      </Link>

      {/* Existing Nav Links Preserved Exactly */}
      <ul className="hidden lg:flex items-center gap-6 text-sm font-medium">
        <li>
          <Link href="/" className={navLinkStyle('/')}>🏛️ Home</Link>
        </li>
        <li>
          <Link href="/gov-jobs" className={navLinkStyle('/gov-jobs')}>🛡️ Govt Jobs</Link>
        </li>
        <li>
          <Link href="/private-jobs" className={navLinkStyle('/private-jobs')}>🏢 Private Jobs</Link>
        </li>
        <li>
          <Link href="/courses" className={navLinkStyle('/courses')}>📚 Courses</Link>
        </li>
        <li>
          <Link href="/videos" className={navLinkStyle('/videos')}>🎥 Videos</Link>
        </li>
        <li>
          <Link href="/govt-forms" className={navLinkStyle('/govt-forms')}>🗂️ Govt Forms</Link>
        </li>
      </ul>

      {/* Auth Action Section */}
      <div className="flex items-center gap-4">
        {isLoggedIn && <NotificationBell />}
        {!isLoggedIn ? (
          /* Unauthenticated State */
          <div className="flex items-center gap-3">
            <Link
              href="/login"
              className="text-slate-700 hover:text-slate-900 text-sm font-semibold px-4 py-2 rounded-xl hover:bg-slate-100 transition-all"
            >
              Login
            </Link>
            <Link
              href="/register"
              className="bg-slate-900 hover:bg-slate-800 text-white px-4 py-2 rounded-xl text-sm font-semibold shadow-md transition-all"
            >
              Register
            </Link>
          </div>
        ) : (
          /* Authenticated State */
          <div className="relative" ref={dropdownRef}>
            <button
              onClick={() => setDropdownOpen(!dropdownOpen)}
              className="flex items-center gap-3 p-1.5 rounded-xl hover:bg-slate-100 transition-all cursor-pointer border border-transparent hover:border-slate-200"
            >
              {/* Avatar Circle */}
              <div className="w-9 h-9 rounded-full bg-slate-900 text-white font-bold flex items-center justify-center text-xs tracking-wider shadow-sm">
                {getInitials(currentUser.name)}
              </div>

              {/* User Name & Role */}
              <div className="hidden sm:flex flex-col text-left">
                <span className="text-sm font-semibold text-slate-900 leading-tight">
                  {currentUser.name}
                </span>
                <span className="text-[10px] uppercase tracking-wider font-semibold text-blue-600 bg-blue-50 px-1.5 py-0.5 rounded w-fit mt-0.5">
                  {currentUser.role}
                </span>
              </div>

              <ChevronDown className={`w-4 h-4 text-slate-500 transition-transform duration-200 ${dropdownOpen ? 'rotate-180' : ''}`} />
            </button>

            {/* Dropdown Menu */}
            {dropdownOpen && (
              <div className="absolute right-0 mt-2 w-56 bg-white rounded-2xl shadow-xl border border-slate-100 py-2 z-50 animate-in fade-in slide-in-from-top-2 duration-150">
                <div className="px-4 py-2.5 border-b border-slate-100 sm:hidden">
                  <p className="text-sm font-semibold text-slate-900">{currentUser.name}</p>
                </div>

                <Link
                  href={dashboardPath}
                  onClick={() => setDropdownOpen(false)}
                  className="flex items-center gap-3 px-4 py-2.5 text-sm text-slate-700 hover:bg-slate-50 hover:text-blue-600 transition-all font-medium"
                >
                  <LayoutDashboard className="w-4 h-4" />
                  Dashboard
                </Link>

                <Link
                  href={profilePath}
                  onClick={() => setDropdownOpen(false)}
                  className="flex items-center gap-3 px-4 py-2.5 text-sm text-slate-700 hover:bg-slate-50 hover:text-blue-600 transition-all font-medium"
                >
                  <User className="w-4 h-4" />
                  Profile
                </Link>

                <div className="my-1 border-t border-slate-100"></div>

                <button
                  onClick={handleLogout}
                  className="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-red-600 hover:bg-red-50 transition-all font-medium text-left cursor-pointer"
                >
                  <LogOut className="w-4 h-4" />
                  Logout
                </button>
              </div>
            )}
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navbar;