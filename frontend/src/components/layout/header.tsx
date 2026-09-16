"use client";

import React, { useState } from "react";
import Link from "next/link";
import { useTheme } from "next-themes";
import { 
  Moon, Sun, Search, Bell, User, Command, 
  Menu, X, Flame, Zap 
} from "lucide-react";
import { cn } from "@/lib/utils";

export function Header() {
  const { theme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);
  const [isSearchOpen, setIsSearchOpen] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  React.useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) {
    return null;
  }

  return (
    <header className="sticky top-0 z-50 bg-white/80 dark:bg-stone-900/80 backdrop-blur-md border-b border-stone-200 dark:border-stone-800">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <div className="flex items-center">
            <Link href="/" className="flex items-center gap-2">
              <div className="p-2 bg-gradient-to-br from-cyan-500 to-stone-500 rounded-lg">
                <Zap className="w-5 h-5 text-white" />
              </div>
              <span className="text-xl font-bold text-stone-900 dark:text-white">
                RojgarSetu <span className="text-cyan-600 dark:text-cyan-400">2.0</span>
              </span>
            </Link>
          </div>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center space-x-1">
            <Link 
              href="/jobs" 
              className="px-4 py-2 rounded-lg text-stone-600 dark:text-stone-300 hover:bg-stone-100 dark:hover:bg-stone-800 hover:text-cyan-600 dark:hover:text-cyan-400 transition-all flex items-center gap-2"
            >
              <Search className="w-4 h-4" />
              Jobs
            </Link>
            <Link 
              href="/courses" 
              className="px-4 py-2 rounded-lg text-stone-600 dark:text-stone-300 hover:bg-stone-100 dark:hover:bg-stone-800 hover:text-cyan-600 dark:hover:text-cyan-400 transition-all flex items-center gap-2"
            >
              <Flame className="w-4 h-4" />
              Courses
            </Link>
            <Link 
              href="/profile" 
              className="px-4 py-2 rounded-lg text-stone-600 dark:text-stone-300 hover:bg-stone-100 dark:hover:bg-stone-800 hover:text-cyan-600 dark:hover:text-cyan-400 transition-all flex items-center gap-2"
            >
              <User className="w-4 h-4" />
              Profile
            </Link>
          </nav>

          {/* Right Actions */}
          <div className="flex items-center space-x-2">
            {/* Search Trigger */}
            <button
              onClick={() => setIsSearchOpen(true)}
              className="hidden md:flex items-center gap-2 px-3 py-1.5 rounded-lg bg-stone-100 dark:bg-stone-800 text-stone-600 dark:text-stone-300 hover:bg-stone-200 dark:hover:bg-stone-700 transition-colors"
            >
              <Command className="w-4 h-4" />
              <span className="text-sm">K</span>
            </button>

            {/* Notifications */}
            <button className="relative p-2 rounded-lg bg-stone-100 dark:bg-stone-800 text-stone-600 dark:text-stone-300 hover:bg-stone-200 dark:hover:bg-stone-700 transition-colors">
              <Bell className="w-5 h-5" />
              <span className="absolute top-1 right-1 w-2 h-2 bg-stone-500 rounded-full"></span>
            </button>

            {/* Theme Toggle */}
            <button
              onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
              className="p-2 rounded-lg bg-stone-100 dark:bg-stone-800 text-stone-600 dark:text-stone-300 hover:bg-stone-200 dark:hover:bg-stone-700 transition-colors"
            >
              {theme === "dark" ? <Sun className="w-5 h-5" /> : <Moon className="w-5 h-5" />}
            </button>

            {/* Mobile Menu Toggle */}
            <button
              onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
              className="md:hidden p-2 rounded-lg bg-stone-100 dark:bg-stone-800 text-stone-600 dark:text-stone-300 hover:bg-stone-200 dark:hover:bg-stone-700 transition-colors"
            >
              {isMobileMenuOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
            </button>
          </div>
        </div>
      </div>

      {/* Mobile Menu */}
      {isMobileMenuOpen && (
        <div className="md:hidden border-t border-stone-200 dark:border-stone-800 bg-white dark:bg-stone-900">
          <div className="px-4 py-4 space-y-2">
            <Link 
              href="/jobs" 
              className="block px-4 py-2 rounded-lg text-stone-600 dark:text-stone-300 hover:bg-stone-100 dark:hover:bg-stone-800"
              onClick={() => setIsMobileMenuOpen(false)}
            >
              Jobs
            </Link>
            <Link 
              href="/courses" 
              className="block px-4 py-2 rounded-lg text-stone-600 dark:text-stone-300 hover:bg-stone-100 dark:hover:bg-stone-800"
              onClick={() => setIsMobileMenuOpen(false)}
            >
              Courses
            </Link>
            <Link 
              href="/profile" 
              className="block px-4 py-2 rounded-lg text-stone-600 dark:text-stone-300 hover:bg-stone-100 dark:hover:bg-stone-800"
              onClick={() => setIsMobileMenuOpen(false)}
            >
              Profile
            </Link>
          </div>
        </div>
      )}

      {/* Search Modal */}
      {isSearchOpen && (
        <div className="fixed inset-0 z-50 flex items-start justify-center pt-20 px-4">
          <div className="absolute inset-0 bg-black/50 backdrop-blur-sm" onClick={() => setIsSearchOpen(false)} />
          <div className="relative w-full max-w-2xl bg-white dark:bg-stone-900 rounded-bento-lg shadow-2xl border border-stone-200 dark:border-stone-800">
            <div className="flex items-center p-4 border-b border-stone-200 dark:border-stone-800">
              <Search className="w-5 h-5 text-stone-400 mr-3" />
              <input
                type="text"
                placeholder="Search jobs, courses, skills..."
                className="flex-1 bg-transparent text-stone-900 dark:text-white placeholder-stone-400 focus:outline-none"
                autoFocus
              />
              <button
                onClick={() => setIsSearchOpen(false)}
                className="p-2 rounded-lg hover:bg-stone-100 dark:hover:bg-stone-800 transition-colors"
              >
                <X className="w-5 h-5 text-stone-400" />
              </button>
            </div>
            <div className="p-4">
              <p className="text-sm text-stone-500 dark:text-stone-400">
                Press <kbd className="px-2 py-1 bg-stone-100 dark:bg-stone-800 rounded text-xs">Enter</kbd> to search
              </p>
            </div>
          </div>
        </div>
      )}
    </header>
  );
}