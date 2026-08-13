import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

const SampleJobs = [
  { id: 1, title: 'SSC CGL Executive Officer 2026', dept: 'Staff Selection Commission', loc: 'New Delhi', date: '2026-08-30' },
  { id: 2, title: 'UPSC Civil Services Examination', dept: 'Union Public Service Commission', loc: 'All India', date: '2026-09-15' },
  { id: 3, title: 'RRB NTPC Senior Clerk', dept: 'Indian Railways', loc: 'Multiple Zones', date: '2026-08-25' },
  { id: 4, title: 'IBPS PO Assistant Manager', dept: 'Institute of Banking Personnel', loc: 'Pan India', date: '2026-09-01' },
  { id: 5, title: 'DRDO Senior Technical Assistant', dept: 'Defence Research Organisation', loc: 'Bengaluru', date: '2026-09-10' },
  { id: 6, title: 'ISRO Scientist / Engineer', dept: 'Indian Space Research Org', loc: 'Sriharikota', date: '2026-09-20' },
];

const GovJobsPage = () => {
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
        <h1 className="text-3xl font-bold text-slate-900 mb-2">Government Jobs Portal</h1>
        <p className="text-slate-600 mb-8 text-sm">
          Browse official notifications, exam schedules, and government sector vacancies.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {SampleJobs.map((item, index) => {
            const cardJSX = (
              <div className="bg-white rounded-2xl p-6 border border-slate-200 shadow-sm flex flex-col justify-between h-full">
                <div>
                  <span className="text-[10px] uppercase font-bold tracking-wider text-blue-600 bg-blue-50 px-2 py-1 rounded">
                    {item.dept}
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

export default GovJobsPage;