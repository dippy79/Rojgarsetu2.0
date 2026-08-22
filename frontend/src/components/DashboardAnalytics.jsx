import React from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Cell, PieChart, Pie } from 'recharts';
import { TrendingUp, Award, Clock, Target } from 'lucide-react';

const DashboardAnalytics = ({ applicationStats = [], skillMatch = 75 }) => {
  // Sample data for charts if none provided
  const chartData = applicationStats.length > 0 ? applicationStats : [
    { name: 'Mon', apps: 4 },
    { name: 'Tue', apps: 7 },
    { name: 'Wed', apps: 5 },
    { name: 'Thu', apps: 9 },
    { name: 'Fri', apps: 12 },
    { name: 'Sat', apps: 3 },
    { name: 'Sun', apps: 2 },
  ];

  const pieData = [
    { name: 'Shortlisted', value: 4, color: '#10B981' },
    { name: 'Pending', value: 8, color: '#3B82F6' },
    { name: 'Rejected', value: 2, color: '#EF4444' },
  ];

  return (
    <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
      {/* Activity Chart */}
      <div className="lg:col-span-8 bg-white border border-slate-200/60 rounded-[2rem] p-8 shadow-sm overflow-hidden relative">
        <div className="flex items-center justify-between mb-8">
          <div>
            <h3 className="text-xl font-black text-slate-900 tracking-tight">Application Velocity</h3>
            <p className="text-slate-400 text-xs font-bold uppercase tracking-widest mt-1">Activity past 7 days</p>
          </div>
          <div className="p-3 bg-emerald-50 text-emerald-600 rounded-2xl border border-emerald-100 flex items-center gap-2">
            <TrendingUp className="w-4 h-4" />
            <span className="text-xs font-black">+24% Increase</span>
          </div>
        </div>

        <div className="h-64 w-full">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#f1f5f9" />
              <XAxis
                dataKey="name"
                axisLine={false}
                tickLine={false}
                tick={{fill: '#94a3b8', fontSize: 10, fontWeight: 700}}
                dy={10}
              />
              <YAxis hide />
              <Tooltip
                cursor={{fill: '#f8fafc'}}
                contentStyle={{borderRadius: '16px', border: 'none', boxShadow: '0 10px 15px -3px rgba(0,0,0,0.1)', fontWeight: 'bold'}}
              />
              <Bar dataKey="apps" radius={[6, 6, 0, 0]}>
                {chartData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={index === 4 ? '#3b82f6' : '#cbd5e1'} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Profile Score Widget */}
      <div className="lg:col-span-4 space-y-6">
        <div className="bg-slate-900 rounded-[2rem] p-8 text-white relative overflow-hidden group">
          <div className="absolute top-0 right-0 w-32 h-32 bg-blue-500/20 rounded-full -mr-16 -mt-16 blur-3xl group-hover:bg-blue-500/40 transition-all duration-700"></div>

          <div className="relative z-10 space-y-6">
            <div className="flex items-center gap-3">
              <div className="p-3 bg-white/10 backdrop-blur-md rounded-2xl border border-white/10">
                <Award className="w-6 h-6 text-blue-400" />
              </div>
              <h3 className="font-black text-lg">AI Match Score</h3>
            </div>

            <div className="flex items-end gap-2">
              <span className="text-6xl font-black tracking-tighter text-blue-400">{skillMatch}%</span>
              <span className="text-blue-100/50 font-bold mb-2">Overall Fit</span>
            </div>

            <div className="w-full bg-white/10 h-3 rounded-full overflow-hidden">
              <div
                className="h-full bg-gradient-to-r from-blue-600 to-indigo-400 rounded-full transition-all duration-1000 ease-out shadow-[0_0_20px_rgba(59,130,246,0.5)]"
                style={{ width: `${skillMatch}%` }}
              ></div>
            </div>

            <p className="text-xs font-medium text-slate-400 leading-relaxed">
              Your profile is stronger than <span className="text-white font-bold">82%</span> of other candidates in the Tech industry.
            </p>
          </div>
        </div>

        {/* Small Status breakdown */}
        <div className="bg-white border border-slate-200/60 rounded-[2rem] p-6 shadow-sm">
          <div className="flex items-center justify-between mb-4">
            <h4 className="text-sm font-black text-slate-900 uppercase tracking-tighter">Status Mix</h4>
            <Clock className="w-4 h-4 text-slate-400" />
          </div>
          <div className="space-y-3">
            {pieData.map((item) => (
              <div key={item.name} className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <div className="w-2 h-2 rounded-full" style={{ backgroundColor: item.color }}></div>
                  <span className="text-xs font-bold text-slate-500">{item.name}</span>
                </div>
                <span className="text-xs font-black text-slate-900">{item.value}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default DashboardAnalytics;
