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
            <div className="general-information id-field">
              <span className="field-label">ID:</span>
              <span className="field-value">30128127</span>
            </div>
            .field-label {
  font-size: 0.7vw;
  color: #777;
  margin-right: 0.5vw;
  font-weight: 600;
  min-width: 5vw; /* Fixed width for labels to align values */
}

/* Value for info fields */
.field-value {
  font-size: 0.8vw;
  color: #333;
  font-weight: 400;
  flex-grow: 1;
}

          </div>
      </div>
    );
  };

export default Profile;