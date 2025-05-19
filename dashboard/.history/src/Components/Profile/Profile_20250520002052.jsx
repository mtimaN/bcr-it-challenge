import profileIconDay from '../../assets/profileIconDay.png'
import profileIconNight from '../../assets/profileIconNight.png'
import React from 'react';
import './Profile.css';

const Profile = ({ theme, setTheme }) => {
    // Theme toggle logic
    const toggle_mode = () => {
      theme === 'light' ? setTheme('dark') : setTheme('light');
    };
  
    return (
      <div className="profile-container">
        <h1>Welcome to the User Profile Page</h1>
        <p>This is a test to verify navigation is working.</p>
  
        {/* Show the correct icon based on theme */}
        <img
          src={theme === 'light' ? profileIconDay : profileIconNight}
          alt="Profile Icon"
          className="profile-image"
        />
  
        {/* Theme toggle button */}
        <button className="theme-toggle-btn" onClick={toggle_mode}>
          Switch to {theme === 'light' ? 'Dark' : 'Light'} Mode
        </button>
      </div>
    );
  };

export default Profile;