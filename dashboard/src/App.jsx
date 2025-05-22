import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';

import NavBar from './Components/NavBar/NavBar';
import Home from './Components/Home/Home';
import Discover from './Components/Discover/Discover';
import Profile from './Components/Profile/Profile';
import Login from './Components/Auth/Login';
import Register from './Components/Auth/Register';

const App = () => {
  const current_theme = localStorage.getItem('current_theme');
  const [theme, setTheme] = useState(current_theme ? current_theme : 'light');
  const [lang, setLang] = useState(localStorage.getItem('lang') || 'RO');
  const [loggedIn, setLoggedIn] = useState(false);

  // Add state to store user data
  const [userData, setUserData] = useState({
    firstName: '',
    lastName: '',
    email: '',
    username: '',
    password: ''
  });

  useEffect(() => {
    localStorage.setItem('current_theme', theme);
    localStorage.setItem('lang', lang);
  }, [theme, lang]);

  return (
    <div className={`container ${theme}`}>
      <Router>
        {loggedIn && <NavBar theme={theme} setTheme={setTheme} lang={lang} setLang={setLang} />}
        <Routes>
          {!loggedIn ? (
            <>
              <Route path="/" element={<Login lang={lang} setLang={setLang} setLoggedIn={setLoggedIn} setUserData={setUserData} />} />
              <Route path="/register" element={<Register lang={lang} setLang={setLang} setUserData={setUserData} />} />
              <Route path="*" element={<Navigate to="/" replace />} />
            </>
          ) : (
            <>
              <Route path="/home" element={<Home theme={theme} setTheme={setTheme} />} />
              <Route path="/discover" element={<Discover theme={theme} setTheme={setTheme} />} />
              <Route path="/profile" element={<Profile theme={theme} setTheme={setTheme} setLoggedIn={setLoggedIn} userData={userData} />} />
              <Route path="*" element={<Navigate to="/home" replace />} />
            </>
          )}
        </Routes>
      </Router>
    </div>
  );
};

export default App;
