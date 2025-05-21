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
          <div className="square">
            <div className="icon-wrapper">
              <img
                src={theme === 'light' ? profileIconNight : profileIconDay}
                alt="Profile"
                className="profile-icon"
              />
            </div>
            <p className="square-username">luktechmech</p>
            
            {/* Search-bar style information section */}
            <div className="info-section">
              {/* ID Field */}
              <div className="info-field id-field">
                <span className="field-label">ID:</span>
                <span className="field-value">30128127</span>
              </div>
              
              {/* Job Field */}
              <div className="info-field job-field">
                <span className="field-label">Job:</span>
                <span className="field-value">Frontend Developer</span>
              </div>
              
              {/* Additional fields - examples */}
              <div className="info-field">
                <span className="field-label">Email:</span>
                <span className="field-value">luktechmech@example.com</span>
              </div>
              
              <div className="info-field">
                <span className="field-label">Location:</span>
                <span className="field-value">San Francisco, CA</span>
              </div>
              
              <div className="info-field">
                <span className="field-label">Joined:</span>
                <span className="field-value">March 2023</span>
              </div>
            </div>
          </div>
      </div>
    );
  };

export default Profile;