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

import settingsDay from '../../assets/profile_page/settingsDay.png'
import logoutDay from '../../assets/profile_page/logoutDay.png'
import bankDay from '../../assets/profile_page/bankDay.png'
import helpDay from '../../assets/profile_page/helpDay.png'
import userAgreementDay from '../../assets/profile_page/userAgreementDay.png'

import LogoutModal from './LogoutModal/LogoutModal';
import EditModal from './EditModal/EditModal';
import DeleteModal from './DeleteModal/DeleteModal';

import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';

import './Profile.css';
import profileIconDay from '../../assets/profile_page/profileIconDay.png'
import profileIconNight from '../../assets/profile_page/profileIconNight.png'

// Settings Modal Component
const SettingsModal = ({ onClose }) => {
  const [newPassword, setNewPassword] = useState('');

  const handlePasswordChange = () => {
    if (newPassword.trim()) {
      console.log('Password changed to:', newPassword);

      //! Add your password change logic here
      // Update endpoint
      
      alert('Password changed successfully!');
      onClose();
    } else {
      alert('Please enter a new password');
    }
  };

  return (
    <div className="logout-modal-overlay" onClick={onClose}>
      <div className="logout-modal-container" onClick={(e) => e.stopPropagation()}>
        <h2 className="logout-modal-title">Schimbare parolă</h2>
        <p className="logout-modal-body">Introduceți noua parolă:</p>
        <input
          type="password"
          value={newPassword}
          onChange={(e) => setNewPassword(e.target.value)}
          placeholder="Noua parolă"
          style={{
            width: '100%',
            padding: '10px',
            marginBottom: '20px',
            border: '1px solid #ccc',
            borderRadius: '4px',
            fontSize: '16px'
          }}
        />
        <div className="logout-modal-buttons">
          <button className="logout-modal-confirm" onClick={handlePasswordChange}>
            Salvează
          </button>
          <button className="logout-modal-cancel" onClick={onClose}>Anulare</button>
        </div>
      </div>
    </div>
  );
};

const Profile = ({ theme, setTheme, setLoggedIn, userData }) => {
  const navigate = useNavigate();

  const [showLogoutModal, setShowLogoutModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [showSettingsModal, setShowSettingsModal] = useState(false);

  // Use userData from props instead of location.state
  const { firstName = '', lastName = '', email = '', username = '', password = '' } = userData || {};

  // Generate random ID with same number of digits (8 digits)
  const [userID] = useState(() => Math.floor(10000000 + Math.random() * 90000000));

  // handle logout logic
  const handleLogout = () => {
    localStorage.removeItem('jwtToken');
    setLoggedIn(false);
    setShowLogoutModal(false);
    navigate('/');
  };

  useEffect(() => {  
    const handleClickOutside = (event) => {
      // Remove gender and married dropdown logic
    };
  
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  const [editMode, setEditMode] = useState(false);

  // handle edit logic - check against registration password
  const handleEditConfirm = (inputPassword) => {
    if (inputPassword === password) {
      setShowEditModal(false);
      setEditMode(true);
      return true;
    } else {
      return false;
    }
  };

  const handleDelete = () => {
    console.log('Account deleted'); // Replace with real delete logic
    setShowDeleteModal(false);
  };

  // info for each field - initialize with userData
  const [firstNameState, setFirstName] = useState(firstName);
  const [lastNameState, setLastName] = useState(lastName);
  const [phoneNumber, setPhoneNumber] = useState('');
  const [emailState, setEmail] = useState(email);
  const [usernameState] = useState(username);
  const [county, setCounty] = useState('');
  const [address, setAddress] = useState('');
  const [gender, setGender] = useState('');
  const [marriedStatus, setMarriedStatus] = useState('');

  // Update state when userData changes
  useEffect(() => {
    if (userData) {
      setFirstName(userData.firstName || '');
      setLastName(userData.lastName || '');
      setEmail(userData.email || '');
    }
  }, [userData]);

  /* change theme logic */
  const toggle_mode = () => {
    theme === 'light' ? setTheme('dark') : setTheme('light');
  };

  // the magic
  return (
    // wrapper for everything
    <div className="profile-container">  
        {/* profile card */}
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
          <p className="profile-card-username">{usernameState}</p>

          {/* account ID */}
          <p className="profile-card-ID">ID: {userID}</p>

          <img
              src={frame}
              alt=""
              className="qr-code"
          />

          {/* joined date */}
          <p className="profile-card-join-date">S-a alăturat în mai 2025</p>
          </div>

        {/* general info section */}
        {editMode && <div className="edit-overlay"></div>}
        <div className={`general-information ${editMode ? 'edit-mode' : ''}`}>
          {/* save edit changes button */}
          {editMode && (
            <button
              className="save-button"
              onClick={() => setEditMode(false)}
            >
              Salvare
            </button>
          )}

          <p className="general-information-title">INFORMAȚII PERSONALE</p>

          <img
              src={gradientLightBlue}
              alt=""
              className="gradient-blue-info"
          />

          {/* edit pencil icon */}
          <img
            src={theme === 'light' ? editPencilDay : editPencilNight}
            alt="Edit"
            className="edit-icon"
            onClick={() => {
              if (!editMode) setShowEditModal(true);
            }}
            style={{ cursor: editMode ? 'default' : 'pointer', opacity: editMode ? 0.5 : 1 }}
          />

          {/* first name */}
          <div className="search-icon-first-name">
            <p className="search-icon-text">Prenume *</p>
            {editMode ? (
              <input
                type="text"
                value={firstNameState}
                onChange={(e) => setFirstName(e.target.value)}
                className="search-icon-info editable"
              />
            ) : (
              <p className="search-icon-info">{firstNameState}</p>
            )}
          </div>
          
          {/* last name */}
          <div className="search-icon-last-name">
            <p className="search-icon-text">Nume *</p>
            {editMode ? (
              <input
                type="text"
                value={lastNameState}
                onChange={(e) => setLastName(e.target.value)}
                className="search-icon-info editable"
              />
            ) : (
              <p className="search-icon-info">{lastNameState}</p>
            )}
          </div>

          {/* phone number */}
          <div className="search-icon-phone-number">
            <p className="search-icon-text">Număr de telefon *</p>
            {editMode ? (
              <input
                type="text"
                value={phoneNumber}
                onChange={(e) => setPhoneNumber(e.target.value)}
                className="search-icon-info editable"
              />
            ) : (
              <p className="search-icon-info">{phoneNumber}</p>
            )}
          </div>

          {/* email address */}
          <div className="search-icon-email">
            <p className="search-icon-text">Adresă de email *</p>
            {editMode ? (
              <input
                type="text"
                value={emailState}
                onChange={(e) => setEmail(e.target.value)}
                className="search-icon-info editable"
              />
            ) : (
              <p className="search-icon-info">{emailState}</p>
            )}
          </div>

          {/* gender */}
          <div className="search-icon-gender">
            <p className="search-icon-text">Gen</p>
            {editMode ? (
              <input
                type="text"
                value={gender}
                onChange={(e) => setGender(e.target.value)}
                className="search-icon-info editable"
              />
            ) : (
              <p className="search-icon-info">{gender}</p>
            )}
          </div>
          
          {/* married status */}
          <div className="search-icon-married">
            <p className="search-icon-text">Căsătorit/ă</p>
            {editMode ? (
              <input
                type="text"
                value={marriedStatus}
                onChange={(e) => setMarriedStatus(e.target.value)}
                className="search-icon-info editable"
              />
            ) : (
              <p className="search-icon-info">{marriedStatus}</p>
            )}
          </div>

          {/* county */}
          <div className="search-icon-county">
            <p className="search-icon-text">Județ *</p>
            {editMode ? (
              <input
                type="text"
                value={county}
                onChange={(e) => setCounty(e.target.value)}
                className="search-icon-info editable"
              />
            ) : (
              <p className="search-icon-info">{county}</p>
            )}
          </div>

          {/* address of residence */}
          <div className="search-icon-address">
            <p className="search-icon-text">Adresă *</p>
            {editMode ? (
              <input
                type="text"
                value={address}
                onChange={(e) => setAddress(e.target.value)}
                className="search-icon-info editable"
              />
            ) : (
              <p className="search-icon-info">{address}</p>
            )}
          </div>

          {/* delete option - only show in edit mode */}
          {editMode && (
            <img
              src={trashCan}
              alt="Delete"
              className="delete-icon"
              onClick={() => setShowDeleteModal(true)}
              style={{ cursor: 'pointer' }}
            />
          )}
        </div>

        
        <div className="logout-button" onClick={() => setShowLogoutModal(true)}>
          <img src={logoutDay} alt="Logout" className="logout-icon" />
          <p className="logout-title">Ieșire din cont</p>
        </div>

        <div className="settings-button" onClick={() => setShowSettingsModal(true)}>
          <img
              src={settingsDay}
              alt=""
              className="settings-icon"
          />

          <p className="settings-title">Setări</p>
        </div>

        <a
          href="https://www.bcr.ro/en/about-us/our-company"
          target="_blank"
          rel="noopener noreferrer"
          style={{ textDecoration: 'none' }}
        >
          <div className="about-button">
            <img
              src={bankDay}
              alt=""
              className="about-icon"
            />
            <p className="about-title">Despre BCR</p>
          </div>
        </a>
        
        <a
          href="https://www.bcr.ro/en/individuals/help-center"
          target="_blank"
          rel="noopener noreferrer"
          style={{ textDecoration: 'none' }}
        >
          <div className="help-button">
            <img
                src={helpDay}
                alt=""
                className="help-icon"
            />

            <p className="help-title">Ajutor</p>
          </div>
        </a>

        <a
          href="https://www.bcr.ro/en/terms-conditions"
          target="_blank"
          rel="noopener noreferrer"
          style={{ textDecoration: 'none' }}
        >
        <div className="user-agreement-button">
          <img
              src={userAgreementDay}
              alt=""
              className="user-agreement-icon"
          />

          <p className="user-agreement-title">Condiții</p>
        </div>
        </a>

        {showSettingsModal && (
          <SettingsModal onClose={() => setShowSettingsModal(false)} />
        )}

        {showLogoutModal && (
          <LogoutModal onClose={() => setShowLogoutModal(false)} onConfirm={handleLogout} />
        )}

        {showEditModal && (
          <EditModal onClose={() => setShowEditModal(false)} onConfirm={handleEditConfirm} />
        )}

        {showDeleteModal && (
          <DeleteModal onClose={() => setShowDeleteModal(false)} onConfirm={handleDelete} />
        )}
    </div>

  );
};

export default Profile;
