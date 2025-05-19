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
            <img src={theme =='light' ? profileIconDay : profileIconNight} alt=""/>
        </div>
  
        <img onClick={()=>{toggle_mode()}} src={theme =='light' ? moon : sun}
        alt="" className='toggle-icon'/>
      </div>
    );
  };

export default Profile;