import React, { useState, useEffect } from 'react';
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';
import { Play } from 'lucide-react';

export default function AdminDashboard() {
  const [tab, setTab] = useState('overview');
  const [stats, setStats] = useState(null);
  const [users, setUsers] = useState([]);
  const [crawlerStats, setCrawlerStats] = useState({ health: 'IDLE', logs: [] });
  const [emails, setEmails] = useState([]);

  useEffect(() => {
    const token = localStorage.getItem('rojgar_token');
    const headers = { Authorization: `Bearer ${token}` };

    if (tab === 'overview') {
      fetch('http://localhost:3001/api/v1/stats', { headers }).then(r => r.ok ? r.json() : null).then(setStats);
    } else if (tab === 'users') {
      fetch('http://localhost:3001/api/v1/admin/users', { headers }).then(r => r.ok ? r.json() : []).then(setUsers);
    } else if (tab === 'crawler') {
      fetch('http://localhost:3001/api/crawler/health').then(r => r.ok ? r.json() : { status: 'UNKNOWN' }).then(h => {
        setCrawlerStats(prev => ({ ...prev, health: h.status }));
      });
      fetch('http://localhost:3001/api/crawler/stats').then(r => r.ok ? r.json() : { logs: [] }).then(s => {
        setCrawlerStats(prev => ({ ...prev, logs: s.logs || [] }));
      });
    } else if (tab === 'email') {
      fetch('http://localhost:3001/api/v1/admin/emails', { headers }).then(r => r.ok ? r.json() : []).then(setEmails);
    }
  }, [tab]);

  const triggerCrawl = async () => {
    await fetch('http://localhost:3001/api/crawler/crawl', { method: 'POST' });
    setTab('crawler');
  };

  return (
    <div className="p-8 bg-slate-50 min-h-screen">
      <h1 className="text-2xl font-bold text-slate-900 mb-6">System Administration</h1>

      <div className="flex border-b border-slate-200 mb-6 gap-6">
        {['overview', 'users', 'crawler', 'email'].map(t => (
          <button
            key={t}
            onClick={() => setTab(t)}
            className={`pb-3 font-medium capitalize text-sm border-b-2 transition ${tab === t ? 'border-indigo-600 text-indigo-600' : 'border-transparent text-slate-500 hover:text-slate-800'}`}
          >
            {t}
          </button>
        ))}
      </div>

      {tab === 'overview' && (
        <div className="space-y-6">
          <div className="grid grid-cols-4 gap-6">
            <div className="p-6 bg-white rounded-2xl shadow-sm border border-slate-100">
              <p className="text-sm text-slate-500">Total Users</p>
              <p className="text-3xl font-bold">{(stats?.total_candidates || 0) + (stats?.total_companies || 0)}</p>
            </div>
            <div className="p-6 bg-white rounded-2xl shadow-sm border border-slate-100">
              <p className="text-sm text-slate-500">Total Jobs</p>
              <p className="text-3xl font-bold">{stats?.total_jobs || 0}</p>
            </div>
            <div className="p-6 bg-white rounded-2xl shadow-sm border border-slate-100">
              <p className="text-sm text-slate-500">Total Companies</p>
              <p className="text-3xl font-bold">{stats?.total_companies || 0}</p>
            </div>
            <div className="p-6 bg-white rounded-2xl shadow-sm border border-slate-100">
              <p className="text-sm text-slate-500">Placements</p>
              <p className="text-3xl font-bold">{stats?.total_placements || 0}</p>
            </div>
          </div>
          <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 h-80">
            <h3 className="font-semibold text-slate-800 mb-4">Application Trends (7 Days)</h3>
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={stats?.chart_data || []}>
                <XAxis dataKey="day" />
                <YAxis />
                <Tooltip />
                <Line type="monotone" dataKey="applications" stroke="#4f46e5" strokeWidth={3} />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>
      )}

      {tab === 'users' && (
        <div className="bg-white rounded-2xl p-6 shadow-sm border border-slate-100">
          <table className="w-full text-left text-sm">
            <thead>
              <tr className="border-b text-slate-400 uppercase text-xs">
                <th className="py-3">Email</th><th className="py-3">Role</th><th className="py-3">Created At</th><th className="py-3">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {users.map(u => (
                <tr key={u.id}>
                  <td className="py-3 font-medium">{u.email}</td>
                  <td className="py-3 capitalize">{u.role}</td>
                  <td className="py-3">{new Date(u.created_at).toLocaleDateString()}</td>
                  <td className="py-3 flex gap-2">
                    <button className="px-3 py-1 bg-rose-50 text-rose-600 rounded-lg text-xs font-semibold">Ban</button>
                    <button className="px-3 py-1 bg-emerald-50 text-emerald-600 rounded-lg text-xs font-semibold">Verify</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {tab === 'crawler' && (
        <div className="space-y-6">
          <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 flex items-center justify-between">
            <div>
              <span className={`inline-block px-3 py-1 rounded-full text-xs font-bold ${crawlerStats.health === 'RUNNING' ? 'bg-emerald-100 text-emerald-700' : 'bg-slate-100 text-slate-700'}`}>
                STATUS: {crawlerStats.health}
              </span>
            </div>
            <button onClick={triggerCrawl} className="flex items-center gap-2 bg-indigo-600 text-white px-5 py-2.5 rounded-xl font-medium">
              <Play className="w-4 h-4" /> Run Now
            </button>
          </div>
          <div className="bg-slate-900 text-slate-100 p-4 rounded-2xl font-mono text-xs h-64 overflow-y-auto">
            {crawlerStats.logs.map((log, index) => <div key={index}>{log}</div>)}
          </div>
        </div>
      )}

      {tab === 'email' && (
        <div className="bg-white rounded-2xl p-6 shadow-sm border border-slate-100">
          <table className="w-full text-left text-sm">
            <thead>
              <tr className="border-b text-slate-400 uppercase text-xs">
                <th className="py-3">Recipient</th><th className="py-3">Subject</th><th className="py-3">Status</th><th className="py-3">Attempts</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {emails.map(e => (
                <tr key={e.id}>
                  <td className="py-3">{e.to_email}</td>
                  <td className="py-3">{e.subject}</td>
                  <td className="py-3">
                    <span className={`px-2.5 py-0.5 rounded-full text-xs font-bold ${e.status === 'SENT' ? 'bg-emerald-100 text-emerald-700' : e.status === 'PENDING' ? 'bg-amber-100 text-amber-700' : 'bg-rose-100 text-rose-700'}`}>
                      {e.status}
                    </span>
                  </td>
                  <td className="py-3">{e.attempts}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
