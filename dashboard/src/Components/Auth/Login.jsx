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

  // sync theme with DOM class
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

  const handleLogin = () => {
    setLoggedIn(true);
    navigate('/home');
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

      {/* Auth box */}
      <div className="auth-box">
        <input type="text" placeholder={lang === 'RO' ? 'Utilizator' : 'Username'} />
        <input type="password" placeholder={lang === 'RO' ? 'Parolă' : 'Password'} />
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
