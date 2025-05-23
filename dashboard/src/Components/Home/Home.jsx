import Balance from './Balance/Balance'
import Transactions from './Transactions/Transactions'
import SpendingVisual from './SpendingVisual/SpendingVisual';
import SpentChart from './SpentChart/SpentChart'

import React, { useState } from 'react';

import './Home.css'

const Home = ({ theme, lang }) => {

  return (
    <div className="home-container">
      <Balance theme={theme} lang={lang}/>
      <Transactions theme={theme} lang={lang}/>
      <SpendingVisual theme={theme} lang={lang}/>
      <SpentChart theme={theme} lang={lang}/>
    </div>
  )
}

export default Home
