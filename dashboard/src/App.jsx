import React, {useEffect, useState} from 'react'

import NavBar from './Components/NavBar/NavBar'
import Home from './Components/Home/Home'
import Discover from './Components/Discover/Discover'
import Products from './Components/Products/Products'
import Profile from './Components/Profile/Profile'
import Settings from './Components/Settings/Settings'

/* to integrate page browsing */
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

const App = () => {
  /* keep current theme on refresh */
  const current_theme = localStorage.getItem('current_theme');
  const [theme, setTheme] = useState(current_theme?
    current_theme : 'light');

  /* when theme updates, this function calls */
  useEffect(()=>{
    localStorage.setItem('current_theme', theme);
  }, [theme])

  return (
    <div className={`container ${theme}`}>
      <Router>
        <NavBar theme={theme} setTheme={setTheme} />
        <Routes>
          <Route path="/home" element={<Home theme={theme} setTheme={setTheme} />} />
          <Route path="/discover" element={<Discover theme={theme} setTheme={setTheme} />} />
          <Route path="/products" element={<Products theme={theme} setTheme={setTheme} />} />
          <Route path="/profile" element={<Profile theme={theme} setTheme={setTheme} />} />
          <Route path="/settings" element={<Settings theme={theme} setTheme={setTheme} />} />
        </Routes>
      </Router>
    </div>
  )
}

export default App
