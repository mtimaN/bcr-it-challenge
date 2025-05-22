import React from 'react';
import './ServiceTiles.css';
import { allServices, serviceOrderByCluster } from '../Data';

const ServiceTiles = ({ userCluster }) => {
  const serviceOrder = serviceOrderByCluster[userCluster] || Object.keys(allServices);

  return (
    <div className="services-container">
      {serviceOrder.map((key) => {
        const service = allServices[key];
        return (
          <div key={key} className="service-wrapper">
            <div className="service-tile" style={{ backgroundColor: service.color }}>
              <img src={service.image} alt={service.label} className="service-image" />
            </div>
            <div className="service-label">{service.label}</div>
          </div>
        );
      })}
    </div>
  );
};

export default ServiceTiles;
``