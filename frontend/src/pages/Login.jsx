import axios from 'axios';
import React, { useState } from 'react';
import { useAuth } from '../components/auth';

const Login = () => {
    const {login} = useAuth();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

 const handleSubmit = async (e) => {
     e.preventDefault();
    try {
      const response = await axios.post('http://localhost:8080/login', {
          email: email,
          password: password,
        });
           login(response.data.token, response.data.user)
    } catch (error) {
      console.error("Erro ao fazer login:", error);
        }
  }

    return (
        <div>
           <h2>Login</h2>
             <form onSubmit={handleSubmit}>
                <div>
                    <label htmlFor="email">Email:</label>
                    <input id="email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
                </div>
                 <div>
                    <label htmlFor="password">Password:</label>
                    <input id="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} required />
               </div>
               <button type="submit">Login</button>
            </form>
       </div>
     );
};

 export default Login;