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
        <div className='profile-icon'>
          <img src={theme =='light' ? profileIconNight : profileIconDay} alt=""/>
        </div>

        <div>
          <div className="square"></div>
        </div>
  
      </div>
    );
  };

export default Profile;