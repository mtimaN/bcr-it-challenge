import React from 'react';
import './AdsManager.css';
import { adsByCluster } from '../Data';


const AdsManager = ({ userCluster }) => {
  const ads = adsByCluster[userCluster] || [];

  return (
    <div className="ads-manager">
      <div className="ads-container">
        {ads.map((ad, index) => (
          <div key={index} className="ad-card">
            {ad.image && (
              <img
                src={ad.image}
                alt={ad.title}
                className="ad-image"
              />
            )}
            <h3>{ad.title}</h3>
            <p>{ad.description}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default AdsManager;
