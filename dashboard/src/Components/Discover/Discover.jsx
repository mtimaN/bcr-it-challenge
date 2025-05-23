import { useEffect, useState } from 'react';
import AdsManager from './AdsManager/AdsManager';
import ServiceTiles from './ServiceTiles/ServiceTiles';
import './Discover.css';

const Discover = ({ theme, lang, setTheme }) => {
  const [userCluster, setUserCluster] = useState(null);

  useEffect(() => {
    const fetchCategory = async () => {
      try {
        const token = localStorage.getItem('jwtToken');

        const response = await fetch('https://localhost:8443/v1/get_ads', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          }
        });

        if (response.ok) {
          const data = await response.json();
          setUserCluster(data.category);
        } else {
          const errorData = await response.json();
          console.error('Failed to fetch user category:', errorData.message || 'Unknown error');
        }
      } catch (error) {
        console.error('Error fetching user category:', error.message);
      }
    };

    fetchCategory();
  }, []);

  return (
    <div>
      <h2 className="discover-title">{lang === 'RO' ? 'Doar pentru tine:' : 'Only for you:'}</h2>
      <AdsManager userCluster={userCluster} lang={lang} />
      <h3 className="discover-subtitle"> {lang === 'RO' ? 'Servicii' : 'Services'}</h3><ServiceTiles userCluster={userCluster} lang={lang} />
    </div>
  )
}

export default Discover;
