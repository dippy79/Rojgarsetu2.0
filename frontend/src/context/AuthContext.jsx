import React, { createContext, useContext, useState } from 'react';

export const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [initialized, setInitialized] = useState(false);

  React.useEffect(() => {
    if (typeof window !== 'undefined') {
      const savedName = localStorage.getItem('userName') || localStorage.getItem('name');
      const savedRole = localStorage.getItem('userRole') || localStorage.getItem('role') || 'candidate';
      const savedEmail = localStorage.getItem('userEmail') || '';
      const token = localStorage.getItem('token');

      if (token) {
        setUser({ name: savedName || 'User', role: savedRole, email: savedEmail });
      }
      setInitialized(true);
    }
  }, []);

  const login = async (email, password, role = 'candidate') => {
    const mockUser = { name: email.split('@')[0] || 'User', email, role };
    if (typeof window !== 'undefined') {
      localStorage.setItem('token', 'mock-jwt-token');
      localStorage.setItem('userName', mockUser.name);
      localStorage.setItem('userRole', role);
      localStorage.setItem('userEmail', email);
    }
    setUser(mockUser);
    return mockUser;
  };

  const logout = () => {
    if (typeof window !== 'undefined') {
      localStorage.clear();
    }
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, login, logout, isAuthenticated: !!user, initialized }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    // Fallback for cases outside AuthProvider, but still SSR-safe
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
    const name = typeof window !== 'undefined' ? (localStorage.getItem('userName') || localStorage.getItem('name') || 'Simranjeet Singh') : 'Simranjeet Singh';
    const role = typeof window !== 'undefined' ? (localStorage.getItem('userRole') || localStorage.getItem('role') || 'candidate') : 'candidate';
    
    return {
      user: token ? { name, role } : null,
      isAuthenticated: !!token,
      login: async () => {},
      logout: () => { if (typeof window !== 'undefined') localStorage.clear(); }
    };
  }
  return context;
};

export default AuthContext;