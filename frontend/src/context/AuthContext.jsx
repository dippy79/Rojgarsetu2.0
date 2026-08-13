import React, { createContext, useContext, useState } from 'react';

export const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(() => {
    const savedName = localStorage.getItem('userName') || localStorage.getItem('name');
    const savedRole = localStorage.getItem('userRole') || localStorage.getItem('role') || 'candidate';
    const savedEmail = localStorage.getItem('userEmail') || '';
    const token = localStorage.getItem('token');
    
    if (token) {
      return { name: savedName || 'User', role: savedRole, email: savedEmail };
    }
    return null;
  });

  const login = async (email, password, role = 'candidate') => {
    const mockUser = { name: email.split('@')[0] || 'User', email, role };
    localStorage.setItem('token', 'mock-jwt-token');
    localStorage.setItem('userName', mockUser.name);
    localStorage.setItem('userRole', role);
    localStorage.setItem('userEmail', email);
    setUser(mockUser);
    return mockUser;
  };

  const logout = () => {
    localStorage.clear();
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, login, logout, isAuthenticated: !!user }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    const token = localStorage.getItem('token');
    const name = localStorage.getItem('userName') || localStorage.getItem('name') || 'Simranjeet Singh';
    const role = localStorage.getItem('userRole') || localStorage.getItem('role') || 'candidate';
    
    return {
      user: token ? { name, role } : null,
      isAuthenticated: !!token,
      login: async () => {},
      logout: () => { localStorage.clear(); }
    };
  }
  return context;
};

export default AuthContext;