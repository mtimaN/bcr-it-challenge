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

  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');

  const [email, setEmail] = useState('');

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

  const handleSubmit = async () => {
    try {
      const response = await fetch('https://localhost:8443/v1/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        credentials: 'include',
        body: JSON.stringify({
          first_name: firstName,
          last_name: lastName,
          email: email,
          username: username,
          password: password
        })
      });

      if (response.ok) {
        navigate('/', {
          state: {
            firstName,
            lastName,
            email,
            username
          }
        });

      } else {
        const errorData = await response.json();
        alert('Registration failed: ' + (errorData.message || 'Unknown error'));
      }

    } catch (error) {
      alert('Error: ' + error.message);
    }
  };

  const getLangIcon = () => {
    if (lang === 'RO') {
      return theme === 'light' ? roD : roL;
    }

    return theme === 'light' ? enD : enL;
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
        <input
          type="text"
          placeholder={lang === 'RO' ? 'Prenume' : 'First Name'}
          value={firstName}
          onChange={(e) => setFirstName(e.target.value)}
        />
        <input
          type="text"
          placeholder={lang === 'RO' ? 'Nume' : 'Last Name'}
          value={lastName}
          onChange={(e) => setLastName(e.target.value)}
        />
        <input
          type="text"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
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
