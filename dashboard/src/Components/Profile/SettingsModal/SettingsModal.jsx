import React, { useState } from 'react';
import './SettingsModal.css';
import { useNavigate } from 'react-router-dom';

const SettingsModal = ({ onClose, userData, setLoggedIn, lang, theme }) => {
  const navigate = useNavigate();
  const [step, setStep] = useState(1);
  const [oldPasswordInput, setOldPasswordInput] = useState('');

  const [newPassword, setNewPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const { password = '' } = userData || {};

  const handleVerifyOldPassword = () => {
    if (oldPasswordInput === password) {
      setStep(2);
    } else {
      {lang === 'RO' ? alert('Parola veche introdusă este incorectă.') : alert('Old password is incorrect.')};   
    }
  };

  const handlePasswordChange = async () => {
    if (!newPassword.trim()) {
      {lang === 'RO' ? alert('Introduceți noua parolă') : alert('Please enter your new password')};
      return;
    }

    setIsLoading(true);
    try {
      const token = localStorage.getItem('jwtToken');

      const response = await fetch('https://localhost:8443/v1/update', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          password: oldPasswordInput,
          new_password: newPassword
        })
      });

      if (response.ok) {
        alert('Parola a fost schimbată cu succes!');
        onClose();

        localStorage.removeItem('jwtToken');
        setLoggedIn(false);
        onClose();
        navigate('/');

      } else {
        const errorData = await response.json();
        console.error('Failed to update password:', errorData.error || 'Unknown error');
        {lang === 'RO' ? alert('Eroare la schimbarea parolei. Încercați din nou.') : alert('Error while changing the password. Please try again.')};
      }
    } catch (error) {
      console.error('Error updating password:', error.message);
      {lang === 'RO' ? alert('A apărut o eroare. Încercați din nou.') : alert('An error has occured. Please try again.')};
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="settings-modal-overlay" onClick={onClose}>
      <div className="settings-modal-container" onClick={(e) => e.stopPropagation()}>
        <h2 className="settings-modal-title">
          {step === 1 ? (lang === 'RO' ? 'Confirmare parolă actuală' : 'Confirm current password') : (lang === 'RO' ? 'Schimbare parolă' : 'Change password')}
        </h2>
        <p className="settings-modal-body">
          {step === 1
            ? (lang === 'RO' ? 'Introduceți parola actuală pentru a continua.' : 'Insert current password to continue')
            : (lang === 'RO' ? 'Introduceți noua parolă.' : 'Insert new password.')}
        </p>

        {step === 1 ? (
          <input
            type="password"
            value={oldPasswordInput}
            onChange={(e) => setOldPasswordInput(e.target.value)}
            placeholder={lang === 'RO' ? 'Parola veche' : 'Old password'}
            className="settings-password-input"
            disabled={isLoading}
          />
        ) : (
          <input
            type="password"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            placeholder={lang === 'RO' ? 'Noua parolă' : 'New password'}
            className="settings-password-input"
            disabled={isLoading}
          />
        )}

        <div className="settings-modal-buttons">
          <button
            className="settings-modal-confirm"
            onClick={step === 1 ? handleVerifyOldPassword : handlePasswordChange}
            disabled={isLoading}
          >
            {isLoading
              ? (lang === 'RO' ? 'Se salvează...' : 'Saving...')
              : step === 1
              ? 'Continuă'
              : 'Salvează'}
          </button>
          <button
            className="settings-modal-cancel"
            onClick={onClose}
            disabled={isLoading}
          >
            Anulare
          </button>
        </div>
      </div>
    </div>
  );
};

export default SettingsModal;
