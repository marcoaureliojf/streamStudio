import React from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { AuthProvider } from './components/auth';
import Dashboard from './pages/Dashboard';
import Login from './pages/Login';
import './styles/styles.css';
function App() {
  return (
  <AuthProvider>
  <BrowserRouter>
  <Routes>
  <Route path="/login" element={<Login />} />
  <Route path="/dashboard" element={<Dashboard />} />
  <Route path="/" element={<Login />} />
  </Routes>
  </BrowserRouter>
  </AuthProvider>
  );
  }
  
  export default App;