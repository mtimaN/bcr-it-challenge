import React, { useState } from 'react';
import './EditModal.css';

const EditModal = ({ onClose, onConfirm, lang, theme }) => {
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = () => {
    if (password.trim() === '') {
      {lang === 'RO' ? setError('Vă rugăm introduceți o parolă.') : setError('Please enter a valid password.')};
    } else {
      const success = onConfirm(password); // get success/failure
      if (success === false) {
        {lang === 'RO' ? setError('Parolă greșită. Vă rugăm încercați din nou.') : setError('Wrong password. Please try again.')};
        setPassword('');
      }
    }
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      handleSubmit();
    }
  };

  return (
    <div className="edit-modal-overlay" onClick={onClose}>
      <div className="edit-modal-container" onClick={(e) => e.stopPropagation()}>
        <h2 className="edit-modal-title"> {lang === 'RO' ? 'Confirmare identitate' : 'Confirm identity'}</h2>
        <p className="edit-modal-body">{lang === 'RO' ? 'Introduceți parola pentru a continua:' : 'To continue, please enter your password:'}</p>

        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          onKeyDown={handleKeyDown}
          className="edit-modal-input"
        />

        {error && <p className="edit-modal-error">{error}</p>}

        <div className="edit-modal-buttons">
          <button className="confirm" onClick={handleSubmit}>{lang === 'RO' ? 'Confirmare' : 'Confirm'}</button>
          <button className="cancel" onClick={onClose}>{lang === 'RO' ? 'Anulare' : 'Cancel'}</button>
        </div>
      </div>
    </div>
  );
};

export default EditModal;
