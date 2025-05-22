import React, { useState } from 'react';
import './SettingsModal.css';

const SettingsModal = ({ onClose, userData }) => {
  const [step, setStep] = useState(1);
  const [oldPasswordInput, setOldPasswordInput] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const { password = '' } = userData || {};

  const handleVerifyOldPassword = () => {
    if (oldPasswordInput === password) {
      setStep(2);
    } else {
      alert('Parola veche introdusă este incorectă.');
    }
  };

  const handlePasswordChange = async () => {
    if (!newPassword.trim()) {
      alert('Introduceți noua parolă');
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
      } else {
        const errorData = await response.json();
        console.error('Failed to update password:', errorData.error || 'Unknown error');
        alert('Eroare la schimbarea parolei. Încercați din nou.');
      }
    } catch (error) {
      console.error('Error updating password:', error.message);
      alert('A apărut o eroare. Încercați din nou.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="settings-modal-overlay" onClick={onClose}>
      <div className="settings-modal-container" onClick={(e) => e.stopPropagation()}>
        <h2 className="settings-modal-title">
          {step === 1 ? 'Confirmare parolă actuală' : 'Schimbare parolă'}
        </h2>
        <p className="settings-modal-body">
          {step === 1
            ? 'Introduceți parola actuală pentru a continua.'
            : 'Introduceți noua parolă.'}
        </p>

        {step === 1 ? (
          <input
            type="password"
            value={oldPasswordInput}
            onChange={(e) => setOldPasswordInput(e.target.value)}
            placeholder="Parola veche"
            className="settings-password-input"
            disabled={isLoading}
          />
        ) : (
          <input
            type="password"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            placeholder="Noua parolă"
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
              ? 'Se salvează...'
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
