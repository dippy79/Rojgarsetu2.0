import React, { createContext, useContext, useState, useEffect } from 'react';
import { authAPI } from '../lib/api';

export const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [initialized, setInitialized] = useState(false);

  useEffect(() => {
    const initAuth = async () => {
      try {
        // Attempt to fetch profile to see if session exists (via HttpOnly cookie)
        const response = await authAPI.getProfile();
        if (response.data && response.data.success) {
          setUser(response.data.data);
        }
      } catch (err) {
        console.log("No active session found");
      } finally {
        setInitialized(true);
      }
    };

    initAuth();
  }, []);

  const login = async (email, password) => {
    try {
      const response = await authAPI.login({ email, password });
      if (response.data && response.data.data) {
        setUser(response.data.data.user);
        return response.data.data.user;
      }
      throw new Error("Login failed");
    } catch (err) {
      console.error("Login Error:", err);
      throw err;
    }
  };

  const logout = async () => {
    try {
      await authAPI.logout();
    } catch (err) {
      console.error("Logout failed:", err);
    } finally {
      setUser(null);
      if (typeof window !== 'undefined') {
        window.location.href = '/login';
      }
    }
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
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

export default AuthContext;
