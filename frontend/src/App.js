import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Navbar from './components/Navbar';
import GovJobsPage from './pages/gov-jobs/GovJobsPage';
import PrivateJobsPage from './pages/private-jobs/PrivateJobsPage';
import CoursesPage from './pages/courses/CoursesPage';
import VideosPage from './pages/videos/VideosPage';
import GovtFormsDashboard from './components/GovtFormsDashboard';
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

const NotFound = () => (
  <div style={{ textAlign: 'center', padding: '4rem 2rem' }}>
    <h1 style={{ fontSize: '3rem', color: '#e53935' }}>404</h1>
    <p style={{ fontSize: '1.2rem', color: '#666' }}>Page not found</p>
    <a href="/" style={{ color: '#1a237e', fontWeight: 600, marginTop: '1rem', display: 'inline-block' }}>Go Home</a>
  </div>
);

function App() {
  return (
    <Router>
      <div className="App">
        <Navbar />
        <main style={{ flex: 1 }}>
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/gov-jobs" element={<GovJobsPage />} />
            <Route path="/private-jobs" element={<PrivateJobsPage />} />
            <Route path="/courses" element={<CoursesPage />} />
<Route path="/videos" element={<VideosPage />} />
            <Route path="/govt-forms" element={<GovtFormsDashboard />} />
            <Route path="*" element={<NotFound />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;
