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
        <div className="profile-wrapper">
          <div className="square">
          <img
            src={theme === 'light' ? profileIconNight : profileIconDay}
            alt=""
            className="profile-icon"
          />
</div>
        </div>
  
      </div>
    );
  };

export default Profile;