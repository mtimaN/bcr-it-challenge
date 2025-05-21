import profileIconDay from '../../assets/profile_page/profileIconDay.png'
import profileIconNight from '../../assets/profile_page/profileIconNight.png'
import editPencilDay from '../../assets/profile_page/editPencilDay.png'
import editPencilNight from '../../assets/profile_page/editPencilNight.png'
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
            <p className="profile-card-location">Bacau, Romania</p>
          </div>

          <div className="general-information">
            <p className="general-information-title">GENERAL INFORMATION</p>

            {/* edit pencil icon */}
            <img
              src={theme === 'light' ? editPencilDay : editPencilNight}
              alt=""
              className="edit-icon"
            />

            {/* first name */}
              <div className="search-icon-first-name">
              <p className="search-icon-text">First name*</p>
              <p className="search-icon-info">Luca</p>
            </div>
            
            {/* last name */}
            {/* <div className="search-icon-last-name">
              <p className="search-icon-text">Last name*</p>
              <p className="search-icon-info">Botez</p>
            </div> */}

            {/* email address */}
            <div className="search-icon-email">
              <p className="search-icon-text">Email*</p>
              <p className="search-icon-info">luktechmech@gmail.com</p>
            </div>

            {/* gender */}
            <div className="search-icon-gender">
              <p className="search-icon-text">Gender</p>
              <p className="search-icon-info">Undisclosed</p>
            </div>
            
            {/* married status */}
            <div className="search-icon-married">
              <p className="search-icon-text">Married</p>
              <p className="search-icon-info">Yes</p>
            </div>

            {/* county */}
            <div className="search-icon-county">
              <p className="search-icon-text">County*</p>
              <p className="search-icon-info">Bacau</p>
            </div>

            {/* address of residence */}
            <div className="search-icon-address">
              <p className="search-icon-text">Address*</p>
              <p className="search-icon-info">Str Milcov, 4</p>
            </div>

          </div>
      </div>
    );
  };

export default Profile;