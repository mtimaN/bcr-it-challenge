import React from 'react';
import './AdsManager.css';
import { adsByCluster } from '../Data';


const AdsManager = ({ userCluster, lang }) => {
  const ads = adsByCluster[userCluster] || [];
  const adsToDisplay = [...ads, ...ads];

  return (
    <div className="ads-manager">
      <div className="ads-container">
        {adsToDisplay.map((ad, index) => (
          <div key={index} className="ad-card">
            {ad.image && (
              <img
                src={ad.image}
                alt=""
                className="ad-image"
              />
            )}
            <h3>{lang === 'RO' ? ad.title_ro: ad.title_eng}</h3>
            <p>{lang === 'RO' ? ad.description_ro: ad.description_eng}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default AdsManager;
