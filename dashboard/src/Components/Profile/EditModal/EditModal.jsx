import React, { useState } from 'react';
import './EditModal.css';

const EditModal = ({ onClose, onConfirm }) => {
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = () => {
    if (password.trim() === '') {
      setError('Please enter a password.');
    } else {
      const success = onConfirm(password); // get success/failure
      if (success === false) {
        setError('Wrong password. Please try again.');
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
        <h2 className="edit-modal-title">Confirmare identitate</h2>
        <p className="edit-modal-body">Introduceți parola pentru a continua:</p>

        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          onKeyDown={handleKeyDown}
          className="edit-modal-input"
        />

        {error && <p className="edit-modal-error">{error}</p>}

        <div className="edit-modal-buttons">
          <button className="confirm" onClick={handleSubmit}>Confirmare</button>
          <button className="cancel" onClick={onClose}>Anulare</button>
        </div>
      </div>
    </div>
  );
};

export default EditModal;