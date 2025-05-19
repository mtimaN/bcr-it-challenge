import React, {useEffect, useState} from 'react'
import NavBar from './Components/NavBar/NavBar'

/* to integrate page browsing */
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

/* profile page */
import Profile from './Components/Profile/Profile'

const App = () => {
  /* keep current theme on refresh */
  const current_theme = localStorage.getItem('current_theme');
  const [theme, setTheme] = useState(current_theme?
    current_theme : 'light');

  /* when theme updates, function calls */
  useEffect(()=>{
    localStorage.setItem('current_theme', theme);
  }, [theme])

  return (
    <div className={`container ${theme}`}>
      <Router>
        <NavBar theme={theme} setTheme={setTheme} />
        <Routes>
          <Route path="/" element={<h1>Welcome Home</h1>} />
          <Route path="/profile" element={<Profile />} />
        </Routes>
      </Router>
    </div>
  )
}

export default App
