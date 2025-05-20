import React from 'react';
import './AdsManager.css';

import smartCard from '../../../assets/ads/'

const adsByCluster = {
  0: [
    {
      title: "Economisește mai ușor",
      description: "Vezi cele mai bune conturi de economii.",
      image: ""
    },
    {
      title: "Termen lung, câștig mare",
      description: "Depozite avantajoase pentru tine.",
      image: ""
    },
    {
      title: "Analizează-ți cheltuielile",
      description: "Grafice lunare pentru control financiar.",
      image: ""
    }
  ],
  1: [
    {
      title: "Ai cheltuieli mari?",
      description: "Descoperă carduri de credit smart.",
      image: ""
    },
    {
      title: "Credit rapid, fără griji",
      description: "Aplică 100% online pentru împrumut.",
      image: ""
    },
    {
      title: "Folosește la maxim overdraft-ul",
      description: "Vezi limitele disponibile.",
      image: ""
    }
  ],
  2: [
    {
      title: "George te poate ajuta",
      description: "Descoperă beneficiile contului digital.",
      image: ""
    },
    {
      title: "Transferuri mai rapide",
      description: "Încearcă plățile instant.",
      image: ""
    },
    {
      title: "Economisește fără efort",
      description: "Setează un plan automat.",
      image: ""
    }
  ],
  3: [
    {
      title: "Trimite bani rapid",
      description: "Transfer instant către prieteni.",
      image: ""
    },
    {
      title: "Cheltuie smart",
      description: "Vezi unde se duc banii tăi.",
      image: ""
    },
    {
      title: "Economii simple",
      description: "Economisește automat când cheltui.",
      image: ""
    }
  ]
};


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
