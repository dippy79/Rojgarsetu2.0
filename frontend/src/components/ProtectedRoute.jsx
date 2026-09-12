import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import { useAuth } from '../context/AuthContext';

const ProtectedRoute = ({ children, allowedRoles }) => {
  const router = useRouter();
  const { user, initialized } = useAuth();
  const [isAuthorized, setIsAuthorized] = useState(false);

  useEffect(() => {
    if (initialized) {
      const role = user?.role;

      if (!user) {
        router.replace('/login');
      } else if (allowedRoles && !allowedRoles.includes(role)) {
        router.replace('/unauthorized');
      } else {
        setIsAuthorized(true);
      }
    }
  }, [initialized, user, allowedRoles, router]);

  if (!initialized || !isAuthorized) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-stone-50">
        <div className="flex flex-col items-center">
          <div className="w-12 h-12 border-4 border-stone-200 border-t-stone-800 rounded-full animate-spin mb-4"></div>
          <p className="text-stone-500 font-medium">Verifying access...</p>
        </div>
      </div>
    );
  }

  return children;
};

export default ProtectedRoute;
