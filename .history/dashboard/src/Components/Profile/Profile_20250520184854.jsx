import profileIconDay from '../../assets/profileIconDay.png'
import profileIconNight from '../../assets/profileIconNight.png'
import editPencilDay from '../../assets/editPencilDay.png'
import editPencilNight from '../../assets/editPencilNight.png'
import React from 'react';
import './Profile.css';

const Profile = ({ theme, setTheme }) => {
    
    /* change theme logic */
    const toggle_mode = () => {
      theme === 'light' ? setTheme('dark') : setTheme('light');
    };
  
    return (
      <div className="profile-container">  
          <div className="profile-card">
            {/* Added icon wrapper div */}
            <div className="icon-wrapper">
              <img
                src={theme === 'light' ? profileIconNight : profileIconDay}
                alt="Profile"
                className="profile-icon"
              />
            </div>
            <p className="profile-card-username">luktechmech</p>
            <p className="profile-card-ID">ID: 30128127</p>
            <p className="profile-card-job">Underwater ceramic expert</p>
            <p className="profile-card-location">Bucharest, Romania</p>
          </div>

          <div className="general-information">
          <p className="general-information-title">GENERAL INFORMATION</p>
          <img
                src={theme === 'light' ? profileIconNight : profileIconDay}
                alt="Profile"
                className="profile-icon"
              />
            {/* <div className="search-icon">
            </div> */}
            

          </div>
      </div>
    );
  };

export default Profile;