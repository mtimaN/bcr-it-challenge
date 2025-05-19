import React, {useEffect, useState} from 'react'
import NavBar from './Components/NavBar/NavBar'

/* to integrate page browsing */
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'; 

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
      <NavBar theme={theme} setTheme={setTheme}/>
    </div>
  )
}

export default App
