import React, { useState, useEffect, useCallback } from 'react';
import {
  LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid, AreaChart, Area
} from 'recharts';
import {
  Play, Users, Briefcase, Mail, Cpu, Zap, Search,
  ShieldCheck, ShieldAlert, Clock, ArrowUpRight, Loader2, Globe, Activity
} from 'lucide-react';
import { useAuth } from '../../hooks/useAuth';
import { apiUrl } from '../../apiConfig';

export default function AdminDashboard() {
  const [tab, setTab] = useState('overview');
  const [stats, setStats] = useState(null);
  const [users, setUsers] = useState([]);
  const [crawlerStats, setCrawlerStats] = useState({ health: 'IDLE', logs: [] });
  const [emails, setEmails] = useState([]);
  const [loading, setLoading] = useState(false);

  const fetchData = useCallback(async () => {
    setLoading(true);
    const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');
    const headers = { Authorization: token ? `Bearer ${token}` : '' };

    try {
      if (tab === 'overview') {
        const res = await fetch(apiUrl('/api/v1/stats'), { headers });
        if (res.ok) setStats(await res.json());
      } else if (tab === 'users') {
        const res = await fetch(apiUrl('/api/v1/admin/users'), { headers });
        if (res.ok) setUsers(await res.json());
      } else if (tab === 'crawler') {
        const hRes = await fetch(apiUrl('/api/crawler/health'));
        const sRes = await fetch(apiUrl('/api/crawler/stats'));
        if (hRes.ok) setCrawlerStats(prev => ({ ...prev, health: (hRes.json()).status || 'IDLE' }));
        if (sRes.ok) {
           const sData = await sRes.json();
           setCrawlerStats(prev => ({ ...prev, logs: sData.logs || [] }));
        }
      } else if (tab === 'email') {
        const res = await fetch(apiUrl('/api/v1/admin/emails'), { headers });
        if (res.ok) setEmails(await res.json());
      }
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [tab]);

  useEffect(() => { fetchData(); }, [fetchData]);

  const triggerCrawl = async () => {
    await fetch(apiUrl('/api/crawler/crawl'), { method: 'POST' });
    setTab('crawler');
    fetchData();
  };

  return (
    <div className="min-h-screen bg-[#FBFBFB] font-sans">
      <header className="bg-white border-b border-slate-200 pt-16 pb-8 px-10 sticky top-0 z-40">
        <div className="max-w-7xl mx-auto flex flex-col md:flex-row md:items-center justify-between gap-8">
           <div>
              <div className="flex items-center gap-3">
                 <div className="p-2.5 bg-slate-900 rounded-2xl shadow-xl">
                    <Cpu className="w-6 h-6 text-indigo-400" />
                 </div>
                 <h1 className="text-3xl font-black text-slate-900 tracking-tight">System Node.</h1>
              </div>
              <p className="text-slate-400 text-[10px] font-black uppercase tracking-[0.3em] mt-3 ml-1">Platform Orchestration Interface</p>
           </div>

           <div className="flex bg-slate-50 p-1.5 rounded-[1.5rem] border border-slate-100">
              {['overview', 'users', 'crawler', 'email'].map(t => (
                <button
                  key={t}
                  onClick={() => setTab(t)}
                  className={`px-6 py-2.5 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all ${
                    tab === t ? 'bg-white text-indigo-600 shadow-sm border border-slate-200' : 'text-slate-400 hover:text-slate-600'
                  }`}
                >
                  {t}
                </button>
              ))}
           </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto p-10">
        {loading ? (
          <div className="flex flex-col items-center justify-center py-40">
             <Loader2 className="w-12 h-12 text-slate-900 animate-spin mb-4" />
             <p className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Polling System Metrics...</p>
          </div>
        ) : (
          <>
            {tab === 'overview' && (
              <div className="space-y-12 animate-in fade-in duration-500">
                <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
                  {[
                    { label: 'Cloud Nodes', val: (stats?.total_candidates || 0) + (stats?.total_companies || 0), icon: Users, color: 'blue' },
                    { label: 'Active Jobs', val: stats?.total_jobs || 0, icon: Briefcase, color: 'indigo' },
                    { label: 'Deployments', val: stats?.total_placements || 0, icon: Zap, color: 'emerald' },
                    { label: 'Partners', val: stats?.total_companies || 0, icon: Globe, color: 'amber' },
                  ].map(s => (
                    <div key={s.label} className="bg-white border border-slate-200 p-8 rounded-[2.5rem] shadow-sm relative overflow-hidden group">
                       <div className={`absolute top-0 right-0 w-24 h-24 bg-${s.color}-500/5 rounded-full -mr-12 -mt-12 blur-2xl transition-transform group-hover:scale-150`}></div>
                       <s.icon className={`w-5 h-5 text-${s.color}-500 mb-6`} />
                       <p className="text-4xl font-black text-slate-900 tracking-tighter">{s.val}</p>
                       <p className="text-[10px] font-black text-slate-400 uppercase tracking-widest mt-2">{s.label}</p>
                    </div>
                  ))}
                </div>

                <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
                   <div className="lg:col-span-8 bg-white border border-slate-200 rounded-[3rem] p-10 shadow-sm">
                      <div className="flex items-center justify-between mb-10">
                         <h3 className="text-xl font-black text-slate-900 tracking-tight flex items-center gap-3">
                            <Activity className="w-5 h-5 text-indigo-600" />
                            Throughput Analytics
                         </h3>
                         <div className="px-3 py-1 bg-slate-50 text-slate-400 rounded-lg text-[9px] font-black uppercase">Live: 7D Window</div>
                      </div>
                      <div className="h-80 w-full">
                        <ResponsiveContainer width="100%" height="100%">
                          <AreaChart data={stats?.chart_data || []}>
                            <defs>
                              <linearGradient id="colorApps" x1="0" y1="0" x2="0" y2="1">
                                <stop offset="5%" stopColor="#6366f1" stopOpacity={0.1}/>
                                <stop offset="95%" stopColor="#6366f1" stopOpacity={0}/>
                              </linearGradient>
                            </defs>
                            <XAxis dataKey="day" axisLine={false} tickLine={false} tick={{fontSize: 10, fontWeight: 700, fill: '#94a3b8'}} dy={10} />
                            <YAxis hide />
                            <Tooltip contentStyle={{borderRadius: '16px', border: 'none', boxShadow: '0 20px 25px -5px rgba(0,0,0,0.1)', fontWeight: 'bold'}} />
                            <Area type="monotone" dataKey="applications" stroke="#6366f1" strokeWidth={4} fillOpacity={1} fill="url(#colorApps)" />
                          </AreaChart>
                        </ResponsiveContainer>
                      </div>
                   </div>

                   <div className="lg:col-span-4 bg-slate-900 rounded-[3rem] p-10 text-white relative overflow-hidden flex flex-col justify-between">
                      <div className="absolute top-0 right-0 w-64 h-64 bg-indigo-500/10 rounded-full -mr-32 -mt-32 blur-[100px]"></div>
                      <div className="relative z-10">
                         <h3 className="text-xl font-black tracking-tight">System Status</h3>
                         <p className="text-slate-500 text-xs font-bold uppercase tracking-widest mt-2">All Nodes Operational</p>
                      </div>

                      <div className="relative z-10 space-y-6">
                         <div className="flex items-center justify-between p-4 bg-white/5 rounded-2xl border border-white/5">
                            <span className="text-[10px] font-black uppercase text-slate-400">Database Latency</span>
                            <span className="text-xs font-black text-emerald-400">12ms</span>
                         </div>
                         <div className="flex items-center justify-between p-4 bg-white/5 rounded-2xl border border-white/5">
                            <span className="text-[10px] font-black uppercase text-slate-400">AI Cache Hit</span>
                            <span className="text-xs font-black text-blue-400">94.2%</span>
                         </div>
                      </div>
                   </div>
                </div>
              </div>
            )}

            {tab === 'users' && (
              <div className="bg-white border border-slate-200 rounded-[2.5rem] overflow-hidden shadow-sm animate-in slide-in-from-bottom-4 duration-500">
                <div className="p-8 border-b border-slate-100 flex items-center justify-between bg-slate-50/30">
                   <h2 className="text-2xl font-black text-slate-900 tracking-tight">Node Directory</h2>
                   <div className="relative">
                      <Search className="absolute left-3 top-2.5 w-4 h-4 text-slate-400" />
                      <input type="text" placeholder="Search user hash..." className="pl-10 pr-4 py-2 bg-white border border-slate-200 rounded-xl text-xs font-bold outline-none focus:ring-2 focus:ring-indigo-500/20 w-64" />
                   </div>
                </div>
                <div className="overflow-x-auto">
                  <table className="w-full text-left">
                    <thead>
                      <tr className="text-[10px] font-black text-slate-400 uppercase tracking-widest border-b border-slate-50">
                        <th className="px-8 py-6">Identity</th>
                        <th className="px-8 py-6">Privileges</th>
                        <th className="px-8 py-6">First Seen</th>
                        <th className="px-8 py-6 text-right">Ops</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-50 text-sm font-bold text-slate-600">
                      {users.length === 0 ? (
                        <tr><td colSpan="4" className="py-20 text-center uppercase tracking-widest text-[10px] text-slate-400">Registry is empty</td></tr>
                      ) : (
                        users.map(u => (
                          <tr key={u.id} className="hover:bg-slate-50/50 transition-colors">
                            <td className="px-8 py-6 text-slate-900 font-black">{u.email}</td>
                            <td className="px-8 py-6">
                              <span className="px-3 py-1 bg-slate-100 border border-slate-200 rounded-lg text-[10px] uppercase tracking-tighter">{u.role}</span>
                            </td>
                            <td className="px-8 py-6 text-slate-400 font-mono text-xs">{new Date(u.created_at).toLocaleDateString()}</td>
                            <td className="px-8 py-6 text-right space-x-2">
                              <button className="px-4 py-2 bg-rose-50 text-rose-600 rounded-xl text-[10px] uppercase tracking-tighter hover:bg-rose-600 hover:text-white transition-all">Revoke</button>
                              <button className="px-4 py-2 bg-emerald-50 text-emerald-600 rounded-xl text-[10px] uppercase tracking-tighter hover:bg-emerald-600 hover:text-white transition-all">Validate</button>
                            </td>
                          </tr>
                        ))
                      )}
                    </tbody>
                  </table>
                </div>
              </div>
            )}

            {tab === 'crawler' && (
              <div className="space-y-8 animate-in zoom-in-95 duration-500">
                <div className="bg-white border border-slate-200 p-10 rounded-[3rem] shadow-sm flex items-center justify-between">
                  <div className="flex items-center gap-6">
                    <div className={`w-4 h-4 rounded-full ${crawlerStats.health === 'RUNNING' ? 'bg-emerald-500 animate-pulse' : 'bg-slate-300 shadow-inner shadow-slate-400'}`}></div>
                    <div>
                       <h3 className="text-2xl font-black text-slate-900 tracking-tight">Aggregator Core</h3>
                       <p className="text-[10px] font-black text-slate-400 uppercase tracking-widest mt-1">Status: {crawlerStats.health}</p>
                    </div>
                  </div>
                  <button onClick={triggerCrawl} className="flex items-center gap-3 bg-slate-900 text-white px-8 py-4 rounded-2xl font-black text-xs uppercase tracking-widest hover:bg-indigo-600 transition-all shadow-xl shadow-slate-900/10">
                    <Play className="w-4 h-4 fill-current" /> Initialize Sequence
                  </button>
                </div>

                <div className="bg-slate-950 rounded-[3rem] p-10 shadow-2xl relative overflow-hidden">
                   <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-transparent via-indigo-500 to-transparent"></div>
                   <div className="flex items-center justify-between mb-6">
                      <h4 className="text-indigo-400 text-xs font-black uppercase tracking-[0.2em]">Kernel Logs</h4>
                      <button className="text-[9px] font-black text-slate-600 uppercase hover:text-white transition-colors">Clear Buffer</button>
                   </div>
                   <div className="font-mono text-[11px] text-slate-300 h-96 overflow-y-auto space-y-2 custom-scrollbar">
                     {crawlerStats.logs.length === 0 ? (
                       <div className="text-slate-700 italic">// Awaiting input...</div>
                     ) : (
                       crawlerStats.logs.map((log, index) => <div key={index} className="flex gap-4 border-l border-white/5 pl-4"><span className="text-slate-700">{index+1}</span> {log}</div>)
                     )}
                   </div>
                </div>
              </div>
            )}

            {tab === 'email' && (
              <div className="bg-white border border-slate-200 rounded-[2.5rem] overflow-hidden shadow-sm animate-in slide-in-from-right-4 duration-500">
                <div className="p-8 border-b border-slate-100 bg-slate-50/30">
                   <h2 className="text-2xl font-black text-slate-900 tracking-tight">Signal Queue</h2>
                </div>
                <div className="overflow-x-auto">
                  <table className="w-full text-left">
                    <thead>
                      <tr className="text-[10px] font-black text-slate-400 uppercase tracking-widest border-b border-slate-50">
                        <th className="px-8 py-6">Target</th>
                        <th className="px-8 py-6">Signal Subject</th>
                        <th className="px-8 py-6">State</th>
                        <th className="px-8 py-6 text-right">Timestamp</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-50 text-sm font-bold text-slate-600">
                      {emails.length === 0 ? (
                         <tr><td colSpan="4" className="py-20 text-center uppercase tracking-widest text-[10px] text-slate-400">No outgoing signals in buffer</td></tr>
                      ) : (
                        emails.map(e => (
                          <tr key={e.id} className="hover:bg-slate-50/50 transition-colors">
                            <td className="px-8 py-6 text-slate-900">{e.to_email}</td>
                            <td className="px-8 py-6">{e.subject}</td>
                            <td className="px-8 py-6">
                              <div className={`inline-flex items-center gap-2 px-3 py-1 rounded-full text-[9px] font-black uppercase tracking-tighter ${
                                e.status === 'SENT' ? 'bg-emerald-50 text-emerald-600 border border-emerald-100' :
                                e.status === 'PENDING' ? 'bg-amber-50 text-amber-600 border border-amber-100' :
                                'bg-rose-50 text-rose-600 border border-rose-100'
                              }`}>
                                <div className={`w-1.5 h-1.5 rounded-full ${e.status === 'SENT' ? 'bg-emerald-500' : e.status === 'PENDING' ? 'bg-amber-500' : 'bg-rose-500'}`}></div>
                                {e.status}
                              </div>
                            </td>
                            <td className="px-8 py-6 text-right text-slate-400 font-mono text-xs">{new Date(e.created_at).toLocaleString()}</td>
                          </tr>
                        ))
                      )}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
          </>
        )}
      </main>
    </div>
  );
}
