import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import { useAuth } from '../context/AuthContext';

const ProtectedRoute = ({ children, allowedRoles }) => {
  const router = useRouter();
  const { user, loading } = useAuth();
  const [isAuthorized, setIsAuthorized] = useState(false);

  useEffect(() => {
    if (!loading) {
      const role = user?.role;

      if (!user) {
        router.replace('/login');
      } else if (allowedRoles && !allowedRoles.includes(role)) {
        router.replace('/unauthorized');
      } else {
        setIsAuthorized(true);
      }
    }
  }, [loading, user, allowedRoles, router]);

  if (loading || !isAuthorized) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-slate-50">
        <div className="animate-pulse flex flex-col items-center">
          <div className="w-12 h-12 bg-blue-600 rounded-full mb-4"></div>
          <p className="text-slate-500 font-medium">Verifying access...</p>
        </div>
      </div>
    );
  }

  return children;
};

export default ProtectedRoute;