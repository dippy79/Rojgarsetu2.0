import React from 'react';
import Link from 'next/link';
import { Briefcase, ShieldCheck, GraduationCap, PlayCircle, FileText, Sparkles } from 'lucide-react';

import PublicCounter from '../components/PublicCounter';

export default function Home() {
  return (
    <div className="min-h-screen bg-slate-50 font-sans">
      {/* Hero Section */}
      <section className="relative overflow-hidden pt-20 pb-32 bg-white">
        <div className="absolute inset-0 bg-[radial-gradient(#e2e8f0_1px,transparent_1px)] [background-size:16px_16px] [mask-image:radial-gradient(ellipse_50%_50%_at_50%_50%,#000_70%,transparent_100%)]"></div>

        <div className="relative max-w-7xl mx-auto px-6 text-center">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-50 border border-blue-100 text-blue-600 text-xs font-bold mb-6 animate-bounce">
            <Sparkles className="w-3.5 h-3.5" />
            <span>AI-Powered Career Portal for India</span>
          </div>

          <h1 className="text-5xl md:text-7xl font-black text-slate-900 tracking-tight mb-6">
            Your Gateway to <br />
            <span className="text-blue-600">Digital Empowerment</span>
          </h1>

          <p className="text-lg md:text-xl text-slate-500 max-w-2xl mx-auto mb-10 font-medium">
            Sourcing verified government jobs, premium private sector roles, and world-class educational content in one unified platform.
          </p>

          <div className="flex flex-wrap justify-center gap-4">
            <Link href="/gov-jobs" className="px-8 py-4 bg-slate-900 text-white font-bold rounded-2xl hover:bg-slate-800 transition-all shadow-xl shadow-slate-200">
              Browse Gov Jobs
            </Link>
            <Link href="/private-jobs" className="px-8 py-4 bg-white text-slate-900 border-2 border-slate-100 font-bold rounded-2xl hover:border-blue-400 transition-all">
              Private Sector
            </Link>
          </div>
        </div>
      </section>

      {/* Platform Stats Counter */}
      <section className="max-w-7xl mx-auto px-6 mb-20">
         <PublicCounter />
      </section>

      {/* Categories Grid */}
      <section className="max-w-7xl mx-auto px-6 -mt-20 relative z-10">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <Link href="/gov-jobs" className="group bg-white/80 backdrop-blur-xl p-8 rounded-3xl border border-slate-200 shadow-sm hover:border-blue-400 hover:shadow-xl transition-all">
            <div className="w-12 h-12 bg-blue-50 text-blue-600 rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform">
              <ShieldCheck className="w-6 h-6" />
            </div>
            <h3 className="text-xl font-bold text-slate-900 mb-2">Govt Jobs</h3>
            <p className="text-sm text-slate-500 font-medium">SSC, UPSC, RRB, and State PSC notifications.</p>
          </Link>

          <Link href="/private-jobs" className="group bg-white/80 backdrop-blur-xl p-8 rounded-3xl border border-slate-200 shadow-sm hover:border-emerald-400 hover:shadow-xl transition-all">
            <div className="w-12 h-12 bg-emerald-50 text-emerald-600 rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform">
              <Briefcase className="w-6 h-6" />
            </div>
            <h3 className="text-xl font-bold text-slate-900 mb-2">Private Jobs</h3>
            <p className="text-sm text-slate-500 font-medium">Top tech and corporate openings across India.</p>
          </Link>

          <Link href="/courses" className="group bg-white/80 backdrop-blur-xl p-8 rounded-3xl border border-slate-200 shadow-sm hover:border-purple-400 hover:shadow-xl transition-all">
            <div className="w-12 h-12 bg-purple-50 text-purple-600 rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform">
              <GraduationCap className="w-6 h-6" />
            </div>
            <h3 className="text-xl font-bold text-slate-900 mb-2">Courses</h3>
            <p className="text-sm text-slate-500 font-medium">Skill-based learning from industry experts.</p>
          </Link>

          <Link href="/videos" className="group bg-white/80 backdrop-blur-xl p-8 rounded-3xl border border-slate-200 shadow-sm hover:border-rose-400 hover:shadow-xl transition-all">
            <div className="w-12 h-12 bg-rose-50 text-rose-600 rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform">
              <PlayCircle className="w-6 h-6" />
            </div>
            <h3 className="text-xl font-bold text-slate-900 mb-2">Video Prep</h3>
            <p className="text-sm text-slate-500 font-medium">Interview guides and exam preparation videos.</p>
          </Link>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-32 max-w-7xl mx-auto px-6">
        <div className="flex flex-col lg:flex-row items-center gap-16">
          <div className="flex-1 space-y-8">
            <h2 className="text-4xl font-black text-slate-900 leading-tight">
              Empowering your career <br />
              with <span className="text-blue-600">AI Intelligence.</span>
            </h2>
            <ul className="space-y-4">
              {[
                { icon: Sparkles, text: 'AI-based job matching tailored to your skills' },
                { icon: FileText, text: 'Smart form deadline tracking system' },
                { icon: ShieldCheck, text: 'Verified and spam-free job listings' },
              ].map((item, i) => (
                <li key={i} className="flex items-center gap-3 text-slate-700 font-bold">
                  <div className="w-6 h-6 bg-blue-100 text-blue-600 rounded-full flex items-center justify-center shrink-0">
                    <item.icon className="w-3.5 h-3.5" />
                  </div>
                  {item.text}
                </li>
              ))}
            </ul>
            <div className="pt-4">
              <Link href="/register" className="text-blue-600 font-black hover:underline underline-offset-4">
                Join thousands of candidates today →
              </Link>
            </div>
          </div>

          <div className="flex-1 grid grid-cols-2 gap-4">
            <div className="space-y-4 pt-8">
              <div className="bg-white p-6 rounded-3xl border border-slate-200 shadow-sm h-48 flex flex-col justify-center text-center">
                <span className="text-4xl font-black text-slate-900">500+</span>
                <span className="text-xs font-bold text-slate-400 uppercase tracking-widest mt-2">Active Jobs</span>
              </div>
              <div className="bg-blue-600 p-6 rounded-3xl shadow-xl h-40 flex flex-col justify-center text-center text-white">
                <span className="text-3xl font-black">Verified</span>
                <span className="text-xs font-bold text-blue-100 uppercase tracking-widest mt-2">Source Check</span>
              </div>
            </div>
            <div className="space-y-4">
              <div className="bg-slate-900 p-6 rounded-3xl shadow-xl h-40 flex flex-col justify-center text-center text-white">
                <span className="text-3xl font-black">200+</span>
                <span className="text-xs font-bold text-slate-400 uppercase tracking-widest mt-2">Companies</span>
              </div>
              <div className="bg-white p-6 rounded-3xl border border-slate-200 shadow-sm h-48 flex flex-col justify-center text-center">
                <span className="text-4xl font-black text-slate-900">10k+</span>
                <span className="text-xs font-bold text-slate-400 uppercase tracking-widest mt-2">Candidates</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="bg-slate-900 py-20">
        <div className="max-w-4xl mx-auto px-6 text-center text-white space-y-8">
          <h2 className="text-3xl md:text-5xl font-black tracking-tight">Ready to find your next opportunity?</h2>
          <p className="text-slate-400 font-medium">Join RojgarSetu and let our AI handle the search for you.</p>
          <div className="pt-4 flex flex-wrap justify-center gap-4">
            <Link href="/login" className="px-10 py-4 bg-blue-600 hover:bg-blue-500 text-white font-bold rounded-2xl transition-all shadow-xl shadow-blue-900/20">
              Get Started Now
            </Link>
          </div>
        </div>
      </section>

      <footer className="py-12 border-t border-slate-200 text-center">
        <p className="text-sm font-bold text-slate-400 uppercase tracking-widest">
          © 2026 RojgarSetu 2.0 • Made with ❤️ for India
        </p>
      </footer>
    </div>
  );
}
