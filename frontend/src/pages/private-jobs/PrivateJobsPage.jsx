import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

const SampleJobs = [
  { id: 1, title: 'Senior Full Stack Developer', company: 'TechCorp Solutions', loc: 'Bengaluru / Remote', date: '2026-08-28' },
  { id: 2, title: 'AI & Data Pipeline Engineer', company: 'DataScale Systems', loc: 'Hyderabad', date: '2026-09-05' },
  { id: 3, title: 'Product UI/UX Designer', company: 'DesignCraft Studio', loc: 'Mumbai', date: '2026-08-31' },
  { id: 4, title: 'Backend Node.js Architect', company: 'CloudScale Inc', loc: 'Gurugram', date: '2026-09-12' },
  { id: 5, title: 'Cybersecurity Operations Lead', company: 'SecureNet India', loc: 'Pune', date: '2026-09-18' },
  { id: 6, title: 'DevOps & SRE Specialist', company: 'InfraScale Global', loc: 'Remote', date: '2026-09-22' },
];

const PrivateJobsPage = () => {
  const auth = useAuth();
  const token = localStorage.getItem('token');
  const isLoggedIn = Boolean(token || auth?.isAuthenticated || auth?.user);
  const total = SampleJobs.length;

  return (
    <div className="min-h-screen bg-slate-50">
      {!isLoggedIn && (
        <div className="bg-slate-900 text-white text-center py-3 text-sm">
          <span>Showing demo view — </span>
          <Link to="/login" className="underline font-semibold">
            Login to access all {total} jobs →
          </Link>
        </div>
      )}

      <div className="max-w-6xl mx-auto px-6 py-8">
        <h1 className="text-3xl font-bold text-slate-900 mb-2">Private Sector Job Openings</h1>
        <p className="text-slate-600 mb-8 text-sm">
          Discover verified technology, design, product, and management roles across top companies.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {SampleJobs.map((item, index) => {
            const cardJSX = (
              <div className="bg-white rounded-2xl p-6 border border-slate-200 shadow-sm flex flex-col justify-between h-full">
                <div>
                  <span className="text-[10px] uppercase font-bold tracking-wider text-teal-600 bg-teal-50 px-2 py-1 rounded">
                    {item.company}
                  </span>
                  <h3 className="text-lg font-bold text-slate-800 mt-3">{item.title}</h3>
                  <p className="text-xs text-slate-500 mt-1">📍 {item.loc}</p>
                </div>
                <div className="mt-6 pt-4 border-t border-slate-100 flex items-center justify-between">
                  <span className="text-xs text-slate-400">Apply by {item.date}</span>
                  <button className="px-3 py-1.5 text-xs font-semibold bg-slate-900 text-white rounded-lg">
                    View Details
                  </button>
                </div>
              </div>
            );

            if (!isLoggedIn && index >= 3) {
              return (
                <div key={item.id} className="relative">
                  <div className="blur-sm pointer-events-none">{cardJSX}</div>
                  <div className="absolute inset-0 flex flex-col items-center justify-center bg-white/60 backdrop-blur-sm rounded-2xl p-4 text-center z-10">
                    <p className="text-slate-700 font-semibold text-sm mb-2">Login to see more jobs</p>
                    <Link to="/login" className="bg-slate-900 text-white text-xs px-4 py-2 rounded-lg font-medium shadow-sm hover:bg-slate-800 transition-all">
                      Login / Register
                    </Link>
                  </div>
                </div>
              );
            }

            return <div key={item.id}>{cardJSX}</div>;
          })}
        </div>
      </div>
    </div>
  );
};

export default PrivateJobsPage;