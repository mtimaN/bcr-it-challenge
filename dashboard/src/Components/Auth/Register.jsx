import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './Auth.css';

import roD from '../../assets/nav_bar/roDark.png';
import enD from '../../assets/nav_bar/enDark.png';
import roL from '../../assets/nav_bar/roLight.png';
import enL from '../../assets/nav_bar/enLight.png';
import sun from '../../assets/nav_bar/dayLogo.png';
import moon from '../../assets/nav_bar/nightLogo.png';

const Register = ({ lang, setLang }) => {
  const navigate = useNavigate();

  const [theme, setTheme] = useState(localStorage.getItem('current_theme') || 'light');

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

  const handleSubmit = () => {
    navigate('/'); // back to login
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

      {/* Register box */}
      <div className="auth-box">
        <input type="text" placeholder={lang === 'RO' ? 'Nume complet' : 'Full Name'} />
        <input type="text" placeholder="Email" />
        <input type="text" placeholder={lang === 'RO' ? 'Utilizator' : 'Username'} />
        <input type="password" placeholder={lang === 'RO' ? 'Parolă' : 'Password'} />
        <button className="auth-button" onClick={handleSubmit}>
          {lang === 'RO' ? 'Înregistrează-te' : 'Register'}
        </button>
        <p className="switch-link" onClick={() => navigate('/')}>
          {lang === 'RO' ? 'Ai deja cont? Autentifică-te' : 'Already have an account? Login'}
        </p>
      </div>
    </div>
  );
};

export default Register;
