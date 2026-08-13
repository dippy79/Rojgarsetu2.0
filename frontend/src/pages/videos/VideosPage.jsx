import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

const SampleVideos = [
  { id: 1, title: 'How to Crack Technical Interviews in 2026', host: 'Industry Expert Series', views: '14.2k views', duration: '45 mins' },
  { id: 2, title: 'UPSC Preparation Strategy & Booklist', host: 'IAS Mentorship Channel', views: '28.9k views', duration: '60 mins' },
  { id: 3, title: 'System Design Interview Essentials', host: 'Architecture Monthly', views: '19.5k views', duration: '50 mins' },
  { id: 4, title: 'Aptitude & Reasoning Shortcut Tricks', host: 'Govt Prep Hub', views: '32.1k views', duration: '40 mins' },
  { id: 5, title: 'Resume Building & LinkedIn Optimization', host: 'Career Guidance Lab', views: '22.8k views', duration: '30 mins' },
  { id: 6, title: 'Cybersecurity OSINT Live Demonstration', host: 'Ethical Hacking Workshop', views: '11.4k views', duration: '75 mins' },
];

const VideosPage = () => {
  const auth = useAuth();
  const token = localStorage.getItem('token');
  const isLoggedIn = Boolean(token || auth?.isAuthenticated || auth?.user);
  const total = SampleVideos.length;

  return (
    <div className="min-h-screen bg-slate-50">
      {!isLoggedIn && (
        <div className="bg-slate-900 text-white text-center py-3 text-sm">
          <span>Showing demo view — </span>
          <Link to="/login" className="underline font-semibold">
            Login to access all {total} videos →
          </Link>
        </div>
      )}

      <div className="max-w-6xl mx-auto px-6 py-8">
        <h1 className="text-3xl font-bold text-slate-900 mb-2">Career Prep Videos & Webinars</h1>
        <p className="text-slate-600 mb-8 text-sm">
          Watch interactive webinars, interview preparation guides, and exam strategies.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {SampleVideos.map((item, index) => {
            const cardJSX = (
              <div className="bg-white rounded-2xl p-6 border border-slate-200 shadow-sm flex flex-col justify-between h-full">
                <div>
                  <span className="text-[10px] uppercase font-bold tracking-wider text-red-600 bg-red-50 px-2 py-1 rounded">
                    {item.host}
                  </span>
                  <h3 className="text-lg font-bold text-slate-800 mt-3">{item.title}</h3>
                  <p className="text-xs text-slate-500 mt-1">📺 {item.views}</p>
                </div>
                <div className="mt-6 pt-4 border-t border-slate-100 flex items-center justify-between">
                  <span className="text-xs text-slate-400">Length: {item.duration}</span>
                  <button className="px-3 py-1.5 text-xs font-semibold bg-slate-900 text-white rounded-lg">
                    Watch Now
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

export default VideosPage;