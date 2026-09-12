import React, { createContext, useContext, useState, useEffect } from 'react';
import { authAPI } from '../lib/api';
import { toast } from 'react-hot-toast';

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
      const { success, data, message } = response.data;

      if (success && data) {
        setUser(data.user);
        toast.success("Welcome back!");
        return data.user;
      }
      throw new Error(message || "Login failed");
    } catch (err) {
      const errMsg = err.response?.data?.message || err.message || "An unexpected error occurred during login";
      console.error("Login Error:", errMsg);
      toast.error(errMsg);
      throw new Error(errMsg);
    }
  };

  const register = async (registerData) => {
    try {
      const response = await authAPI.register(registerData);
      const { success, message, error: apiError } = response.data;

      if (success) {
        toast.success("Account created successfully!");
        return login(registerData.email, registerData.password);
      }
      throw new Error(apiError || message || "Registration failed");
    } catch (err) {
      const errMsg = err.response?.data?.error || err.response?.data?.message || err.message || "Registration failed";
      console.error("Registration Error:", errMsg);
      toast.error(errMsg);
      throw new Error(errMsg);
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
        localStorage.removeItem('rojgar_user');
        window.location.href = '/login';
      }
    }
  };

  return (
    <AuthContext.Provider value={{ user, login, register, logout, isAuthenticated: !!user, initialized }}>
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
