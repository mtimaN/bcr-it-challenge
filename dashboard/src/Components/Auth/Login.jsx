import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './Auth.css';

import roD from '../../assets/nav_bar/roDark.png';
import enD from '../../assets/nav_bar/enDark.png';
import roL from '../../assets/nav_bar/roLight.png';
import enL from '../../assets/nav_bar/enLight.png';
import sun from '../../assets/nav_bar/dayLogo.png';
import moon from '../../assets/nav_bar/nightLogo.png';

const Login = ({ lang, setLang, setLoggedIn }) => {
  const navigate = useNavigate();
  const [theme, setTheme] = useState(localStorage.getItem('current_theme') || 'light');

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  useEffect(() => {
    const container = document.querySelector('.container');
    if (theme === 'dark') {
      container?.classList.add('dark');
    } else {
      container?.classList.remove('dark');
    }
    localStorage.setItem('current_theme', theme);
  }, [theme]);

  const toggleTheme = () => {
    const newTheme = theme === 'light' ? 'dark' : 'light';
    setTheme(newTheme);
  };

  const toggleLang = () => {
    setLang(lang === 'RO' ? 'EN' : 'RO');
  };

  const handleLogin = async () => {
    setLoggedIn(true);
    navigate('/home');

    // try {
    //   const response = await fetch('https://localhost:8443/v1/login', {
    //     method: 'POST',
    //     headers: {
    //       'Content-Type': 'application/json'
    //     },
    //     credentials: 'include', // if backend sets cookies
    //     body: JSON.stringify({
    //       username: username,
    //       password: password
    //     })
    //   });

    //   if (response.ok) {
    //     setLoggedIn(true);
    //     navigate('/home');
    //   } else {
    //     const errorData = await response.json();
    //     alert('Login failed: ' + (errorData.error || 'Invalid credentials'));
    //   }
    // } catch (error) {
    //   alert('Error: ' + error.message);
    // }
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
      <div className="lang-toggle">
        <img
          src={getLangIcon()}
          alt="Lang"
          onClick={toggleLang}
        />
      </div>

      {/* Theme toggle */}
      <div className="theme-toggle">
        <img
          src={theme === 'light' ? moon : sun}
          alt="Theme"
          onClick={toggleTheme}
        />
      </div>

      <p className="login-title">GEORG.IO</p>

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
