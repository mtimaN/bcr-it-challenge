import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './Auth.css';

import roD from '../../assets/nav_bar/roDark.png';
import enD from '../../assets/nav_bar/enDark.png';
import roL from '../../assets/nav_bar/roLight.png';
import enL from '../../assets/nav_bar/enLight.png';
import sun from '../../assets/nav_bar/dayLogo.png';
import moon from '../../assets/nav_bar/nightLogo.png';

const Login = ({ lang, setLang, setLoggedIn, setUserData }) => {
  const navigate = useNavigate();
  const [theme, setTheme] = useState(localStorage.getItem('current_theme') || 'light');

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  useEffect(() => {
    // Force theme to light
    setTheme('light');
    localStorage.setItem('current_theme', 'light');
  
    // Force language to RO
    setLang('RO');
  
    // Remove dark class from container if it exists
    const container = document.querySelector('.container');
    container?.classList.remove('dark');
  }, []);

  const handleLogin = async () => {
    try {
      const response = await fetch('https://localhost:8443/v1/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        credentials: 'include',
        body: JSON.stringify({
          username: username,
          password: password
        })
      });

      if (response.ok) {
        const data = await response.json();
        const token = data.token;

        localStorage.setItem('jwtToken', token);

        const localUserData = JSON.parse(localStorage.getItem('userData') || '{}');

        setUserData({
          firstName: localUserData.firstName || '',
          lastName: localUserData.lastName || '',
          email: localUserData.email || '',
          username,
          password
        });

        setLoggedIn(true);
        navigate('/home');

      } else {
        const errorData = await response.json();
        alert('Login failed: ' + (errorData.error || 'Invalid credentials'));
      }
    } catch (error) {
      alert('Error: ' + error.message);
    }
  };


  const getLangIcon = () => {
    if (lang === 'RO') {
      return theme === 'light' ? roD : roL;
    } else {
      return theme === 'light' ? enD : enL;
    }
  };


  return (
    <div className="auth-container">
      {/* Lang toggle */}

      <p className="login-title">GEORG.IO</p>

      <p className="login-description">Next-gen Personalized Banking Experience</p>

      {/* Auth box */}
      <div className="auth-box">
        <input
          type="text"
          placeholder={lang === 'RO' ? 'Utilizator' : 'Username'}
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <input
          type="password"
          placeholder={lang === 'RO' ? 'Parolă' : 'Password'}
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <button className="auth-button" onClick={handleLogin}>
          {lang === 'RO' ? 'Autentificare' : 'Login'}
        </button>
        <p className="switch-link" onClick={() => navigate('/register')}>
          {lang === 'RO' ? 'Nu ai cont? Înregistrează-te' : "Don't have an account? Register"}
        </p>
      </div>
    </div>
  );
};

export default Login;
