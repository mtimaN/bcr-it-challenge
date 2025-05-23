import React from 'react';
import './LogoutModal.css';

const LogoutModal = ({ onClose, onConfirm, lang, theme }) => {
  return (
    <div className={`logout-modal-overlay ${theme === 'dark' ? 'dark' : ''}`} onClick={onClose}>
      <div
        className={`logout-modal-container ${theme === 'dark' ? 'dark' : ''}`}
        onClick={(e) => e.stopPropagation()}
      >
        <h2 className={`logout-modal-title ${theme === 'dark' ? 'dark' : ''}`}>
          {lang === 'RO' ? 'Confirmare ieșire din cont' : 'Confirm logout'}
        </h2>
        <p className={`logout-modal-body ${theme === 'dark' ? 'dark' : ''}`}>
          {lang === 'RO' ? 'Sunteți sigur/ă că doriți să ieșiți din cont?' : 'Are you sure you want to log out?'}
        </p>
        <div className="logout-modal-buttons">
          <button className="logout-modal-confirm" onClick={onConfirm}>
            {lang === 'RO' ? 'Ieșire din cont' : 'Log out'}
          </button>
          <button className="logout-modal-cancel" onClick={onClose}>
            {lang === 'RO' ? 'Anulare' : 'Cancel'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default LogoutModal;
