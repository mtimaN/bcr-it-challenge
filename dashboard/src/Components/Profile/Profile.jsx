// Asset Imports
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
import profileIconDay from '../../assets/profile_page/profileIconDay.png'
import profileIconNight from '../../assets/profile_page/profileIconNight.png'

// Component Imports
import LogoutModal from './LogoutModal/LogoutModal';
import EditModal from './EditModal/EditModal';

import DeleteModal from './DeleteModal/DeleteModal';
import SettingsModal from './SettingsModal/SettingsModal';

import { useLocation } from 'react-router-dom';

// Library Imports
import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import './Profile.css';

const Profile = ({ theme, setTheme, setLoggedIn, userData }) => {
  const navigate = useNavigate();

  // STATE MANAGEMENT
  
  // Modal States
  const [showLogoutModal, setShowLogoutModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [showSettingsModal, setShowSettingsModal] = useState(false);

  // User Data States
  const { firstName = '', lastName = '', email = '', username = '', password = '' } = userData || {};
  const [userID] = useState(() => Math.floor(10000000 + Math.random() * 90000000));
  const [editMode, setEditMode] = useState(false);

  // Form Data States
  const [firstNameState, setFirstName] = useState(firstName);
  const [lastNameState, setLastName] = useState(lastName);
  const [phoneNumber, setPhoneNumber] = useState('');
  const [emailState, setEmail] = useState(email);
  const [usernameState] = useState(username);
  const [county, setCounty] = useState('');
  const [address, setAddress] = useState('');
  const [gender, setGender] = useState('');
  const [marriedStatus, setMarriedStatus] = useState('');

  // API FUNCTIONS

  const handleDeleteAccount = async () => {
    try {
      const token = localStorage.getItem('jwtToken');

      const response = await fetch('https://localhost:8443/v1/delete', {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Accept': 'application/json'
        }
      });

      if (response.ok) {
        // Account deleted successfully, now logout
        localStorage.removeItem('jwtToken');
        setLoggedIn(false);
        setShowDeleteModal(false);
        navigate('/');
        alert('Account deleted successfully');
      } else {
        const errorData = await response.json();
        console.error('Failed to delete account:', errorData.message || 'Unknown error');
        alert('Failed to delete account. Please try again.');
      }
    } catch (error) {
      console.error('Error deleting account:', error.message);
      alert('An error occurred. Please try again.');
    }
  };

  // EVENT HANDLERS

  const handleLogout = () => {
    localStorage.removeItem('jwtToken');
    setLoggedIn(false);
    setShowLogoutModal(false);
    navigate('/');
  };

  const handleEditConfirm = (inputPassword) => {
    if (inputPassword === password) {
      setShowEditModal(false);
      setEditMode(true);
      return true;
    } else {
      return false;
    }
  };

  const toggleTheme = () => {
    theme === 'light' ? setTheme('dark') : setTheme('light');
  };

  // EFFECTS

  useEffect(() => {  
    const handleClickOutside = (event) => {
      // Remove gender and married dropdown logic
    };
  
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  // Update state when userData changes
  useEffect(() => {
    if (userData) {
      setFirstName(userData.firstName || '');
      setLastName(userData.lastName || '');
      setEmail(userData.email || '');
    }
  }, [userData]);

  // RENDER HELPER COMPONENTS

  const FormField = ({ className, label, value, setValue, required = false, disabled = false }) => (
    <div className={className}>
      <p className="search-icon-text">{label} {required && '*'}</p>
      {editMode && !disabled ? (
        <input
          type="text"
          value={value}
          onChange={(e) => setValue(e.target.value)}
          className="search-icon-info editable"
        />
      ) : (
        <p className="search-icon-info">{value}</p>
      )}
    </div>
  );

  const ProfileCard = () => (
    <div className="profile-card">
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
      <p className="profile-card-username">{usernameState}</p>
      <p className="profile-card-ID">ID: {userID}</p>
      <img
        src={frame}
        alt=""
        className="qr-code"
      />
      <p className="profile-card-join-date">S-a alăturat în mai 2025</p>
    </div>
  );

  const PersonalInformation = () => (
    <div className={`general-information ${editMode ? 'edit-mode' : ''}`}>
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

      <img
        src={theme === 'light' ? editPencilDay : editPencilNight}
        alt="Edit"
        className="edit-icon"
        onClick={() => {
          if (!editMode) setShowEditModal(true);
        }}
        style={{ cursor: editMode ? 'default' : 'pointer', opacity: editMode ? 0.5 : 1 }}
      />

      <FormField
        className="search-icon-first-name"
        label="Prenume"
        value={firstNameState}
        setValue={setFirstName}
        required={true}
      />

      <FormField
        className="search-icon-last-name"
        label="Nume"
        value={lastNameState}
        setValue={setLastName}
        required={true}
      />

      <FormField
        className="search-icon-phone-number"
        label="Număr de telefon"
        value={phoneNumber}
        setValue={setPhoneNumber}
        required={true}
      />

      <FormField
        className="search-icon-email"
        label="Adresă de email"
        value={emailState}
        setValue={setEmail}
        required={true}
      />

      <FormField
        className="search-icon-gender"
        label="Gen"
        value={gender}
        setValue={setGender}
      />

      <FormField
        className="search-icon-married"
        label="Căsătorit/ă"
        value={marriedStatus}
        setValue={setMarriedStatus}
      />

      <FormField
        className="search-icon-county"
        label="Județ"
        value={county}
        setValue={setCounty}
        required={true}
      />

      <FormField
        className="search-icon-address"
        label="Adresă"
        value={address}
        setValue={setAddress}
        required={true}
      />

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
  );

  const ActionButtons = () => (
    <>
      <div className="logout-button" onClick={() => setShowLogoutModal(true)}>
        <img src={logoutDay} alt="Logout" className="logout-icon" />
        <p className="logout-title">Ieșire din cont</p>
      </div>

      <div className="settings-button" onClick={() => setShowSettingsModal(true)}>
        <img src={settingsDay} alt="" className="settings-icon" />
        <p className="settings-title">Setări</p>
      </div>

      <a
        href="https://www.bcr.ro/en/about-us/our-company"
        target="_blank"
        rel="noopener noreferrer"
        style={{ textDecoration: 'none' }}
      >
        <div className="about-button">
          <img src={bankDay} alt="" className="about-icon" />
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
          <img src={helpDay} alt="" className="help-icon" />
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
          <img src={userAgreementDay} alt="" className="user-agreement-icon" />
          <p className="user-agreement-title">Condiții</p>
        </div>
      </a>
    </>
  );

  const Modals = () => (
    <>
      {showSettingsModal && (
        <SettingsModal onClose={() => setShowSettingsModal(false)} userData={userData} setLoggedIn={setLoggedIn}/>
      )}

      {showLogoutModal && (
        <LogoutModal onClose={() => setShowLogoutModal(false)} onConfirm={handleLogout} />
      )}

      {showEditModal && (
        <EditModal onClose={() => setShowEditModal(false)} onConfirm={handleEditConfirm} />
      )}

      {showDeleteModal && (
        <DeleteModal onClose={() => setShowDeleteModal(false)} onConfirm={handleDeleteAccount} />
      )}
    </>
  );

  // MAIN RENDER

  return (
    <div className="profile-container">
      {editMode && <div className="edit-overlay"></div>}
      
      <ProfileCard />
      <PersonalInformation />
      <ActionButtons />
      <Modals />
    </div>
  );
};

export default Profile;
