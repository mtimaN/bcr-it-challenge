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
        <h1>Welcome to the User Profile Page</h1>
        <p>This is a test to verify navigation is working.</p>
  
        
        <div className='profile-icon'>
            <input type="text" placeholder='Search'/>
            <img src={theme =='light' ? magGlassLight : magGlassDark} alt=""/>
        </div>
  
        {/* Theme toggle button */}
        <button className="theme-toggle-btn" onClick={toggle_mode}>
          Switch to {theme === 'light' ? 'Dark' : 'Light'} Mode
        </button>
      </div>
    );
  };

export default Profile;