import React from 'react'
import AdsManager from './AdsManager/AdsManager';
import './Discover.css'

const Discover = () => {
  const userCluster = 1;

  return (
    <div>
      <AdsManager userCluster={userCluster} />
    </div>
  )
}

export default Discover
