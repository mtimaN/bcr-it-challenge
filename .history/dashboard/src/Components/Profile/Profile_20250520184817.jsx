import profileIconDay from '../../assets/profileIconDay.png'
import profileIconNight from '../../assets/profileIconNight.png'
import editPencil_day from '../../assets/edit_pencil_day.png'
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
            {/* <div className="search-icon">
            </div> */}
            

          </div>
      </div>
    );
  };

export default Profile;