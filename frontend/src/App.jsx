import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';

// Pages & Components
import CompanyDashboard from './pages/company/CompanyDashboard';
import JobApplicants from './pages/company/JobApplicants';
import PostJob from './pages/company/PostJob';
import AdminDashboard from './pages/admin/AdminDashboard';
import AIJobMatches from './pages/candidate/AIJobMatches';
import VideoCallPage from './pages/VideoCallPage';
import NotificationBell from './components/NotificationBell';

const ProtectedRoute = ({ allowedRoles, children }) => {
  const token = localStorage.getItem('rojgar_token');
  if (!token) return <Navigate to="/login" replace />;
  return children;
};

export default function App() {
  return (
    <BrowserRouter>
      <div className="min-h-screen bg-slate-50 text-slate-900 font-sans">
        <header className="bg-white border-b border-slate-200 px-8 py-4 flex justify-between items-center">
          <h1 className="font-bold text-xl text-indigo-600">ROJGARSETU</h1>
          <NotificationBell />
        </header>

        <Routes>
          {/* Company Routes */}
          <Route path="/dashboard/company" element={<ProtectedRoute allowedRoles={['company', 'admin']}><CompanyDashboard /></ProtectedRoute>} />
          <Route path="/company/applicants" element={<ProtectedRoute allowedRoles={['company', 'admin']}><JobApplicants /></ProtectedRoute>} />
          <Route path="/company/post-job" element={<ProtectedRoute allowedRoles={['company', 'admin']}><PostJob /></ProtectedRoute>} />

          {/* Admin Routes */}
          <Route path="/dashboard/admin" element={<ProtectedRoute allowedRoles={['admin']}><AdminDashboard /></ProtectedRoute>} />

          {/* Candidate AI Matches & Video Calls */}
          <Route path="/candidate/ai-matches" element={<ProtectedRoute allowedRoles={['candidate']}><AIJobMatches /></ProtectedRoute>} />
          <Route path="/interview/:id" element={<ProtectedRoute allowedRoles={['candidate', 'company', 'admin']}><VideoCallPage /></ProtectedRoute>} />

          {/* Default Fallback */}
          <Route path="*" element={<Navigate to="/dashboard/company" replace />} />
        </Routes>
      </div>
    </BrowserRouter>
  );
}
