import React from 'react';
import Stream from '../components/Stream';
import { useAuth } from '../components/auth';
const Dashboard = () => {
    const { logout, user } = useAuth();

     return (
        <div>
            <h2>Dashboard</h2>
            {user && <p>Bem-vindo, {user.name}</p>}
             <Stream />
           <button type="button" onClick={logout}>Logout</button>
        </div>
    );
};

export default Dashboard;