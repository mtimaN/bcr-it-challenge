import React from 'react';
import './LogoutModal.css';

const LogoutModal = ({ onClose, onConfirm }) => {
  return (
    <div className="logout-modal-overlay" onClick={onClose}>
      <div className="logout-modal-container" onClick={(e) => e.stopPropagation()}>
        <h2 className="logout-modal-title">Confirmare ieșire din cont</h2>
        <p className="logout-modal-body">Sunteți sigur/ă că doriți să ieșiți din cont?</p>
        <div className="logout-modal-buttons">
          <button className="logout-modal-confirm" onClick={onConfirm}>
            Ieșire din cont
          </button>
          <button className="logout-modal-cancel" onClick={onClose}>Anulare</button>
        </div>
      </div>
    </div>
  );
};

export default LogoutModal;
