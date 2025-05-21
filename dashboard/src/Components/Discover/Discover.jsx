import React from 'react'
import AdsManager from './AdsManager/AdsManager';
import './Discover.css'

const Discover = () => {
  const userCluster = 0;

  return (
    <div>
      <AdsManager userCluster={userCluster} />
    </div>
  )
}

export default Discover
