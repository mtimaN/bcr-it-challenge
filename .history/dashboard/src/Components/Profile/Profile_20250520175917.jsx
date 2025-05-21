import profileIconDay from '../../assets/profileIconDay.png'
import profileIconNight from '../../assets/profileIconNight.png'
import React from 'react';
import './Profile.css';

const Profile = ({ theme, setTheme }) => {
    
    /* change theme logic */
    const toggle_mode = () => {
      theme === 'light' ? setTheme('dark') : setTheme('light');
    };
  
    return (
      <div className="profile-container">  
          <div className="square">
            {/* Added icon wrapper div */}
            <div className="icon-wrapper">
              <img
                src={theme === 'light' ? profileIconNight : profileIconDay}
                alt="Profile"
                className="profile-icon"
              />
            </div>
            <p className="square-username">luktechmech</p>
            <p className="square-ID">ID: 30128127</p>
            <p className="square-ID">ID: 30128127</p>
          </div>
      </div>
    );
  };

export default Profile;