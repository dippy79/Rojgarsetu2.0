import Link from "next/link";
import { Search, Flame, Sparkles, BarChart3 } from "lucide-react";

export default function Home() {
  return (
    <div className="bg-off-white dark:bg-slate-navy text-slate-900 dark:text-off-white transition-colors">
      {/* Hero Section */}
      <section className="py-20 px-4">
        <div className="max-w-4xl mx-auto text-center">
          <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-cyan-100 dark:bg-cyan-900/30 text-cyan-700 dark:text-cyan-300 text-sm font-semibold mb-6">
            <Sparkles className="w-4 h-4" />
            <span>AI-Powered Job & Course Aggregator</span>
          </div>
          <h1 className="text-4xl md:text-6xl font-bold mb-6 bg-gradient-to-r from-cyan-600 to-emerald-500 bg-clip-text text-transparent">
            RojgarSetu 2.0
          </h1>
          <p className="text-xl text-slate-600 dark:text-slate-300 mb-8">
            Your Gateway to Government Jobs & Professional Courses with intelligent matching and gamified career tracking
          </p>
          <div className="flex justify-center gap-4">
            <Link 
              href="/jobs" 
              className="btn-primary px-8 py-3 text-lg flex items-center gap-2"
            >
              <Search className="w-5 h-5" />
              Find Jobs
            </Link>
            <Link 
              href="/courses" 
              className="btn-secondary px-8 py-3 text-lg flex items-center gap-2"
            >
              <Sparkles className="w-5 h-5" />
              Explore Courses
            </Link>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-16 px-4 bg-slate-100 dark:bg-slate-800/50">
        <div className="max-w-7xl mx-auto">
          <h2 className="text-3xl font-bold text-center mb-12">Why RojgarSetu 2.0?</h2>
          <div className="grid md:grid-cols-3 gap-8">
            <div className="bento-card">
              <div className="text-cyan-600 dark:text-cyan-400 mb-4">
                <Sparkles className="w-12 h-12" />
              </div>
              <h3 className="text-xl font-semibold mb-2">AI-Powered Matching</h3>
              <p className="text-slate-600 dark:text-slate-300">
                Intelligent job and course recommendations based on your skills and preferences
              </p>
            </div>
            <div className="bento-card">
              <div className="text-emerald-600 dark:text-emerald-400 mb-4">
                <BarChart3 className="w-12 h-12" />
              </div>
              <h3 className="text-xl font-semibold mb-2">Verified Opportunities</h3>
              <p className="text-slate-600 dark:text-slate-300">
                All jobs and courses are verified and updated regularly from trusted sources
              </p>
            </div>
            <div className="bento-card">
              <div className="text-amber-600 dark:text-amber-400 mb-4">
                <Flame className="w-12 h-12" />
              </div>
              <h3 className="text-xl font-semibold mb-2">Gamified Experience</h3>
              <p className="text-slate-600 dark:text-slate-300">
                Track your application streaks, earn badges, and level up your career profile
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="py-8 px-4 bg-slate-900 text-slate-300">
        <div className="max-w-7xl mx-auto text-center">
          <p>&copy; 2024 RojgarSetu 2.0. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
}