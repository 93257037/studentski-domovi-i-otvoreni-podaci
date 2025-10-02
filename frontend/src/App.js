import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import Login from './components/Login';
import Register from './components/Register';
import Dashboard from './components/Dashboard';
import StDomDetail from './components/StDomDetail';
import RoomDetail from './components/RoomDetail';
import AdvancedRoomSearch from './components/AdvancedRoomSearch';
import AcademicYearApplications from './components/AcademicYearApplications';
import ProtectedRoute from './components/ProtectedRoute';
import './App.css';

function AppContent() {
  const { isAuthenticated, loading } = useAuth();

  if (loading) {
    return (
      <div className="loading-container">
        <div className="loading-spinner"></div>
        <p>Učitavanje...</p>
      </div>
    );
  }

  return (
    <Router>
      <div className="App">
        <Routes>
          <Route 
            path="/login" 
            element={isAuthenticated ? <Navigate to="/dashboard" replace /> : <Login />} 
          />
          <Route 
            path="/register" 
            element={isAuthenticated ? <Navigate to="/dashboard" replace /> : <Register />} 
          />
          <Route 
            path="/dashboard" 
            element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/st-dom/:id" 
            element={
              <ProtectedRoute>
                <StDomDetail />
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/room/:id" 
            element={
              <ProtectedRoute>
                <RoomDetail />
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/advanced-search" 
            element={
              <ProtectedRoute>
                <AdvancedRoomSearch />
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/academic-year-applications" 
            element={
              <ProtectedRoute>
                <AcademicYearApplications />
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/" 
            element={<Navigate to={isAuthenticated ? "/dashboard" : "/login"} replace />} 
          />
        </Routes>
      </div>
    </Router>
  );
}

function App() {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider>
  );
}

export default App;
