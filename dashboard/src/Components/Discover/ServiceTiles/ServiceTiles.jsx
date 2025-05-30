import React from 'react';
import './ServiceTiles.css';
import { allServices, serviceOrderByCluster } from '../Data';

import { Link } from 'react-router-dom';

const ServiceTiles = ({ userCluster, lang }) => {
  const serviceOrder = serviceOrderByCluster[userCluster] || Object.keys(allServices);

  return (
    <div className="services-container">
      {serviceOrder.map((key) => {
        const service = allServices[key];
        return (
          <div key={key} className="service-wrapper">
            <Link to="/" key={key} className="service-wrapper">
              <div className="service-tile" style={{ backgroundColor: service.color }}>
                <img src={service.image} alt="" className="service-image" />
              </div>
              <div className="service-label">{lang === 'RO' ? service.label_ro : service.label_eng}</div>
            </Link>
          </div>
        );
      })}
    </div>
  );
};

export default ServiceTiles;
