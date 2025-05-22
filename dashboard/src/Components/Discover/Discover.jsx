import React from 'react'
import AdsManager from './AdsManager/AdsManager';
import ServiceTiles from './ServiceTiles/ServiceTiles';
import './Discover.css'

const Discover = ({theme, setTheme}) => {
  const userCluster = 0;

  /* change theme logic */
  const toggle_mode = () => {
    theme === 'light' ? setTheme('dark') : setTheme('light');
  }; 


  return (
    <div>
      <h2 className="discover-title">Descoperă</h2>
      <AdsManager userCluster={userCluster} />
      <h3 className="discover-subtitle"> Servicii</h3><ServiceTiles userCluster={userCluster} />
    </div>
  )
}

export default Discover
