import newProfileIconDay from '../../assets/profile_page/newProfileIconDay.png'
import newProfileIconNight from '../../assets/profile_page/newProfileIconNight.png'

import editPencilDay from '../../assets/profile_page/editPencilDay.png'
import editPencilNight from '../../assets/profile_page/editPencilNight.png'

import trashCan from '../../assets/profile_page/trashCan.png'

import downCollapseArrowDay from '../../assets/profile_page/downCollapseArrowDay.png'
import downCollapseArrowNight from '../../assets/profile_page/downCollapseArrowNight.png'

import upCollapseArrowDay from '../../assets/profile_page/upCollapseArrowDay.png'
import upCollapseArrowNight from '../../assets/profile_page/upCollapseArrowNight.png'

import gradientLightBlue from '../../assets/profile_page/gradientLightBlue.jpg'

import frame from '../../assets/profile_page/frame.png'

import React, { useState, useEffect, useRef } from 'react';
import './Profile.css';


const Profile = ({ theme, setTheme }) => {
  const [isGenderOpen, setIsGenderOpen] = useState(false);
  const [selectedGender, setSelectedGender] = useState('Undisclosed');

  const [isMarriedOpen, setIsMarriedOpen] = useState(false);
  const [selectedMarried, setSelectedMarried] = useState('Yes');

  const toggleGenderDropdown = () => setIsGenderOpen(prev => !prev);
  const toggleMarriedDropdown = () => setIsMarriedOpen(prev => !prev);

  const genderRef = useRef(null);
  const marriedRef = useRef(null);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (genderRef.current && !genderRef.current.contains(event.target)) {
        setIsGenderOpen(false);
      }
      if (marriedRef.current && !marriedRef.current.contains(event.target)) {
        setIsMarriedOpen(false);
      }
    };
  
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);
    
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
                src={theme === 'light' ? newProfileIconDay : newProfileIconNight}
                alt="Profile"
                className="profile-icon"
              />
            </div>
            <img
                src={gradientLightBlue}
                alt=""
                className="gradient-blue"
            />

            {/* profile username */}
            <p className="profile-card-username">luktechmech</p>

            {/* account ID */}
            <p className="profile-card-ID">ID: 30128127</p>

            <img
                src={frame}
                alt=""
                className=""
            />

            {/* occupation */}
            <p className="profile-card-job">Underwater ceramic expert</p>

            {/* current address */}
            <p className="profile-card-location">Bacau, Romania</p>

            {/* joined date */}
            <p className="profile-card-join-date">Joined May 2025</p>
            </div>

          <div className="general-information">
            <p className="general-information-title">PERSONAL INFORMATION</p>

            <img
                src={gradientLightBlue}
                alt=""
                className="gradient-blue-info"
            />

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
            <div className="search-icon-last-name">
              <p className="search-icon-text">Last name*</p>
              <p className="search-icon-info">Botez</p>
            </div>

            {/* phone number */}
            <div className="search-icon-phone-number">
              <p className="search-icon-text">Phone number*</p>
              <p className="search-icon-info">+40721628090</p>
            </div>

            {/* email address */}
            <div className="search-icon-email">
              <p className="search-icon-text">Email*</p>
              <p className="search-icon-info">luktechmech@gmail.com</p>
            </div>

            {/* gender */}
            <div className="search-icon-gender" ref={genderRef}>
              <p className="search-icon-text">Gender</p>
              <p className="search-icon-info">{selectedGender}</p>
              <img
                src={theme === 'light'
                  ? (isGenderOpen ? upCollapseArrowDay : downCollapseArrowDay)
                  : (isGenderOpen ? upCollapseArrowNight : downCollapseArrowNight)
                }
                alt=""
                className="arrow-icon"
                onClick={toggleGenderDropdown}
              />
              {isGenderOpen && (
                <ul className="dropdown-menu">
                  {['Male', 'Female', 'Undisclosed'].map(option => (
                    <li
                      key={option}
                      className={`dropdown-item ${option.toLowerCase()}`}
                      onClick={() => {
                        setSelectedGender(option);
                        setIsGenderOpen(false);
                      }}
                    >
                      {option}
                    </li>
                  ))}
                </ul>
              )}
            </div>
            
            {/* married status */}
            <div className="search-icon-married" ref={marriedRef}>
              <p className="search-icon-text">Married</p>
              <p className="search-icon-info">{selectedMarried}</p>
              <img
                src={theme === 'light'
                  ? (isMarriedOpen ? upCollapseArrowDay : downCollapseArrowDay)
                  : (isMarriedOpen ? upCollapseArrowNight : downCollapseArrowNight)
                }
                alt=""
                className="arrow-icon"
                onClick={toggleMarriedDropdown}
              />
              {isMarriedOpen && (
                <ul className="married-dropdown-menu">
                  {['Yes', 'No'].map(option => (
                    <li
                      key={option}
                      className={`married-dropdown-item ${option.toLowerCase()}`}
                      onClick={() => {
                        setSelectedMarried(option);
                        setIsMarriedOpen(false);
                      }}
                    >
                      {option}
                    </li>
                  ))}
                </ul>
              )}
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

            {/* delete option */}
            <img
              src={trashCan}
              alt=""
              className="delete-icon"
            />
          </div>
      </div>
    );
  };

export default Profile;