import profileIconDay from '../../assets/profileIconDay.png'
import profileIconNight from '../../assets/profileIconNight.png'
import React from 'react';
import './Profile.css';

const Profile = () => {

    /* change theme logic */
  c onst toggle_mode = ()=>{
    theme == 'light' ? setTheme('dark') : setTheme('light');
    }
    return (
      <div className="profile-container">
        <h1>Welcome to the User Profile Page</h1>
        <p>This is a test to verify navigation is working.</p>
        <img src={profileIconDay} alt="Profile Icon" className="profile-image" />
      </div>
    );
  };

export default Profile;