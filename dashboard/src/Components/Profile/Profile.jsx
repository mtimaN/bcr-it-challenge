import React from 'react';
import './Profile.css';
import profileIconDay from '../../assets/profile_page/profileIconDay.png'
import profileIconNight from '../../assets/profile_page/profileIconNight.png'

const Profile = ({ theme, setTheme }) => {
    /* change theme logic */
    const toggle_mode = () => {
      theme === 'light' ? setTheme('dark') : setTheme('light');
    };
  
    return (
      <div className="profile-container">
        <h1>Welcome to the User Profile Page</h1>
        <p>This is a test to verify navigation is working.</p>
  
        
        <div className='profile-icon'>
            <img src={theme =='light' ? profileIconNight : profileIconDay} alt=""/>
        </div>
  
      </div>
    );
  };

export default Profile;