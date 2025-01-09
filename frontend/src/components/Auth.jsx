import React from 'react';
import { useNavigate } from 'react-router-dom';
const AuthContext = React.createContext(null);
const AuthProvider = ({children}) => {
    const navigate = useNavigate();
  const [user, setUser] = React.useState(null);
  const [token, setToken] = React.useState(localStorage.getItem('token'));
      
  const login = async (token, user) => {
     setToken(token);
     setUser(user);
       localStorage.setItem('token', token);
       navigate('/dashboard');
      }
      
  const logout = () => {
        setToken(null);
        setUser(null);
        localStorage.removeItem('token')
        navigate('/login');
   }

  const value = {
      user,
      token,
      login,
      logout
  }
    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

const useAuth = () => {
  return React.useContext(AuthContext)
}
export { AuthProvider, useAuth };
