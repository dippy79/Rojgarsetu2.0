import React, { useState, useEffect } from 'react';
import { Star } from 'lucide-react';

const COLUMNS = ['Applied', 'Shortlisted', 'Interview', 'Hired', 'Rejected'];

export default function JobApplicants() {
  const [applications, setApplications] = useState([]);

  const fetchApplications = () => {
    const token = localStorage.getItem('rojgar_token');
    fetch('http://localhost:3001/api/v1/company/applications', {
      headers: { Authorization: `Bearer ${token}` }
    })
      .then(res => res.ok ? res.json() : [])
      .then(data => setApplications(Array.isArray(data) ? data : []))
      .catch(err => console.error("Error fetching applications:", err));
  };

  useEffect(() => { fetchApplications(); }, []);

  const updateStatus = async (id, status) => {
    const token = localStorage.getItem('rojgar_token');
    await fetch(`http://localhost:3001/api/v1/applications/${id}/status`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify({ status, recruiter_notes: `Status changed to ${status}` })
    });
    fetchApplications();
  };

  return (
    <div className="p-8 bg-slate-50 min-h-screen">
      <h1 className="text-2xl font-bold text-slate-900 mb-6">Candidate Hiring Pipeline</h1>
      <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
        {COLUMNS.map(col => {
          const colApps = applications.filter(a => (a.status || 'Applied').toLowerCase() === col.toLowerCase());
          return (
            <div key={col} className="bg-slate-100/70 p-4 rounded-2xl min-h-[600px] flex flex-col">
              <h2 className="font-semibold text-slate-700 text-sm mb-4 uppercase tracking-wider flex justify-between">
                {col}
                <span className="bg-slate-200 text-slate-700 px-2 py-0.5 rounded-full text-xs">
                  {colApps.length}
                </span>
              </h2>
              <div className="space-y-3 flex-1 overflow-y-auto">
                {colApps.map(app => (
                  <div key={app.id} className="bg-white p-4 rounded-xl shadow-sm border border-slate-200/80 hover:shadow-md transition">
                    <div className="flex justify-between items-start mb-2">
                      <h3 className="font-semibold text-slate-900">{app.candidate_name || "Candidate"}</h3>
                      <div className="flex text-amber-400">
                        {[...Array(5)].map((_, i) => (
                          <Star key={i} className={`w-3 h-3 ${i < (app.star_rating || 3) ? 'fill-amber-400' : 'text-slate-200'}`} />
                        ))}
                      </div>
                    </div>
                    <p className="text-xs text-slate-500 mb-2">Applied: {app.applied_date ? new Date(app.applied_date).toLocaleDateString() : 'Recently'}</p>
                    <div className="flex flex-wrap gap-1 mb-4">
                      {app.skills?.map(skill => (
                        <span key={skill} className="bg-slate-100 text-slate-600 text-[10px] px-2 py-0.5 rounded-md font-medium">
                          {skill}
                        </span>
                      ))}
                    </div>
                    <div className="flex gap-1 justify-between border-t border-slate-100 pt-3 text-xs flex-wrap">
                      {col !== 'Shortlisted' && <button onClick={() => updateStatus(app.id, 'Shortlisted')} className="text-blue-600 hover:underline">Shortlist</button>}
                      {col !== 'Interview' && <button onClick={() => updateStatus(app.id, 'Interview')} className="text-amber-600 hover:underline">Interview</button>}
                      {col !== 'Hired' && <button onClick={() => updateStatus(app.id, 'Hired')} className="text-emerald-600 hover:underline">Hire</button>}
                      {col !== 'Rejected' && <button onClick={() => updateStatus(app.id, 'Rejected')} className="text-rose-600 hover:underline">Reject</button>}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
