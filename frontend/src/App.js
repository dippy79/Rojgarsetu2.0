import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link, useLocation, Navigate } from 'react-router-dom';
import Navbar from './components/Navbar';
import GovJobsPage from './pages/gov-jobs/GovJobsPage';
import PrivateJobsPage from './pages/private-jobs/PrivateJobsPage';
import CoursesPage from './pages/courses/CoursesPage';
import VideosPage from './pages/videos/VideosPage';
import GovtFormsDashboard from './components/GovtFormsDashboard';
import LoginPage from './pages/auth/LoginPage';
import ProtectedRoute from './components/ProtectedRoute';
import { AuthProvider } from './context/AuthContext';
import './index.css';

const HomePage = () => (
  <div style={{ textAlign: 'center', padding: '4rem 2rem', maxWidth: 800, margin: '0 auto' }}>
    <h1 style={{ fontSize: '2.5rem', color: '#1a237e', marginBottom: '1rem' }}>RojgarSetu 2.0</h1>
    <p style={{ fontSize: '1.2rem', color: '#555', marginBottom: '2rem' }}>
      Your comprehensive job portal for government and private sector opportunities, courses, and educational videos.
    </p>
    <div style={{ display: 'flex', gap: '1rem', justifyContent: 'center', flexWrap: 'wrap' }}>
      <a href="/gov-jobs" style={{
        padding: '0.75rem 2rem',
        background: 'linear-gradient(135deg, #1a237e, #283593)',
        color: '#fff',
        borderRadius: 8,
        fontWeight: 600,
        fontSize: '1rem',
        textDecoration: 'none'
      }}>Browse Government Jobs</a>
      <a href="/private-jobs" style={{
        padding: '0.75rem 2rem',
        background: 'linear-gradient(135deg, #00897b, #26a69a)',
        color: '#fff',
        borderRadius: 8,
        fontWeight: 600,
        fontSize: '1rem',
        textDecoration: 'none'
      }}>Browse Private Jobs</a>
      <a href="/courses" style={{
        padding: '0.75rem 2rem',
        background: 'linear-gradient(135deg, #6a1b9a, #8e24aa)',
        color: '#fff',
        borderRadius: 8,
        fontWeight: 600,
        fontSize: '1rem',
        textDecoration: 'none'
      }}>Browse Courses</a>
      <a href="/videos" style={{
        padding: '0.75rem 2rem',
        background: 'linear-gradient(135deg, #c62828, #e53935)',
        color: '#fff',
        borderRadius: 8,
        fontWeight: 600,
        fontSize: '1rem',
        textDecoration: 'none'
      }}>Browse Videos</a>
    </div>
  </div>
);

const UnauthorizedPage = () => (
  <div className="flex flex-col items-center justify-center min-h-screen bg-slate-50 text-slate-900 p-6 text-center">
    <div className="w-16 h-16 bg-red-100 text-red-600 rounded-full flex items-center justify-center mb-4 text-2xl font-bold">
      !
    </div>
    <h1 className="text-3xl font-extrabold mb-2">403 - Access Denied</h1>
    <p className="text-slate-600 max-w-md mb-6 text-sm">
      You do not have the required permissions to access this page. Please log in with an authorized account.
    </p>
    <Link
      to="/login"
      className="px-6 py-3 bg-slate-900 text-white font-semibold rounded-xl hover:bg-slate-800 transition-all text-sm shadow-sm"
    >
      Return to Login
    </Link>
  </div>
);

const NotFound = () => (
  <div style={{ textAlign: 'center', padding: '4rem 2rem' }}>
    <h1 style={{ fontSize: '3rem', color: '#e53935' }}>404</h1>
    <p style={{ fontSize: '1.2rem', color: '#666' }}>Page not found</p>
    <a href="/" style={{ color: '#1a237e', fontWeight: 600, marginTop: '1rem', display: 'inline-block' }}>Go Home</a>
  </div>
);

function AppContent() {
  const location = useLocation();
  const hideNavbar = ['/login'].includes(location.pathname);

  return (
    <div className="App">
      {!hideNavbar && <Navbar />}
      <main style={{ flex: 1 }}>
        <Routes>
          {/* Public Routes */}
          <Route path="/" element={<HomePage />} />
          <Route path="/gov-jobs" element={<GovJobsPage />} />
          <Route path="/private-jobs" element={<PrivateJobsPage />} />
          <Route path="/courses" element={<CoursesPage />} />
          <Route path="/videos" element={<VideosPage />} />
          <Route path="/govt-forms" element={<GovtFormsDashboard />} />

          {/* Auth Routes */}
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<Navigate to="/login" replace />} />
          <Route path="/unauthorized" element={<UnauthorizedPage />} />

          {/* Protected Dashboards */}
          <Route
            path="/dashboard/candidate"
            element={
              <ProtectedRoute allowedRoles={['candidate']}>
                <div className="p-8 text-2xl font-bold">Candidate Dashboard</div>
              </ProtectedRoute>
            }
          />
          <Route
            path="/dashboard/company"
            element={
              <ProtectedRoute allowedRoles={['company']}>
                <div className="p-8 text-2xl font-bold">Employer Dashboard</div>
              </ProtectedRoute>
            }
          />
          <Route
            path="/dashboard/admin"
            element={
              <ProtectedRoute allowedRoles={['admin']}>
                <div className="p-8 text-2xl font-bold">Admin Dashboard</div>
              </ProtectedRoute>
            }
          />

          <Route path="*" element={<NotFound />} />
        </Routes>
      </main>
    </div>
  );
}

function App() {
  return (
    <AuthProvider>
      <Router>
        <AppContent />
      </Router>
    </AuthProvider>
  );
}

export default App;