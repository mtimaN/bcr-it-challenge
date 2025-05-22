import Balance from './Balance/Balance'
import Transactions from './Transactions/Transactions'
import SpendingVisual from './SpendingVisual/SpendingVisual';
import SpentChart from './SpentChart/SpentChart'

import React, { useState } from 'react';

import './Home.css'

const Home = ({ theme, setTheme }) => {

  /* change theme logic */
  const toggle_mode = () => {
    theme === 'light' ? setTheme('dark') : setTheme('light');
  };

  return (
    <div className="home-container">
      <Balance theme={theme} setTheme={setTheme}/>
      <Transactions theme={theme} setTheme={setTheme}/>
      <SpendingVisual theme={theme} setTheme={setTheme}/>
      <SpentChart theme={theme} setTheme={setTheme}/>
    </div>
  )
}

export default Home
