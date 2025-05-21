import React from 'react';
import './AdsManager.css';

const adsByCluster = {
  0: [
    {
      title: "Economisește mai ușor",
      description: "Vezi cele mai bune conturi de economii.",
      image: "images/ads/economyAccount.jpg"
    },
    {
      title: "Termen lung, câștig mare",
      description: "Depozite avantajoase pentru tine.",
      image: "/images/ads/depositFavorable.jpg"
    },
    {
      title: "Analizează-ți cheltuielile",
      description: "Grafice lunare pentru control financiar.",
      image: "/images/ads/financialGraphic.jpg"
    }
  ],
  1: [
    {
      title: "Ai cheltuieli mari?",
      description: "Descoperă cardurile de credit smart.",
      image: "images/ads/smartCard.jpg"
    },
    {
      title: "Credit rapid, fără griji",
      description: "Aplică 100% online pentru împrumut.",
      image: "images/ads/loan.jpg "
    },
    {
      title: "Folosește la maxim overdraft-ul",
      description: "Vezi limitele disponibile.",
      image: "images/ads/overdraft.jpg"
    }
  ],
  2: [
    {
      title: "George te poate ajuta",
      description: "Descoperă beneficiile contului digital.",
      image: "images/ads/digitalAccount.jpg"
    },
    {
      title: "Transferuri mai rapide",
      description: "Încearcă plățile instant.",
      image: "images/ads/contactless.jpg"
    },
    {
      title: "Economisește fără efort",
      description: "Setează un plan automat.",
      image: "images/ads/financePlan.jpg"
    }
  ],
  3: [
    {
      title: "Trimite bani rapid",
      description: "Transfer instant către prieteni.",
      image: "images/ads/fastTransfers.jpg"
    },
    {
      title: "Economii simple",
      description: "Economisește automat când cheltui.",
      image: "images/ads/economyAuto.jpg"
    },
    {
      title: "Cheltuie smart",
      description: "Vezi unde se duc banii tăi.",
      image: "images/ads/moneyLocation.jpg"
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
