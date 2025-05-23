import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './Auth.css';

import roD from '../../assets/nav_bar/roDark.png';
import enD from '../../assets/nav_bar/enDark.png';
import roL from '../../assets/nav_bar/roLight.png';
import enL from '../../assets/nav_bar/enLight.png';
import sun from '../../assets/nav_bar/dayLogo.png';
import moon from '../../assets/nav_bar/nightLogo.png';

const Register = ({ lang, setLang, setUserData }) => {
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
        // Store user data in parent component state for later use
        setUserData({
          firstName,
          lastName,
          email,
          username,
          password
        });

        localStorage.setItem('userData', JSON.stringify({
          firstName,
          lastName,
          email,
          username,
          password
        }));
        
        // Navigate back to login page after successful registration
        navigate('/');

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

      <p className="login-title">GEORG.IO</p>

      <p className="login-description">Next-gen Personalized Banking Experience</p>

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
