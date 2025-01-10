import React from 'react';
import { useAuth } from '../components/Auth';
import Stream from '../components/Stream';
const Dashboard = () => {
    const { logout, user } = useAuth();

     return (
        <div>
            <h2>Dashboard</h2>
            {user && <p>Bem-vindo, {user.name}</p>}
             <Stream />
           {/* biome-ignore lint/a11y/useButtonType: <explanation> */}
           <button onClick={logout}>Logout</button>
        </div>
    );
};

export default Dashboard;