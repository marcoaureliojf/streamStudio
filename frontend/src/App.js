import React from 'react';
import { Route, Routes } from 'react-router-dom';
import { AuthProvider } from './components/Auth';
import Dashboard from './pages/Dashboard';
import Login from './pages/Login';
import './styles/styles.css';

function App() {
    return (
        <AuthProvider>
                <Routes>
                   <Route path="/login" element={<Login />} />
                    <Route path="/dashboard" element={<Dashboard />} />
                    <Route path="/" element={<Login />} />
               </Routes>
       </AuthProvider>
   );
}

export default App;