import React from 'react'
import AdsManager from './AdsManager/AdsManager';
import ServiceTiles from './ServiceTiles/ServiceTiles';
import './Discover.css'

const Discover = ({theme, lang}) => {
  const userCluster = 3;

  return (
    <div>
      <h2 className="discover-title">Descoperă</h2>
      <AdsManager userCluster={userCluster} lang={lang} />
      <h3 className="discover-subtitle"> Servicii</h3><ServiceTiles userCluster={userCluster} lang={lang} />
    </div>
  )
}

export default Discover
