import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';

// Auth Pages (Using YOUR Original Files)
import LoginPage from '../pages/auth/LoginPage';
import RegisterPage from '../pages/auth/RegisterPage';

// Candidate Pages
import JobSearch from '../pages/candidate/JobSearch';
import CandidateApplications from '../pages/candidate/CandidateApplications';
import CandidateProfile from '../pages/candidate/CandidateProfile';

// Employer Pages
import EmployerDashboard from '../pages/employer/EmployerDashboard';
import PostJob from '../pages/employer/PostJob';
import EmployerApplicants from '../pages/employer/EmployerApplicants';

export const AppRoutes = () => {
  return (
    <Routes>
      {/* Auth Routes */}
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />

      {/* Candidate Routes */}
      <Route path="/jobs" element={<JobSearch />} />
      <Route path="/candidate/applications" element={<CandidateApplications />} />
      <Route path="/candidate/profile" element={<CandidateProfile />} />

      {/* Employer Routes */}
      <Route path="/employer/dashboard" element={<EmployerDashboard />} />
      <Route path="/employer/post-job" element={<PostJob />} />
      <Route path="/employer/applicants" element={<EmployerApplicants />} />

      {/* Default Fallback */}
      <Route path="/" element={<Navigate to="/jobs" replace />} />
      <Route path="*" element={<Navigate to="/jobs" replace />} />
    </Routes>
  );
};

export default AppRoutes;